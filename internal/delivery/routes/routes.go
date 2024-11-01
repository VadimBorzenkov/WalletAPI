package routes

import (
	"github.com/VadimBorzenkov/WalletAPI/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes регистрирует маршруты приложения.
func SetupRoutes(app *fiber.App, h handler.WalletHandler) *fiber.App {
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Настройка CORS, чтобы разрешить доступ со всех доменов
	}))

	api := app.Group("/api/v1/wallets")
	api.Get("/:walletID", h.HandleBalance)
	api.Patch("/", h.HandleTransaction)

	return app
}
