package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/config"
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

	usersHandler := user_handlers.New(cfg.Log, cfg.UsersService)

	cfg.Router.Route(APIVersion1URL, func(r chi.Router) {
		// User management endpoints.
		r.Route("/users", func(r chi.Router) {
			r.Get("/", usersHandler.GetUsers)
			r.Post("/", usersHandler.CreateUser)
			r.Get("/{userID}", usersHandler.GetUser)
			r.Put("/{userID}", usersHandler.UpdateUser)
			r.Delete("/{userID}", usersHandler.DeleteUser)
			r.Post("/{userID}/activate", usersHandler.ActivateUser)
			r.Post("/{userID}/deactivate", usersHandler.DeactivateUser)
			r.Post("/{userID}/suspend", usersHandler.SuspendUser)
			r.Post("/{userID}/unsuspend", usersHandler.UnSuspendUser)
		})
	})
}
