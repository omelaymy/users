package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/omelaymy/users/internal/api"
	"github.com/rs/zerolog"
)

type HttpErrorHandler struct {
	log *zerolog.Logger
}

func NewHttpErrorHandler(
	log *zerolog.Logger,
) *HttpErrorHandler {
	return &HttpErrorHandler{
		log: log,
	}
}

func (h *HttpErrorHandler) Handler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return ctx.Status(code).JSON(api.ErrorResponse{
		Message: err.Error(),
	})
}
