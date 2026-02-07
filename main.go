// @title Notes API
// @version 1.0
// @description API для работы с заметками
// @host localhost:3000
// @BasePath /api/v1

package main

import (
	"log"
	"notes_project/handlers"
	"notes_project/routes"
	"notes_project/services"

	fiberSwagger "github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"

	_ "notes_project/docs"
)

func main() {
	app := fiber.New()

	noteService := services.NewNoteService("notes.csv")

	noteHandler := handlers.NewNoteHandler(noteService)

	api := app.Group("/api/v1")
	routes.RegisterNoteRoutes(api, noteHandler)

	app.Get("/swagger/*", fiberSwagger.New())

	log.Fatal(app.Listen(":3000"))

}
