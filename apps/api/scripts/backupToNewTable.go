package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	client := dynamodb.NewFromConfig(cfg)

	source := "sidekick-temp"
	target := "SidekickApi-ApiWithDynamoSidekickBD73AF7B-RKXQLMN63VYD"

	var lastKey map[string]types.AttributeValue
	for {
		out, _ := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:         &source,
			ExclusiveStartKey: lastKey,
		})

		for i := 0; i < len(out.Items); i += 25 {
			end := min(i+25, len(out.Items))
			requests := make([]types.WriteRequest, 0, end-i)
			for _, item := range out.Items[i:end] {
				requests = append(requests, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
			}
			client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{target: requests},
			})
		}

		if out.LastEvaluatedKey == nil {
			break
		}
		lastKey = out.LastEvaluatedKey
	}
	fmt.Println("Done")
}
