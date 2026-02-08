package habits

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jimvid/sidekick/internal/auth"
)

type HabitHandler struct {
	service   *HabitService
	getUserId func(r *http.Request) (string, error)
}

func NewHabitHandler(service *HabitService) *HabitHandler {
	return &HabitHandler{
		service:   service,
		getUserId: auth.GetUserId,
	}
}

func (h *HabitHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *HabitHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *HabitHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	habits, err := h.service.GetAll(userId)
	if err != nil {
		slog.Error("Failed to get all habits", "error", err, "userId", userId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get all habits")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, habits)
}

func (h *HabitHandler) Create(w http.ResponseWriter, r *http.Request) {
	var habitReq HabitReq
	err := json.NewDecoder(r.Body).Decode(&habitReq)
	if err != nil {
		slog.Error("Failed to parse JSON", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not parse JSON")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	habit, err := h.service.Create(userId, habitReq)
	if err != nil {
		slog.Error("Failed to create habit", "error", err, "userId", userId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not create habit")
		return
	}

	slog.Info("Habit created", "habitId", habit.ID, "userId", userId)
	h.writeSuccessResponse(w, http.StatusCreated, habit)
}

func (h *HabitHandler) FindById(w http.ResponseWriter, r *http.Request) {
	habitId := chi.URLParam(r, "habitId")
	if habitId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	habit, err := h.service.FindById(userId, habitId)
	if err != nil {
		slog.Error("Could not find habit by ID", "error", err, "habitId", habitId)
		h.writeErrorResponse(w, http.StatusNotFound, "Could not find habit by ID")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, habit)
}

func (h *HabitHandler) Delete(w http.ResponseWriter, r *http.Request) {
	habitId := chi.URLParam(r, "habitId")
	if habitId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	err = h.service.Delete(userId, habitId)
	if err != nil {
		slog.Error("Could not delete habit", "error", err, "habitId", habitId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not delete habit")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Successfully deleted habit"})
}

func (h *HabitHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req HabitReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("Failed to parse JSON", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not parse JSON")
		return
	}

	habitId := chi.URLParam(r, "habitId")
	if habitId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	updatedHabit, err := h.service.Update(userId, habitId, req)
	if err != nil {
		slog.Error("Could not update habit", "error", err, "habitId", habitId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not update habit")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, updatedHabit)
}
