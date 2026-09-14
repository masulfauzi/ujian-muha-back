package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:       "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		Done: func(c *fiber.Ctx, logString []byte) {
			log.Println(string(logString))
		},
	})
}
