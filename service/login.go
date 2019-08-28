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
)

// LoginObject ...
type LoginObject struct {
	Identifier  string                   `json:"identifier"`
	Permissions []permissions.Permission `json:"permissions"`
}

// LoginHandler ...
func LoginHandler(body string) (string, error) {
	r := login.LoginRequest{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't unmarshall login: %v, %v", err, body))
		return "", err
	}

	rf, err := Login(r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't get login: %v, %v", err, r))
		return "", err
	}

	rfb, err := json.Marshal(rf)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't marshall login: %v, %v", err, rf))
		return "", err
	}

	return string(rfb), nil
}

// Login ...
func Login(l login.LoginRequest) (LoginObject, error) {
	lo, err := LoginUser(l)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't get login for user: %v, %v", err, l))
		return LoginObject{}, err
	}

	resp, err := LoginPermissions(lo)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't get permissions for user: %v, %v", err, lo))
		return LoginObject{}, err
	}

	return LoginObject{
		Identifier:  lo.Identifier,
		Permissions: resp,
	}, nil
}

// LoginUser ...
func LoginUser(l login.LoginRequest) (login.Login, error) {
	lr := login.Login{}

	j, err := json.Marshal(&l)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't unmarshall login: %v, %v", err, l))
		return lr, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", os.Getenv("SERVICE_LOGIN")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
		return lr, err
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("login client err: %v", err))
		return lr, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("login resp err: %v", err))
			return lr, err
		}

		err = json.Unmarshal(body, &lr)
		if err != nil {
			fmt.Println(fmt.Sprintf("can't unmarshal login service: %v, %v", err, string(body)))
			return lr, err
		}

		if lr.Error != "" {
			return lr, errors.New(lr.Error)
		}
	}

	return lr, nil
}

// LoginPermissions ...
func LoginPermissions(l login.Login) ([]permissions.Permission, error) {
	p := permissions.Permissions{
		Identifier: l.Identifier,
	}

	j, err := json.Marshal(&p)
	if err != nil {
		return []permissions.Permission{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/retrieve", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
		return p.Permissions, err
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("permissions client err: %v", err))
		return p.Permissions, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("permissions resp err: %v", err))
			return p.Permissions, err
		}

		err = json.Unmarshal(body, &p)
		if err != nil {
			fmt.Println(fmt.Sprintf("can't unmarshall permissions body: %v, %v", err, string(body)))
			return p.Permissions, err
		}

		if p.Status != "" {
			return p.Permissions, errors.New(p.Status)
		}
	}

	return p.Permissions, nil
}
