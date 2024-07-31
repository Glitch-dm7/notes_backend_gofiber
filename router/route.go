package router

import (
	"notes_api/internal/routes/noteRoutes"
	"notes_api/internal/routes/userRoutes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	userRoutes.SetupUserRoutes(api)
	noteRoutes.SetupNoteRoutes(api)
}