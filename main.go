package main

import (
	"notes_project/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.RegisterNoteRoutes(app)
	app.Listen(":3000")

}
