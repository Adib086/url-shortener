package routes

import (
	"github.com/Adib086/url-shortener/database"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.Client(0)

	defer func() {
		_ = r.Close()
	}()
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "URL not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Internal Server Error"})
	}
	rInr := database.Client(1)
	_ = rInr.Close()
	rInr.Incr(database.Ctx, "counter:"+url)

	return c.Redirect(value, fiber.StatusTemporaryRedirect)
}
