//go:build lambda

package main

import (
	"context"
	"funkey-grab-and-bite/funkey-bite-api/internal/app"
	"funkey-grab-and-bite/funkey-bite-api/internal/realtime"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

// LambdaSilentBroadcaster drops broadcast requests silently to prevent hanging loops in serverless space
type LambdaSilentBroadcaster struct{}

func (s *LambdaSilentBroadcaster) Broadcast(event string, data interface{}) {
	// In Lambda, real-time broadcasts are ignored internally to maintain ephemerality.
	// (Can be piped out to AWS EventBridge or API Gateway WebSockets here if needed later).
}

var ginLambda *ginadapter.GinLambda

func init() {
	log.Println("Initializing Funkey Grab-and-Bite API in AWS LAMBDA Mode...")

	// 1. Swap the stateful Hub with our lightweight, serverless-safe broadcaster stub
	realtime.GlobalBroadcaster = &LambdaSilentBroadcaster{}

	// 2. Initialize the exact same engine setup safely
	router, _ := app.SetupEngine()

	ginLambda = ginadapter.New(router)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
