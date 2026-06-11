//go:build lambda

package main

import (
	"context"
	"funkey-grab-and-bite/funkey-bite-api/internal/app"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Println("Initializing Funkey Grab-and-Bite API in AWS LAMBDA Mode...")

	// Initialize the exact same engine setup used locally
	router, _ := app.SetupEngine()

	// Wrap the Gin engine inside the API Gateway proxy adapter
	ginLambda = ginadapter.New(router)
}

// Handler intercepts API Gateway proxy events and routes them through Gin
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Hand off execution control to the AWS serverless runtime
	lambda.Start(Handler)
}
