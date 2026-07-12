package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"zuhroai/internal/domain"
)

type productRepository struct {
	db *sql.DB
}

// NewProductRepository создает экземпляр с инъекцией зависимости БД
func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	const query = `
		INSERT INTO products (id, seller_id, title, description, price, category, status, stock_quantity, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	
	_, err := r.db.ExecContext(ctx, query, 
		p.ID, p.SellerID, p.Title, p.Description, p.Price, 
		p.Category, p.Status, p.StockQuantity, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	const query = `
		SELECT id, seller_id, title, description, price, category, status, stock_quantity, created_at, updated_at 
		FROM products WHERE id = $1`
	
	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, 
		&p.Category, &p.Status, &p.StockQuantity, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return p, nil
}

