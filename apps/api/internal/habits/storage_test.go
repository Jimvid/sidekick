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
			{AttributeName: aws.String("userId"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("itemId"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("userId"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("itemId"), AttributeType: aws.String("S")},
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

func makeLog(id, habitId, date string) HabitLogModel {
	return HabitLogModel{
		ID:      id,
		HabitId: habitId,
		Date:    date,
		Note:    "test note",
	}
}

func TestStorageCreateHabit(t *testing.T) {
	storage := setupTestDB(t)
	habit := makeHabit("habit-1", "Exercise")

	err := storage.CreateHabit("user-1", habit)
	if err != nil {
		t.Fatalf("CreateHabit failed: %v", err)
	}

	// Verify it was persisted
	result, err := storage.FindHabitById("user-1", "habit-1")
	if err != nil {
		t.Fatalf("FindHabitById after CreateHabit failed: %v", err)
	}
	if result.Name != "Exercise" {
		t.Errorf("expected name %q, got %q", "Exercise", result.Name)
	}
}

func TestStorageFindHabitById(t *testing.T) {
	storage := setupTestDB(t)
	habit := makeHabit("habit-1", "Read")
	storage.CreateHabit("user-1", habit)

	t.Run("existing habit", func(t *testing.T) {
		result, err := storage.FindHabitById("user-1", "habit-1")
		if err != nil {
			t.Fatalf("FindHabitById failed: %v", err)
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
		_, err := storage.FindHabitById("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent habit, got nil")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		_, err := storage.FindHabitById("user-999", "habit-1")
		if err == nil {
			t.Fatal("expected error when querying with wrong user, got nil")
		}
	})
}

func TestStorageGetAllHabits(t *testing.T) {
	storage := setupTestDB(t)

	t.Run("empty result", func(t *testing.T) {
		habits, err := storage.GetAllHabits("user-1")
		if err != nil {
			t.Fatalf("GetAllHabits failed: %v", err)
		}
		if len(habits) != 0 {
			t.Errorf("expected 0 habits, got %d", len(habits))
		}
	})

	t.Run("returns only habits for the given user", func(t *testing.T) {
		storage.CreateHabit("user-1", makeHabit("h1", "Exercise"))
		storage.CreateHabit("user-1", makeHabit("h2", "Read"))
		storage.CreateHabit("user-2", makeHabit("h3", "Meditate"))

		habits, err := storage.GetAllHabits("user-1")
		if err != nil {
			t.Fatalf("GetAllHabits failed: %v", err)
		}
		if len(habits) != 2 {
			t.Fatalf("expected 2 habits for user-1, got %d", len(habits))
		}

		habits2, err := storage.GetAllHabits("user-2")
		if err != nil {
			t.Fatalf("GetAllHabits failed: %v", err)
		}
		if len(habits2) != 1 {
			t.Fatalf("expected 1 habit for user-2, got %d", len(habits2))
		}
	})
}

func TestStorageDeleteHabit(t *testing.T) {
	storage := setupTestDB(t)
	storage.CreateHabit("user-1", makeHabit("habit-1", "Exercise"))

	t.Run("existing habit", func(t *testing.T) {
		err := storage.DeleteHabit("user-1", "habit-1")
		if err != nil {
			t.Fatalf("DeleteHabit failed: %v", err)
		}

		_, err = storage.FindHabitById("user-1", "habit-1")
		if err == nil {
			t.Fatal("expected error after delete, got nil")
		}
	})

	t.Run("non-existent habit", func(t *testing.T) {
		err := storage.DeleteHabit("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent habit, got nil")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		storage.CreateHabit("user-1", makeHabit("habit-2", "Read"))

		err := storage.DeleteHabit("user-999", "habit-2")
		if err == nil {
			t.Fatal("expected error when deleting with wrong user, got nil")
		}
	})
}

func TestStorageUpdateHabit(t *testing.T) {
	storage := setupTestDB(t)
	original := makeHabit("habit-1", "Exercise")
	storage.CreateHabit("user-1", original)

	updated := original
	updated.Name = "Morning Exercise"
	updated.Color = "#00ff00"
	updated.UpdatedAt = time.Now().Unix() + 100

	err := storage.UpdateHabit("user-1", "habit-1", updated)
	if err != nil {
		t.Fatalf("UpdateHabit failed: %v", err)
	}

	result, err := storage.FindHabitById("user-1", "habit-1")
	if err != nil {
		t.Fatalf("FindHabitById after UpdateHabit failed: %v", err)
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

// Habit logs
func TestStorageCreateHabitLog(t *testing.T) {
	storage := setupTestDB(t)
	log := makeLog("log-1", "habit-1", "2026-02-08")

	err := storage.CreateHabitLog("user-1", log)
	if err != nil {
		t.Fatalf("CreateHabitLog failed: %v", err)
	}

	result, err := storage.FindHabitLogById("user-1", "log-1")
	if err != nil {
		t.Fatalf("FindHabitLogById after CreateHabitLog failed: %v", err)
	}
	if result.HabitId != "habit-1" {
		t.Errorf("expected habitId %q, got %q", "habit-1", result.HabitId)
	}
	if result.Date != "2026-02-08" {
		t.Errorf("expected date %q, got %q", "2026-02-08", result.Date)
	}
}

func TestStorageFindHabitLogById(t *testing.T) {
	storage := setupTestDB(t)
	storage.CreateHabitLog("user-1", makeLog("log-1", "habit-1", "2026-02-08"))

	t.Run("existing log", func(t *testing.T) {
		result, err := storage.FindHabitLogById("user-1", "log-1")
		if err != nil {
			t.Fatalf("FindHabitLogById failed: %v", err)
		}
		if result.ID != "log-1" {
			t.Errorf("expected ID %q, got %q", "log-1", result.ID)
		}
		if result.HabitId != "habit-1" {
			t.Errorf("expected habitId %q, got %q", "habit-1", result.HabitId)
		}
	})

	t.Run("non-existent log", func(t *testing.T) {
		_, err := storage.FindHabitLogById("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent log, got nil")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		_, err := storage.FindHabitLogById("user-999", "log-1")
		if err == nil {
			t.Fatal("expected error when querying with wrong user, got nil")
		}
	})
}

func TestStorageGetAllHabitLogs(t *testing.T) {
	storage := setupTestDB(t)

	t.Run("empty result", func(t *testing.T) {
		logs, err := storage.GetAllHabitLogs("user-1")
		if err != nil {
			t.Fatalf("GetAllHabitLogs failed: %v", err)
		}
		if len(logs) != 0 {
			t.Errorf("expected 0 logs, got %d", len(logs))
		}
	})

	t.Run("returns only logs for the given user", func(t *testing.T) {
		storage.CreateHabitLog("user-1", makeLog("l1", "habit-1", "2026-02-08"))
		storage.CreateHabitLog("user-1", makeLog("l2", "habit-1", "2026-02-09"))
		storage.CreateHabitLog("user-2", makeLog("l3", "habit-2", "2026-02-08"))

		logs, err := storage.GetAllHabitLogs("user-1")
		if err != nil {
			t.Fatalf("GetAllHabitLogs failed: %v", err)
		}
		if len(logs) != 2 {
			t.Fatalf("expected 2 logs for user-1, got %d", len(logs))
		}

		logs2, err := storage.GetAllHabitLogs("user-2")
		if err != nil {
			t.Fatalf("GetAllHabitLogs failed: %v", err)
		}
		if len(logs2) != 1 {
			t.Fatalf("expected 1 log for user-2, got %d", len(logs2))
		}
	})
}

func TestStorageDeleteHabitLog(t *testing.T) {
	storage := setupTestDB(t)
	storage.CreateHabitLog("user-1", makeLog("log-1", "habit-1", "2026-02-08"))

	t.Run("existing log", func(t *testing.T) {
		err := storage.DeleteHabitLog("user-1", "log-1")
		if err != nil {
			t.Fatalf("DeleteHabitLog failed: %v", err)
		}

		_, err = storage.FindHabitLogById("user-1", "log-1")
		if err == nil {
			t.Fatal("expected error after delete, got nil")
		}
	})

	t.Run("non-existent log", func(t *testing.T) {
		err := storage.DeleteHabitLog("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent log, got nil")
		}
	})

	t.Run("wrong user", func(t *testing.T) {
		storage.CreateHabitLog("user-1", makeLog("log-2", "habit-1", "2026-02-09"))

		err := storage.DeleteHabitLog("user-999", "log-2")
		if err == nil {
			t.Fatal("expected error when deleting with wrong user, got nil")
		}
	})
}

func TestStorageUpdateHabitLog(t *testing.T) {
	storage := setupTestDB(t)
	original := makeLog("log-1", "habit-1", "2026-02-08")
	storage.CreateHabitLog("user-1", original)

	updated := original
	updated.Date = "2026-02-09"
	updated.Note = "updated note"

	err := storage.UpdateHabitLog("user-1", "log-1", updated)
	if err != nil {
		t.Fatalf("UpdateLog failed: %v", err)
	}

	result, err := storage.FindHabitLogById("user-1", "log-1")
	if err != nil {
		t.Fatalf("FindHabitLogById after UpdateLog failed: %v", err)
	}
	if result.Date != "2026-02-09" {
		t.Errorf("expected date %q, got %q", "2026-02-09", result.Date)
	}
	if result.Note != "updated note" {
		t.Errorf("expected note %q, got %q", "updated note", result.Note)
	}
	if result.HabitId != "habit-1" {
		t.Errorf("expected habitId %q, got %q", "habit-1", result.HabitId)
	}
}
