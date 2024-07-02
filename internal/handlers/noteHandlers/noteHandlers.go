package noteHandler

import (
	"net/http"
	"notes_api/database"
	"notes_api/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler for fetching notes
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

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "some issue occured",
				"error": err,
			})
		return err
	}
		
	if len(*notes) == 0 {
		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message": "No notes found",
				"data" : notes,
			})
	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Notes found",
			"data": notes,
		})
}

// Handler for creating note
func CreateNotes(c *fiber.Ctx) error {
	db := database.DB
	note := model.Note{}

	err := c.BodyParser(&note)

	if err!=nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message" : "request failed, review your input",
				"data" : err,
			})
		return err
	}

	var missingFields []string
	if note.Title == "" {
		missingFields = append(missingFields, "Title")
	}
	if note.Text == "" {
		missingFields = append(missingFields, "Text")
	}

	if len(missingFields) > 0{
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message" : "The following fields are required and cannot be empty",
				"data" : missingFields,
		})
	}

	note.ID = uuid.New()

	err = db.Create(&note).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message" : "could not create your note",
				"data" : err,
			})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message" : "note created successfully",
		"data" : note,
	})
}
