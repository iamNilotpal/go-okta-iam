package group_handlers

import (
	group_service "github.com/iamNilotpal/iam/internal/services/group"
	"go.uber.org/zap"
)

type Handler struct {
	log       *zap.SugaredLogger
	groupsSvc *group_service.Service
}

func New(log *zap.SugaredLogger, svc *group_service.Service) *Handler {
	return &Handler{log: log, groupsSvc: svc}
}
