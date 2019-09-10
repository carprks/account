package service_test

import (
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
					t.Errorf("godotenv err: %w", err)
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
		t.Errorf("login register failed: %w", err)
	}

	tests := []struct {
		name string
		request login.LoginRequest
		expect  service.LoginObject
		err     error
	}{
		{
			name: "login: tester@carpark.ninja",
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
						Identifier: resp.Identifier,
					},
					{
						Name:       "account",
						Action:     "edit",
						Identifier: resp.Identifier,
					},
					{
						Name:       "account",
						Action:     "view",
						Identifier: resp.Identifier,
					},
					{
						Name:       "payments",
						Action:     "create",
						Identifier: resp.Identifier,
					},
					{
						Name:       "payments",
						Action:     "view",
						Identifier: resp.Identifier,
					},
					{
						Name:       "payments",
						Action:     "report",
						Identifier: resp.Identifier,
					},
					{
						Name:       "carparks",
						Action:     "book",
						Identifier: "*",
					},
					{
						Name:       "carparks",
						Action:     "report",
						Identifier: "*",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.Login(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("login test err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("login test not equal: %v", test.request)
			}
		})
	}

	deleteAccount(resp.Identifier)
}

func TestLoginUser(t *testing.T) {
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
}

func TestLoginPermissions(t *testing.T) {
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
}
