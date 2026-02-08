package habits

import (
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jimvid/sidekick/internal/config"
)

const (
	itemPrefixHabit = "habit#"
)

type HabitStorage struct {
	db  *dynamodb.DynamoDB
	cfg *config.Config
}

func NewHabitStorage(db *dynamodb.DynamoDB, cfg *config.Config) *HabitStorage {
	return &HabitStorage{
		db:  db,
		cfg: cfg,
	}
}

func (c *HabitStorage) Create(userId string, habit HabitModel) error {
	newItem := habitItem{
		UserId: userId,
		ItemId: itemPrefixHabit + habit.ID,
		HabitModel: habit,
	}

	attributeValue, err := dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		slog.Error("Failed to marshal habit", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	slog.Debug("Writing to DynamoDB", "table", c.cfg.TABLE_NAME, "userId", newItem.UserId, "itemId", newItem.ItemId)

	_, err = c.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "table", c.cfg.TABLE_NAME)
		return err
	}

	return nil
}

func (c *HabitStorage) GetAll(userId string) ([]HabitModel, error) {
	var items []habitItem

	input := &dynamodb.QueryInput{
		TableName:              aws.String(c.cfg.TABLE_NAME),
		KeyConditionExpression: aws.String("UserId = :userId AND begins_with(ItemId, :itemId)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {S: aws.String(userId)},
			":itemId": {S: aws.String(itemPrefixHabit)},
		},
	}

	result, err := c.db.Query(input)
	if err != nil {
		slog.Error("DynamoDB Query failed", "error", err, "userId", userId)
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		slog.Error("Failed to unmarshal habits", "error", err)
		return nil, err
	}

	habits := make([]HabitModel, len(items))
	for i, item := range items {
		habits[i] = item.HabitModel
	}

	return habits, nil
}

func (c *HabitStorage) FindById(userId, habitId string) (HabitModel, error) {
	var item habitItem

	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {S: aws.String(userId)},
			"ItemId": {S: aws.String(itemPrefixHabit + habitId)},
		},
	}

	result, err := c.db.GetItem(input)
	if err != nil {
		slog.Error("DynamoDB GetItem failed", "error", err, "userId", userId, "habitId", habitId)
		return HabitModel{}, err
	}

	if result.Item == nil {
		return HabitModel{}, errors.New("could not find a habit with that ID")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		slog.Error("Failed to unmarshal habit", "error", err)
		return HabitModel{}, err
	}

	return item.HabitModel, nil
}

func (c *HabitStorage) Delete(userId, habitId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {S: aws.String(userId)},
			"ItemId": {S: aws.String(itemPrefixHabit + habitId)},
		},
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := c.db.DeleteItem(input)
	if err != nil {
		slog.Error("DynamoDB DeleteItem failed", "error", err, "userId", userId, "habitId", habitId)
		return err
	}

	if result.Attributes == nil {
		return errors.New("could not find a habit with that ID")
	}

	slog.Info("Habit deleted", "habitId", habitId, "userId", userId)
	return nil
}

func (c *HabitStorage) Update(userId, habitId string, habit HabitModel) error {
	item := habitItem{
		UserId: userId,
		ItemId: itemPrefixHabit + habitId,
		HabitModel: habit,
	}

	attributeValue, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		slog.Error("Failed to marshal habit", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	_, err = c.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "userId", userId, "habitId", habitId)
		return err
	}

	slog.Info("Habit updated", "habitId", habitId)
	return nil
}
