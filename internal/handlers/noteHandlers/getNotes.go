package noteHandler

import (
	"net/http"
	"notes_api/database"
	"notes_api/internal/model"

	"github.com/gofiber/fiber/v2"
)

func GetNotes(c *fiber.Ctx) error {
	db := database.DB
	
	if db == nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message" : "db is not initialized",
				"log" : db,
			})
		}

	notes := &[]model.Note{}
	err := db.Find(notes).Error

	if len(*notes) == 0 {
		return c.Status(404).JSON(
			&fiber.Map{
				"message": "No notes found",
				"date" : notes,
			})
	}

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "some issue occured",
				"error": err,
			})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "notes found",
			"data": notes,
		})
	return nil
}
