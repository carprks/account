package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/carprks/account/service"
)

func main() {
	lambda.Start(service.Handler)
}
