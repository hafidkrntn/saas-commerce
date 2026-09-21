package entities

import (
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Shared Request Params
// =============================================================================

// ListParams is the standard query parameters for list/pagination endpoints.
// All filter params from frontend go here — shared across all modules.
type ListParams struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	Search     string `form:"search"`
	Status     string `form:"status"`
	CategoryID string `form:"category_id"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
}

// =============================================================================
// Shared Response DTOs
// =============================================================================

// PostResponse is returned after a successful POST (create).
type PostResponse struct {
	Id            uuid.UUID  `json:"id"`
	NoTransaction string     `json:"no_transaction"`
	Version       int        `json:"version"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// PutResponse is returned after a successful PUT (update).
type PutResponse struct {
	Id            uuid.UUID  `json:"id"`
	NoTransaction string     `json:"no_transaction"`
	Version       int        `json:"version"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
