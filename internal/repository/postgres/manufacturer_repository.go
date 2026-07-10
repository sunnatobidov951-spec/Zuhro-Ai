package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/domain"
)

type ManufacturerRepository struct {
	db *sql.DB
}

func NewManufacturerRepository(db *sql.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

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

func (r *ManufacturerRepository) GetAll(ctx context.Context) ([]domain.Manufacturer, error) {
	const query = `SELECT id, name, created_at FROM manufacturers`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get manufacturers: %w", err)
	}
	defer rows.Close()

	var manufacturers []domain.Manufacturer
	for rows.Next() {
		var m domain.Manufacturer
		if err := rows.Scan(&m.ID, &m.Name, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository: failed to scan manufacturer: %w", err)
		}
		manufacturers = append(manufacturers, m)
	}

	return manufacturers, nil
}
