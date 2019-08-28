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

func TestRegister(t *testing.T) {
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
		request login.RegisterRequest
		expect  service.RegisterObject
		err     error
	}{
		{
			request: login.RegisterRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
				Verify:   "tester",
			},
			expect: service.RegisterObject{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Email:      "tester@carpark.ninja",
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
		response, err := service.Register(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("register create register test err: %v, request: %v", err, test.request))
		}
		passed = assert.Equal(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("register create register test not equal: %v", test.request))
		}

		r := Remove{
			Identifier: test.expect.Identifier,
		}
		err = r.deleteLogin()
		if err != nil {
			fmt.Println(fmt.Sprintf("delete login: %v", err))
		}
		err = r.deletePermission()
		if err != nil {
			fmt.Println(fmt.Sprintf("delete permission: %v", err))
		}
	}
}

func TestCreatePermissions(t *testing.T) {
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
		request login.Register
		expect  []permissions.Permission
		err     error
	}{
		{
			request: login.Register{
				Identifier: "tester",
				Email:      "tester@carpark.ninja",
			},
			expect: []permissions.Permission{
				{
					Name:       "account",
					Action:     "login",
					Identifier: "tester",
				},
				{
					Name:       "carparks",
					Action:     "book",
					Identifier: "*",
				},
			},
		},
	}

	for _, test := range tests {
		response, err := service.CreatePermissions(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("register create permissions test err: %v, request: %v", err, test.request))
		}
		passed = assert.Equal(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("register create permissions test not equal: %v", test.request))
		}

		r := Remove{
			Identifier: test.request.Identifier,
		}
		err = r.deletePermission()
		if err != nil {
			fmt.Println(fmt.Sprintf("delete permission: %v", err))
		}
	}
}

func TestCreateLogin(t *testing.T) {
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
		request login.RegisterRequest
		expect  login.Register
		err     error
	}{
		{
			request: login.RegisterRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
				Verify:   "tester",
			},
			expect: login.Register{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Email:      "tester@carpark.ninja",
				Error:      "",
			},
		},
	}

	for _, test := range tests {
		response, err := service.CreateLogin(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("register create login test err: %v, request: %v", err, test.request))
		}

		passed = assert.Equal(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("register create login test not equal: %v", test.request))
		}

		r := Remove{
			Identifier: test.expect.Identifier,
		}
		err = r.deleteLogin()
		if err != nil {
			fmt.Println(fmt.Sprintf("delete login: %v", err))
		}
	}
}
