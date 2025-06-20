package user_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/models"
	user_service "github.com/iamNilotpal/iam/internal/services/user"
	"github.com/iamNilotpal/iam/pkg/response"
)

type Handler struct {
	log      *zap.SugaredLogger
	usersSvc *user_service.Service
}

func New(log *zap.SugaredLogger, svc *user_service.Service) *Handler {
	return &Handler{log: log, usersSvc: svc}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("Create user request received")

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("Failed to decode create user request", zap.Error(err))
		h.respondWithError(w, "Invalid request body - please check your JSON format", http.StatusBadRequest)
		return
	}

	user, err := h.usersSvc.CreateUser(r.Context(), &req)
	if err != nil {
		h.log.Infow("Failed to create user", zap.Error(err), "email", req.Email)
		h.respondWithError(w, "Failed to create user - please try again", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User created successfully", "userId", user.ID, "email", user.Email)
	response.RespondSuccess(
		w, http.StatusCreated, fmt.Sprintf("User %s created successfully", user.Email), user,
	)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("Get users request received")

	users, err := h.usersSvc.GetUsers(r.Context())
	if err != nil {
		h.log.Infow("Failed to get users", zap.Error(err))
		h.respondWithError(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Users retrieved successfully", "count", len(users))
	response.RespondSuccess(w, http.StatusOK, "Success", users)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Get user request received", zap.String("userId", userID))

	user, err := h.usersSvc.GetUser(r.Context(), userID)
	if err != nil {
		h.log.Infow("Failed to get user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User retrieved successfully", "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "Success", user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Update user request received", "userId", userID)

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("Failed to decode update user request", zap.Error(err))
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.usersSvc.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		h.log.Infow("Failed to update user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User updated successfully", zap.String("userId", userID))
	response.RespondSuccess(w, http.StatusOK, "User updated successfully", user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Delete user request received", "userId", userID)

	err := h.usersSvc.DeleteUser(r.Context(), userID)
	if err != nil {
		h.log.Infow("Failed to delete user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User deleted successfully", "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User deleted successfully", nil)
}

func (h *Handler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Activate user request received", "userId", userID)

	err := h.usersSvc.ActivateUser(r.Context(), userID)
	if err != nil {
		h.log.Infow("Failed to activate user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to activate user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User activated successfully", "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User activated successfully", nil)
}

func (h *Handler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Deactivate user request received", "userId", userID)

	if err := h.usersSvc.DeactivateUser(r.Context(), userID); err != nil {
		h.log.Infow("Failed to deactivate user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to deactivate user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User deactivated successfully", "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User deactivated successfully", nil)
}

func (h *Handler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Suspend user request received", "userId", userID)

	if err := h.usersSvc.SuspendUser(r.Context(), userID); err != nil {
		h.log.Infow("Failed to suspend user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to suspend user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User suspended successfully", "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User suspended successfully", nil)
}

func (h *Handler) UnSuspendUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Unsuspend user request received", "userId", userID)

	if err := h.usersSvc.UnsuspendUser(r.Context(), userID); err != nil {
		h.log.Infow("Failed to unsuspend user", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to unsuspend user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User unsuspended successfully", "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User unsuspended successfully", nil)
}

func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response.RespondError(w, statusCode, "API_ERROR", message, nil)
}
