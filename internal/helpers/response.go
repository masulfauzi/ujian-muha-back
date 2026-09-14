package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return ctx.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(ctx *fiber.Ctx, statusCode int, message string, errors interface{}) error {
	return ctx.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
