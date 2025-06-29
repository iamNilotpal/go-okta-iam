package group_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/iamBelugaa/iam/internal/models"
	group_service "github.com/iamBelugaa/iam/internal/services/group"
	"github.com/iamBelugaa/iam/pkg/response"
)

type Handler struct {
	log       *zap.SugaredLogger
	groupsSvc *group_service.Service
}

func New(log *zap.SugaredLogger, svc *group_service.Service) *Handler {
	return &Handler{log: log, groupsSvc: svc}
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("Create group request received")

	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("Failed to decode create group request", zap.Error(err))
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group, err := h.groupsSvc.CreateGroup(r.Context(), &req)
	if err != nil {
		h.log.Infow("Failed to create group", zap.Error(err), "name", req.Name)
		h.respondWithError(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Group created successfully", "groupId", group.ID, "name", group.Name)
	response.RespondSuccess(
		w, http.StatusCreated, fmt.Sprintf("Group '%s' created successfully", group.Name), group,
	)
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	h.log.Infow("Get groups request received")

	groups, err := h.groupsSvc.GetGroups(r.Context())
	if err != nil {
		h.log.Infow("Failed to get groups", zap.Error(err))
		h.respondWithError(w, "Failed to retrieve groups", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Groups retrieved successfully", zap.Int("count", len(groups)))
	response.RespondSuccess(w, http.StatusOK, "Success", groups)
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		h.respondWithError(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Get group request received", "groupId", groupID)

	group, err := h.groupsSvc.GetGroup(r.Context(), groupID)
	if err != nil {
		h.log.Infow("Failed to get group", zap.Error(err), "groupId", groupID)
		h.respondWithError(w, "Failed to retrieve group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Group retrieved successfully", zap.String("groupId", groupID))
	response.RespondSuccess(w, http.StatusOK, "Success", group)
}

func (h *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		h.respondWithError(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Update group request received", "groupId", groupID)

	var req models.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Infow("Failed to decode update group request", zap.Error(err))
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group, err := h.groupsSvc.UpdateGroup(r.Context(), groupID, &req)
	if err != nil {
		h.log.Infow("Failed to update group", zap.Error(err), "groupId", groupID)
		h.respondWithError(w, "Failed to update group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Group updated successfully", "groupId", groupID)
	response.RespondSuccess(w, http.StatusOK, "Group updated successfully", group)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		h.respondWithError(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Delete group request received", "groupId", groupID)

	if err := h.groupsSvc.DeleteGroup(r.Context(), groupID); err != nil {
		h.log.Infow("Failed to delete group", zap.Error(err), "groupId", groupID)
		h.respondWithError(w, "Failed to delete group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Group deleted successfully", "groupId", groupID)
	response.RespondSuccess(w, http.StatusOK, "Group deleted successfully", nil)
}

func (h *Handler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		h.respondWithError(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Get group members request received", "groupId", groupID)

	members, err := h.groupsSvc.GetGroupMembers(r.Context(), groupID)
	if err != nil {
		h.log.Infow("Failed to get group members", zap.Error(err), "groupId", groupID)
		h.respondWithError(w, "Failed to retrieve group members", http.StatusInternalServerError)
		return
	}

	h.log.Infow("Group members retrieved successfully", "groupId", groupID, "memberCount", len(members))
	response.RespondSuccess(w, http.StatusOK, "Success", members)
}

func (h *Handler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	userID := chi.URLParam(r, "userID")

	if groupID == "" || userID == "" {
		h.respondWithError(w, "Both Group ID and User ID are required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Add user to group request received", "groupId", groupID, "userId", userID)

	if err := h.groupsSvc.AddUserToGroup(r.Context(), groupID, userID); err != nil {
		h.log.Infow("Failed to add user to group", zap.Error(err), "groupId", groupID, "userId", userID)
		h.respondWithError(w, "Failed to add user to group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User added to group successfully", "groupId", groupID, "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User added to group successfully", nil)
}

func (h *Handler) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	userID := chi.URLParam(r, "userID")

	if groupID == "" || userID == "" {
		h.respondWithError(w, "Both Group ID and User ID are required", http.StatusBadRequest)
		return
	}

	h.log.Infow("Remove user from group request received", "groupId", groupID, "userId", userID)

	if err := h.groupsSvc.RemoveUserFromGroup(r.Context(), groupID, userID); err != nil {
		h.log.Infow("Failed to remove user from group", zap.Error(err), "groupId", groupID, "userId", userID)
		h.respondWithError(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	h.log.Infow("User removed from group successfully", "groupId", groupID, "userId", userID)
	response.RespondSuccess(w, http.StatusOK, "User removed from group successfully", nil)
}

func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response.RespondError(w, statusCode, "API_ERROR", message, nil)
}
