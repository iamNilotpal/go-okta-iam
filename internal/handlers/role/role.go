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

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	if roleID == "" {
		h.respondWithError(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Update role request received", "roleId", roleID)

	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("Failed to decode update role request", zap.Error(err))
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role, err := h.rolesSvc.UpdateRole(r.Context(), roleID, &req)
	if err != nil {
		h.log.Infow("Failed to update role", zap.Error(err), "roleId", roleID)
		h.respondWithError(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role updated successfully", "roleId", roleID)
	response.RespondSuccess(w, http.StatusOK, "Role updated successfully", role)
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	if roleID == "" {
		h.respondWithError(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Delete role request received", "roleId", roleID)

	if err := h.rolesSvc.DeleteRole(r.Context(), roleID); err != nil {
		h.log.Infow("Failed to delete role", zap.Error(err), "roleId", roleID)
		h.respondWithError(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role deleted successfully", "roleId", roleID)
	response.RespondSuccess(w, http.StatusOK, "Role deleted successfully", nil)
}

func (h *Handler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	userID := chi.URLParam(r, "userID")

	if roleID == "" || userID == "" {
		h.respondWithError(w, "Both Role ID and User ID are required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Assign role to user request received", "roleId", roleID, "userId", userID)

	if err := h.rolesSvc.AssignRoleToUser(r.Context(), userID, roleID); err != nil {
		h.log.Infow("Failed to assign role to user", zap.Error(err), "roleId", roleID, "userId", userID)
		h.respondWithError(w, "Failed to assign role to user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role assigned to user successfully", "roleId", roleID, "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "Role assigned to user successfully", nil)
}

func (h *Handler) UnassignRoleFromUser(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	userID := chi.URLParam(r, "userID")

	if roleID == "" || userID == "" {
		h.respondWithError(w, "Both Role ID and User ID are required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Unassign role from user request received", "roleId", roleID, "userId", userID)

	if err := h.rolesSvc.UnassignRoleFromUser(r.Context(), userID, roleID); err != nil {
		h.log.Infow("Failed to unassign role from user", zap.Error(err), "roleId", roleID, "userId", userID)
		h.respondWithError(w, "Failed to unassign role from user", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role unassigned from user successfully", "roleId", roleID, "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "Role unassigned from user successfully", nil)
}

func (h *Handler) AssignRoleToGroup(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	groupID := chi.URLParam(r, "groupID")

	if roleID == "" || groupID == "" {
		h.respondWithError(w, "Both Role ID and Group ID are required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Assign role to group request received", "roleId", roleID, "groupId", groupID)

	if err := h.rolesSvc.AssignRoleToGroup(r.Context(), groupID, roleID); err != nil {
		h.log.Infow("Failed to assign role to group", zap.Error(err), "roleId", roleID, "groupId", groupID)
		h.respondWithError(w, "Failed to assign role to group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role assigned to group successfully", "roleId", roleID, "groupId", groupID)
	response.RespondSuccess(w, http.StatusOK, "Role assigned to group successfully", nil)
}

func (h *Handler) UnassignRoleFromGroup(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	groupID := chi.URLParam(r, "groupID")

	if roleID == "" || groupID == "" {
		h.respondWithError(w, "Both Role ID and Group ID are required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Unassign role from group request received", "roleId", roleID, "groupId", groupID)

	if err := h.rolesSvc.UnassignRoleFromGroup(r.Context(), groupID, roleID); err != nil {
		h.log.Infow("Failed to unassign role from group", zap.Error(err), "roleId", roleID, "groupId", groupID)
		h.respondWithError(w, "Failed to unassign role from group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Role unassigned from group successfully", "roleId", roleID, "groupId", groupID)
	response.RespondSuccess(w, http.StatusOK, "Role unassigned from group successfully", nil)
}

func (h *Handler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondWithError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Get user roles request received", "userId", userID)

	roles, err := h.rolesSvc.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.log.Infow("Failed to get user roles", zap.Error(err), "userId", userID)
		h.respondWithError(w, "Failed to retrieve user roles", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User roles retrieved successfully", "userId", userID, "roleCount", len(roles))
	response.RespondSuccess(w, http.StatusOK, "Success", roles)
}

func (h *Handler) GetGroupRoles(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		h.respondWithError(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Get group roles request received", "groupId", groupID)

	roles, err := h.rolesSvc.GetGroupRoles(r.Context(), groupID)
	if err != nil {
		h.log.Infow("Failed to get group roles", zap.Error(err), "groupId", groupID)
		h.respondWithError(w, "Failed to retrieve group roles", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Group roles retrieved successfully", "groupId", groupID, "roleCount", len(roles))
	response.RespondSuccess(w, http.StatusOK, "Success", roles)
}

func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response.RespondError(w, statusCode, "API_ERROR", message, nil)
}
