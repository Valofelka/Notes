package routes

import (
	"notes_project/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(app *fiber.App, handler *handlers.NoteHandler) {
	app.Post("/notes", handler.CreateNote)
	app.Get("/notes", handler.GetAllNotes)
	app.Get("/notes/:id", handler.GetNoteByID)
	app.Put("notes/:id", handler.UpdateNote)
	app.Delete("/notes/:id", handler.DeleteNote)

}
