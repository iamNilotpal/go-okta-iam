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

func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response.RespondError(w, statusCode, "API_ERROR", message, nil)
}
