package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/carprks/account/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

type resp struct {
	Identifier string `json:"identifier"`
	Status     string `json:"status"`
}

type Remove struct {
	Identifier string `json:"identifier"`
}

func deleteAccount(ident string) error {
	r := Remove{
		Identifier: ident,
	}

	err := r.deleteLogin()
	if err != nil {
		fmt.Println(fmt.Sprintf("Remove Login: %v", err))
		return err
	}

	err = r.deletePermission()
	if err != nil {
		fmt.Println(fmt.Sprintf("Remove Permissions: %v", err))
		return err
	}

	return nil
}

func (r Remove) deleteLogin() error {
	j, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", os.Getenv("SERVICE_LOGIN")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("login remove req err: %v", err))
		return err
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_LOGIN"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("login remove client err: %v", err))
	}
	defer resp.Body.Close()

	return err
}

func (r Remove) deletePermission() error {
	j, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("permissions remove req err: %v", err))
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_PERMISSIONS"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("permissions remove client err: %v", err))
		return err
	}
	defer resp.Body.Close()

	return err
}

func TestHandler(t *testing.T) {
	if len(os.Args) >= 1 {
		for _, env := range os.Args {
			if env == "localDev" {
				err := godotenv.Load()
				if err != nil {
					fmt.Println(fmt.Sprintf("godotenv err: %v", err))
				}
			}
		}
	}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  events.APIGatewayProxyResponse
		err     error
	}{
		// Register
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/register",
				Body:     `{"email":"tester@carpark.ninja","password":"tester","verify":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250","email":"tester@carpark.ninja","permissions":[{"name":"account","action":"login","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"account","action":"edit","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"account","action":"view","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"payments","action":"create","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"payments","action":"view","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"payments","action":"report","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"carparks","action":"book","identifier":"*"},{"name":"carparks","action":"report","identifier":"*"}]}`,
			},
			err: nil,
		},
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/register",
				Body:     `{"email":"tester@carpark.ninja","password":"tester","verify":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `can't register: can't get create login: login already exists`,
			},
		},

		// Login
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/login",
				Body:     `{"email":"tester@carpark.ninja","password":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250","permissions":[{"name":"account","action":"login","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"account","action":"edit","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"account","action":"view","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"payments","action":"create","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"payments","action":"view","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"payments","action":"report","identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250"},{"name":"carparks","action":"book","identifier":"*"},{"name":"carparks","action":"report","identifier":"*"}]}`,
			},
		},
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/login",
				Body:     `{"email":"tester@carpark.ninja","password":"failure"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `can't get login: can't get login for user: login response err: invalid password`,
			},
		},
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/login",
				Body:     `{"email":"failure@carpark.ninja","password":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `can't get login: can't get login for user: login response err: no identity`,
			},
		},

		// Allowed
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/allowed",
				Body:     `{"identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250","permissions":[{"action":"login","name":"account"}]}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250","status":"allowed"}`,
			},
		},
		{
			request: events.APIGatewayProxyRequest{
				Resource: "/allowed",
				Body:     `{"identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250","permissions":[{"name":"carparks","action":"create"}]}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"identifier":"5f46cf19-5399-55e3-aa62-0e7c19382250","status":"denied"}`,
			},
		},
	}

	for _, test := range tests {
		response, err := service.Handler(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("service test type err: %v, request: %v", err, test.request))
		}
		passed = assert.Equal(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("service test value response: %v, request: %v", response, test.request))
		}
	}

	deleteAccount("5f46cf19-5399-55e3-aa62-0e7c19382250")
}
