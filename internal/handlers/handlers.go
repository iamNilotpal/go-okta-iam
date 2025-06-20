package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/config"
	group_handlers "github.com/iamNilotpal/iam/internal/handlers/group"
	user_handlers "github.com/iamNilotpal/iam/internal/handlers/user"
	group_service "github.com/iamNilotpal/iam/internal/services/group"
	user_service "github.com/iamNilotpal/iam/internal/services/user"
)

const (
	APIVersion1URL = "/api/v1"
)

type Config struct {
	Router        *chi.Mux
	Config        *config.Config
	Log           *zap.SugaredLogger
	UsersService  *user_service.Service
	GroupsService *group_service.Service
}

func Setup(cfg *Config) {
	// Standard middleware for RealIP, RequestID, Logger, Recoverer etc.
	cfg.Router.Use(middleware.RealIP)
	cfg.Router.Use(middleware.RequestID)
	cfg.Router.Use(middleware.Logger)
	cfg.Router.Use(middleware.Recoverer)

	userHandlers := user_handlers.New(cfg.Log, cfg.UsersService)
	groupHandlers := group_handlers.New(cfg.Log, cfg.GroupsService)

	cfg.Router.Route(APIVersion1URL, func(r chi.Router) {
		// User management endpoints.
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandlers.GetUsers)
			r.Post("/", userHandlers.CreateUser)
			r.Get("/{userID}", userHandlers.GetUser)
			r.Put("/{userID}", userHandlers.UpdateUser)
			r.Delete("/{userID}", userHandlers.DeleteUser)
			r.Post("/{userID}/activate", userHandlers.ActivateUser)
			r.Post("/{userID}/deactivate", userHandlers.DeactivateUser)
			r.Post("/{userID}/suspend", userHandlers.SuspendUser)
			r.Post("/{userID}/unsuspend", userHandlers.UnSuspendUser)
		})

		// Group management endpoints.
		r.Route("/groups", func(r chi.Router) {
			r.Get("/", groupHandlers.GetGroups)
			r.Post("/", groupHandlers.CreateGroup)
			r.Get("/{groupID}", groupHandlers.GetGroup)
			r.Put("/{groupID}", groupHandlers.UpdateGroup)
			r.Delete("/{groupID}", groupHandlers.DeleteGroup)
			r.Get("/{groupID}/members", groupHandlers.GetGroupMembers)
			r.Post("/{groupID}/users/{userID}", groupHandlers.AddUserToGroup)
			r.Delete("/{groupID}/users/{userID}", groupHandlers.RemoveUserFromGroup)
		})
	})
}
