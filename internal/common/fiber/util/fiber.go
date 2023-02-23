package util

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status int  `json:"status"`
	Data   any  `json:"data"`
	Ok     bool `json:"ok"`
}

func Send(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(Response{
		Status: status,
		Data:   data,
		Ok:     status < http.StatusBadRequest,
	})
}
