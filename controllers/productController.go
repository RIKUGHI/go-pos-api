package controllers

import (
	"encoding/json"

	"github.com/RIKUGHI/go-pos-api/helper"
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/RIKUGHI/go-pos-api/models/dto"
	"github.com/RIKUGHI/go-pos-api/services"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	*services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (pc *ProductController) FindAll(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	products, err := pc.ProductService.FindAll(user.ID)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to retrieve products")
	}

	return c.JSON(fiber.Map{
		"data": helper.MapToProductResponses(products),
	})
}

func (pc *ProductController) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)
	product, err := pc.ProductService.FindByID(id, user.ID)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Product not found")
	}

	return c.JSON(fiber.Map{
		"data": helper.MapToProductResponse(product),
	})
}

func (pc *ProductController) Create(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*models.User)
	request := new(dto.ProductDTO)

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	product, err := pc.ProductService.Create(request, user.ID)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to create product")
	}

	return c.JSON(fiber.Map{
		"data": helper.MapToProductResponse(product),
	})
}

func (pc *ProductController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)

	request := new(dto.ProductDTO)

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	product, err := pc.ProductService.Update(id, user.ID, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to update product")
	}

	return c.JSON(fiber.Map{
		"data": helper.MapToProductResponse(product),
	})
}

func (pc *ProductController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	user, _ := c.Locals("user").(*models.User)

	if err := pc.ProductService.Delete(id, user.ID); err != nil {
		return err
	}

	return c.JSON(fiber.Map{})
}
