package domain

import (
	"time"

	"github.com/google/uuid"
)

type Manufacturer struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	LegalName            string    `json:"legal_name"`
	Country              string    `json:"country"`
	City                 string    `json:"city"`
	Industries           []string  `json:"industries"`
	Products             []string  `json:"products"`
	FactorySize          int       `json:"factory_size"`
	WorkersCount         int       `json:"workers_count"`
	ProductionCapacity   int       `json:"production_capacity"`
	MinimumOrderQuantity int       `json:"minimum_order_quantity"`
	AverageProductionDays int      `json:"average_production_days"`
	Currency             string    `json:"currency"`
	Status               string    `json:"status"`
	AIProfile            string    `json:"ai_profile"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
