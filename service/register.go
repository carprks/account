package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	login "github.com/carprks/login/service"
	permissions "github.com/carparks/permissions/service"
	"io/ioutil"
	"net/http"
	"os"
)

// RegisterObject ...
type RegisterObject struct {
	ID          string       `json:"id"`
	Email       string       `json:"email"`
	Permissions []Permission `json:"permissions"`
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
		ID:          ro.ID,
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
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("account resp err: %v", err))
			return rr, err
		}

		err = json.Unmarshal(body, &rr)
		if err != nil {
			return rr, err
		}
	}

	return rr, nil
}

// CreatePermissions ...
func CreatePermissions(r login.Register) ([]Permission, error) {
	// rr := perm.

	return []Permission{}, nil
}
