package services

import (
	"github.com/RIKUGHI/go-pos-api/initializers"
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/RIKUGHI/go-pos-api/models/dto"
	"github.com/gofiber/fiber/v2"
)

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (p ProductService) FindAll(userID uint) ([]models.Product, error) {
	var products []models.Product
	err := initializers.DB.Where("user_id = ?", userID).Order("id desc").Find(&products).Error
	return products, err
}

func (p ProductService) FindByID(id string, userID uint) (*models.Product, error) {
	var product models.Product
	err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).Take(&product).Error
	return &product, err
}

func (p ProductService) Create(request *dto.ProductDTO, userID uint) (*models.Product, error) {
	product := models.Product{Name: request.Name, Price: request.Price, UserId: userID}
	result := initializers.DB.Create(&product)
	return &product, result.Error
}

func (p ProductService) Update(id string, userID uint, request *dto.ProductDTO) (*models.Product, error) {
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

func (p ProductService) Delete(id string, userID uint) error {
	var product models.Product
	result := initializers.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&product)
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to delete product")
	}
	return nil
}
