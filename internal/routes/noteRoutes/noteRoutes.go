package noteRoutes

import (
	noteHandler "notes_api/internal/handlers/noteHandlers"

	"github.com/gofiber/fiber/v2"
)

func SetupNoteRoutes(router fiber.Router){
	note := router.Group("/note")

	// Get all notes
	note.Get("/", noteHandler.GetNotes)

	// Create a note
	note.Post("/", noteHandler.CreateNotes)

	// Get a single note
	note.Get("/:id", noteHandler.GetNote)
}