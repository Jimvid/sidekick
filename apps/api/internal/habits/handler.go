package habits

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jimvid/sidekick/internal/user"
)

type HabitHandler struct {
	service   *HabitService
	getUserId func(r *http.Request) (string, error)
}

func NewHabitHandler(service *HabitService) *HabitHandler {
	return &HabitHandler{
		service:   service,
		getUserId: user.GetUserId,
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

func (h *HabitHandler) GetAllHabits(w http.ResponseWriter, r *http.Request) {
	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	habits, err := h.service.GetAllHabits(userId)
	if err != nil {
		slog.Error("Failed to get all habits", "error", err, "userId", userId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get all habits")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, habits)
}

func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
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

	habit, err := h.service.CreateHabit(userId, habitReq)
	if err != nil {
		slog.Error("Failed to create habit", "error", err, "userId", userId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not create habit")
		return
	}

	slog.Info("Habit created", "habitId", habit.ID, "userId", userId)
	h.writeSuccessResponse(w, http.StatusCreated, habit)
}

func (h *HabitHandler) FindHabitById(w http.ResponseWriter, r *http.Request) {
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

	habit, err := h.service.FindHabitById(userId, habitId)
	if err != nil {
		slog.Error("Could not find habit by ID", "error", err, "habitId", habitId)
		h.writeErrorResponse(w, http.StatusNotFound, "Could not find habit by ID")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, habit)
}

func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeleteHabit(userId, habitId)
	if err != nil {
		slog.Error("Could not delete habit", "error", err, "habitId", habitId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not delete habit")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Successfully deleted habit"})
}

func (h *HabitHandler) UpdateHabit(w http.ResponseWriter, r *http.Request) {
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

	updatedHabit, err := h.service.UpdateHabit(userId, habitId, req)
	if err != nil {
		slog.Error("Could not update habit", "error", err, "habitId", habitId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not update habit")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, updatedHabit)
}

// Habit logs
func (h *HabitHandler) CreateHabitLog(w http.ResponseWriter, r *http.Request) {
	var logReq HabitLogReq
	err := json.NewDecoder(r.Body).Decode(&logReq)
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

	log, err := h.service.CreateHabitLog(userId, logReq)
	if err != nil {
		slog.Error("Failed to create log", "error", err, "userId", userId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not create log")
		return
	}

	slog.Info("Log created", "logId", log.ID, "userId", userId)
	h.writeSuccessResponse(w, http.StatusCreated, log)
}

func (h *HabitHandler) GetAllHabitLogs(w http.ResponseWriter, r *http.Request) {
	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	logs, err := h.service.GetAllHabitLogs(userId)
	if err != nil {
		slog.Error("Failed to get all logs", "error", err, "userId", userId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get all logs")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, logs)
}

func (h *HabitHandler) FindHabitLogById(w http.ResponseWriter, r *http.Request) {
	logId := chi.URLParam(r, "id")
	if logId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	log, err := h.service.FindHabitLogById(userId, logId)
	if err != nil {
		slog.Error("Could not find log by ID", "error", err, "logId", logId)
		h.writeErrorResponse(w, http.StatusNotFound, "Could not find log by ID")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, log)
}

func (h *HabitHandler) DeleteHabitLog(w http.ResponseWriter, r *http.Request) {
	logId := chi.URLParam(r, "id")
	if logId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	err = h.service.DeleteHabitLog(userId, logId)
	if err != nil {
		slog.Error("Could not delete log", "error", err, "logId", logId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not delete log")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Successfully deleted log"})
}

func (h *HabitHandler) UpdateHabitLog(w http.ResponseWriter, r *http.Request) {
	var req HabitLogReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("Failed to parse JSON", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not parse JSON")
		return
	}

	logId := chi.URLParam(r, "id")
	if logId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	userId, err := h.getUserId(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	updatedLog, err := h.service.UpdateHabitLog(userId, logId, req)
	if err != nil {
		slog.Error("Could not update log", "error", err, "logId", logId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not update log")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, updatedLog)
}
