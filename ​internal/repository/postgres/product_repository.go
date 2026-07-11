package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"zuhtroai/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type productRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{
		db:     db,
		tracer: otel.Tracer("product.repository"),
	}
}

// WithTx — оставляем без изменений
func (r *productRepository) WithTx(ctx context.Context, fn func(repo domain.ProductRepository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	txRepo := &productRepository{db: tx, tracer: r.tracer}
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.Create")
	defer span.End()

	query := `INSERT INTO products (id, seller_id, title, price, description, category) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at`

	err := r.db.QueryRowContext(ctx, query,
		p.ID, p.SellerID, p.Title, p.Price, p.Description, p.Category,
	).Scan(&p.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("create product failed: %w", err)
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.GetByID")
	defer span.End()

	query := `SELECT id, seller_id, title, price, description, category, created_at, updated_at 
	          FROM products WHERE id = $1`
	
	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SellerID, &p.Title, &p.Price, &p.Description, &p.Category, &p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func (r *productRepository) Update(ctx context.Context, p *domain.Product) error {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.Update")
	defer span.End()

	query := `UPDATE products SET title = $1, price = $2, description = $3, category = $4, updated_at = NOW() 
	          WHERE id = $5 RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, p.Title, p.Price, p.Description, p.Category, p.ID).Scan(&p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (r *productRepository) Archive(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE products SET status = 'archived' WHERE id = $1", id)
	return err
}

func (r *productRepository) Restore(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE products SET status = 'active' WHERE id = $1", id)
	return err
}

// Метод GetBySeller (бывший GetByManufacturer)
func (r *productRepository) GetBySeller(ctx context.Context, opts domain.ProductListOptions) ([]*domain.Product, int, error) {
    // Здесь должна быть логика поиска по SellerID вместо ManufacturerID
    return nil, 0, nil 
}
