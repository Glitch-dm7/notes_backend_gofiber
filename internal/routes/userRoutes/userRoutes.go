package userRoutes

import (
	userhandlers "notes_api/internal/handlers/userHandlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/user")

	user.Post("/register", userhandlers.RegisterUser)
	
	user.Post("/login", userhandlers.LoginUser)
}