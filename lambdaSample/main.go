package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func loadAreaFromZip(zipCode string) {
	fmt.Println("------")
	fmt.Println(zipCode)
	fmt.Println("------")
}

//urlのパターンは　/area/274-0077(郵便番号)
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
	zipCode := request.PathParameters["zipCode"]
	loadAreaFromZip(zipCode)
	return "", nil
}

func main() {
	lambda.Start(handler)
}
