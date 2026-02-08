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
	itemPrefixHabit    = "habit#"
	itemPrefixHabitLog = "habit-log#"
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

func (s *HabitStorage) CreateHabit(userId string, habit HabitModel) error {
	newItem := habitItem{
		UserId:     userId,
		ItemId:     itemPrefixHabit + habit.ID,
		HabitModel: habit,
	}

	attributeValue, err := dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		slog.Error("Failed to marshal habit", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	slog.Debug("Writing to DynamoDB", "table", s.cfg.TABLE_NAME, "userId", newItem.UserId, "itemId", newItem.ItemId)

	_, err = s.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "table", s.cfg.TABLE_NAME)
		return err
	}

	return nil
}

func (s *HabitStorage) GetAllHabits(userId string) ([]HabitModel, error) {
	var items []habitItem

	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.cfg.TABLE_NAME),
		KeyConditionExpression: aws.String("userId = :userId AND begins_with(itemId, :itemId)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {S: aws.String(userId)},
			":itemId": {S: aws.String(itemPrefixHabit)},
		},
	}

	result, err := s.db.Query(input)
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

func (s *HabitStorage) FindHabitById(userId, habitId string) (HabitModel, error) {
	var item habitItem

	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {S: aws.String(userId)},
			"itemId": {S: aws.String(itemPrefixHabit + habitId)},
		},
	}

	result, err := s.db.GetItem(input)
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

func (s *HabitStorage) DeleteHabit(userId, habitId string) error {
	// TODO: Delete all habit logs related to the habit
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {S: aws.String(userId)},
			"itemId": {S: aws.String(itemPrefixHabit + habitId)},
		},
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := s.db.DeleteItem(input)
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

func (s *HabitStorage) UpdateHabit(userId, habitId string, habit HabitModel) error {
	item := habitItem{
		UserId:     userId,
		ItemId:     itemPrefixHabit + habitId,
		HabitModel: habit,
	}

	attributeValue, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		slog.Error("Failed to marshal habit", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	_, err = s.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "userId", userId, "habitId", habitId)
		return err
	}

	slog.Info("Habit updated", "habitId", habitId)
	return nil
}

func (s *HabitStorage) CreateHabitLog(userId string, log HabitLogModel) error {
	newItem := habitLogItem{
		UserId:        userId,
		ItemId:        itemPrefixHabitLog + log.ID,
		HabitLogModel: log,
	}

	attributeValue, err := dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		slog.Error("Failed to marshal habit log", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	slog.Debug("Writing habit log to DynamoDB", "table", s.cfg.TABLE_NAME, "userId", newItem.UserId, "itemId", newItem.ItemId)

	_, err = s.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "table", s.cfg.TABLE_NAME)
		return err
	}

	return nil
}

func (s *HabitStorage) GetAllHabitLogs(userId string) ([]HabitLogModel, error) {
	var items []habitLogItem

	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.cfg.TABLE_NAME),
		KeyConditionExpression: aws.String("userId = :userId AND begins_with(itemId, :itemId)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userId": {S: aws.String(userId)},
			":itemId": {S: aws.String(itemPrefixHabitLog)},
		},
	}

	result, err := s.db.Query(input)
	if err != nil {
		slog.Error("DynamoDB Query failed", "error", err, "userId", userId)
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		slog.Error("Failed to unmarshal habit logs", "error", err)
		return nil, err
	}

	logs := make([]HabitLogModel, len(items))
	for i, item := range items {
		logs[i] = item.HabitLogModel
	}

	return logs, nil
}

func (s *HabitStorage) FindHabitLogById(userId, logId string) (HabitLogModel, error) {
	var item habitLogItem

	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {S: aws.String(userId)},
			"itemId": {S: aws.String(itemPrefixHabitLog + logId)},
		},
	}

	result, err := s.db.GetItem(input)
	if err != nil {
		slog.Error("DynamoDB GetItem failed", "error", err, "userId", userId, "logId", logId)
		return HabitLogModel{}, err
	}

	if result.Item == nil {
		return HabitLogModel{}, errors.New("could not find a habit log with that ID")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		slog.Error("Failed to unmarshal habit log", "error", err)
		return HabitLogModel{}, err
	}

	return item.HabitLogModel, nil
}

func (s *HabitStorage) DeleteHabitLog(userId, logId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {S: aws.String(userId)},
			"itemId": {S: aws.String(itemPrefixHabitLog + logId)},
		},
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := s.db.DeleteItem(input)
	if err != nil {
		slog.Error("DynamoDB DeleteItem failed", "error", err, "userId", userId, "logId", logId)
		return err
	}

	if result.Attributes == nil {
		return errors.New("could not find a habit log with that ID")
	}

	slog.Info("Habit log deleted", "logId", logId, "userId", userId)
	return nil
}

func (s *HabitStorage) UpdateHabitLog(userId, logId string, log HabitLogModel) error {
	item := habitLogItem{
		UserId:        userId,
		ItemId:        itemPrefixHabitLog + logId,
		HabitLogModel: log,
	}

	attributeValue, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		slog.Error("Failed to marshal habit log", "error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.cfg.TABLE_NAME),
		Item:      attributeValue,
	}

	_, err = s.db.PutItem(input)
	if err != nil {
		slog.Error("DynamoDB PutItem failed", "error", err, "userId", userId, "logId", logId)
		return err
	}

	slog.Info("Habit log updated", "logId", logId)
	return nil
}
