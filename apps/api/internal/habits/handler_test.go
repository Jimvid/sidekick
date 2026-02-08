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
	r.Post("/habits", handler.Create)
	r.Get("/habits", handler.GetAll)
	r.Get("/habits/{habitId}", handler.FindById)
	r.Delete("/habits/{habitId}", handler.Delete)
	r.Put("/habits/{habitId}", handler.Update)

	return handler, r
}

func TestHandlerCreate(t *testing.T) {
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

func TestHandlerGetAll(t *testing.T) {
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

func TestHandlerFindById(t *testing.T) {
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

func TestHandlerDelete(t *testing.T) {
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

func TestHandlerUpdate(t *testing.T) {
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
