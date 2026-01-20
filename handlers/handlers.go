package handlers

import (
	"github.com/gofiber/fiber/v2"
	"notes_project/functions"
)

type CreateNoteRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Обработчики HTTP-запросов
func CreateNote(c *fiber.Ctx) error { //c *fiber.Ctx - контекст запроса
	var req CreateNoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest). //HTTP- статус ответа
								JSON(fiber.Map{"error": "invalid request body"})
	}

	note := &functions.Note{}
	functions.Create(note, req.Title, req.Text)

	if err := functions.AddNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

func GetNotes(c *fiber.Ctx) error {
	notes, err := functions.ReadNote()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notes)
}
