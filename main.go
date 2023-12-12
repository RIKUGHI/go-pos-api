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
		// ErrorHandler: func(c *fiber.Ctx, err error) error {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "test error : " + err.Error()})
		// },
	})

	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)

	authGroup := app.Group("/auth", middleware.Auth, middleware.EnsureUser)

	authGroup.Get("/validate", controllers.Validate)
	authGroup.Get("/logout", controllers.Logout)

	productGroup := authGroup.Group("/products")

	productGroup.Get("/", controllers.Products)
	productGroup.Get("/:id", controllers.ByID)
	productGroup.Post("/", controllers.Create)
	productGroup.Put("/:id", controllers.Update)
	productGroup.Delete("/:id", controllers.Delete)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
