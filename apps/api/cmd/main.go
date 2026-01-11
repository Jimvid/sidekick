package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jimvid/sidekick/config"
	"github.com/jimvid/sidekick/internal/router"
)

var chiLambda *chiadapter.ChiLambda
var cfg = config.NewConfig()

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func init() {
	clerk.SetKey(cfg.CLERK_SECRET)
	router := router.NewRouter(cfg)
	chiLambda = chiadapter.New(router)
}

func main() {
	lambda.Start(handler)
}
