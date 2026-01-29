package application

import (
	"go-architecture/internal/product/domain"
)

func ToProductResponseDTO(product *domain.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price.Value(),
		Stock:       product.Stock,
		Category:    product.Category,
		Active:      product.Active,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToProductResponseDTOList(products []*domain.Product) []ProductResponseDTO {
	dtos := make([]ProductResponseDTO, len(products))
	for i, product := range products {
		dtos[i] = ToProductResponseDTO(product)
	}
	return dtos
}
