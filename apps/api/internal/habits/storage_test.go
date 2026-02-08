package habits

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jimvid/sidekick/internal/config"
)

const testTableName = "test-habits"

func setupTestDB(t *testing.T) *HabitStorage {
	t.Helper()

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("fake", "fake", ""),
	}))
	db := dynamodb.New(sess)

	// Create table
	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("UserId"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("ItemId"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("UserId"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("ItemId"), AttributeType: aws.String("S")},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	})
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	t.Cleanup(func() {
		db.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(testTableName),
		})
	})

	cfg := &config.Config{TABLE_NAME: testTableName}
	return NewHabitStorage(db, cfg)
}

func makeHabit(id, name string) HabitModel {
	now := time.Now().Unix()
	return HabitModel{
		ID:          id,
		Name:        name,
		Description: "test description",
		Color:       "#ff0000",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestCreate(t *testing.T) {
	storage := setupTestDB(t)
	habit := makeHabit("habit-1", "Exercise")

	err := storage.Create("user-1", habit)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify it was persisted
	result, err := storage.FindById("user-1", "habit-1")
	if err != nil {
		t.Fatalf("FindById after Create failed: %v", err)
	}
	if result.Name != "Exercise" {
		t.Errorf("expected name %q, got %q", "Exercise", result.Name)
	}
}

func TestFindById(t *testing.T) {
	storage := setupTestDB(t)
	habit := makeHabit("habit-1", "Read")
	storage.Create("user-1", habit)

	t.Run("existing habit", func(t *testing.T) {
		result, err := storage.FindById("user-1", "habit-1")
		if err != nil {
			t.Fatalf("FindById failed: %v", err)
		}
		if result.ID != "habit-1" {
			t.Errorf("expected ID %q, got %q", "habit-1", result.ID)
		}
		if result.Name != "Read" {
			t.Errorf("expected name %q, got %q", "Read", result.Name)
		}
		if result.Color != "#ff0000" {
			t.Errorf("expected color %q, got %q", "#ff0000", result.Color)
		}
	})

	t.Run("non-existent habit", func(t *testing.T) {
		_, err := storage.FindById("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent habit, got nil")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		_, err := storage.FindById("user-999", "habit-1")
		if err == nil {
			t.Fatal("expected error when querying with wrong user, got nil")
		}
	})
}

func TestGetAll(t *testing.T) {
	storage := setupTestDB(t)

	t.Run("empty result", func(t *testing.T) {
		habits, err := storage.GetAll("user-1")
		if err != nil {
			t.Fatalf("GetAll failed: %v", err)
		}
		if len(habits) != 0 {
			t.Errorf("expected 0 habits, got %d", len(habits))
		}
	})

	t.Run("returns only habits for the given user", func(t *testing.T) {
		storage.Create("user-1", makeHabit("h1", "Exercise"))
		storage.Create("user-1", makeHabit("h2", "Read"))
		storage.Create("user-2", makeHabit("h3", "Meditate"))

		habits, err := storage.GetAll("user-1")
		if err != nil {
			t.Fatalf("GetAll failed: %v", err)
		}
		if len(habits) != 2 {
			t.Fatalf("expected 2 habits for user-1, got %d", len(habits))
		}

		habits2, err := storage.GetAll("user-2")
		if err != nil {
			t.Fatalf("GetAll failed: %v", err)
		}
		if len(habits2) != 1 {
			t.Fatalf("expected 1 habit for user-2, got %d", len(habits2))
		}
	})
}

func TestDelete(t *testing.T) {
	storage := setupTestDB(t)
	storage.Create("user-1", makeHabit("habit-1", "Exercise"))

	t.Run("existing habit", func(t *testing.T) {
		err := storage.Delete("user-1", "habit-1")
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = storage.FindById("user-1", "habit-1")
		if err == nil {
			t.Fatal("expected error after delete, got nil")
		}
	})

	t.Run("non-existent habit", func(t *testing.T) {
		err := storage.Delete("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent habit, got nil")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		storage.Create("user-1", makeHabit("habit-2", "Read"))

		err := storage.Delete("user-999", "habit-2")
		if err == nil {
			t.Fatal("expected error when deleting with wrong user, got nil")
		}
	})
}

func TestUpdate(t *testing.T) {
	storage := setupTestDB(t)
	original := makeHabit("habit-1", "Exercise")
	storage.Create("user-1", original)

	updated := original
	updated.Name = "Morning Exercise"
	updated.Color = "#00ff00"
	updated.UpdatedAt = time.Now().Unix() + 100

	err := storage.Update("user-1", "habit-1", updated)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	result, err := storage.FindById("user-1", "habit-1")
	if err != nil {
		t.Fatalf("FindById after Update failed: %v", err)
	}
	if result.Name != "Morning Exercise" {
		t.Errorf("expected name %q, got %q", "Morning Exercise", result.Name)
	}
	if result.Color != "#00ff00" {
		t.Errorf("expected color %q, got %q", "#00ff00", result.Color)
	}
	if result.UpdatedAt != updated.UpdatedAt {
		t.Errorf("expected updatedAt %d, got %d", updated.UpdatedAt, result.UpdatedAt)
	}
	// CreatedAt should be unchanged
	if result.CreatedAt != original.CreatedAt {
		t.Errorf("expected createdAt %d to be unchanged, got %d", original.CreatedAt, result.CreatedAt)
	}
}
