package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type resp struct {
	ID string `json:"id"`
}

type Remove struct {
	ID string `json:"id"`
}

func removeAccount(ident string) error {
	r := Remove{
		ID: ident,
	}

	err := r.removeLogin()
	if err != nil {
		fmt.Println(fmt.Sprintf("Remove Login: %v", err))
		return err
	}

	err = r.removePermission()
	if err != nil {
		fmt.Println(fmt.Sprintf("Remove Permissions: %v", err))
		return err
	}

	return nil
}

func (r Remove) removeLogin() error {
	j, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/remove", os.Getenv("SERVICE_LOGIN")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("login remove req err: %v", err))
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("login remove client err: %v", err))
	}
	defer resp.Body.Close()

	return err
}

func (r Remove) removePermission() error {
	j, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/permissions/%s", os.Getenv("SERVICE_PERMISSIONS"), r.ID), nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("permissions remove req err: %v", err))
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("permissions remove client err: %v", err))
	}
	defer resp.Body.Close()

	return err
}

// func TestHandler(t *testing.T) {
// 	if len(os.Args) >= 1 {
// 		for _, env := range os.Args {
// 			if env == "localDev" {
// 				err := godotenv.Load()
// 				if err != nil {
// 					fmt.Println(fmt.Sprintf("godotenv err: %v", err))
// 				}
// 			}
// 		}
// 	}
//
// 	tests := []struct {
// 		request events.APIGatewayProxyRequest
// 		expect  events.APIGatewayProxyResponse
// 		err     error
// 	}{
// 		{
// 			request: events.APIGatewayProxyRequest{
// 				Resource: "/register",
// 				Body:     `{"email":"tester@carpark.ninja","password":"tester","verify":"tester"}`,
// 			},
// 			expect: events.APIGatewayProxyResponse{
// 				StatusCode: 200,
// 				Body:       `{"id":"5f46cf19-5399-55e3-aa62-0e7c19382250","email":"tester@carpark.ninja","permissions":[]}`,
// 			},
// 			err: nil,
// 		},
// 	}
//
// 	for _, test := range tests {
// 		response, err := service.Handler(test.request)
// 		passed := assert.IsType(t, test.err, err)
// 		if !passed {
// 			fmt.Println(fmt.Sprintf("service test err: %v, request: %v", err, test.request))
// 		}
// 		assert.Equal(t, test.expect, response)
//
// 		s := resp{}
// 		err = json.Unmarshal([]byte(response.Body), &s)
// 		if err != nil {
// 			t.Fail()
// 		}
// 	}
// }
