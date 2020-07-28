package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, snsEvent events.SNSEvent) {

	log.Println("ctx")
	log.Println(ctx)
	log.Println("snsEvent")
	log.Println(snsEvent)
	os.Exit(0)
}

func main() {
	lambda.Start(handler)
}
