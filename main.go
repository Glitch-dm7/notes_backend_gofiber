package main

import (
	"log"
	"notes_api/database"
	"notes_api/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.ConnectDB()

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}