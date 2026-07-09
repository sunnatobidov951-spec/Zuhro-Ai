package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sunnatobidov951-spec/Zuhro-Ai-/internal/domain"
)

// ProductProcessor — главный сервис обработки товаров
type ProductProcessor struct {
	logger        *slog.Logger
	validator     *ProductValidator
	aiEngine      *AIEngine
	seoEngine     *SEOEngine
	marketAdapter MarketplaceAdapter
	cache         Cache
}

func NewProductProcessor(
	logger *slog.Logger,
	aiEngine *AIEngine,
	seoEngine *SEOEngine,
	marketAdapter MarketplaceAdapter,
	cache Cache,
) *ProductProcessor {
	return &ProductProcessor{
		logger:        logger,
		validator:     NewProductValidator(),
		aiEngine:      aiEngine,
		seoEngine:     seoEngine,
		marketAdapter: marketAdapter,
		cache:         cache,
	}
}

// Process — возвращает НОВЫЙ объект, оригинал не меняется
func (p *ProductProcessor) Process(ctx context.Context, original *domain.Product, opts ProcessOptions) (*domain.Product, error) {
	start := time.Now()

	prod := original.Clone()

	if err := p.validator.Validate(prod); err != nil {
		return nil, err
	}

	if cached, ok := p.cache.GetProcessed(prod.ID, opts); ok {
		p.logger.Info("cache hit", "product_id", prod.ID)
		return cached, nil
	}

	processed, err := p.processInternal(ctx, prod, opts)
	if err != nil {
		return nil, err
	}

	p.cache.SetProcessed(prod.ID, opts, processed)

	duration := time.Since(start)
	p.logger.Info("product processed successfully",
		"product_id", processed.ID,
		"marketplace", opts.Channel,
		"duration_ms", duration.Milliseconds(),
		"version", processed.Version,
	)

	return processed, nil
}

func (p *ProductProcessor) processInternal(ctx context.Context, prod *domain.Product, opts ProcessOptions) (*domain.Product, error) {
	if p.aiEngine != nil {
		aiResult, err := p.aiEngine.GenerateDescription(ctx, prod, opts)
		if err != nil {
			p.logger.Warn("ai engine failed, using fallback", "error", err)
		} else {
			prod.Description = aiResult.Description
		}
	}

	prod = p.seoEngine.Enrich(prod, opts)
	prod = p.marketAdapter.Adapt(prod, opts.Channel)

	return prod, nil
}

