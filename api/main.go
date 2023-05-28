package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kaushal9900/url-shortner/configs"
	"github.com/kaushal9900/url-shortner/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	configs.InitEnvConfigs()

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	log.Fatal(app.Listen(configs.EnvConfigs.AppPort))

}
