package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product — основная сущность маркетплейса (AI-ready)
type Product struct {
	ID              uuid.UUID      `json:"id"`
	SellerID        uuid.UUID      `json:"seller_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	ShortDescription string        `json:"short_description"`

	// Цены
	Price           int64          `json:"price"`           // в центах
	OriginalPrice   int64          `json:"original_price"`  // в центах
	Currency        string         `json:"currency"`
	DiscountPercent float64        `json:"discount_percent"`

	// Категории и атрибуты
	Category        string         `json:"category"`
	SubCategory     string         `json:"sub_category"`
	Brand           string         `json:"brand"`
	Tags            []string       `json:"tags"`
	Attributes      map[string]interface{} `json:"attributes"` // цвет, размер, материал и т.д.

	// Медиа
	Images          []ProductImage `json:"images"`
	VideoURL        string         `json:"video_url,omitempty"`

	// Статус и наличие
	Status          ProductStatus  `json:"status"`
	StockQuantity   int            `json:"stock_quantity"`
	IsInStock       bool           `json:"is_in_stock"`
	IsFeatured      bool           `json:"is_featured"`

	// AI данные
	AIScore         float64        `json:"ai_score"`         // качество карточки по AI
	AIKeywords      []string       `json:"ai_keywords"`
	AIRecommendations []string     `json:"ai_recommendations"`

	// Метрики
	Views           int64          `json:"views"`
	SalesCount      int64          `json:"sales_count"`
	Rating          float64        `json:"rating"`
	ReviewCount     int            `json:"review_count"`

	// Timestamps
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	PublishedAt     *time.Time     `json:"published_at,omitempty"`
}

// ProductImage
type ProductImage struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Order     int       `json:"order"`
	IsMain    bool      `json:"is_main"`
}

// ProductStatus
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPending   ProductStatus = "pending"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusPaused    ProductStatus = "paused"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusRejected  ProductStatus = "rejected"
)

// Методы сущности
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.StockQuantity > 0
}

func (p *Product) CalculateDiscount() float64 {
	if p.OriginalPrice > 0 && p.Price < p.OriginalPrice {
		return float64(p.OriginalPrice-p.Price) / float64(p.OriginalPrice) * 100
	}
	return 0
}

func (p *Product) UpdateStock(quantity int) {
	p.StockQuantity = quantity
	p.IsInStock = quantity > 0
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) Publish() {
	now := time.Now().UTC()
	p.Status = ProductStatusActive
	p.PublishedAt = &now
	p.UpdatedAt = now
}
