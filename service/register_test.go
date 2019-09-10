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

func TestRegister(t *testing.T) {
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

	tests := []struct {
		name string
		request login.RegisterRequest
		expect  service.RegisterObject
		err     error
	}{
		{
			name: "register: tester@carpark.ninja",
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
						Name:       "account",
						Action:     "edit",
						Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
					},
					{
						Name:       "account",
						Action:     "view",
						Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
					},
					{
						Name:       "payments",
						Action:     "create",
						Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
					},
					{
						Name:       "payments",
						Action:     "view",
						Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
					},
					{
						Name:       "payments",
						Action:     "report",
						Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
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
			response, err := service.Register(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("register create register test err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("register create register test not equal: %v", test.request)
			}

			r := Remove{
				Identifier: test.expect.Identifier,
			}
			err = r.deleteLogin()
			if err != nil {
				t.Errorf("delete login: %w", err)
			}
			err = r.deletePermission()
			if err != nil {
				t.Errorf("delete permission: %w", err)
			}
		})
	}
}

func TestCreatePermissions(t *testing.T) {
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

	tests := []struct {
		name string
		request login.Register
		expect  []permissions.Permission
		err     error
	}{
		{
			name: "create permissions: tester@carpark.ninja",
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
					Name:       "account",
					Action:     "edit",
					Identifier: "tester",
				},
				{
					Name:       "account",
					Action:     "view",
					Identifier: "tester",
				},
				{
					Name:       "payments",
					Action:     "create",
					Identifier: "tester",
				},
				{
					Name:       "payments",
					Action:     "view",
					Identifier: "tester",
				},
				{
					Name:       "payments",
					Action:     "report",
					Identifier: "tester",
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.CreatePermissions(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("register create permissions test err: %w, request: %v", err, test.request)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("register create permissions test not equal: %v", test.request)
			}

			r := Remove{
				Identifier: test.request.Identifier,
			}
			err = r.deletePermission()
			if err != nil {
				t.Errorf("delete permission: %w", err)
			}
		})
	}
}

func TestCreateLogin(t *testing.T) {
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

	tests := []struct {
		name string
		request login.RegisterRequest
		expect  login.Register
		err     error
	}{
		{
			name: "create login tester@carpark.ninia",
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
		t.Run(test.name, func(t *testing.T) {
			response, err := service.CreateLogin(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("register create login test err: %w, request: %v", err, test.request)
			}

			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("register create login test not equal: %v", test.request)
			}

			r := Remove{
				Identifier: test.expect.Identifier,
			}
			err = r.deleteLogin()
			if err != nil {
				t.Errorf("delete login: %w", err)
			}
		})
	}
}
