package role_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/models"
	role_service "github.com/iamNilotpal/iam/internal/services/role"
	"github.com/iamNilotpal/iam/pkg/response"
)

type Handler struct {
	log      *zap.SugaredLogger
	rolesSvc *role_service.Service
}

func New(log *zap.SugaredLogger, svc *role_service.Service) *Handler {
	return &Handler{log: log, rolesSvc: svc}
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("Create role request received")

	var req models.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("Failed to decode create role request", zap.Error(err))
		h.respondWithError(w, "Invalid request body - please check your JSON format", http.StatusBadRequest)
		return
	}

	role, err := h.rolesSvc.CreateRole(r.Context(), &req)
	if err != nil {
		h.log.Infow("Failed to create role", zap.Error(err), "name", req.Name)
		h.respondWithError(w, "Failed to create role - please try again", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role created successfully", "roleId", role.ID, "name", role.Name)
	response.RespondSuccess(
		w, http.StatusCreated, fmt.Sprintf("Role '%s' created successfully", role.Name), role,
	)
}

func (h *Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("Get roles request received")

	roles, err := h.rolesSvc.GetRoles(r.Context())
	if err != nil {
		h.log.Infow("Failed to get roles", zap.Error(err))
		h.respondWithError(w, "Failed to retrieve roles", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Roles retrieved successfully", "count", len(roles))
	response.RespondSuccess(w, http.StatusOK, "Success", roles)
}

func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	if roleID == "" {
		h.respondWithError(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Get role request received", "roleId", roleID)

	role, err := h.rolesSvc.GetRole(r.Context(), roleID)
	if err != nil {
		h.log.Infow("Failed to get role", zap.Error(err), "roleId", roleID)
		h.respondWithError(w, "Failed to retrieve role", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role retrieved successfully", "roleId", roleID)
	response.RespondSuccess(w, http.StatusOK, "Success", role)
}

func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response.RespondError(w, statusCode, "API_ERROR", message, nil)
}
