package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/kaushal9900/url-shortner/database"
)

// ResolveURL handles the resolution of a shortened URL and redirects the user to the original URL.
// It uses Redis as the database to store the mappings between short URLs and their corresponding original URLs.
// The Redis client is created twice in this function to connect to different databases.

func ResolveURL(c *fiber.Ctx) error {
	// Get the short URL parameter from the request
	url := c.Params("url")

	// Create a Redis client to connect to database 0
	r := database.CreateClient(0)
	defer r.Close()

	// Retrieve the original URL associated with the short URL from the Redis database
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		// If the short URL is not found in the database, return a not found error response
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short URL not found in database"})
	}
	if err != nil {
		// If there is an error connecting to the database, return an internal server error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to database: " + err.Error()})
	}

	// Create a Redis client to connect to database 1
	rInr := database.CreateClient(1)
	defer rInr.Close()

	// Increment a counter in the second Redis database
	_ = rInr.Incr(database.Ctx, "counter")

	// Redirect the user to the original URL
	return c.Redirect(value, 301)
}
