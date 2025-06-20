package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/config"
	group_handlers "github.com/iamNilotpal/iam/internal/handlers/group"
	role_handlers "github.com/iamNilotpal/iam/internal/handlers/role"
	user_handlers "github.com/iamNilotpal/iam/internal/handlers/user"
	group_service "github.com/iamNilotpal/iam/internal/services/group"
	role_service "github.com/iamNilotpal/iam/internal/services/role"
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
	RolesService  *role_service.Service
}

func Setup(cfg *Config) {
	// Standard middleware for RealIP, RequestID, Logger, Recoverer etc.
	cfg.Router.Use(middleware.RealIP)
	cfg.Router.Use(middleware.RequestID)
	cfg.Router.Use(middleware.Logger)
	cfg.Router.Use(middleware.Recoverer)

	userHandlers := user_handlers.New(cfg.Log, cfg.UsersService)
	groupHandlers := group_handlers.New(cfg.Log, cfg.GroupsService)
	roleHandlers := role_handlers.New(cfg.Log, cfg.RolesService)

	cfg.Router.Route(APIVersion1URL, func(r chi.Router) {
		// User management endpoints.
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandlers.GetUsers)
			r.Post("/", userHandlers.CreateUser)

			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", userHandlers.GetUser)
				r.Put("/", userHandlers.UpdateUser)
				r.Delete("/", userHandlers.DeleteUser)

				// User lifecycle actions.
				r.Post("/activate", userHandlers.ActivateUser)
				r.Post("/deactivate", userHandlers.DeactivateUser)
				r.Post("/suspend", userHandlers.SuspendUser)
				r.Post("/unsuspend", userHandlers.UnSuspendUser)

				// User roles sub-resource.
				r.Route("/roles", func(r chi.Router) {
					r.Get("/", roleHandlers.GetUserRoles)
					r.Put("/{roleID}", roleHandlers.AssignRoleToUser)
					r.Delete("/{roleID}", roleHandlers.UnassignRoleFromUser)
				})
			})
		})

		// Group management endpoints.
		r.Route("/groups", func(r chi.Router) {
			r.Get("/", groupHandlers.GetGroups)
			r.Post("/", groupHandlers.CreateGroup)

			r.Route("/{groupID}", func(r chi.Router) {
				r.Get("/", groupHandlers.GetGroup)
				r.Put("/", groupHandlers.UpdateGroup)
				r.Delete("/", groupHandlers.DeleteGroup)

				// Group members sub-resource.
				r.Route("/members", func(r chi.Router) {
					r.Get("/", groupHandlers.GetGroupMembers)
					r.Put("/{userID}", groupHandlers.AddUserToGroup)
					r.Delete("/{userID}", groupHandlers.RemoveUserFromGroup)
				})

				// Group roles sub-resource.
				r.Route("/roles", func(r chi.Router) {
					r.Get("/", roleHandlers.GetGroupRoles)
					r.Put("/{roleID}", roleHandlers.AssignRoleToGroup)
					r.Delete("/{roleID}", roleHandlers.UnassignRoleFromGroup)
				})
			})
		})

		// Role management endpoints.
		r.Route("/roles", func(r chi.Router) {
			r.Get("/", roleHandlers.GetRoles)
			r.Post("/", roleHandlers.CreateRole)

			r.Route("/{roleID}", func(r chi.Router) {
				r.Get("/", roleHandlers.GetRole)
				r.Put("/", roleHandlers.UpdateRole)
				r.Delete("/", roleHandlers.DeleteRole)
			})
		})
	})
}
