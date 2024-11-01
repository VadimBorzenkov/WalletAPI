package app

import (
	"github.com/VadimBorzenkov/WalletAPI/config"
	"github.com/VadimBorzenkov/WalletAPI/internal/db"
	"github.com/VadimBorzenkov/WalletAPI/internal/delivery/handler"
	"github.com/VadimBorzenkov/WalletAPI/internal/delivery/routes"
	"github.com/VadimBorzenkov/WalletAPI/internal/repository"
	"github.com/VadimBorzenkov/WalletAPI/internal/service"
	"github.com/VadimBorzenkov/WalletAPI/pkg/logger"
	"github.com/VadimBorzenkov/WalletAPI/pkg/migrator"
	"github.com/gofiber/fiber/v2"
)

func Run() {
	logger := logger.InitLogger()

	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	dbase := db.Init(config)
	defer func() {
		if err := db.Close(dbase); err != nil {
			logger.Errorf("Failed to close database: %v", err)
		}
	}()

	if err := migrator.RunDatabaseMigrations(dbase); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewApiWalletRepository(dbase, logger)

	svc := service.NewApiWalletService(repo, logger)

	handler := handler.NewApiWalletHandler(svc, logger)

	app := fiber.New()

	routes.SetupRoutes(app, handler)

	logger.Infof("Starting server on port %s", config.Port)
	if err := app.Listen(config.Port); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
