package main

import (
	"log"
	"notes_project/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.RegisterNoteRoutes(app)
	log.Fatal(app.Listen(":3000"))

}
