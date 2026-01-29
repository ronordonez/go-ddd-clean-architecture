package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"go-architecture/internal/product/domain"
	apperrors "go-architecture/internal/shared/errors"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type productModel struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Stock       int       `db:"stock"`
	Category    string    `db:"category"`
	Active      bool      `db:"active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, stock, category, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price.Value(),
		product.Stock,
		product.Category,
		product.Active,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category, active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var model productModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return r.toDomain(&model)
}

func (r *ProductRepository) FindAll(ctx context.Context, filters domain.ProductFilters) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category, active, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if filters.Category != "" {
		query += ` AND category = $` + string(rune(argCount+48))
		args = append(args, filters.Category)
		argCount++
	}

	if filters.Active != nil {
		query += ` AND active = $` + string(rune(argCount+48))
		args = append(args, *filters.Active)
		argCount++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + string(rune(argCount+48)) + ` OFFSET $` + string(rune(argCount+49))
	args = append(args, filters.Limit, filters.Offset)

	var models []productModel
	err := r.db.SelectContext(ctx, &models, query, args...)
	if err != nil {
		return nil, err
	}

	products := make([]*domain.Product, len(models))
	for i, model := range models {
		product, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		products[i] = product
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, category = $5, active = $6, updated_at = $7
		WHERE id = $8
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price.Value(),
		product.Stock,
		product.Category,
		product.Active,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *ProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE name = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, name)
	return exists, err
}

func (r *ProductRepository) toDomain(model *productModel) (*domain.Product, error) {
	price, err := domain.NewPrice(model.Price)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       price,
		Stock:       model.Stock,
		Category:    model.Category,
		Active:      model.Active,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}
