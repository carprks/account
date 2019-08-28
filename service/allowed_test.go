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

func TestAllowed(t *testing.T) {
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
		request permissions.Permissions
		expect  permissions.Permissions
		err     error
	}{
		{
			request: permissions.Permissions{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Permissions: []permissions.Permission{
					{
						Name:   "account",
						Action: "login",
					},
				},
			},
			expect: permissions.Permissions{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Status:     "allowed",
			},
		},
		{
			request: permissions.Permissions{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Permissions: []permissions.Permission{
					{
						Name:   "carparks",
						Action: "create",
					},
				},
			},
			expect: permissions.Permissions{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Status:     "denied",
			},
		},
	}

	for _, test := range tests {
		response, err := service.Allowed(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("verify test type err: %v, request: %v", err, test.request))
		}
		passed = assert.Equal(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("verify test value response: %v, request: %v", response, test.request))
		}
	}

	deleteAccount(resp.Identifier)
}
