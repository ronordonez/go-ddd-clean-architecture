package application

import (
	"context"
	"errors"

	"go-architecture/internal/product/domain"
	apperrors "go-architecture/internal/shared/errors"
)

type ProductService struct {
	repo      domain.ProductRepository
	validator *Validator
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo:      repo,
		validator: NewValidator(),
	}
}

func (s *ProductService) Create(ctx context.Context, dto CreateProductDTO) (*ProductResponseDTO, error) {
	// Validate DTO
	if err := s.validator.Validate(dto); err != nil {
		return nil, err
	}

	// Check if product with same name exists
	exists, err := s.repo.ExistsByName(ctx, dto.Name)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check product existence", err)
	}
	if exists {
		return nil, apperrors.NewAppError(409, "Product with this name already exists", apperrors.ErrConflict)
	}

	// Create domain entity
	product, err := domain.NewProduct(dto.Name, dto.Description, dto.Price, dto.Stock, dto.Category)
	if err != nil {
		return nil, apperrors.NewValidationError(err.Error(), nil)
	}

	// Save to repository
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, apperrors.NewInternalError("Failed to create product", err)
	}

	response := ToProductResponseDTO(product)
	return &response, nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*ProductResponseDTO, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFoundError("Product not found")
		}
		return nil, apperrors.NewInternalError("Failed to get product", err)
	}

	response := ToProductResponseDTO(product)
	return &response, nil
}

func (s *ProductService) GetAll(ctx context.Context, filtersDTO ProductListFiltersDTO) ([]ProductResponseDTO, error) {
	// Validate filters
	if err := s.validator.Validate(filtersDTO); err != nil {
		return nil, err
	}

	// Set default pagination
	if filtersDTO.Limit == 0 {
		filtersDTO.Limit = 20
	}

	filters := domain.ProductFilters{
		Category: filtersDTO.Category,
		Active:   filtersDTO.Active,
		Limit:    filtersDTO.Limit,
		Offset:   filtersDTO.Offset,
	}

	products, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get products", err)
	}

	return ToProductResponseDTOList(products), nil
}

func (s *ProductService) Update(ctx context.Context, id string, dto UpdateProductDTO) (*ProductResponseDTO, error) {
	// Validate DTO
	if err := s.validator.Validate(dto); err != nil {
		return nil, err
	}

	// Get existing product
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFoundError("Product not found")
		}
		return nil, apperrors.NewInternalError("Failed to get product", err)
	}

	// Update domain entity
	if err := product.Update(dto.Name, dto.Description, dto.Price, dto.Stock, dto.Category); err != nil {
		return nil, apperrors.NewValidationError(err.Error(), nil)
	}

	// Save changes
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, apperrors.NewInternalError("Failed to update product", err)
	}

	response := ToProductResponseDTO(product)
	return &response, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	// Check if product exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NewNotFoundError("Product not found")
		}
		return apperrors.NewInternalError("Failed to get product", err)
	}

	// Delete product
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.NewInternalError("Failed to delete product", err)
	}

	return nil
}
