package climbs

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/go-chi/chi/v5"
)

type ClimbHandler struct {
	service *ClimbService
}

func NewClimbHandler(service *ClimbService) *ClimbHandler {
	return &ClimbHandler{
		service: service,
	}
}
func (h *ClimbHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *ClimbHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *ClimbHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get claims from clerk
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Warn("No session claims in context")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user from claims
	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		slog.Error("Failed to get user from Clerk", "error", err, "subject", claims.Subject)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	// Get all climbs for user
	climbs, err := h.service.GetAll(usr.ID)
	if err != nil {
		slog.Error("Failed to get all climbs", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get all climbs")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(climbs)
}

func (h *ClimbHandler) Create(w http.ResponseWriter, r *http.Request) {
	var climb ClimbModel
	err := json.NewDecoder(r.Body).Decode(&climb)
	if err != nil {
		slog.Error("Failed to parse JSON", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not parse JSON")
		return
	}

	// Get claims from clerk
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Warn("No session claims in context")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user from claims
	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		slog.Error("Failed to get user from Clerk", "error", err, "subject", claims.Subject)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}
	climb.UserId = usr.ID

	out, err := h.service.Create(climb)
	if err != nil {
		slog.Error("Failed to create climb", "error", err, "userId", climb.UserId)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not create climb")
		return
	}

	slog.Info("Climb created", "climbId", out.ID, "userId", out.UserId)
	h.writeSuccessResponse(w, http.StatusCreated, out)
}

func (h *ClimbHandler) FindById(w http.ResponseWriter, r *http.Request) {
	// get the climb id from params
	climbId := chi.URLParam(r, "climbId")

	// If we get no ID from url
	if climbId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get ID from URL")
		return
	}

	// Get claims from clerk
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Warn("No session claims in context")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user from claims
	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		slog.Error("Failed to get user from Clerk", "error", err, "subject", claims.Subject)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	// Get climb
	climb, err := h.service.FindById(usr.ID, climbId)
	if err != nil {
		slog.Error("Could not find climb by ID", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not find climb by ID")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, climb)
}

func (h *ClimbHandler) Delete(w http.ResponseWriter, r *http.Request) {

	// get the climb id from params
	climbId := chi.URLParam(r, "climbId")

	// If we get no ID from url
	if climbId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get ID from URL")
		return
	}

	// Get claims from clerk
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Warn("No session claims in context")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user from claims
	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		slog.Error("Failed to get user from Clerk", "error", err, "subject", claims.Subject)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	// Get climb
	err = h.service.Delete(usr.ID, climbId)
	if err != nil {
		slog.Error("Could not delete climb", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not delete climb")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, "Successfully deleted climb")
}

func (h *ClimbHandler) Update(w http.ResponseWriter, r *http.Request) {
	var data ClimbUpdateReq
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("Failed to parse JSON", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not parse JSON")
		return
	}

	// get the climb id from params
	climbId := chi.URLParam(r, "climbId")

	// If we get no ID from url
	if climbId == "" {
		slog.Warn("Could not get ID from URL")
		h.writeErrorResponse(w, http.StatusBadRequest, "Could not get ID from URL")
		return
	}

	// Get claims from clerk
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Warn("No session claims in context")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user from claims
	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		slog.Error("Failed to get user from Clerk", "error", err, "subject", claims.Subject)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not get user")
		return
	}

	updatedClimb, err := h.service.Update(usr.ID, climbId, data)
	if err != nil {
		slog.Error("Could not update climb", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Could not update climb")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, updatedClimb)
}
