package main

import (
	"github.com/carprks/account/service"
)

func main() {
	lambda.Start(service.Handler)
}
