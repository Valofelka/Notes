package handlers

import (
	"notes_project/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	service *services.NoteService
}

type CreateNoteRequest struct {
	Title string `json: "title" example:"Заголовок"`
	Text  string `json: "text" example:"Текст заметки"`
}

type UpdateNoteRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// CreateNote godoc
// @Summary Создать заметку
// @Description Создает новую заметку
// @Tags notes
// @Accept json
// @Produce json
// @Param note body models.CreateNoteRequest true "Данные заметки"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Router /notes [post]
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

	var req UpdateNoteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})

	}

	note, err := h.service.UpdateNote(id, req.Title, req.Text)

	if err != nil {
		return c.Status(fiber.StatusNotFound).
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
