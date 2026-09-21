package setting

import (
	"backend-go/internal/entities"
	"time"

	"github.com/google/uuid"
)

// StoreSetting represents the `store_settings` table (1 row per tenant).
type StoreSetting struct {
	TenantID            uuid.UUID `gorm:"column:tenant_id;type:uuid;primaryKey"`
	Name                string    `gorm:"column:name"`
	Description         string    `gorm:"column:description"`
	Currency            string    `gorm:"column:currency;not null;default:'IDR'"`
	Timezone            string    `gorm:"column:timezone;not null;default:'Asia/Jakarta'"`
	OrderPrefix         string    `gorm:"column:order_prefix;not null;default:'ORD'"`
	LowStockThreshold   int       `gorm:"column:low_stock_threshold;not null;default:5"`
	EnableReviews       bool      `gorm:"column:enable_reviews;not null;default:true"`
	EnableWishlist      bool      `gorm:"column:enable_wishlist;not null;default:true"`
	EnableGiftCards     bool      `gorm:"column:enable_gift_cards;not null;default:false"`
	DefaultWeightUnit   string    `gorm:"column:default_weight_unit;not null;default:'kg'"`
	DefaultDimensionUnit string   `gorm:"column:default_dimension_unit;not null;default:'cm'"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (StoreSetting) TableName() string {
	return "store_settings"
}

// CompanySetting represents the `company_settings` table (1 row per tenant).
type CompanySetting struct {
	TenantID   uuid.UUID       `gorm:"column:tenant_id;type:uuid;primaryKey"`
	Name       string          `gorm:"column:name"`
	LegalName  string          `gorm:"column:legal_name"`
	Logo       string          `gorm:"column:logo"`
	Email      string          `gorm:"column:email"`
	Phone      string          `gorm:"column:phone"`
	Website    string          `gorm:"column:website"`
	Address    entities.JSONB  `gorm:"column:address;type:jsonb"`
	TaxID      string          `gorm:"column:tax_id"`
	Currency   string          `gorm:"column:currency;not null;default:'IDR'"`
	Timezone   string          `gorm:"column:timezone;not null;default:'Asia/Jakarta'"`
	DateFormat string          `gorm:"column:date_format;not null;default:'d MMM yyyy'"`
	CreatedAt  time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;autoUpdateTime"`
}

func (CompanySetting) TableName() string {
	return "company_settings"
}

// TaxRate represents the `tax_rates` table.
type TaxRate struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  uuid.UUID `gorm:"column:tenant_id;type:uuid;not null;index"`
	Name      string    `gorm:"column:name;not null"`
	Rate      float64   `gorm:"column:rate;not null;default:0"`
	Region    string    `gorm:"column:region;not null;default:'all'"`
	Type      string    `gorm:"column:type;not null;default:'vat'"`
	AppliesTo string    `gorm:"column:applies_to;not null;default:'all'"`
	IsActive  bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (TaxRate) TableName() string {
	return "tax_rates"
}

// ShippingZone represents the `shipping_zones` table.
type ShippingZone struct {
	ID            uuid.UUID         `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID      uuid.UUID         `gorm:"column:tenant_id;type:uuid;not null;index"`
	Name          string            `gorm:"column:name;not null"`
	Regions       entities.StringList `gorm:"column:regions;type:jsonb"`
	Method        string            `gorm:"column:method;not null;default:'flat_rate'"`
	Rate          float64           `gorm:"column:rate;not null;default:0"`
	FreeAbove     *float64          `gorm:"column:free_above"`
	EstimatedDays string            `gorm:"column:estimated_days"`
	IsActive      bool              `gorm:"column:is_active;not null;default:true"`
	CreatedAt     time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time         `gorm:"column:updated_at;autoUpdateTime"`
}

func (ShippingZone) TableName() string {
	return "shipping_zones"
}

// NotificationSetting represents the `notification_settings` table.
type NotificationSetting struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    uuid.UUID `gorm:"column:tenant_id;type:uuid;not null;index"`
	Event       string    `gorm:"column:event;not null"`
	Email       bool      `gorm:"column:email;not null;default:false"`
	SMS         bool      `gorm:"column:sms;not null;default:false"`
	InApp       bool      `gorm:"column:in_app;not null;default:true"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (NotificationSetting) TableName() string {
	return "notification_settings"
}

// =============================================================================
// Request DTOs
// =============================================================================

type UpdateStoreSettingRequest struct {
	Name                *string `json:"name"`
	Description         *string `json:"description"`
	Currency            *string `json:"currency"`
	Timezone            *string `json:"timezone"`
	OrderPrefix         *string `json:"order_prefix"`
	LowStockThreshold   *int    `json:"low_stock_threshold"`
	EnableReviews       *bool   `json:"enable_reviews"`
	EnableWishlist      *bool   `json:"enable_wishlist"`
	EnableGiftCards     *bool   `json:"enable_gift_cards"`
	DefaultWeightUnit   *string `json:"default_weight_unit"`
	DefaultDimensionUnit *string `json:"default_dimension_unit"`
}

type UpdateCompanySettingRequest struct {
	Name       *string        `json:"name"`
	LegalName  *string        `json:"legal_name"`
	Logo       *string        `json:"logo"`
	Email      *string        `json:"email"`
	Phone      *string        `json:"phone"`
	Website    *string        `json:"website"`
	Address    *entities.JSONB `json:"address"`
	TaxID      *string        `json:"tax_id"`
	Currency   *string        `json:"currency"`
	Timezone   *string        `json:"timezone"`
	DateFormat *string        `json:"date_format"`
}

type CreateTaxRateRequest struct {
	Name      string  `json:"name" binding:"required,min=2,max=100"`
	Rate      float64 `json:"rate" binding:"gte=0,lte=100"`
	Region    string  `json:"region"`
	Type      string  `json:"type" binding:"omitempty,oneof=vat sales_tax gst"`
	AppliesTo string  `json:"applies_to" binding:"omitempty,oneof=all digital physical services"`
	IsActive  *bool   `json:"is_active"`
}

type UpdateTaxRateRequest struct {
	Name      *string  `json:"name" binding:"omitempty,min=2,max=100"`
	Rate      *float64 `json:"rate" binding:"omitempty,gte=0,lte=100"`
	Region    *string  `json:"region"`
	Type      *string  `json:"type" binding:"omitempty,oneof=vat sales_tax gst"`
	AppliesTo *string  `json:"applies_to" binding:"omitempty,oneof=all digital physical services"`
	IsActive  *bool    `json:"is_active"`
}

type CreateShippingZoneRequest struct {
	Name          string    `json:"name" binding:"required,min=2,max=100"`
	Regions       []string  `json:"regions"`
	Method        string    `json:"method" binding:"omitempty,oneof=flat_rate free calculated"`
	Rate          float64   `json:"rate"`
	FreeAbove     *float64  `json:"free_above"`
	EstimatedDays string    `json:"estimated_days"`
	IsActive      *bool     `json:"is_active"`
}

type UpdateShippingZoneRequest struct {
	Name          *string   `json:"name" binding:"omitempty,min=2,max=100"`
	Regions       []string  `json:"regions"`
	Method        *string   `json:"method" binding:"omitempty,oneof=flat_rate free calculated"`
	Rate          *float64  `json:"rate"`
	FreeAbove     *float64  `json:"free_above"`
	EstimatedDays *string   `json:"estimated_days"`
	IsActive      *bool     `json:"is_active"`
}

type UpdateNotificationSettingRequest struct {
	Email       *bool   `json:"email"`
	SMS         *bool   `json:"sms"`
	InApp       *bool   `json:"in_app"`
	Description *string `json:"description"`
}

// =============================================================================
// Response DTOs
// =============================================================================

type StoreSettingResponse struct {
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Currency            string    `json:"currency"`
	Timezone            string    `json:"timezone"`
	OrderPrefix         string    `json:"order_prefix"`
	LowStockThreshold   int       `json:"low_stock_threshold"`
	EnableReviews       bool      `json:"enable_reviews"`
	EnableWishlist      bool      `json:"enable_wishlist"`
	EnableGiftCards     bool      `json:"enable_gift_cards"`
	DefaultWeightUnit   string    `json:"default_weight_unit"`
	DefaultDimensionUnit string   `json:"default_dimension_unit"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CompanySettingResponse struct {
	Name       string        `json:"name"`
	LegalName  string        `json:"legal_name"`
	Logo       string        `json:"logo"`
	Email      string        `json:"email"`
	Phone      string        `json:"phone"`
	Website    string        `json:"website"`
	Address    entities.JSONB `json:"address"`
	TaxID      string        `json:"tax_id"`
	Currency   string        `json:"currency"`
	Timezone   string        `json:"timezone"`
	DateFormat string        `json:"date_format"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type TaxRateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Rate      float64   `json:"rate"`
	Region    string    `json:"region"`
	Type      string    `json:"type"`
	AppliesTo string    `json:"applies_to"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *TaxRate) ToResponse() TaxRateResponse {
	return TaxRateResponse{
		ID:        t.ID,
		Name:      t.Name,
		Rate:      t.Rate,
		Region:    t.Region,
		Type:      t.Type,
		AppliesTo: t.AppliesTo,
		IsActive:  t.IsActive,
		CreatedAt: t.CreatedAt,
	}
}

type ShippingZoneResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Regions       []string  `json:"regions"`
	Method        string    `json:"method"`
	Rate          float64   `json:"rate"`
	FreeAbove     *float64  `json:"free_above"`
	EstimatedDays string    `json:"estimated_days"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

func (z *ShippingZone) ToResponse() ShippingZoneResponse {
	return ShippingZoneResponse{
		ID:            z.ID,
		Name:          z.Name,
		Regions:       z.Regions,
		Method:        z.Method,
		Rate:          z.Rate,
		FreeAbove:     z.FreeAbove,
		EstimatedDays: z.EstimatedDays,
		IsActive:      z.IsActive,
		CreatedAt:     z.CreatedAt,
	}
}

type NotificationSettingResponse struct {
	ID          uuid.UUID `json:"id"`
	Event       string    `json:"event"`
	Email       bool      `json:"email"`
	SMS         bool      `json:"sms"`
	InApp       bool      `json:"in_app"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (n *NotificationSetting) ToResponse() NotificationSettingResponse {
	return NotificationSettingResponse{
		ID:          n.ID,
		Event:       n.Event,
		Email:       n.Email,
		SMS:         n.SMS,
		InApp:       n.InApp,
		Description: n.Description,
		CreatedAt:   n.CreatedAt,
	}
}
