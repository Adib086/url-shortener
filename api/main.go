package main

import (
	"fmt"
	"os"

	"github.com/Adib086/url-shortener/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {

	app.Get("/:url", routes.ResolveUrl)
	app.Post("/api/v1", routes.ShortenUrl)

}

func main() {
	// Load .env file if it exists (optional for Docker environments)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
