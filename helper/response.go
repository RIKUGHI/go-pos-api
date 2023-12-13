package helper

import (
	"github.com/RIKUGHI/go-pos-api/models"
	"github.com/RIKUGHI/go-pos-api/models/api"
)

func MapToProductResponses(products []models.Product) []api.ProductResponse {
	productResponses := []api.ProductResponse{}
	for _, p := range products {
		productResponses = append(productResponses, MapToProductResponse(&p))
	}
	return productResponses
}

func MapToProductResponse(product *models.Product) api.ProductResponse {
	return api.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
