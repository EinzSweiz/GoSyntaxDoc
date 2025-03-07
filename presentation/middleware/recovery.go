package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				Log.WithFields(logrus.Fields{
					"method": c.Method(),
					"url":    c.Path(),
					"ip":     c.IP(),
					"error":  r,
				}).Error("Recovered from panic")

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
			}
		}()
		return c.Next()
	}
}
