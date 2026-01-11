package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewDynamoDB() *dynamodb.DynamoDB {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return db
}
