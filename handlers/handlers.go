package handlers

import (
	"notes_project/models"
	"notes_project/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	service *services.NoteService
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{}
}

// Обработчики HTTP-запросов
func CreateNote(c *fiber.Ctx) error { //c *fiber.Ctx - контекст запроса
	var req models.CreateNoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest). //HTTP- статус ответа
								JSON(fiber.Map{"error": "invalid request body"})
	}

	note := &models.Note{}
	services.Create(note, req.Title, req.Text)

	if err := services.AddNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

func GetNotes(c *fiber.Ctx) error {
	notes, err := services.ReadNote()
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
	note, err := services.ReadNoteId(id)
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

	var req models.UpdateNoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})

	}

	note := &models.Note{
		Id:    id,
		Title: req.Title,
		Text:  req.Text,
	}

	if err := services.UpdateNote(note); err != nil {
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
	note := &models.Note{Id: id}

	if err := services.DeleteNote(note); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)

}
