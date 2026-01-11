package climbs

import (
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/jimvid/sidekick/config"
)

type ClimbStorage struct {
	db  *dynamodb.DynamoDB
	cfg *config.Config
}

func NewClimbStorage(db *dynamodb.DynamoDB, cfg *config.Config) *ClimbStorage {
	return &ClimbStorage{
		db:  db,
		cfg: cfg,
	}
}

func (c *ClimbStorage) Create(climb ClimbModel) (ClimbModel, error) {
	climb.ID = uuid.New().String()
	climb.CreatedAt = time.Now().Unix()
	climb.UpdatedAt = time.Now().Unix()
	climb.PK = "USER#" + climb.UserId
	climb.SK = "CLIMB#" + climb.ID
	climb.EntityType = "CLIMB"

	attributeValue, err := dynamodbattribute.MarshalMap(climb)
	if err != nil {
		slog.Error("Failed to marshal climb", "error", err)
		return ClimbModel{}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	slog.Debug("Writing to DynamoDB", "table", c.cfg.TABLE_NAME, "PK", climb.PK, "SK", climb.SK)
	_, err = c.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "table", c.cfg.TABLE_NAME)
		return ClimbModel{}, err
	}

	return climb, nil
}

func (c *ClimbStorage) GetAll(userId string) ([]ClimbModel, error) {
	var climbs []ClimbModel

	input := &dynamodb.QueryInput{
		TableName:              aws.String(c.cfg.TABLE_NAME),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String("USER#" + userId)},
			":sk": {S: aws.String("CLIMB#")},
		},
	}

	result, err := c.db.Query(input)
	if err != nil {
		slog.Error("DynamoDB Query failed", "error", err, "userId", userId)
		return climbs, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &climbs)
	if err != nil {
		slog.Error("Failed to unmarshal climbs", "error", err)
		return climbs, err
	}

	return climbs, nil
}

func (c *ClimbStorage) FindById(userId, climbId string) (ClimbModel, error) {
	var climb ClimbModel

	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String("USER#" + userId)},
			"SK": {S: aws.String("CLIMB#" + climbId)},
		},
	}

	result, err := c.db.GetItem(input)
	if err != nil {
		slog.Error("DynamoDB GetItem failed", "error", err, "userId", userId, "climbId", climbId)
		return climb, err
	}

	if result.Item == nil {
		return climb, errors.New("Could not find a climb with that ID")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &climb)
	if err != nil {
		slog.Error("Failed to unmarshal climb", "error", err)
		return climb, err
	}

	return climb, nil
}

func (c *ClimbStorage) Delete(userId, climbId string) error {
	climb, err := c.FindById(userId, climbId)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(climb.PK)},
			"SK": {S: aws.String(climb.SK)},
		},
	}

	_, err = c.db.DeleteItem(input)
	if err != nil {
		slog.Error("DynamoDB DeleteItem failed", "error", err, "PK", climb.PK, "SK", climb.SK)
		return err
	}

	slog.Info("Climb deleted", "climbId", climbId, "userId", userId)
	return nil
}

func (c *ClimbStorage) Update(userId, climbId string, data ClimbUpdateReq) (ClimbModel, error) {
	var updatedClimb ClimbModel

	climb, err := c.FindById(userId, climbId)
	if err != nil {
		return updatedClimb, err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(climb.PK)},
			"SK": {S: aws.String(climb.SK)},
		},
		UpdateExpression: aws.String("Set Grade = :grade, PerceivedDifficulty = :pd, TypeOfClimb = :tyc, GradeSystem = :gs, UpdatedAt = :updatedAt"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":grade":     {S: aws.String(data.Grade)},
			":pd":        {S: aws.String(data.PerceivedDifficulty)},
			":tyc":       {S: aws.String(data.TypeOfClimb)},
			":gs":        {S: aws.String(data.GradeSystem)},
			":updatedAt": {N: aws.String(strconv.FormatInt(time.Now().Unix(), 10))},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}

	result, err := c.db.UpdateItem(input)
	if err != nil {
		slog.Error("DynamoDB UpdateItem failed", "error", err, "PK", climb.PK, "SK", climb.SK)
		return updatedClimb, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedClimb)
	if err != nil {
		slog.Error("Failed to unmarshal updated climb", "error", err)
		return updatedClimb, err
	}

	slog.Info("Climb updated", "climbId", climbId)
	return updatedClimb, nil
}
