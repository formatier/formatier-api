package main

import (
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := newApp()
	app.Listen(":3000")
}
