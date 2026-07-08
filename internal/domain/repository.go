package domain

import (
"context"

"github.com/google/uuid"

)

type ProductRepository interface {
// CRUD
Create(ctx context.Context, product *Product) error
GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
Update(ctx context.Context, product *Product) error
Delete(ctx context.Context, id uuid.UUID) error

// Search  
Find(  
	ctx context.Context,  
	filter ProductFilter,  
	params PaginationParams,  
) ([]*Product, int64, error)  

// Seller  
GetBySellerID(  
	ctx context.Context,  
	sellerID uuid.UUID,  
	params PaginationParams,  
) ([]*Product, int64, error)  

// Category  
GetByCategoryID(  
	ctx context.Context,  
	categoryID uuid.UUID,  
	params PaginationParams,  
) ([]*Product, int64, error)  

// Inventory  
UpdateStock(  
	ctx context.Context,  
	productID uuid.UUID,  
	quantity int,  
) error  

// Price  
UpdatePrice(  
	ctx context.Context,  
	update PriceUpdate,  
) error  

// Analytics  
UpdateStats(  
	ctx context.Context,  
	productID uuid.UUID,  
	stats ProductStats,  
) error  

IncrementViews(ctx context.Context, productID uuid.UUID) error  
IncrementSales(ctx context.Context, productID uuid.UUID, count int64) error  

// AI  
FindSimilar(  
	ctx context.Context,  
	productID uuid.UUID,  
	limit int,  
) ([]*Product, error)  

// Bulk  
BulkCreate(ctx context.Context, products []*Product) error  
BulkUpdate(ctx context.Context, products []*Product) error  

// Utility  
Exists(ctx context.Context, id uuid.UUID) (bool, error)  
Count(ctx context.Context, filter ProductFilter) (int64, error)

}
