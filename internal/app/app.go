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

// Run инициализирует и запускает сервер приложения
func Run() {
	// Инициализация логгера для логирования событий приложения
	logger := logger.InitLogger()

	// Загрузка конфигурации из файла config
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация подключения к базе данных
	dbase := db.Init(config)
	defer func() {
		// Закрытие подключения к базе данных при завершении работы приложения
		if err := db.Close(dbase); err != nil {
			logger.Errorf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	// Выполнение миграций базы данных для настройки необходимых таблиц
	if err := migrator.RunDatabaseMigrations(dbase); err != nil {
		logger.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	// Создание нового репозитория для работы с базой данных
	repo := repository.NewApiWalletRepository(dbase, logger)

	// Инициализация сервисного уровня с репозиторием и логгером
	svc := service.NewApiWalletService(repo, logger)

	// Настройка обработчиков API для обработки запросов
	handler := handler.NewApiWalletHandler(svc, logger)

	// Инициализация приложения Fiber для маршрутизации
	app := fiber.New()

	// Регистрация маршрутов API в приложении
	routes.SetupRoutes(app, handler)

	// Запуск сервера на указанном порту из конфигурации
	logger.Infof("Запуск сервера на порту %s", config.Port)
	if err := app.Listen(config.Port); err != nil {
		logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
