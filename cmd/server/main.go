package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/config"
	"github.com/iamNilotpal/iam/internal/handlers"
	group_service "github.com/iamNilotpal/iam/internal/services/group"
	user_service "github.com/iamNilotpal/iam/internal/services/user"
	"github.com/iamNilotpal/iam/pkg/logger"
	"github.com/iamNilotpal/iam/pkg/okta"
)

func main() {
	log := logger.New("flexera-iam")
	defer func() {
		if err := log.Sync(); err != nil {
			log.Infow("sync error", "error", err)
		}
	}()

	if err := godotenv.Load(); err != nil {
		log.Fatalw("error loading envs", "error", err)
	}

	log.Infow("Starting Flexera IAM Platform...")

	if err := run(log); err != nil {
		log.Infow("startup error", "error", err)
		if err := log.Sync(); err != nil {
			log.Infow("sync error", "error", err)
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log.Infow("Configuration loaded successfully")

	oktaClient, err := okta.NewClient(cfg.Okta)
	if err != nil {
		return err
	}

	if err := oktaClient.TestConnection(context.Background()); err != nil {
		return err
	}
	log.Infow("Okta service initialized successfully")

	router := chi.NewRouter()
	usersService := user_service.New(log, oktaClient.SDK())
	groupsService := group_service.New(log, oktaClient.SDK())

	handlers.Setup(&handlers.Config{
		Config: cfg,
		Log:    log,
		Router: router,

		UsersService:  usersService,
		GroupsService: groupsService,
	})

	server := http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	serverErrors := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infow("server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutting down server", "signal", sig)
		defer log.Infow("shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
