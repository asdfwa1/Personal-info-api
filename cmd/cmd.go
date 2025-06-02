package cmd

import (
	"Test_Task_EffMob/internal/config"
	"Test_Task_EffMob/internal/database"
	"Test_Task_EffMob/internal/database/migrations"
	"Test_Task_EffMob/internal/enricher"
	"Test_Task_EffMob/internal/handlers"
	mymiddleware "Test_Task_EffMob/internal/middleware"
	"Test_Task_EffMob/internal/repository"
	"Test_Task_EffMob/internal/service"
	"Test_Task_EffMob/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func StartServer() {
	logger.InitLogger(slog.LevelDebug)
	cfg := config.LoadCfg()

	err := migrations.RunMigrations(cfg, cfg.MigrationFlag)
	if err != nil {
		slog.Error("Failed execute migrations", "error", err)
		panic(err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		panic(err)
	}
	defer db.Close()

	repo, err := repository.NewPersonRepository(db.DB)
	if err != nil {
		slog.Error("Failed to create Person Pepository", "error", err)
	}
	enr := enricher.NewEnricher()
	svc := service.NewPersonService(repo, enr)
	handler := handlers.NewPersonHandler(svc)

	r := chi.NewRouter()
	r.Use(mymiddleware.RequestIDMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/people", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})

	slog.Info("Server starting", "port", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		slog.Error("Server failed", "error", err)
	}
}
