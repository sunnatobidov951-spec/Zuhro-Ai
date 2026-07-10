packge postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"zuhroai/internal/domain"
)

var allowedSortFields = map[string]string{
	"price":     "p.price",
	"rating":    "p.rating",
	"sales":     "p.sales_count",
	"views":     "p.views",
	"ai_score":  "p.ai_score",
	"newest":    "p.created_at",
}

func (r *ProductRepository) Find(
	ctx context.Context,
	filter domain.ProductFilter,
	params domain.PaginationParams,
) ([]*domain.Product, int64, error) {

	where := []string{
		"p.status = 'active'",
	}

	args := make([]any, 0)
	argIndex := 1


	// Category filter
	if filter.CategoryID != uuid.Nil {

		where = append(
			where,
			fmt.Sprintf(
				"p.category_id = $%d",
				argIndex,
			),
		)

		args = append(
			args,
			filter.CategoryID,
		)

		argIndex++
	}


	// Price range
	if filter.MinPrice > 0 {

		where = append(
			where,
			fmt.Sprintf(
				"p.price >= $%d",
				argIndex,
			),
		)

		args = append(
			args,
			filter.MinPrice,
		)

		argIndex++
	}


	if filter.MaxPrice > 0 {

		where = append(
			where,
			fmt.Sprintf(
				"p.price <= $%d",
				argIndex,
			),
		)

		args = append(
			args,
			filter.MaxPrice,
		)

		argIndex++
	}


	order := "p.created_at DESC"


	// Safe sorting
	if field, ok := allowedSortFields[params.SortBy]; ok {

		direction := "DESC"

		if strings.ToLower(params.SortDir) == "asc" {
			direction = "ASC"
		}

		order = field + " " + direction
	}


	query := fmt.Sprintf(`
	SELECT
		p.id,
		p.seller_id,
		p.title,
		p.description,
		p.short_description,
		p.price,
		p.original_price,
		p.currency,
		p.category,
		p.brand,
		p.status,
		p.stock_quantity,
		p.ai_score,
		p.views,
		p.sales_count,
		p.rating,
		p.review_count,
		p.created_at,
		p.updated_at

	FROM products p

	WHERE %s

	ORDER BY %s

	LIMIT $%d
	OFFSET $%d
	`,
		strings.Join(where, " AND "),
		order,
		argIndex,
		argIndex+1,
	)


	args = append(
		args,
		params.Limit,
		params.Offset,
	)


	rows, err := r.db.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil,0,err
	}

	defer rows.Close()


	products := make(
		[]*domain.Product,
		0,
	)


	for rows.Next() {

		var p domain.Product


		err := rows.Scan(
			&p.ID,
			&p.SellerID,
			&p.Title,
			&p.Description,
			&p.ShortDescription,
			&p.Price,
			&p.OriginalPrice,
			&p.Currency,
			&p.Category,
			&p.Brand,
			&p.Status,
			&p.StockQuantity,
			&p.AIScore,
			&p.Views,
			&p.SalesCount,
			&p.Rating,
			&p.ReviewCount,
			&p.CreatedAt,
			&p.UpdatedAt,
		)


		if err != nil {
			return nil,0,err
		}


		products = append(
			products,
			&p,
		)
	}


	if err := rows.Err(); err != nil {
		return nil,0,err
	}


	total, err := r.Count(
		ctx,
		filter,
	)

	if err != nil {
		return nil,0,err
	}


	return products,total,nil
}
