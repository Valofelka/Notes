package handlers

import (
	"notes_project/services"
	"strconv"

	fiber "github.com/gofiber/fiber/v3"
)

type NoteHandler struct {
	service *services.NoteService
}

type CreateNoteRequest struct {
	Title string `json:"title" example:"Заголовок"`
	Text  string `json:"text" example:"Текст заметки"`
}

type UpdateNoteRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

// CreateNote godoc
// @Summary Создать заметку
// @Description Создает новую заметку
// @Tags notes
// @Accept json
// @Produce json
// @Param note body handlers.CreateNoteRequest true "Данные заметки"
// @Success 201 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notes [post]
func (h *NoteHandler) CreateNote(c fiber.Ctx) error { //c fiber.Ctx - контекст запроса
	var req CreateNoteRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})
	}

	note := h.service.NewNote(req.Title, req.Text)

	if err := h.service.AddNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

// GetAllNotes godoc
// @Summary Вывод всех заметок
// @Description Возвращает список все заметок
// @Tags notes
// @Produce json
// @Success 200 {array} models.Note
// @Failure 500 {object} map[string]string
// @Router /notes [get]
func (h *NoteHandler) GetAllNotes(c fiber.Ctx) error {
	notes, err := h.service.GetAllNotes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notes)
}

// GetNoteByID godoc
// @Summary Вывод заметки
// @Description Возвращает одну заметку
// @Tags notes
// @Produce json
// @Param id path int true "ID заметки"
// @Success 200 {array} models.Note
// @Failure 500 {object} map[string]string
// @Router /notes/{id} [get]
func (h *NoteHandler) GetNoteByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	note, err := h.service.GetNoteByID(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(note)
}

// UpdateNote godoc
// @Summary Обновить заметку
// @Description Обновить заметку по ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID заметки"
// @Param note body handlers.CreateNoteRequest true "Данные заметки"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [put]
func (h *NoteHandler) UpdateNote(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid note id"})
	}

	var req UpdateNoteRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body"})

	}

	note, err := h.service.UpdateNote(uint(id), req.Title, req.Text)

	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(note)

}

// CreateNote godoc
// @Summary Удалить заметку
// @Description Удаляет заметку по id
// @Tags notes
// @Param id path int true "ID заметки"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Router /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid id"})
	}

	err = h.service.DeleteNote(uint(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)

}
