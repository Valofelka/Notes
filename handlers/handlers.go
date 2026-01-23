package handlers

import (
	"notes_project/functions"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CreateNoteRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type UpdateNoteRequest struct {
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

func GetNoteID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}
	note, err := functions.ReadNoteId(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(note)
}

func UpdateNote(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid note id"})
	}

	var req UpdateNoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})

	}

	note := &functions.Note{
		Id:    id,
		Title: req.Title,
		Text:  req.Text,
	}

	if err := functions.UpdateNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(note)

}

func DeleteNote(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}
	note := &functions.Note{Id: id}

	if err := functions.DeleteNote(note); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)

}
