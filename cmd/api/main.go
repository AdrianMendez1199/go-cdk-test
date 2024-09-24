package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Hello, world!",
		StatusCode: 200,
	}, nil

}

func main() {
	// log.Printf("Event: %v", event)
	println("Hello, world! 2")
	lambda.Start(handler)
}
