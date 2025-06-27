package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *fiber.Ctx, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

func InternalError(c *fiber.Ctx, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Message: message,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Message: message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Message: message,
	})
}
