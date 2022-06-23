package domain

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  *string         `json:"description"`
	DurationDays int             `json:"duration_days"`
	Price        decimal.Decimal `json:"price"`
}
