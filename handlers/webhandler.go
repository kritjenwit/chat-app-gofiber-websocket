package handlers

import "github.com/gofiber/fiber/v2"

func IndexHandler(c *fiber.Ctx) error {
	return c.SendString("<h1>Chat APP | Backend</h1>")
}
