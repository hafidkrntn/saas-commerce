package setting

import (
	"backend-go/pkg/apperror"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	GetStore(ctx context.Context, tenantID uuid.UUID) (StoreSettingResponse, error)
	UpdateStore(ctx context.Context, tenantID uuid.UUID, req *UpdateStoreSettingRequest) (StoreSettingResponse, error)
	GetCompany(ctx context.Context, tenantID uuid.UUID) (CompanySettingResponse, error)
	UpdateCompany(ctx context.Context, tenantID uuid.UUID, req *UpdateCompanySettingRequest) (CompanySettingResponse, error)

	ListTaxRates(ctx context.Context, tenantID uuid.UUID) ([]TaxRateResponse, error)
	CreateTaxRate(ctx context.Context, tenantID uuid.UUID, req *CreateTaxRateRequest) (TaxRateResponse, error)
	UpdateTaxRate(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateTaxRateRequest) (TaxRateResponse, error)
	DeleteTaxRate(ctx context.Context, tenantID uuid.UUID, id string) error

	ListShippingZones(ctx context.Context, tenantID uuid.UUID) ([]ShippingZoneResponse, error)
	CreateShippingZone(ctx context.Context, tenantID uuid.UUID, req *CreateShippingZoneRequest) (ShippingZoneResponse, error)
	UpdateShippingZone(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateShippingZoneRequest) (ShippingZoneResponse, error)
	DeleteShippingZone(ctx context.Context, tenantID uuid.UUID, id string) error

	ListNotifications(ctx context.Context, tenantID uuid.UUID) ([]NotificationSettingResponse, error)
	UpdateNotification(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateNotificationSettingRequest) (NotificationSettingResponse, error)
}

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

// --- Store ---

func (s *Service) GetStore(ctx context.Context, tenantID uuid.UUID) (StoreSettingResponse, error) {
	setting, err := s.repo.GetStoreSetting(s.db.WithContext(ctx), tenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = StoreSetting{TenantID: tenantID}
			if err := s.repo.UpsertStoreSetting(s.db.WithContext(ctx), &setting); err != nil {
				return StoreSettingResponse{}, apperror.Internal("gagal menyiapkan store settings", err)
			}
			return storeToResponse(setting), nil
		}
		return StoreSettingResponse{}, apperror.Internal("gagal mengambil store settings", err)
	}
	return storeToResponse(setting), nil
}

func (s *Service) UpdateStore(ctx context.Context, tenantID uuid.UUID, req *UpdateStoreSettingRequest) (StoreSettingResponse, error) {
	var updated StoreSetting

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		setting, err := s.repo.GetStoreSetting(tx, tenantID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.Internal("gagal mengambil store settings", err)
			}
			setting = StoreSetting{TenantID: tenantID}
		}

		if req.Name != nil {
			setting.Name = *req.Name
		}
		if req.Description != nil {
			setting.Description = *req.Description
		}
		if req.Currency != nil {
			setting.Currency = *req.Currency
		}
		if req.Timezone != nil {
			setting.Timezone = *req.Timezone
		}
		if req.OrderPrefix != nil {
			setting.OrderPrefix = *req.OrderPrefix
		}
		if req.LowStockThreshold != nil {
			setting.LowStockThreshold = *req.LowStockThreshold
		}
		if req.EnableReviews != nil {
			setting.EnableReviews = *req.EnableReviews
		}
		if req.EnableWishlist != nil {
			setting.EnableWishlist = *req.EnableWishlist
		}
		if req.EnableGiftCards != nil {
			setting.EnableGiftCards = *req.EnableGiftCards
		}
		if req.DefaultWeightUnit != nil {
			setting.DefaultWeightUnit = *req.DefaultWeightUnit
		}
		if req.DefaultDimensionUnit != nil {
			setting.DefaultDimensionUnit = *req.DefaultDimensionUnit
		}

		if err := s.repo.UpsertStoreSetting(tx, &setting); err != nil {
			return apperror.Internal("gagal memperbarui store settings", err)
		}

		updated = setting
		return nil
	})

	if txErr != nil {
		return StoreSettingResponse{}, txErr
	}

	return storeToResponse(updated), nil
}

func storeToResponse(s StoreSetting) StoreSettingResponse {
	return StoreSettingResponse{
		Name:                 s.Name,
		Description:          s.Description,
		Currency:             s.Currency,
		Timezone:             s.Timezone,
		OrderPrefix:          s.OrderPrefix,
		LowStockThreshold:    s.LowStockThreshold,
		EnableReviews:        s.EnableReviews,
		EnableWishlist:       s.EnableWishlist,
		EnableGiftCards:      s.EnableGiftCards,
		DefaultWeightUnit:    s.DefaultWeightUnit,
		DefaultDimensionUnit: s.DefaultDimensionUnit,
		UpdatedAt:            s.UpdatedAt,
	}
}

// --- Company ---

func (s *Service) GetCompany(ctx context.Context, tenantID uuid.UUID) (CompanySettingResponse, error) {
	setting, err := s.repo.GetCompanySetting(s.db.WithContext(ctx), tenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = CompanySetting{TenantID: tenantID}
			if err := s.repo.UpsertCompanySetting(s.db.WithContext(ctx), &setting); err != nil {
				return CompanySettingResponse{}, apperror.Internal("gagal menyiapkan company settings", err)
			}
			return companyToResponse(setting), nil
		}
		return CompanySettingResponse{}, apperror.Internal("gagal mengambil company settings", err)
	}
	return companyToResponse(setting), nil
}

func (s *Service) UpdateCompany(ctx context.Context, tenantID uuid.UUID, req *UpdateCompanySettingRequest) (CompanySettingResponse, error) {
	var updated CompanySetting

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		setting, err := s.repo.GetCompanySetting(tx, tenantID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.Internal("gagal mengambil company settings", err)
			}
			setting = CompanySetting{TenantID: tenantID}
		}

		if req.Name != nil {
			setting.Name = *req.Name
		}
		if req.LegalName != nil {
			setting.LegalName = *req.LegalName
		}
		if req.Logo != nil {
			setting.Logo = *req.Logo
		}
		if req.Email != nil {
			setting.Email = *req.Email
		}
		if req.Phone != nil {
			setting.Phone = *req.Phone
		}
		if req.Website != nil {
			setting.Website = *req.Website
		}
		if req.Address != nil {
			setting.Address = *req.Address
		}
		if req.TaxID != nil {
			setting.TaxID = *req.TaxID
		}
		if req.Currency != nil {
			setting.Currency = *req.Currency
		}
		if req.Timezone != nil {
			setting.Timezone = *req.Timezone
		}
		if req.DateFormat != nil {
			setting.DateFormat = *req.DateFormat
		}

		if err := s.repo.UpsertCompanySetting(tx, &setting); err != nil {
			return apperror.Internal("gagal memperbarui company settings", err)
		}

		updated = setting
		return nil
	})

	if txErr != nil {
		return CompanySettingResponse{}, txErr
	}

	return companyToResponse(updated), nil
}

func companyToResponse(c CompanySetting) CompanySettingResponse {
	return CompanySettingResponse{
		Name:       c.Name,
		LegalName:  c.LegalName,
		Logo:       c.Logo,
		Email:      c.Email,
		Phone:      c.Phone,
		Website:    c.Website,
		Address:    c.Address,
		TaxID:      c.TaxID,
		Currency:   c.Currency,
		Timezone:   c.Timezone,
		DateFormat: c.DateFormat,
		UpdatedAt:  c.UpdatedAt,
	}
}

// --- Tax rates ---

func (s *Service) ListTaxRates(ctx context.Context, tenantID uuid.UUID) ([]TaxRateResponse, error) {
	rates, err := s.repo.ListTaxRates(s.db.WithContext(ctx), tenantID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar pajak", err)
	}

	responses := make([]TaxRateResponse, len(rates))
	for i, t := range rates {
		responses[i] = t.ToResponse()
	}
	return responses, nil
}

func (s *Service) CreateTaxRate(ctx context.Context, tenantID uuid.UUID, req *CreateTaxRateRequest) (TaxRateResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	taxType := req.Type
	if taxType == "" {
		taxType = "vat"
	}
	appliesTo := req.AppliesTo
	if appliesTo == "" {
		appliesTo = "all"
	}
	region := req.Region
	if region == "" {
		region = "all"
	}

	rate := TaxRate{
		TenantID:  tenantID,
		Name:      req.Name,
		Rate:      req.Rate,
		Region:    region,
		Type:      taxType,
		AppliesTo: appliesTo,
		IsActive:  isActive,
	}

	if err := s.repo.CreateTaxRate(s.db.WithContext(ctx), &rate); err != nil {
		return TaxRateResponse{}, apperror.Internal("gagal membuat pajak", err)
	}

	return rate.ToResponse(), nil
}

func (s *Service) UpdateTaxRate(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateTaxRateRequest) (TaxRateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return TaxRateResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	rate, err := s.repo.GetTaxRate(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TaxRateResponse{}, apperror.NotFound("pajak tidak ditemukan", err)
		}
		return TaxRateResponse{}, apperror.Internal("gagal mengambil pajak", err)
	}

	if req.Name != nil {
		rate.Name = *req.Name
	}
	if req.Rate != nil {
		rate.Rate = *req.Rate
	}
	if req.Region != nil {
		rate.Region = *req.Region
	}
	if req.Type != nil {
		rate.Type = *req.Type
	}
	if req.AppliesTo != nil {
		rate.AppliesTo = *req.AppliesTo
	}
	if req.IsActive != nil {
		rate.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateTaxRate(s.db.WithContext(ctx), &rate); err != nil {
		return TaxRateResponse{}, apperror.Internal("gagal memperbarui pajak", err)
	}

	return rate.ToResponse(), nil
}

func (s *Service) DeleteTaxRate(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	if err := s.repo.DeleteTaxRate(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus pajak", err)
	}

	return nil
}

// --- Shipping zones ---

func (s *Service) ListShippingZones(ctx context.Context, tenantID uuid.UUID) ([]ShippingZoneResponse, error) {
	zones, err := s.repo.ListShippingZones(s.db.WithContext(ctx), tenantID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar zona pengiriman", err)
	}

	responses := make([]ShippingZoneResponse, len(zones))
	for i, z := range zones {
		responses[i] = z.ToResponse()
	}
	return responses, nil
}

func (s *Service) CreateShippingZone(ctx context.Context, tenantID uuid.UUID, req *CreateShippingZoneRequest) (ShippingZoneResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	method := req.Method
	if method == "" {
		method = "flat_rate"
	}

	zone := ShippingZone{
		TenantID:      tenantID,
		Name:          req.Name,
		Regions:       req.Regions,
		Method:        method,
		Rate:          req.Rate,
		FreeAbove:     req.FreeAbove,
		EstimatedDays: req.EstimatedDays,
		IsActive:      isActive,
	}

	if err := s.repo.CreateShippingZone(s.db.WithContext(ctx), &zone); err != nil {
		return ShippingZoneResponse{}, apperror.Internal("gagal membuat zona pengiriman", err)
	}

	return zone.ToResponse(), nil
}

func (s *Service) UpdateShippingZone(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateShippingZoneRequest) (ShippingZoneResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ShippingZoneResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	zone, err := s.repo.GetShippingZone(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ShippingZoneResponse{}, apperror.NotFound("zona pengiriman tidak ditemukan", err)
		}
		return ShippingZoneResponse{}, apperror.Internal("gagal mengambil zona pengiriman", err)
	}

	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Regions != nil {
		zone.Regions = req.Regions
	}
	if req.Method != nil {
		zone.Method = *req.Method
	}
	if req.Rate != nil {
		zone.Rate = *req.Rate
	}
	if req.FreeAbove != nil {
		zone.FreeAbove = req.FreeAbove
	}
	if req.EstimatedDays != nil {
		zone.EstimatedDays = *req.EstimatedDays
	}
	if req.IsActive != nil {
		zone.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateShippingZone(s.db.WithContext(ctx), &zone); err != nil {
		return ShippingZoneResponse{}, apperror.Internal("gagal memperbarui zona pengiriman", err)
	}

	return zone.ToResponse(), nil
}

func (s *Service) DeleteShippingZone(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	if err := s.repo.DeleteShippingZone(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus zona pengiriman", err)
	}

	return nil
}

// --- Notifications ---

func (s *Service) ListNotifications(ctx context.Context, tenantID uuid.UUID) ([]NotificationSettingResponse, error) {
	notifications, err := s.repo.ListNotifications(s.db.WithContext(ctx), tenantID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar notifikasi", err)
	}

	responses := make([]NotificationSettingResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = n.ToResponse()
	}
	return responses, nil
}

func (s *Service) UpdateNotification(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateNotificationSettingRequest) (NotificationSettingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return NotificationSettingResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	notification, err := s.repo.GetNotification(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotificationSettingResponse{}, apperror.NotFound("notifikasi tidak ditemukan", err)
		}
		return NotificationSettingResponse{}, apperror.Internal("gagal mengambil notifikasi", err)
	}

	if req.Email != nil {
		notification.Email = *req.Email
	}
	if req.SMS != nil {
		notification.SMS = *req.SMS
	}
	if req.InApp != nil {
		notification.InApp = *req.InApp
	}
	if req.Description != nil {
		notification.Description = *req.Description
	}

	if err := s.repo.UpdateNotification(s.db.WithContext(ctx), &notification); err != nil {
		return NotificationSettingResponse{}, apperror.Internal("gagal memperbarui notifikasi", err)
	}

	return notification.ToResponse(), nil
}
