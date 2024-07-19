package userRoutes

import (
	userhandlers "notes_api/internal/handlers/userHandlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/user")

	// Route to register a user
	user.Post("/register", userhandlers.RegisterUser)
	
	// Route to login a user
	user.Post("/login", userhandlers.LoginUser)
}