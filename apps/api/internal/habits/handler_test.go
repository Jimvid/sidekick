package habits

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

const testUserId = "test-user-1"

func setupHandler(t *testing.T) (*HabitHandler, *chi.Mux) {
	t.Helper()

	storage := setupTestDB(t)
	service := NewHabitService(storage)
	handler := &HabitHandler{
		service: service,
		getUserId: func(r *http.Request) (string, error) {
			return testUserId, nil
		},
	}

	r := chi.NewRouter()
	r.Post("/habits", handler.CreateHabit)
	r.Get("/habits", handler.GetAllHabits)
	r.Get("/habits/{habitId}", handler.FindHabitById)
	r.Delete("/habits/{habitId}", handler.DeleteHabit)
	r.Put("/habits/{habitId}", handler.UpdateHabit)

	r.Post("/habit-logs", handler.CreateHabitLog)
	r.Get("/habit-logs", handler.GetAllHabitLogs)
	r.Get("/habit-logs/{id}", handler.FindHabitLogById)
	r.Delete("/habit-logs/{id}", handler.DeleteHabitLog)
	r.Put("/habit-logs/{id}", handler.UpdateHabitLog)

	return handler, r
}

func TestHandlerCreateHabit(t *testing.T) {
	_, router := setupHandler(t)

	t.Run("valid request", func(t *testing.T) {
		body := `{"name":"Exercise","description":"Daily workout","color":"#ff0000"}`
		req := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var habit HabitModel
		json.NewDecoder(w.Body).Decode(&habit)

		if habit.ID == "" {
			t.Error("expected habit ID to be set")
		}
		if habit.Name != "Exercise" {
			t.Errorf("expected name %q, got %q", "Exercise", habit.Name)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandlerGetAllHabits(t *testing.T) {
	_, router := setupHandler(t)

	t.Run("empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/habits", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var habits []HabitModel
		json.NewDecoder(w.Body).Decode(&habits)

		if len(habits) != 0 {
			t.Errorf("expected 0 habits, got %d", len(habits))
		}
	})

	t.Run("returns list", func(t *testing.T) {
		// Create two habits
		for _, name := range []string{"Exercise", "Read"} {
			body := `{"name":"` + name + `","description":"desc","color":"#000"}`
			req := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		req := httptest.NewRequest(http.MethodGet, "/habits", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var habits []HabitModel
		json.NewDecoder(w.Body).Decode(&habits)

		if len(habits) != 2 {
			t.Errorf("expected 2 habits, got %d", len(habits))
		}
	})
}

func TestHandlerFindHabitById(t *testing.T) {
	_, router := setupHandler(t)

	// Create a habit first
	body := `{"name":"Exercise","description":"desc","color":"#000"}`
	createReq := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created HabitModel
	json.NewDecoder(createW.Body).Decode(&created)

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/habits/"+created.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var habit HabitModel
		json.NewDecoder(w.Body).Decode(&habit)

		if habit.Name != "Exercise" {
			t.Errorf("expected name %q, got %q", "Exercise", habit.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/habits/does-not-exist", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestHandlerDeleteHabit(t *testing.T) {
	_, router := setupHandler(t)

	// Create a habit first
	body := `{"name":"Exercise","description":"desc","color":"#000"}`
	createReq := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created HabitModel
	json.NewDecoder(createW.Body).Decode(&created)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/habits/"+created.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/habits/does-not-exist", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestHandlerUpdateHabit(t *testing.T) {
	_, router := setupHandler(t)

	// Create a habit first
	body := `{"name":"Exercise","description":"desc","color":"#000"}`
	createReq := httptest.NewRequest(http.MethodPost, "/habits", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created HabitModel
	json.NewDecoder(createW.Body).Decode(&created)

	t.Run("success", func(t *testing.T) {
		updateBody := `{"name":"Yoga","description":"Evening yoga","color":"#00ff00"}`
		req := httptest.NewRequest(http.MethodPut, "/habits/"+created.ID, strings.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var updated HabitModel
		json.NewDecoder(w.Body).Decode(&updated)

		if updated.Name != "Yoga" {
			t.Errorf("expected name %q, got %q", "Yoga", updated.Name)
		}
		if updated.Description != "Evening yoga" {
			t.Errorf("expected description %q, got %q", "Evening yoga", updated.Description)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/habits/"+created.ID, strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandlerCreateHabitLog(t *testing.T) {
	_, router := setupHandler(t)

	t.Run("valid request", func(t *testing.T) {
		body := `{"habitId":"habit-1","date":"2026-02-08","note":"Morning run"}`
		req := httptest.NewRequest(http.MethodPost, "/habit-logs", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var log HabitLogModel
		json.NewDecoder(w.Body).Decode(&log)

		if log.ID == "" {
			t.Error("expected log ID to be set")
		}
		if log.HabitId != "habit-1" {
			t.Errorf("expected habitId %q, got %q", "habit-1", log.HabitId)
		}
		if log.Date != "2026-02-08" {
			t.Errorf("expected date %q, got %q", "2026-02-08", log.Date)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/habit-logs", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandlerGetAllHabitLogs(t *testing.T) {
	_, router := setupHandler(t)

	t.Run("empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/habit-logs", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var logs []HabitLogModel
		json.NewDecoder(w.Body).Decode(&logs)

		if len(logs) != 0 {
			t.Errorf("expected 0 logs, got %d", len(logs))
		}
	})

	t.Run("returns list", func(t *testing.T) {
		for _, date := range []string{"2026-02-08", "2026-02-09"} {
			body := `{"habitId":"h1","date":"` + date + `","note":"test"}`
			req := httptest.NewRequest(http.MethodPost, "/habit-logs", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		req := httptest.NewRequest(http.MethodGet, "/habit-logs", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var logs []HabitLogModel
		json.NewDecoder(w.Body).Decode(&logs)

		if len(logs) != 2 {
			t.Errorf("expected 2 logs, got %d", len(logs))
		}
	})
}

func TestHandlerFindHabitLogById(t *testing.T) {
	_, router := setupHandler(t)

	// Create a log first
	body := `{"habitId":"habit-1","date":"2026-02-08","note":"test"}`
	createReq := httptest.NewRequest(http.MethodPost, "/habit-logs", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created HabitLogModel
	json.NewDecoder(createW.Body).Decode(&created)

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/habit-logs/"+created.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var log HabitLogModel
		json.NewDecoder(w.Body).Decode(&log)

		if log.HabitId != "habit-1" {
			t.Errorf("expected habitId %q, got %q", "habit-1", log.HabitId)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/habit-logs/does-not-exist", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestHandlerDeleteHabitLog(t *testing.T) {
	_, router := setupHandler(t)

	// Create a log first
	body := `{"habitId":"habit-1","date":"2026-02-08","note":"test"}`
	createReq := httptest.NewRequest(http.MethodPost, "/habit-logs", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created HabitLogModel
	json.NewDecoder(createW.Body).Decode(&created)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/habit-logs/"+created.ID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/habit-logs/does-not-exist", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestHandlerUpdateHabitLog(t *testing.T) {
	_, router := setupHandler(t)

	// Create a log first
	body := `{"habitId":"habit-1","date":"2026-02-08","note":"test"}`
	createReq := httptest.NewRequest(http.MethodPost, "/habit-logs", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var created HabitLogModel
	json.NewDecoder(createW.Body).Decode(&created)

	t.Run("success", func(t *testing.T) {
		updateBody := `{"habitId":"habit-2","date":"2026-02-09","note":"updated"}`
		req := httptest.NewRequest(http.MethodPut, "/habit-logs/"+created.ID, strings.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var updated HabitLogModel
		json.NewDecoder(w.Body).Decode(&updated)

		if updated.HabitId != "habit-2" {
			t.Errorf("expected habitId %q, got %q", "habit-2", updated.HabitId)
		}
		if updated.Date != "2026-02-09" {
			t.Errorf("expected date %q, got %q", "2026-02-09", updated.Date)
		}
		if updated.Note != "updated" {
			t.Errorf("expected note %q, got %q", "updated", updated.Note)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/habit-logs/"+created.ID, strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
