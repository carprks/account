package service_test

import (
	"fmt"
	"github.com/carprks/account/service"
	login "github.com/carprks/login/service"
	permissions "github.com/carprks/permissions/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
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

	r := login.RegisterRequest{
		Email:    "tester@carpark.ninja",
		Password: "tester",
		Verify:   "tester",
	}
	resp, err := service.Register(r)
	if err != nil {
		fmt.Println(fmt.Sprintf("login register failed: %v", err))
	}

	tests := []struct {
		request login.LoginRequest
		expect  service.LoginObject
		err     error
	}{
		{
			request: login.LoginRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
			},
			expect: service.LoginObject{
				Identifier: resp.Identifier,
				Permissions: []permissions.Permission{
					{
						Name:       "account",
						Action:     "login",
						Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
					},
					{
						Name:       "carparks",
						Action:     "book",
						Identifier: "*",
					},
				},
			},
		},
	}

	for _, test := range tests {
		response, err := service.Login(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("login test err: %v, request: %v", err, test.request))
		}
		passed = assert.Equal(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("login test not equal: %v", test.request))
		}
	}

	deleteAccount(resp.Identifier)
}

func TestLoginUser(t *testing.T) {
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
}

func TestLoginPermissions(t *testing.T) {
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
}
