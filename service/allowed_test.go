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

var testsAllowed = []struct {
	name    string
	request permissions.Permissions
	expect  permissions.Permissions
	err     error
}{
	{
		name: "allowed: account-login",
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
		name: "allowed: carparks-create",
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

func TestAllowed(t *testing.T) {
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

	for _, test := range testsAllowed {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.Allowed(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("verify test type err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("verify test value response: %v, request: %v", response, test.request)
			}
		})
	}

	deleteAccount(resp.Identifier)
}

func BenchmarkAllowed(b *testing.B) {
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

	r := login.RegisterRequest{
		Email:    "tester@carpark.ninja",
		Password: "tester",
		Verify:   "tester",
	}
	resp, err := service.Register(r)
	if err != nil {
		b.Errorf("login register failed: %w", err)
	}

	b.ResetTimer()
	for _, test := range testsAllowed {
		b.Run(test.name, func(t *testing.B) {
			t.StopTimer()

			response, err := service.Allowed(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("verify test type err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("verify test value response: %v, request: %v", response, test.request)
			}

			t.StartTimer()
		})
	}

	deleteAccount(resp.Identifier)
}
