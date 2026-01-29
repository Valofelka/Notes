package main

import (
	"fmt"
	"log"
	"notes_project/routes"
	"notes_project/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	noteService, err := services.NewNoteService("notes.csv")
	if err != nil {
		log.Fatalf("Invalid server", err)
	}
	routes.RegisterNoteRoutes(app, handlers)
	log.Fatal(app.Listen(":3000"))

}
