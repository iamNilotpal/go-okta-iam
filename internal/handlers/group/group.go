package group_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/models"
	group_service "github.com/iamNilotpal/iam/internal/services/group"
	"github.com/iamNilotpal/iam/pkg/response"
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

func (h *Handler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response.RespondError(w, statusCode, "API_ERROR", message, nil)
}
