package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	login "github.com/carprks/login/service"
	permissions "github.com/carprks/permissions/service"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// RegisterObject ...
type RegisterObject struct {
	Identifier  string                   `json:"identifier"`
	Email       string                   `json:"email"`
	Permissions []permissions.Permission `json:"permissions"`
}

// RegisterHandler what is used by service
func RegisterHandler(body string) (string, error) {
	r := login.RegisterRequest{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't unmarshall register: %v, %v", err, body))
		return "", fmt.Errorf("can't unmarshall register: %w", err)
	}

	rf, err := Register(r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't register: %v, %v", err, r))
		return "", fmt.Errorf("can't register: %w", err)
	}

	rfb, err := json.Marshal(rf)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't marshall register: %v, %v", err, rf))
		return "", fmt.Errorf("can't unmarshall register: %w", err)
	}

	return string(rfb), nil
}

// Register underlying functions
func Register(r login.RegisterRequest) (RegisterObject, error) {
	ro, err := CreateLogin(r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't get create login: %v, %v", err, r))
		return RegisterObject{}, fmt.Errorf("can't get create login: %w", err)
	}

	resp, err := CreatePermissions(ro)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't create permissions: %v, %v", err, ro))
		return RegisterObject{}, fmt.Errorf("can't create permissions: %w", err)
	}

	return RegisterObject{
		Identifier:  ro.Identifier,
		Email:       ro.Email,
		Permissions: resp,
	}, nil
}

// CreateLogin create the login
func CreateLogin(r login.RegisterRequest) (login.Register, error) {
	rr := login.Register{}

	j, err := json.Marshal(&r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't marshall register: %v, %v", err, r))
		return rr, fmt.Errorf("can't marshall register: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/register", os.Getenv("SERVICE_LOGIN")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
		return rr, fmt.Errorf("create login req err: %w", err)
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_LOGIN"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("client err: %v", err))
		return rr, fmt.Errorf("create login client err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("account resp err: %v", err))
			return rr, fmt.Errorf("account resp err: %w", err)
		}

		rt := login.Register{}

		err = json.Unmarshal(body, &rt)
		if err != nil {
			fmt.Println(fmt.Sprintf("can't unmarshal register body: %v, %v", err, string(body)))
			return rt, fmt.Errorf("can't unmarshal register body: %w", err)
		}

		if rt.Error != "" {
			if strings.Contains(rt.Error, "ErrCodeConditionalCheckFailedException") {
				return rt, fmt.Errorf("login already exists")
			}
			return rt, fmt.Errorf("unknown login error: %w", rt.Error)
		}

		return rt, nil
	}

	return rr, fmt.Errorf("can't create login")
}

// CreatePermissions ...
func CreatePermissions(r login.Register) ([]permissions.Permission, error) {
	p := permissions.Permissions{
		Identifier:  r.Identifier,
		Permissions: getDefaultPerms(r.Identifier),
	}

	j, err := json.Marshal(&p)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't marshall permissions: %v, %v", err, p))
		return []permissions.Permission{}, fmt.Errorf("can't unmarshall permissions: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/create", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
		return []permissions.Permission{}, fmt.Errorf("create permissions req err: %w", err)
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_PERMISSIONS"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("client err: %v", err))
		return []permissions.Permission{}, fmt.Errorf("create permissions client err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return getDefaultPerms(r.Identifier), nil
	}

	return []permissions.Permission{}, fmt.Errorf("can't create permissions")
}
