package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.SetOutput(os.Stdout)
	}

	Log.SetLevel(logrus.InfoLevel)
}

func FiberLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		Log.WithFields(logrus.Fields{
			"method": c.Method(),
			"url":    c.Path(),
			"ip":     c.IP(),
			"status": c.Response().StatusCode(),
			"error":  err,
		}).Info("Incoming request")
		return err
	}
}
