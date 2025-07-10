package routes

import (
	"formatier-api/shared/ierror"

	"github.com/gofiber/fiber/v2"
)

type RouteMiddleware func(*fiber.Ctx) (bool, int, ierror.IErrorSchema)

func NewRouter(app *fiber.App, prefix string) *Router {
	router := app.Group(prefix)
	return &Router{
		router: router,
	}
}

type Router struct {
	router            fiber.Router
	globalMiddlewares []RouteMiddleware
}

func (rr *Router) AddMiddleware(middleware RouteMiddleware) {
	rr.globalMiddlewares = append(rr.globalMiddlewares, middleware)
}

func runMiddleware(middleware RouteMiddleware) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		canActive, status, ierr := middleware(ctx)
		if !canActive {
			return ctx.Status(status).JSON(ierr)
		}
		return ctx.Next()
	}
}

func (rr *Router) Route(path, method string, handler fiber.Handler, middlewares ...RouteMiddleware) {
	handlersList := []fiber.Handler{}
	for _, globalMiddleware := range rr.globalMiddlewares {
		handlersList = append(
			handlersList,
			runMiddleware(globalMiddleware),
		)
	}
	for _, middleware := range middlewares {
		handlersList = append(
			handlersList,
			runMiddleware(middleware),
		)
	}

	handlersList = append(handlersList, handler)

	rr.router.Add(
		method,
		path,
		handlersList...,
	)
}
