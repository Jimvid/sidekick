package habits

import (
	"testing"
	"time"
)

func TestServiceCreate(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	before := time.Now().Unix()
	habit, err := service.Create("user-1", HabitReq{
		Name:        "Exercise",
		Description: "Daily workout",
		Color:       "#ff0000",
	})
	after := time.Now().Unix()

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if habit.ID == "" {
		t.Fatal("expected UUID to be generated, got empty string")
	}
	if habit.Name != "Exercise" {
		t.Errorf("expected name %q, got %q", "Exercise", habit.Name)
	}
	if habit.Description != "Daily workout" {
		t.Errorf("expected description %q, got %q", "Daily workout", habit.Description)
	}
	if habit.Color != "#ff0000" {
		t.Errorf("expected color %q, got %q", "#ff0000", habit.Color)
	}
	if habit.CreatedAt < before || habit.CreatedAt > after {
		t.Errorf("expected CreatedAt between %d and %d, got %d", before, after, habit.CreatedAt)
	}
	if habit.UpdatedAt < before || habit.UpdatedAt > after {
		t.Errorf("expected UpdatedAt between %d and %d, got %d", before, after, habit.UpdatedAt)
	}

	// Verify persisted
	result, err := storage.FindById("user-1", habit.ID)
	if err != nil {
		t.Fatalf("FindById after Create failed: %v", err)
	}
	if result.Name != "Exercise" {
		t.Errorf("persisted name %q, expected %q", result.Name, "Exercise")
	}
}

func TestServiceGetAll(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	t.Run("empty", func(t *testing.T) {
		habits, err := service.GetAll("user-1")
		if err != nil {
			t.Fatalf("GetAll failed: %v", err)
		}
		if len(habits) != 0 {
			t.Errorf("expected 0 habits, got %d", len(habits))
		}
	})

	t.Run("returns results", func(t *testing.T) {
		service.Create("user-1", HabitReq{Name: "A"})
		service.Create("user-1", HabitReq{Name: "B"})

		habits, err := service.GetAll("user-1")
		if err != nil {
			t.Fatalf("GetAll failed: %v", err)
		}
		if len(habits) != 2 {
			t.Errorf("expected 2 habits, got %d", len(habits))
		}
	})
}

func TestServiceFindById(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.Create("user-1", HabitReq{Name: "Read"})

	t.Run("found", func(t *testing.T) {
		habit, err := service.FindById("user-1", created.ID)
		if err != nil {
			t.Fatalf("FindById failed: %v", err)
		}
		if habit.Name != "Read" {
			t.Errorf("expected name %q, got %q", "Read", habit.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.FindById("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent habit, got nil")
		}
	})
}

func TestServiceDelete(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.Create("user-1", HabitReq{Name: "Exercise"})

	err := service.Delete("user-1", created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = service.FindById("user-1", created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestServiceUpdate(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.Create("user-1", HabitReq{
		Name:        "Exercise",
		Description: "Morning run",
		Color:       "#ff0000",
	})

	time.Sleep(time.Second) // ensure UpdatedAt differs

	updated, err := service.Update("user-1", created.ID, HabitReq{
		Name:        "Yoga",
		Description: "Evening yoga",
		Color:       "#00ff00",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Name != "Yoga" {
		t.Errorf("expected name %q, got %q", "Yoga", updated.Name)
	}
	if updated.Description != "Evening yoga" {
		t.Errorf("expected description %q, got %q", "Evening yoga", updated.Description)
	}
	if updated.Color != "#00ff00" {
		t.Errorf("expected color %q, got %q", "#00ff00", updated.Color)
	}
	if updated.CreatedAt != created.CreatedAt {
		t.Errorf("expected CreatedAt to be preserved (%d), got %d", created.CreatedAt, updated.CreatedAt)
	}
	if updated.UpdatedAt <= created.UpdatedAt {
		t.Errorf("expected UpdatedAt (%d) to be greater than original (%d)", updated.UpdatedAt, created.UpdatedAt)
	}
}
