package routes

import (
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kaushal9900/url-shortner/configs"
	"github.com/kaushal9900/url-shortner/database"
	"github.com/kaushal9900/url-shortner/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	ShortURL        string        `json:"shorturl"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL is responsible for handling the shortening of a URL.
// It receives a request body containing the URL to be shortened.
// The function performs various checks and operations before returning an error or success response.

func ShortenURL(c *fiber.Ctx) error {

	var id string
	// Parse the request body into the "body" struct
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		// Return a bad request error if the request body cannot be parsed
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	// Create a Redis client to connect to database 1
	rc := database.CreateClient(1)
	defer rc.Close()

	// Check if the client's IP address has exceeded the rate limit
	noOfRequest, err := rc.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		// If there are no previous requests from the client, set the request count to 1 and expire it after 30 minutes
		_ = rc.Set(database.Ctx, c.IP(), 1, time.Minute*30).Err()
	} else if err == nil {
		// If there are previous requests, compare the count with the rate limit
		rateLimit := configs.EnvConfigs.APIQuota
		noOfRequestInt, _ := strconv.Atoi(noOfRequest)
		if noOfRequestInt > rateLimit {
			// If the rate limit is exceeded, calculate the remaining time until the rate limit resets
			resetTime, _ := rc.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "rate limit exceeded",
				"rate_limit_reset": (resetTime / time.Nanosecond / time.Minute),
			})
		}
	}

	// Check if the URL is valid
	if !govalidator.IsURL(body.URL) {
		// If the URL is invalid, return a bad request error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// Check for circular redirection and infinite loops in the URL
	if !helpers.RemoveDomainError(body.URL) {
		// If a circular redirection or infinite loop is detected, return a service unavailable error
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Circular Redirection Infinite loop detected"})
	}

	// Enforce the usage of HTTP protocol in the URL
	body.URL = helpers.EnforceHTTP(body.URL)

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	urlExist, _ := r.Get(database.Ctx, id).Result()
	if urlExist != "" {
		id = uuid.New().String()[:6]
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	rerr := r.Set(database.Ctx, id, body.URL, body.Expiry*time.Second*3600)

	if rerr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	// Increment the request count for the client's IP address in Redis
	rc.Incr(database.Ctx, c.IP())

	resetTime, _ := rc.TTL(database.Ctx, c.IP()).Result()
	shortURL := configs.EnvConfigs.Domain + "/" + id
	rateRemaining, _ := strconv.Atoi(noOfRequest)
	resp := response{
		URL:             body.URL,
		ShortURL:        shortURL,
		Expiry:          body.Expiry,
		XRateRemaining:  rateRemaining - 1,
		XRateLimitReset: resetTime,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
