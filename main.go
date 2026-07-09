package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	scraper := worker.NewScraper()
	processor := worker.NewProductProcessor(logger, nil, nil, nil, nil) 

	fmt.Println("🚀 Робот начал поиск...")

	product, err := scraper.Fetch("https://factory-example.com/cap")
	if err != nil {
		logger.Error("Ошибка при поиске", "err", err)
		return
	}
	fmt.Printf("✅ Робот нашел: %s\n", product.Name)

	processed, err := processor.Process(context.Background(), product, worker.ProcessOptions{})
	if err != nil {
		logger.Error("Ошибка обработки", "err", err)
		return
	}

	fmt.Printf("✨ Обработка завершена! Новый товар: %s\n", processed.Name)
}

