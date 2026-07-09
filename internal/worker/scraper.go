package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/domain"
)

type ProductSource interface {
	Fetch(ctx context.Context) ([]domain.Product, error)
}

type Scraper struct {
	repo   domain.ProductRepository
	source ProductSource
	logger *slog.Logger
}

func NewScraper(
	repo domain.ProductRepository,
	source ProductSource,
	logger *slog.Logger,
) *Scraper {
	return &Scraper{
		repo:   repo,
		source: source,
		logger: logger,
	}
}

func (s *Scraper) Start(ctx context.Context) error {
	s.logger.Info("scraper started")

	products, err := s.source.Fetch(ctx)
	if err != nil {
		return err
	}

	for _, p := range products {
		if p.ID == uuid.Nil {
			p.ID = uuid.New()
		}

		if p.Name == "" {
			s.logger.Warn("skip product: empty name")
			continue
		}

		if err := s.repo.Create(ctx, p); err != nil {
			s.logger.Error(
				"create product failed",
				"product", p.Name,
				"error", err,
			)

			if errors.Is(err, context.Canceled) {
				return err
			}
			continue
		}
	}

	s.logger.Info("scraper finished")
	return nil
}

func (s *Scraper) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.logger.Error("scraper cycle failed", "error", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

