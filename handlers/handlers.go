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

type CreateNoteRequest struct {
	Title string `json: "title"`
	Text  string `json: "text"`
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

// Обработчики HTTP-запросов
func (h *NoteHandler) CreateNote(c *fiber.Ctx) error { //c *fiber.Ctx - контекст запроса
	var req CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})
	}

	note := h.service.CreateNote(req.Title, req.Text)

	if err := h.service.AddNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

func (h *NoteHandler) GetAllNotes(c *fiber.Ctx) error {
	notes, err := h.service.GetAllNotes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notes)
}

func (h *NoteHandler) GetNoteByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	note, err := h.service.GetNoteByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(note)
}

func (h *NoteHandler) UpdateNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
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

	if err := h.service.UpdateNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(note)

}

func (h *NoteHandler) DeleteNote(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.service.DeleteNote(id); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)

}
