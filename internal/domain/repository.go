package domain

import (
	"context"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Archive(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

type InventoryRepository interface {
	UpdateStock(ctx context.Context, productID uuid.UUID, quantity int) error
	ReserveStock(ctx context.Context, productID uuid.UUID, quantity int) error
	ReleaseStock(ctx context.Context, productID uuid.UUID, quantity int) error
	GetStock(ctx context.Context, productID uuid.UUID) (int, error)
}

type ProductSearchRepository interface {
	Search(ctx context.Context, query string, filter ProductFilter, params PaginationParams) ([]*Product, int64, error)
	FindSimilar(ctx context.Context, productID uuid.UUID, limit int) ([]*Product, error)
}

type ProductAnalyticsRepository interface {
	IncrementViews(ctx context.Context, productID uuid.UUID) error
	IncrementSales(ctx context.Context, productID uuid.UUID, count int64) error
	UpdateStats(ctx context.Context, productID uuid.UUID, stats ProductStats) error
}

type ProductRecommendationRepository interface {
	GetRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]*Product, error)
	UpdateAIScore(ctx context.Context, productID uuid.UUID, score float64) error
}

