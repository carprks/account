package service

import (
	"bytes"
	"encoding/json"
	"errors"
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
		return "", err
	}

	rf, err := Register(r)
	if err != nil {
		return "", err
	}

	rfb, err := json.Marshal(rf)
	if err != nil {
		return "", err
	}

	return string(rfb), nil
}

// Register underlying functions
func Register(r login.RegisterRequest) (RegisterObject, error) {
	ro, err := CreateLogin(r)
	if err != nil {
		return RegisterObject{}, err
	}

	resp, err := CreatePermissions(ro)
	if err != nil {
		return RegisterObject{}, err
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
		return rr, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/register", os.Getenv("SERVICE_LOGIN")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("client err: %v", err))
		return rr, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("account resp err: %v", err))
			return rr, err
		}

		rt := login.Register{}

		err = json.Unmarshal(body, &rt)
		if err != nil {
			return rt, err
		}

		if rt.Error != "" {
			if strings.Contains(rt.Error, "ErrCodeConditionalCheckFailedException") {
				return rt, errors.New("login already exists")
			}
			return rt, errors.New(rt.Error)
		}

		return rt, nil
	}

	return rr, errors.New("can't create login")
}

// CreatePermissions ...
func CreatePermissions(r login.Register) ([]permissions.Permission, error) {
	perms := []permissions.Permission{
		{
			Identifier: r.Identifier,
			Name:       "account",
			Action:     "login",
		},
		{
			Identifier: "*",
			Name:       "carparks",
			Action:     "book",
		},
	}

	p := permissions.Permissions{
		Identifier:  r.Identifier,
		Permissions: perms,
	}

	j, err := json.Marshal(&p)
	if err != nil {
		return []permissions.Permission{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/create", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("client err: %v", err))
		return []permissions.Permission{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return perms, nil
	}

	return []permissions.Permission{}, errors.New("can't create permissions")
}
