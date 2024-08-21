package noteHandler

import (
	"net/http"
	"notes_api/database"
	"notes_api/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Handler for fetching notes
func GetNotes(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	notes := &[]model.Note{}

	// Fetch all the notes of the users in descending order
	if err := db.Where("user_id = ?", userID).Order("updated_at desc").Find(notes).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Failed to get notes"})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{"data": notes})
}

// Handler for creating note
func CreateNotes(c *fiber.Ctx) error {
	db := database.DB
	note := model.Note{}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	// If mismatching fields in the body we return an error
	if err := c.BodyParser(&note); err!=nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message" : "Invalid input"})
		return err
	}

	// Title and text are mandatory fields, without these fields in the body we throw an error
	var missingFields []string
	if note.Title == "" {
		missingFields = append(missingFields, "title")
	}
	if note.Text == "" {
		missingFields = append(missingFields, "text")
	}

	if len(missingFields) > 0{
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message" : "The following fields are required and cannot be empty",
				"data" : missingFields,
		})
	}

	// Create a new id of the note
	note.ID = uuid.New()
	note.UserID = uuid.MustParse(userID)

	if err := db.Create(&note).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message" : "could not create your note"})
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
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	id := c.Params("id")

	// Fetch the note with the id provided in the params
	if err := db.Find(note, "id = ? AND user_id = ?", id, userID).Error; err!=nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message" : "getting some issue fetching the notes",
		})
		return err
	}

	if note.ID == uuid.Nil{
		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message" : "no note present with this id",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{"data" : note})
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
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	id := c.Params("id")

	if err := db.First(note, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Note not found"})
	}

	updateNote := &UpdateNote{}
	if err := c.BodyParser(updateNote); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message" : "Review your intput",
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

	if err := db.Save(note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update note"})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message" : "Updated note successfully",
		"data" : note,
	})
}

// Handler for deleting a note
func DeleteNote(c *fiber.Ctx) error {
	db := database.DB
	note := &model.Note{}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	id := c.Params("id")

	if err := db.First(note, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Note not found"})
	}

	if err := db.Delete(note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete note"})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message" : "Deleted note successfully",
	})
}