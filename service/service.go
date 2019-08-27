package service

import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

func retn() (string, error) {
	return "", nil
}

// Handler ...
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := retn()

	switch request.Resource {
	case "/login":
		resp, err = LoginHandler(request.Body)
	case "/register":
		resp, err = RegisterHandler(request.Body)
	case "/reset":
		resp, err = ResetHandler(request.Body)
	case "/verify":
		resp, err = VerifyHandler(request.Body)
	case "/allowed":
		resp, err = AllowedHandler(request.Body)
	}

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       resp,
	}, nil
}
