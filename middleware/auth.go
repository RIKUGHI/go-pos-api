package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *fiber.Ctx) error {
	tokenString := c.Cookies("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorization"})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorization"})
		}

		user := models.User{}
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorization"})
		}

		c.Locals("user", user)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorization"})
	}

	return c.Next()
}
