package routes

import (
	"notes_project/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(app *fiber.App) error {
	app.Post("/notes", handlers.CreateNote)
	app.Get("/notes", handlers.GetNotes)

	return nil
}
