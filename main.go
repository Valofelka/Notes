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
	"notes_project/services/store"

	fiberSwagger "github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"

	_ "notes_project/docs"
)

func main() {
	app := fiber.New()
	storage := store.NewCSVStore("notes.csv")

	noteService := services.NewNoteService(storage)

	noteHandler := handlers.NewNoteHandler(noteService)

	api := app.Group("/api/v1")
	routes.RegisterNoteRoutes(api, noteHandler)

	app.Get("/swagger/*", fiberSwagger.New())

	log.Fatal(app.Listen(":3000"))

}
