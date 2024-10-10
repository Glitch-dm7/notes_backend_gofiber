package userhandlers

import (
	"notes_api/config"
	"notes_api/database"
	"notes_api/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Handler for registering the user
func RegisterUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(model.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message":"Invalid input"})
	}

	if user.Username == ""{
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Username is required"})
	}

	if user.Password == ""{
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password is required"})
	}

	existingUser := new(model.User)
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err==nil{
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message":"User already exists please user other username"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message" : "Failed to hash password"})
	}

	user.ID = uuid.New()
	user.Password = string(hashedPassword)

	if err:=db.Create(&user).Error; err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message":"Failed to create user"})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"message":"User created successfully"})
}

// Handler for loggin the user
func LoginUser(c *fiber.Ctx) error {
	db := database.DB
	input := new(model.User)

	if err:=c.BodyParser(&input); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message" : "Invalid input"})
	}

	if input.Username == ""{
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Username is required"})
	}

	if input.Password == ""{
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Password is required"})
	}

	user := new(model.User)
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message":"User not found, please register first"})
	}

	if err:= bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message":"Invalid Credentials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims:= token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	secretKey := config.Config("JWT_SECRET_KEY")

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message" : "Could not login"})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"token" : t})
}