package main

import (
	"time"

	"github.com/RIKUGHI/go-pos-api/controllers"
	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/middleware"
	"github.com/RIKUGHI/go-pos-api/services"
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
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error: " + err.Error(),
				"code":  fiber.StatusInternalServerError,
			})
		},
	})

	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)

	authGroup := app.Group("/auth", middleware.Auth, middleware.EnsureUser)

	authGroup.Get("/validate", controllers.Validate)
	authGroup.Get("/logout", controllers.Logout)

	productGroup := authGroup.Group("/products")
	productController := controllers.NewProductController(services.NewProductService())

	productGroup.Get("/", productController.FindAll)
	productGroup.Get("/:id", productController.FindByID)
	productGroup.Post("/", productController.Create)
	productGroup.Put("/:id", productController.Update)
	productGroup.Delete("/:id", productController.Delete)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
