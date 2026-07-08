package domain

import (
	"time"

	"github.com/google/uuid"
)

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

type PaginationParams struct {
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	Offset    int               `json:"offset,omitempty"`
	SortBy    string            `json:"sort_by,omitempty"`
	SortOrder SortDirection     `json:"sort_order,omitempty"`
	Search    string            `json:"search,omitempty"`
	Filters   map[string]string `json:"filters,omitempty"`
}

type ProductFilter struct {
	IDs          []uuid.UUID    `json:"ids,omitempty"`
	SellerID     *uuid.UUID     `json:"seller_id,omitempty"`
	CategoryID   *uuid.UUID     `json:"category_id,omitempty"`

	Brand        string         `json:"brand,omitempty"`
	Status       ProductStatus  `json:"status,omitempty"`

	MinPrice     int64          `json:"min_price,omitempty"`
	MaxPrice     int64          `json:"max_price,omitempty"`
	Currency     string         `json:"currency,omitempty"`

	MinRating    float64        `json:"min_rating,omitempty"`

	InStock      *bool          `json:"in_stock,omitempty"`
	IsFeatured   *bool          `json:"is_featured,omitempty"`

	Tags         []string       `json:"tags,omitempty"`

	Source       string         `json:"source,omitempty"`

	CreatedAfter  *time.Time    `json:"created_after,omitempty"`
	CreatedBefore *time.Time    `json:"created_before,omitempty"`
}

type PriceUpdate struct {
	ProductID uuid.UUID `json:"product_id"`

	OldPrice int64 `json:"old_price"`
	NewPrice int64 `json:"new_price"`

	Currency string `json:"currency"`

	Reason string `json:"reason,omitempty"`

	UpdatedBy uuid.UUID `json:"updated_by"`

	UpdatedAt time.Time `json:"updated_at"`
}

type ProductStats struct {
	ProductID uuid.UUID `json:"product_id"`

	ViewsCount     int64 `json:"views_count"`
	SalesCount     int64 `json:"sales_count"`

	FavoritesCount int64 `json:"favorites_count"`
	ReviewsCount   int64 `json:"reviews_count"`

	Revenue        int64 `json:"revenue"`

	AverageRating  float64 `json:"average_rating"`

	ReturnsCount   int64 `json:"returns_count"`

	LastViewedAt   *time.Time `json:"last_viewed_at,omitempty"`
	LastSoldAt     *time.Time `json:"last_sold_at,omitempty"`

	UpdatedAt      time.Time `json:"updated_at"`
}
