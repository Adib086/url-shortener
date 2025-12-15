package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/Adib086/url-shortener/database"
	"github.com/Adib086/url-shortener/server"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	Url string `json:"url"`

	CustomShort string `json:"custom_short"`

	Expiry time.Duration `json:"expirey"`
}

type response struct {
	Url string `json:"url"`

	CustomShort string `json:"custom_short"`

	Expiry time.Duration `json:"expirey"`

	RateRemianing  int           `json:"rate_remaining"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenUrl(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	r2 := database.Client(1)
	defer func() {
		_ = r2.Close()
	}()

	value, err := r2.Get(database.Ctx, c.IP()).Result()

	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
	} else {
		value, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(value)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "rate limit exceeded", "rate_limit_reset": limit.String()})

		}
	}

	if !govalidator.IsURL(body.Url) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	if !server.RemoveProhibitedUrls(body.Url) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "this URL is not allowed",
		})
	}

	body.Url = server.EnforceHTTPS(body.Url)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.Client(0)
	defer func() {
		_ = r.Close()
	}()

	value, _ = r.Get(database.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "custom short URL is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	}
	err = r.Set(database.Ctx, id, body.Url, body.Expiry).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect to the database",
		})
	}

	resp := response{
		Url:            body.Url,
		CustomShort:    id,
		Expiry:         body.Expiry,
		RateLimitReset: 30,
		RateRemianing:  10,
	}
	r2.Decr(database.Ctx, c.IP())

	val, _ := r2.Get(database.Ctx, id).Result()
	resp.RateRemianing, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp.RateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(200).JSON(resp)
}
