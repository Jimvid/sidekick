package habits

import (
	"testing"
	"time"
)

func TestServiceCreateHabit(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	before := time.Now().Unix()
	habit, err := service.CreateHabit("user-1", HabitReq{
		Name:        "Exercise",
		Description: "Daily workout",
		Color:       "#ff0000",
	})
	after := time.Now().Unix()

	if err != nil {
		t.Fatalf("CreateHabit failed: %v", err)
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
	result, err := storage.FindHabitById("user-1", habit.ID)
	if err != nil {
		t.Fatalf("FindHabitById after CreateHabit failed: %v", err)
	}
	if result.Name != "Exercise" {
		t.Errorf("persisted name %q, expected %q", result.Name, "Exercise")
	}
}

func TestServiceGetAllHabits(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	t.Run("empty", func(t *testing.T) {
		habits, err := service.GetAllHabits("user-1")
		if err != nil {
			t.Fatalf("GetAllHabits failed: %v", err)
		}
		if len(habits) != 0 {
			t.Errorf("expected 0 habits, got %d", len(habits))
		}
	})

	t.Run("returns results", func(t *testing.T) {
		service.CreateHabit("user-1", HabitReq{Name: "A"})
		service.CreateHabit("user-1", HabitReq{Name: "B"})

		habits, err := service.GetAllHabits("user-1")
		if err != nil {
			t.Fatalf("GetAllHabits failed: %v", err)
		}
		if len(habits) != 2 {
			t.Errorf("expected 2 habits, got %d", len(habits))
		}
	})
}

func TestServiceFindHabitById(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.CreateHabit("user-1", HabitReq{Name: "Read"})

	t.Run("found", func(t *testing.T) {
		habit, err := service.FindHabitById("user-1", created.ID)
		if err != nil {
			t.Fatalf("FindHabitById failed: %v", err)
		}
		if habit.Name != "Read" {
			t.Errorf("expected name %q, got %q", "Read", habit.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.FindHabitById("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent habit, got nil")
		}
	})
}

func TestServiceDeleteHabit(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.CreateHabit("user-1", HabitReq{Name: "Exercise"})

	err := service.DeleteHabit("user-1", created.ID)
	if err != nil {
		t.Fatalf("DeleteHabit failed: %v", err)
	}

	_, err = service.FindHabitById("user-1", created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestServiceUpdateHabit(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.CreateHabit("user-1", HabitReq{
		Name:        "Exercise",
		Description: "Morning run",
		Color:       "#ff0000",
	})

	time.Sleep(time.Second) // ensure UpdatedAt differs

	updated, err := service.UpdateHabit("user-1", created.ID, HabitReq{
		Name:        "Yoga",
		Description: "Evening yoga",
		Color:       "#00ff00",
	})
	if err != nil {
		t.Fatalf("UpdateHabit failed: %v", err)
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

func TestServiceCreateHabitLog(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	log, err := service.CreateHabitLog("user-1", HabitLogReq{
		HabitId: "habit-1",
		Date:    "2026-02-08",
		Note:    "Morning run",
	})

	if err != nil {
		t.Fatalf("CreateLog failed: %v", err)
	}
	if log.ID == "" {
		t.Fatal("expected UUID to be generated, got empty string")
	}
	if log.HabitId != "habit-1" {
		t.Errorf("expected habitId %q, got %q", "habit-1", log.HabitId)
	}
	if log.Date != "2026-02-08" {
		t.Errorf("expected date %q, got %q", "2026-02-08", log.Date)
	}
	if log.Note != "Morning run" {
		t.Errorf("expected note %q, got %q", "Morning run", log.Note)
	}

	// Verify persisted
	result, err := storage.FindHabitLogById("user-1", log.ID)
	if err != nil {
		t.Fatalf("FindLogById after CreateLog failed: %v", err)
	}
	if result.HabitId != "habit-1" {
		t.Errorf("persisted habitId %q, expected %q", result.HabitId, "habit-1")
	}
}

func TestServiceGetAllHabitLogs(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	t.Run("empty", func(t *testing.T) {
		logs, err := service.GetAllHabitLogs("user-1")
		if err != nil {
			t.Fatalf("GetAllLogs failed: %v", err)
		}
		if len(logs) != 0 {
			t.Errorf("expected 0 logs, got %d", len(logs))
		}
	})

	t.Run("returns results", func(t *testing.T) {
		service.CreateHabitLog("user-1", HabitLogReq{HabitId: "h1", Date: "2026-02-08"})
		service.CreateHabitLog("user-1", HabitLogReq{HabitId: "h1", Date: "2026-02-09"})

		logs, err := service.GetAllHabitLogs("user-1")
		if err != nil {
			t.Fatalf("GetAllLogs failed: %v", err)
		}
		if len(logs) != 2 {
			t.Errorf("expected 2 logs, got %d", len(logs))
		}
	})
}

func TestServiceFindHabitLogById(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.CreateHabitLog("user-1", HabitLogReq{HabitId: "h1", Date: "2026-02-08", Note: "test"})

	t.Run("found", func(t *testing.T) {
		log, err := service.FindHabitLogById("user-1", created.ID)
		if err != nil {
			t.Fatalf("FindLogById failed: %v", err)
		}
		if log.HabitId != "h1" {
			t.Errorf("expected habitId %q, got %q", "h1", log.HabitId)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.FindHabitLogById("user-1", "does-not-exist")
		if err == nil {
			t.Fatal("expected error for non-existent log, got nil")
		}
	})
}

func TestServiceDeleteHabitLog(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.CreateHabitLog("user-1", HabitLogReq{HabitId: "h1", Date: "2026-02-08"})

	err := service.DeleteHabitLog("user-1", created.ID)
	if err != nil {
		t.Fatalf("DeleteLog failed: %v", err)
	}

	_, err = service.FindHabitLogById("user-1", created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestServiceUpdateHabitLog(t *testing.T) {
	storage := setupTestDB(t)
	service := NewHabitService(storage)

	created, _ := service.CreateHabitLog("user-1", HabitLogReq{
		HabitId: "habit-1",
		Date:    "2026-02-08",
		Note:    "Morning run",
	})

	updated, err := service.UpdateHabitLog("user-1", created.ID, HabitLogReq{
		HabitId: "habit-2",
		Date:    "2026-02-09",
		Note:    "Evening yoga",
	})
	if err != nil {
		t.Fatalf("UpdateLog failed: %v", err)
	}

	if updated.HabitId != "habit-2" {
		t.Errorf("expected habitId %q, got %q", "habit-2", updated.HabitId)
	}
	if updated.Date != "2026-02-09" {
		t.Errorf("expected date %q, got %q", "2026-02-09", updated.Date)
	}
	if updated.Note != "Evening yoga" {
		t.Errorf("expected note %q, got %q", "Evening yoga", updated.Note)
	}
	if updated.ID != created.ID {
		t.Errorf("expected ID to be preserved (%q), got %q", created.ID, updated.ID)
	}
}
