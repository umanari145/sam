package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() {

	fmt.Println("morigami")
	os.Exit(0)
}

func main() {
	lambda.Start(handler)
}
