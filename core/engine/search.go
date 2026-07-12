package engine

import (
	"errors"
	"fmt"
	"strings"
	"github.com/sunnatobidov951-spec/Zuhro-Ai/core/repository"
)

type SearchContext struct {
	Query      string
	UserIntent string
}

type SearchResult struct {
	Query   string
	Intent  string
	Message string
	Product *repository.Product
}

var ErrEmptyQuery = errors.New("query cannot be empty")

func PerformFastSearch(ctx SearchContext) (*SearchResult, error) {
	query := strings.TrimSpace(ctx.Query)
	if query == "" {
		return nil, ErrEmptyQuery
	}

	intent := normalizeIntent(ctx.UserIntent)

	product, err := repository.GetProductByID(1)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить товар: %w", err)
	}

	return &SearchResult{
		Query:   query,
		Intent:  intent,
		Message: fmt.Sprintf("Результат для: %s", query),
		Product: product,
	}, nil
}

func normalizeIntent(intent string) string {
	intent = strings.ToLower(strings.TrimSpace(intent))
	if intent == "" {
		return "general"
	}
	return intent
}
