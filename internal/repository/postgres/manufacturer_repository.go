package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/domain"
)

type ManufacturerRepository struct {
	db *sql.DB
}

func NewManufacturerRepository(db *sql.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

func (r *ManufacturerRepository) Create(ctx context.Context, m *domain.Manufacturer) error {
	query := `
		INSERT INTO manufacturers (
			id, name, legal_name, country, city, industries, products, 
			factory_size, workers_count, production_capacity, 
			minimum_order_quantity, average_production_days, currency, 
			status, ai_profile, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.Name, m.LegalName, m.Country, m.City, m.Industries, m.Products,
		m.FactorySize, m.WorkersCount, m.ProductionCapacity, m.MinimumOrderQuantity,
		m.AverageProductionDays, m.Currency, m.Status, m.AIProfile,
	)
	if err != nil {
		return fmt.Errorf("manufacturer create failed: %w", err)
	}
	return nil
}

func (r *ManufacturerRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Manufacturer, error) {
	query := `
		SELECT 
			id, name, legal_name, country, city, industries, products, 
			factory_size, workers_count, production_capacity, 
			minimum_order_quantity, average_production_days, currency, 
			status, ai_profile, created_at, updated_at
		FROM manufacturers
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*domain.Manufacturer, 0)

	for rows.Next() {
		m := &domain.Manufacturer{}
		err := rows.Scan(
			&m.ID, &m.Name, &m.LegalName, &m.Country, &m.City, &m.Industries, &m.Products,
			&m.FactorySize, &m.WorkersCount, &m.ProductionCapacity, &m.MinimumOrderQuantity,
			&m.AverageProductionDays, &m.Currency, &m.Status, &m.AIProfile, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

