package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/zuhroai/backend/internal/domain"
)

type ManufacturerMatch struct {
	Manufacturer *domain.Manufacturer
	MatchScore   float64
	Reason       string
}

type ManufacturerRepository struct {
	db *sql.DB
}

func NewManufacturerRepository(db *sql.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

func (r *ManufacturerRepository) Create(ctx context.Context, m *domain.Manufacturer) error {
	query := `INSERT INTO manufacturers (id, name, legal_name, country, city, industries, products, factory_size, workers_count, production_capacity, moq, production_days, currency, status, verification, ai_profile) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	
	_, err := r.db.ExecContext(ctx, query, m.ID, m.Name, m.LegalName, m.Country, m.City, m.Industries, m.Products, m.FactorySize, m.WorkersCount, m.ProductionCapacity, m.MinimumOrderQuantity, m.AverageProductionDays, m.Currency, m.Status, m.Verification, m.AIProfile)
	return err
}

func (r *ManufacturerRepository) FindBestManufacturers(ctx context.Context, category string, maxPrice float64, limit int) ([]ManufacturerMatch, error) {
	query := `
	SELECT id, name, legal_name, country, city, 
	       (ai_profile->>'trust_score')::float, (ai_profile->>'quality_score')::float, (ai_profile->>'risk_score')::float, 
	       moq, production_days, status, created_at, updated_at,
	       (( (ai_profile->>'trust_score')::float * 0.35) + ((ai_profile->>'quality_score')::float * 0.35) + ((100 - (ai_profile->>'risk_score')::float) * 0.20)) AS match_score
	FROM manufacturers
	WHERE status = 'active' AND $1 = ANY(industries)
	ORDER BY match_score DESC
	LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, category, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ManufacturerMatch
	for rows.Next() {
		var m domain.Manufacturer
		var score float64
		err := rows.Scan(&m.ID, &m.Name, &m.LegalName, &m.Country, &m.City, &m.AIProfile.TrustScore, &m.AIProfile.QualityScore, &m.AIProfile.RiskScore, &m.MinimumOrderQuantity, &m.AverageProductionDays, &m.Status, &m.CreatedAt, &m.UpdatedAt, &score)
		if err != nil {
			return nil, err
		}
		results = append(results, ManufacturerMatch{Manufacturer: &m, MatchScore: score, Reason: "Подходит по AI анализу"})
	}
	return results, nil
}
