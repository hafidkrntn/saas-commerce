package analytics

import (
	"context"
	"time"
)

// =============================================================================
// Response DTOs
// =============================================================================

type Summary struct {
	Revenue         float64 `json:"revenue"`
	RevenueChange   float64 `json:"revenue_change"`
	Profit          float64 `json:"profit"`
	ProfitChange    float64 `json:"profit_change"`
	Orders          int64   `json:"orders"`
	OrdersChange    float64 `json:"orders_change"`
	ConversionRate  float64 `json:"conversion_rate"`
	ConversionChange float64 `json:"conversion_change"`
	Customers       int64   `json:"customers"`
	CustomersChange float64 `json:"customers_change"`
	AverageOrderValue float64 `json:"average_order_value"`
	AOVChange       float64 `json:"aov_change"`
}

type RevenuePoint struct {
	Month    string  `json:"month"`
	Revenue  float64 `json:"revenue"`
	Profit   float64 `json:"profit"`
	Orders   int64   `json:"orders"`
	Expenses float64 `json:"expenses"`
}

type SalesPoint struct {
	Date   string  `json:"date"`
	Sales  float64 `json:"sales"`
	Orders int64   `json:"orders"`
}

type TopProduct struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Revenue    float64 `json:"revenue"`
	UnitsSold  int64   `json:"units_sold"`
	Percentage float64 `json:"percentage"`
}

type TopCategory struct {
	Name       string  `json:"name"`
	Revenue    float64 `json:"revenue"`
	Percentage float64 `json:"percentage"`
}

// =============================================================================
// Service Interface
// =============================================================================

type ServiceInterface interface {
	Summary(ctx context.Context, tenantID string, from, to time.Time) (Summary, error)
	Revenue(ctx context.Context, tenantID string, months int) ([]RevenuePoint, error)
	Sales(ctx context.Context, tenantID string, days int) ([]SalesPoint, error)
	TopProducts(ctx context.Context, tenantID string, from, to time.Time, limit int) ([]TopProduct, error)
	TopCategories(ctx context.Context, tenantID string, from, to time.Time, limit int) ([]TopCategory, error)
}
