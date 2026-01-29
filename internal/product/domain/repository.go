package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context, filters ProductFilters) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type ProductFilters struct {
	Category string
	Active   *bool
	Limit    int
	Offset   int
}
