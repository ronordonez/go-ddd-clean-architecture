package mssql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"go-architecture/internal/product/domain"
	apperrors "go-architecture/internal/shared/errors"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type productModel struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Price       float64        `db:"price"`
	Stock       int            `db:"stock"`
	Category    string         `db:"category"`
	Active      bool           `db:"active"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (id, name, description, price, stock, category, active, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	q := r.db.Rebind(query)
	_, err := r.db.ExecContext(ctx, q,
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
	query := `SELECT id, name, description, price, stock, category, active, created_at, updated_at FROM products WHERE id = ?`
	q := r.db.Rebind(query)

	var m productModel
	err := r.db.GetContext(ctx, &m, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	desc := ""
	if m.Description.Valid {
		desc = m.Description.String
	}

	price, err := domain.NewPrice(m.Price)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: desc,
		Price:       price,
		Stock:       m.Stock,
		Category:    m.Category,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, filters domain.ProductFilters) ([]*domain.Product, error) {
	var sb strings.Builder
	sb.WriteString("SELECT id, name, description, price, stock, category, active, created_at, updated_at FROM products WHERE 1=1")
	args := []interface{}{}

	if filters.Category != "" {
		sb.WriteString(" AND category = ?")
		args = append(args, filters.Category)
	}
	if filters.Active != nil {
		sb.WriteString(" AND active = ?")
		args = append(args, *filters.Active)
	}

	// SQL Server pagination requires ORDER BY when using OFFSET/FETCH
	sb.WriteString(" ORDER BY created_at DESC OFFSET ? ROWS FETCH NEXT ? ROWS ONLY")
	args = append(args, filters.Offset, filters.Limit)

	q := r.db.Rebind(sb.String())

	var models []productModel
	if err := r.db.SelectContext(ctx, &models, q, args...); err != nil {
		return nil, err
	}

	products := make([]*domain.Product, 0, len(models))
	for _, m := range models {
		desc := ""
		if m.Description.Valid {
			desc = m.Description.String
		}
		price, err := domain.NewPrice(m.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &domain.Product{
			ID:          m.ID,
			Name:        m.Name,
			Description: desc,
			Price:       price,
			Stock:       m.Stock,
			Category:    m.Category,
			Active:      m.Active,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		})
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET name = ?, description = ?, price = ?, stock = ?, category = ?, active = ?, updated_at = ? WHERE id = ?`
	q := r.db.Rebind(query)
	res, err := r.db.ExecContext(ctx, q,
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
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = ?`
	q := r.db.Rebind(query)
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *ProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT CASE WHEN EXISTS(SELECT 1 FROM products WHERE name = ?) THEN 1 ELSE 0 END`
	q := r.db.Rebind(query)
	var existsInt int
	if err := r.db.GetContext(ctx, &existsInt, q, name); err != nil {
		return false, err
	}
	return existsInt == 1, nil
}

// no helper needed; CreatedAt/UpdatedAt are scanned as time.Time by sqlx
