// @title Notes API
// @version 1.0
// @description API для работы с заметками
// @host localhost:3000
// @BasePath /api/v1

package main

import (
	"log"
	"notes_project/handlers"
	"notes_project/models"
	"notes_project/routes"
	"notes_project/services"
	store "notes_project/services/store"
	"notes_project/services/store/config"
	database "notes_project/services/store/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gofiber/contrib/v3/swaggerui"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	_ "notes_project/docs"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := config.New()

	db, err := database.ConnectDB(conf)

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Note{})

	storage := store.NewPostgresStore(db)

	// if err := db.Ping(); err != nil {
	// 	log.Fatal("failed to connect", err)
	// }

	app := fiber.New()
	app.Use(cors.New(), swaggerui.New())
	cfg := swaggerui.Config{
		Next:     nil,
		BasePath: "/",
		FilePath: "./swagger.json",
		Path:     "docs",
		Title:    "Fiber API documentation",
		CacheAge: 3600, // Default to 1 hour
	}
	app.Use(swaggerui.New(cfg))
	// app.Use("/docs", static.New("./docs"))
	// app.Get("/swagger/*", adaptor.HTTPHandler(swagger.New()))
	// storage := store.NewCSVStore("notes.csv")

	service := services.NewNoteService(storage)

	noteHandler := handlers.NewNoteHandler(service)

	api := app.Group("/api/v1")
	routes.RegisterNoteRoutes(api, noteHandler)

	log.Fatal(app.Listen(":3000"))

}
