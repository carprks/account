package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	permissions "github.com/carprks/permissions/service"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// AllowedHandler ...
func AllowedHandler(body string) (string, error) {
	r := permissions.Permissions{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't unmarshall input: %v, %v", err, body))
		return "", fmt.Errorf("can't unmarshall input: %w", err)
	}

	rf, err := Allowed(r)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't get allowed: %v, %v", err, r))
		return "", fmt.Errorf("can't get allowed: %w", err)
	}

	rfb, err := json.Marshal(rf)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't marshal allowed: %v, %v", err, rf))
		return "", fmt.Errorf("can't unmarshal allowed: %w", err)
	}

	return string(rfb), nil
}

// Allowed ...
func Allowed(p permissions.Permissions) (permissions.Permissions, error) {
	pr := permissions.Permissions{}

	j, err := json.Marshal(&p)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't unmarshal permissions: %v, %v", err, p))
		return pr, fmt.Errorf("can't unmarshal permissions: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/allowed", os.Getenv("SERVICE_PERMISSIONS")), bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(fmt.Sprintf("req err: %v", err))
		return pr, fmt.Errorf("req allowed err: %w", err)
	}

	req.Header.Set("X-Authorization", os.Getenv("AUTH_PERMISSIONS"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     2 * time.Minute,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("verify client err: %v", err))
		return pr, fmt.Errorf("verify client err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("verify resp err: %v", err))
			return pr, fmt.Errorf("verify resp err: %w", err)
		}

		err = json.Unmarshal(body, &pr)
		if err != nil {
			fmt.Println(fmt.Sprintf("can't unmarshall permissions: %v, %v", err, string(body)))
			return pr, fmt.Errorf("can't unmarshall permissions: %w", err)
		}

		return pr, nil
	}

	return pr, fmt.Errorf("allowed came back with a different statuscode: %v", resp.StatusCode)
}
