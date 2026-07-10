package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"zuhroai/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

	query := `INSERT INTO products (id, manufacturer_id, name, price, description, category, version)
	VALUES ($1, $2, $3, $4, $5, $6, 1) RETURNING created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, p.ID, p.ManufacturerID, p.Name, p.Price, p.Description, p.Category).Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create product failed: %w", err)
	}
	return nil
}

func (r *productRepository) GetByManufacturer(ctx context.Context, opts domain.ProductListOptions) ([]*domain.Product, int, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.GetByManufacturer")
	defer span.End()

	var query strings.Builder
	var args []interface{}
	argPos := 1
	query.WriteString(`SELECT id, manufacturer_id, name, price, description, category, version, created_at, updated_at, COUNT(*) OVER() as total_count FROM products WHERE manufacturer_id = $1`)
	args = append(args, *opts.Filter.ManufacturerID)
	argPos++

	if opts.Filter.Search != nil {
		query.WriteString(fmt.Sprintf(" AND name ILIKE $%d", argPos))
		args = append(args, "%"+*opts.Filter.Search+"%")
		argPos++
	}
	
	query.WriteString(" AND is_deleted = false ORDER BY name ASC LIMIT $")
	query.WriteString(fmt.Sprintf("%d", argPos))
	query.WriteString(" OFFSET $")
	query.WriteString(fmt.Sprintf("%d", argPos+1))
	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	var totalCount int
	for rows.Next() {
		p := &domain.Product{}
		var total int
		err := rows.Scan(&p.ID, &p.ManufacturerID, &p.Name, &p.Price, &p.Description, &p.Category, &p.Version, &p.CreatedAt, &p.UpdatedAt, &total)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, p)
		totalCount = total
	}

	return products, totalCount, nil
}
func (r *productRepository) Update(ctx context.Context, p *domain.Product) error {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.Update")
	defer span.End()

	query := `
		UPDATE products 
		SET name = $1, price = $2, description = $3, category = $4, 
		    version = version + 1, updated_at = NOW()
		WHERE id = $5 AND version = $6 AND is_deleted = false
		RETURNING version, updated_at`

	var newVersion int
	var updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query,
		p.Name, p.Price, p.Description, p.Category, p.ID, p.Version,
	).Scan(&newVersion, &updatedAt)

	if err == sql.ErrNoRows {
		return domain.ErrOptimisticLock
	}
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	p.Version = newVersion
	p.UpdatedAt = updatedAt
	return nil
}

func (r *productRepository) SemanticSearch(ctx context.Context, vector []float32, limit int) ([]*domain.Product, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.SemanticSearch")
	defer span.End()

	query := `SELECT id, manufacturer_id, name, price, description, category, version, created_at, updated_at 
	FROM products WHERE is_deleted = false ORDER BY embedding <=> $1 LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, vector, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		err := rows.Scan(&p.ID, &p.ManufacturerID, &p.Name, &p.Price, &p.Description, &p.Category, &p.Version, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	ctx, span := r.tracer.Start(ctx, "ProductRepository.GetByID")
	defer span.End()

	query := `SELECT id, manufacturer_id, name, price, description, category, version, created_at, updated_at FROM products WHERE id = $1 AND is_deleted = false`
	p := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.ManufacturerID, &p.Name, &p.Price, &p.Description, &p.Category, &p.Version, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}
