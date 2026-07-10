package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/repository/postgres"
	"github.com/sunnatobidov951-spec/Zuhro-Ai/internal/worker"

	_ "github.com/lib/pq" // не забудь импортировать драйвер postgres
)

func main() {
	// 1. Создаем логгер
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. Подключаемся к БД (замени строку на свою реальную, если нужно)
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/zuhro?sslmode=disable")
	if err != nil {
		logger.Error("ошибка подключения к БД", "error", err)
		return
	}
	defer db.Close()

	// 3. Инициализируем репозиторий
	repo := postgres.NewProductRepository(db)

	// 4. Инициализируем скрейпер (здесь мы передаем реальный репо)
	// ВАЖНО: тебе нужно передать и реальный source вместо nil
	scraper := worker.NewScraper(repo, nil, logger)

	logger.Info("Робот запущен...")

	// Дальше твой код...
	_ = context.Background()
	_ = fmt.Sprintf
	_ = scraper
}
