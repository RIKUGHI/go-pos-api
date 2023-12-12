package controllers

import (
	"encoding/json"
	"time"

	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
)

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Products(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	products := []models.Product{}

	err := initializers.DB.Where("user_id = ?", user.ID).Order("id desc").Find(&products).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var productResponses []ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, ProductResponse{
			ID:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"data": productResponses,
	})
}

func ByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)
	product := models.Product{}

	err := initializers.DB.Where("id = ? AND user_id = ?", id, user.ID).Take(&product).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	productResponse := ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	return c.JSON(fiber.Map{
		"data": productResponse,
	})
}

func Create(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	request := new(struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	})

	err := json.Unmarshal(c.Body(), request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product := models.Product{Name: request.Name, Price: request.Price, UserId: user.ID}
	result := initializers.DB.Create(&product)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.JSON(fiber.Map{
		"data": ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		},
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)
	product := models.Product{}
	request := new(struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	})

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := initializers.DB.Where("id = ? AND user_id = ?", id, user.ID).Take(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	product.Name = request.Name
	product.Price = request.Price

	if err := initializers.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		},
	})
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)
	product := models.Product{}

	isSuccess := initializers.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&product).RowsAffected

	if isSuccess == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to delete"})
	}

	return c.JSON(fiber.Map{})
}
