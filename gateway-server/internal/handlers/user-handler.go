package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UserHandler struct {
	Store *session.Store
}

func (uh *UserHandler) CreateDraftUser(ctx *fiber.Ctx) error {
	return nil
}
