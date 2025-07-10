package routes

import (
	"gateway-server/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(app *fiber.App, userHandler *handlers.UserHandler) {
	userRouter := NewRouter(app, "/user")
	userRouter.Route("/user", fiber.MethodGet, userHandler.CreateDraftUser)
}
