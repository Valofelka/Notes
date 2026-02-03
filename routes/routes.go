package routes

import (
	"notes_project/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(router fiber.Router, handler *handlers.NoteHandler) {
	router.Post("/notes", handler.CreateNote)
	router.Get("/notes", handler.GetAllNotes)
	router.Get("/notes/:id", handler.GetNoteByID)
	router.Put("notes/:id", handler.UpdateNote)
	router.Delete("/notes/:id", handler.DeleteNote)

}
