package main

import (
	"github.com/gofiber/fiber/v2",
	"Notes_project/Functions_Notes/functions/"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hi")
	})

	app.Listen(":3000")
}
