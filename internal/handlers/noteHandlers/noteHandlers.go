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
	err := db.Order("updated_at desc").Find(notes).Error

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

// Handler for getting a single note
func GetNote(c *fiber.Ctx) error {
	db := database.DB
	note := &model.Note{}

	id := c.Params("id")

	err := db.Find(note, "id = ?", id).Error
	if err!=nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message" : "getting some issue fetching the notes",
				"error" : err,
		})
		return err
	}

	if note.ID == uuid.Nil{
		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message" : "no note present",
				"data" : nil,	
		})
	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message" : "note found with the id",
			"data" : note,
	})
}

// Handler for updating a note
func UpdateNote(c *fiber.Ctx) error {
	type UpdateNote struct {
		Title			*string	`json:"title"`  //using pointers to differentiate between missing and empty values
		Subtitle	*string	`json:"subtitle"`
		Text			*string	`json:"text"`
	}
	db := database.DB
	note := &model.Note{}

	id := c.Params("id")

	db.Find(note, "id = ?", id)

	if note.ID == uuid.Nil{
		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message" : "No note found with this id",
		})
	}

	updateNote := &UpdateNote{}
	err := c.BodyParser(updateNote)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message" : "Review your intput",
			"data" : err,
		})
	}

	if updateNote.Title != nil {
		note.Title = *updateNote.Title
	}
	if updateNote.Subtitle != nil {
		note.Subtitle = *updateNote.Subtitle
	}
	if updateNote.Text != nil {
		note.Text = *updateNote.Text
	}

	db.Save(note)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message" : "Updated note successfully",
		"data" : note,
	})
}

// Handler for deleting a note
func DeleteNote(c *fiber.Ctx) error {
	db := database.DB
	note := &model.Note{}

	id := c.Params("id")

	err := db.Find(note, "id = ?", id).Error
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message" : "Issue while finding note",
				"error" : err,
		})
	}

	if note.ID == uuid.Nil {
		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message" : "no note found with this id",
		})
	}

	err = db.Delete(note, "id = ?", id).Error
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message" : "Failed to delete note",
				"data" : err,
			})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message" : "Delete note successfully",
	})
}