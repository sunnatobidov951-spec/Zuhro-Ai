package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zuhroai/backend/internal/domain"
)

type ProductRepository interface {

	// =========================
	// CORE PRODUCT
	// =========================

	Create(
		ctx context.Context,
		product *domain.Product,
	) error


	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Product, error)


	Update(
		ctx context.Context,
		product *domain.Product,
	) error


	// Мягкое удаление
	Archive(
		ctx context.Context,
		id uuid.UUID,
	) error



	// =========================
	// MARKETPLACE SEARCH
	// =========================

	List(
		ctx context.Context,
		params PaginationParams,
	) ([]*domain.Product, int64, error)


	Search(
		ctx context.Context,
		query string,
		params PaginationParams,
	) ([]*domain.Product, int64, error)


	AdvancedSearch(
		ctx context.Context,
		filter domain.ProductFilter,
		params PaginationParams,
	) ([]*domain.Product, int64, error)



	// =========================
	// AI ENGINE
	// =========================

	FindSimilarProducts(
		ctx context.Context,
		productID uuid.UUID,
		limit int,
	) ([]*domain.Product, error)


	SearchByEmbedding(
		ctx context.Context,
		vector []float32,
		limit int,
	) ([]*domain.Product, error)


	UpdateAIData(
		ctx context.Context,
		productID uuid.UUID,
		score float64,
		keywords []string,
	) error



	// =========================
	// SELLER SYSTEM
	// =========================

	GetBySellerID(
		ctx context.Context,
		sellerID uuid.UUID,
		params PaginationParams,
	) ([]*domain.Product, int64, error)



	// =========================
	// CATEGORY
	// =========================

	GetByCategory(
		ctx context.Context,
		categoryID uuid.UUID,
		params PaginationParams,
	) ([]*domain.Product, int64, error)



	// =========================
	// INVENTORY
	// =========================

	UpdateStock(
		ctx context.Context,
		productID uuid.UUID,
		quantity int,
	) error


	CheckAvailability(
		ctx context.Context,
		productID uuid.UUID,
	) (bool, error)



	// =========================
	// IMPORT / PARSER
	// =========================

	GetByExternalID(
		ctx context.Context,
		source string,
		externalID string,
	) (*domain.Product, error)


	ImportProduct(
		ctx context.Context,
		product *domain.Product,
		source string,
	) error



	// =========================
	// ANALYTICS
	// =========================

	IncrementViews(
		ctx context.Context,
		productID uuid.UUID,
	) error


	IncrementSales(
		ctx context.Context,
		productID uuid.UUID,
		count int,
	) error


	GetStats(
		ctx context.Context,
		productID uuid.UUID,
	) (*domain.ProductStats, error)



	// =========================
	// BULK OPERATIONS
	// =========================

	BulkUpdate(
		ctx context.Context,
		products []*domain.Product,
	) error


	BulkArchive(
		ctx context.Context,
		ids []uuid.UUID,
	) error
}
