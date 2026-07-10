package repository

import (
	"context"

	"github.com/google/uuid"
	"zuhroai/internal/domain"
)

type ProductRepository interface {
	// ===============================
	// CORE PRODUCT
	// ===============================
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Archive(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error

	// ===============================
	// MARKETPLACE SEARCH
	// ===============================
	List(ctx context.Context, params domain.PaginationParams) ([]*domain.Product, int64, error)
	Search(ctx context.Context, query string, params domain.PaginationParams) ([]*domain.Product, int64, error)
	AdvancedSearch(ctx context.Context, filter domain.ProductFilter, params domain.PaginationParams) ([]*domain.Product, int64, error)

	// ===============================
	// AI ENGINE
	// ===============================
	FindSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]*domain.Product, error)
	SearchByEmbedding(ctx context.Context, vector []float32, limit int) ([]*domain.Product, error)
	UpdateAIData(ctx context.Context, productID uuid.UUID, score float64, keywords []string) error

	// ===============================
	// SELLER SYSTEM
	// ===============================
	GetBySellerID(ctx context.Context, sellerID uuid.UUID, params domain.PaginationParams) ([]*domain.Product, int64, error)

	// ===============================
	// CATEGORY
	// ===============================
	GetByCategory(ctx context.Context, categoryID uuid.UUID, params domain.PaginationParams) ([]*domain.Product, int64, error)

	// ===============================
	// INVENTORY
	// ===============================
	UpdateStock(ctx context.Context, productID uuid.UUID, quantity int) error
	CheckAvailability(ctx context.Context, productID uuid.UUID) (bool, error)

	// ===============================
	// IMPORT / PARSER
	// ===============================
	GetByExternalID(ctx context.Context, source string, externalID string) (*domain.Product, error)
	ImportProduct(ctx context.Context, product *domain.Product, source string) error

	// ===============================
	// ANALYTICS
	// ===============================
	IncrementViews(ctx context.Context, productID uuid.UUID) error
	IncrementSales(ctx context.Context, productID uuid.UUID, count int64) error
	GetStats(ctx context.Context, productID uuid.UUID) (*domain.ProductStats, error)

	// ===============================
	// BULK OPERATIONS & TRANSACTIONS
	// ===============================
	BulkUpdate(ctx context.Context, products []*domain.Product) error
	BulkArchive(ctx context.Context, ids []uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(ProductRepository) error) error
	LockByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	UpdateWithVersion(ctx context.Context, product *domain.Product, version int64) error
	BulkImport(ctx context.Context, products []*domain.Product) error
	BulkUpdatePrices(ctx context.Context, updates []domain.PriceUpdate) error
	ExistsMany(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]bool, error)
}

