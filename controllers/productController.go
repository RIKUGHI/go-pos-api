package controllers

import (
	"encoding/json"
	"time"

	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/gofiber/fiber/v2"
)

type ProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Products(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	products, err := getProducts(user.ID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to retrieve products"})
	}

	return c.JSON(fiber.Map{
		"data": mapToProductResponses(products),
	})
}

func ByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)
	product, err := getProductByID(id, user.ID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(fiber.Map{
		"data": mapToProductResponse(product),
	})
}

func Create(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	request := new(ProductRequest)

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product, err := createProduct(request, user.ID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.JSON(fiber.Map{
		"data": mapToProductResponse(product),
	})
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)

	request := new(ProductRequest)

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product, err := updateProduct(id, user.ID, request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to update product"})
	}

	return c.JSON(fiber.Map{
		"data": mapToProductResponse(product),
	})
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)

	if err := deleteProduct(id, user.ID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{})
}

func getProducts(userID uint) ([]models.Product, error) {
	var products []models.Product
	err := initializers.DB.Where("user_id = ?", userID).Order("id desc").Find(&products).Error
	return products, err
}

func getProductByID(id string, userID uint) (*models.Product, error) {
	var product models.Product
	err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).Take(&product).Error
	return &product, err
}

func createProduct(request *ProductRequest, userID uint) (*models.Product, error) {
	product := models.Product{Name: request.Name, Price: request.Price, UserId: userID}
	result := initializers.DB.Create(&product)
	return &product, result.Error
}

func updateProduct(id string, userID uint, request *ProductRequest) (*models.Product, error) {
	var product models.Product
	if err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).Take(&product).Error; err != nil {
		return nil, err
	}

	product.Name = request.Name
	product.Price = request.Price

	if err := initializers.DB.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func deleteProduct(id string, userID uint) error {
	var product models.Product
	result := initializers.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&product)
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to delete product")
	}
	return nil
}

func mapToProductResponses(products []models.Product) []ProductResponse {
	productResponses := []ProductResponse{}
	for _, p := range products {
		productResponses = append(productResponses, mapToProductResponse(&p))
	}
	return productResponses
}

func mapToProductResponse(product *models.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
