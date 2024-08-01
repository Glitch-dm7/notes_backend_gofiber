package authMiddleware

import (
	"notes_api/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTProtected() fiber.Handler {
	secretKey := config.Config("JWT_SECRET_KEY")

	return jwtware.New(jwtware.Config{
		SigningKey : []byte(secretKey),
	})
}