package engine

import (
	"errors"
	"fmt"
	"strings"
)

type SearchContext struct {
	Query      string
	UserIntent string
}

type SearchResult struct {
	Query   string
	Intent  string
	Message string
}

var ErrEmptyQuery = errors.New("query cannot be empty")

func PerformFastSearch(ctx SearchContext) (*SearchResult, error) {
	query := strings.TrimSpace(ctx.Query)
	if query == "" {
		return nil, ErrEmptyQuery
	}

	intent := normalizeIntent(ctx.UserIntent)

	var message string
	switch intent {
	case "buy":
		message = fmt.Sprintf("Поиск товаров: %s", query)
	case "info":
		message = fmt.Sprintf("Поиск информации: %s", query)
	case "compare":
		message = fmt.Sprintf("Сравнение вариантов по запросу: %s", query)
	default:
		message = fmt.Sprintf("Общий анализ запроса: %s", query)
	}

	return &SearchResult{
		Query:   query,
		Intent:  intent,
		Message: message,
	}, nil
}

func normalizeIntent(intent string) string {
	intent = strings.ToLower(strings.TrimSpace(intent))
	if intent == "" {
		return "general"
	}
	return intent
}
