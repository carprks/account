package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	permissions "github.com/carprks/permissions/service"
	"io/ioutil"
	"net/http"
	"os"
)

// AllowedHandler ...
func AllowedHandler(body string) (string, error) {
	r := permissions.Permissions{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return "", err
	}

	rf, err := Allowed(r)
	if err != nil {
		return "", err
	}

	rfb, err := json.Marshal(rf)
	if err != nil {
		return "", err
	}

	return string(rfb), nil
}

// Allowed ...
func Allowed(p permissions.Permissions) (permissions.Permissions, error) {
	pr := permissions.Permissions{}

	j, err := json.Marshal(&p)
	if err != nil {
		return pr, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/allowed", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
		return pr, err
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_KEY"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("verify client err: %v", err))
		return pr, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("verify resp err: %v", err))
			return pr, err
		}

		err = json.Unmarshal(body, &pr)
		if err != nil {
			return pr, err
		}
	}

	return pr, nil
}
