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
	"time"
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
		return fmt.Errorf("delete account login: %w", err)
	}

	err = r.deletePermission()
	if err != nil {
		return fmt.Errorf("delete account permssions: %w", err)
	}

	return nil
}

func (r Remove) deleteLogin() error {
	j, err := json.Marshal(&r)
	if err != nil {
		return fmt.Errorf("delete login marshall: %w", err)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", os.Getenv("SERVICE_LOGIN")), bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("login remove req err: %w", err)
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_LOGIN"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     2 * time.Minute,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("login remove client err: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (r Remove) deletePermission() error {
	j, err := json.Marshal(&r)
	if err != nil {
		return fmt.Errorf("delete permissions marshall: %w", err)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("permissions remove req err: %w", err)
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_PERMISSIONS"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     2 * time.Minute,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("permissions remove client err: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

var testsService = []struct {
	name    string
	request events.APIGatewayProxyRequest
	expect  events.APIGatewayProxyResponse
	err     error
}{
	// Register
	{
		name: "register success",
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
		name: "register failed",
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
		name: "login success",
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
		name: "login failed password",
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
		name: "login failed identity",
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
		name: "allowed success",
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
		name: "allowed failed",
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

func TestHandler(t *testing.T) {
	if len(os.Args) >= 1 {
		for _, env := range os.Args {
			if env == "localDev" {
				err := godotenv.Load()
				if err != nil {
					t.Errorf("godotenv err: %w", err)
				}
			}
		}
	}

	for _, test := range testsService {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.Handler(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("service test type err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("service test value response: %v, request: %v", response, test.request)
			}
		})
	}

	deleteAccount("5f46cf19-5399-55e3-aa62-0e7c19382250")
}

func BenchmarkHandler(b *testing.B) {
	b.ReportAllocs()

	if len(os.Args) >= 1 {
		for _, env := range os.Args {
			if env == "localDev" {
				err := godotenv.Load()
				if err != nil {
					b.Errorf("godotenv err: %w", err)
				}
			}
		}
	}

	b.ResetTimer()
	for _, test := range testsService {
		b.Run(test.name, func(t *testing.B) {
			b.StopTimer()

			response, err := service.Handler(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("service test type err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("service test value response: %v, request: %v", response, test.request)
			}

			t.StartTimer()
		})
	}

	deleteAccount("5f46cf19-5399-55e3-aa62-0e7c19382250")
}
