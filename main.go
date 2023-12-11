package main

import (
	"time"

	"github.com/RIKUGHI/go-pos-api/controllers"
	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/middleware"
	"github.com/gofiber/fiber/v2"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnnectToDb()
	// initializers.SyncDb()
}

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Prefork:      true,
	})

	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)
	app.Get("/validate", middleware.Auth, controllers.Validate)
	app.Get("/logout", middleware.Auth, controllers.Logout)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
