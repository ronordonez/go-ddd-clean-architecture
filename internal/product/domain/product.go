package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidProductName  = errors.New("product name must be between 3 and 100 characters")
	ErrInvalidPrice        = errors.New("price must be greater than 0")
	ErrInvalidStock        = errors.New("stock cannot be negative")
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       Price
	Stock       int
	Category    string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(name, description string, price float64, stock int, category string) (*Product, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	priceVO, err := NewPrice(price)
	if err != nil {
		return nil, err
	}

	if err := validateStock(stock); err != nil {
		return nil, err
	}

	now := time.Now()

	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       priceVO,
		Stock:       stock,
		Category:    category,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (p *Product) Update(name, description string, price float64, stock int, category string) error {
	if err := validateName(name); err != nil {
		return err
	}

	priceVO, err := NewPrice(price)
	if err != nil {
		return err
	}

	if err := validateStock(stock); err != nil {
		return err
	}

	p.Name = name
	p.Description = description
	p.Price = priceVO
	p.Stock = stock
	p.Category = category
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Product) Deactivate() {
	p.Active = false
	p.UpdatedAt = time.Now()
}

func (p *Product) Activate() {
	p.Active = true
	p.UpdatedAt = time.Now()
}

func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) IncreaseStock(quantity int) error {
	if quantity < 0 {
		return errors.New("quantity must be positive")
	}
	p.Stock += quantity
	p.UpdatedAt = time.Now()
	return nil
}

func validateName(name string) error {
	if len(name) < 3 || len(name) > 100 {
		return ErrInvalidProductName
	}
	return nil
}

func validateStock(stock int) error {
	if stock < 0 {
		return ErrInvalidStock
	}
	return nil
}
