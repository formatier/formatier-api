package main

import (
	"gateway-server/internal/handlers"
	"gateway-server/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := newApp()
	app.Listen(":3000")
}

func newApp() *fiber.App {
	/*redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})*/
	app := fiber.New()
	app.Get("test", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Hello World")
	})
	/*store := session.New(session.Config{
		Expiration: 12 * time.Hour,
		Storage: cache.RedisStore{
			RedisClient: redisClient,
			Ctx:         context.Background(),
		},
	})*/
	routes.InitUserRoute(app, &handlers.UserHandler{
		//Store: store,
	})

	return app
}
