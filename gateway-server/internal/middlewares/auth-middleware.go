package middlewares

import (
	"formatier-api/shared/ierror"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
}

func (am *AuthMiddleware) Use(ctx *fiber.Ctx) (bool, int, *ierror.IErrorSchema) {
	return true, 0, nil
}
