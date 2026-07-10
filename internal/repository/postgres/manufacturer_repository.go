package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/domain"
)

// ManufacturerRepository handles database operations for manufacturers.
type ManufacturerRepository struct {
	db *sql.DB
}

// NewManufacturerRepository creates a new ManufacturerRepository instance.
func NewManufacturerRepository(db *sql.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

// Create inserts a new manufacturer into the database and returns its generated ID.
func (r *ManufacturerRepository) Create(ctx context.Context, m *domain.Manufacturer) (int64, error) {
	const query = `
		INSERT INTO manufacturers (name, created_at)
		VALUES ($1, NOW())
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(ctx, query, m.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to create manufacturer: %w", err)
	}

	return id, nil
}

// GetAll retrieves all manufacturers from the database.
func (r *ManufacturerRepository) GetAll(ctx context.Context) ([]domain.Manufacturer, error) {
	const query = `
		SELECT id, name, created_at
		FROM manufacturers
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query manufacturers: %w", err)
	}
	defer rows.Close()

	var manufacturers []domain.Manufacturer

	for rows.Next() {
		var m domain.Manufacturer
		if err := rows.Scan(&m.ID, &m.Name, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository: failed to scan manufacturer row: %w", err)
		}
		manufacturers = append(manufacturers, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows iteration error: %w", err)
	}

	return manufacturers, nil
}
