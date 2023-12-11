package controllers

import (
	"encoding/json"
	"os"
	"time"

	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	request := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	err := json.Unmarshal(c.Body(), request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := models.User{Email: request.Email, Password: string(hashed)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.JSON(fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	request := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user := models.User{}
	initializers.DB.First(&user, "email = ?", request.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create token " + err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		Domain:   "",
		Path:     "",
		Secure:   false,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func Validate(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": c.Locals("user"),
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		MaxAge:   -1,
		Domain:   "",
		Path:     "",
		Secure:   false,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.JSON(fiber.Map{"message": "Logout successful"})
}
