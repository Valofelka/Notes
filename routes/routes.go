package routes

import (
	"notes_project/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(app *fiber.App) {
	app.Post("/notes", handlers.CreateNote)
	app.Get("/notes", handlers.GetNotes)
	app.Get("/notes/:id", handlers.GetNoteID)
	app.Put("notes/:id", handlers.UpdateNote)
	app.Delete("/notes/:id", handlers.DeleteNote)

}
