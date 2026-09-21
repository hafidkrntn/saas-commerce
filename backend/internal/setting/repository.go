package setting

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	GetStoreSetting(tx *gorm.DB, tenantID uuid.UUID) (StoreSetting, error)
	UpsertStoreSetting(tx *gorm.DB, s *StoreSetting) error
	GetCompanySetting(tx *gorm.DB, tenantID uuid.UUID) (CompanySetting, error)
	UpsertCompanySetting(tx *gorm.DB, c *CompanySetting) error

	CreateTaxRate(tx *gorm.DB, t *TaxRate) error
	GetTaxRate(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (TaxRate, error)
	ListTaxRates(tx *gorm.DB, tenantID uuid.UUID) ([]TaxRate, error)
	UpdateTaxRate(tx *gorm.DB, t *TaxRate) error
	DeleteTaxRate(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error

	CreateShippingZone(tx *gorm.DB, z *ShippingZone) error
	GetShippingZone(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (ShippingZone, error)
	ListShippingZones(tx *gorm.DB, tenantID uuid.UUID) ([]ShippingZone, error)
	UpdateShippingZone(tx *gorm.DB, z *ShippingZone) error
	DeleteShippingZone(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error

	CreateNotification(tx *gorm.DB, n *NotificationSetting) error
	GetNotification(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (NotificationSetting, error)
	ListNotifications(tx *gorm.DB, tenantID uuid.UUID) ([]NotificationSetting, error)
	UpdateNotification(tx *gorm.DB, n *NotificationSetting) error
	DeleteNotification(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

// --- Store & Company ---

func (r *Repository) GetStoreSetting(tx *gorm.DB, tenantID uuid.UUID) (StoreSetting, error) {
	var s StoreSetting
	err := tx.Where("tenant_id = ?", tenantID).First(&s).Error
	if err != nil {
		return StoreSetting{}, err
	}
	return s, nil
}

func (r *Repository) UpsertStoreSetting(tx *gorm.DB, s *StoreSetting) error {
	return tx.Save(s).Error
}

func (r *Repository) GetCompanySetting(tx *gorm.DB, tenantID uuid.UUID) (CompanySetting, error) {
	var c CompanySetting
	err := tx.Where("tenant_id = ?", tenantID).First(&c).Error
	if err != nil {
		return CompanySetting{}, err
	}
	return c, nil
}

func (r *Repository) UpsertCompanySetting(tx *gorm.DB, c *CompanySetting) error {
	return tx.Save(c).Error
}

// --- Tax rates ---

func (r *Repository) CreateTaxRate(tx *gorm.DB, t *TaxRate) error {
	return tx.Create(t).Error
}

func (r *Repository) GetTaxRate(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (TaxRate, error) {
	var t TaxRate
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&t).Error
	if err != nil {
		return TaxRate{}, err
	}
	return t, nil
}

func (r *Repository) ListTaxRates(tx *gorm.DB, tenantID uuid.UUID) ([]TaxRate, error) {
	var rates []TaxRate
	err := tx.Where("tenant_id = ?", tenantID).Order("name ASC").Find(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *Repository) UpdateTaxRate(tx *gorm.DB, t *TaxRate) error {
	return tx.Save(t).Error
}

func (r *Repository) DeleteTaxRate(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&TaxRate{}).Error
}

// --- Shipping zones ---

func (r *Repository) CreateShippingZone(tx *gorm.DB, z *ShippingZone) error {
	return tx.Create(z).Error
}

func (r *Repository) GetShippingZone(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (ShippingZone, error) {
	var z ShippingZone
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&z).Error
	if err != nil {
		return ShippingZone{}, err
	}
	return z, nil
}

func (r *Repository) ListShippingZones(tx *gorm.DB, tenantID uuid.UUID) ([]ShippingZone, error) {
	var zones []ShippingZone
	err := tx.Where("tenant_id = ?", tenantID).Order("name ASC").Find(&zones).Error
	if err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *Repository) UpdateShippingZone(tx *gorm.DB, z *ShippingZone) error {
	return tx.Save(z).Error
}

func (r *Repository) DeleteShippingZone(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&ShippingZone{}).Error
}

// --- Notifications ---

func (r *Repository) CreateNotification(tx *gorm.DB, n *NotificationSetting) error {
	return tx.Create(n).Error
}

func (r *Repository) GetNotification(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (NotificationSetting, error) {
	var n NotificationSetting
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&n).Error
	if err != nil {
		return NotificationSetting{}, err
	}
	return n, nil
}

func (r *Repository) ListNotifications(tx *gorm.DB, tenantID uuid.UUID) ([]NotificationSetting, error) {
	var notifications []NotificationSetting
	err := tx.Where("tenant_id = ?", tenantID).Order("event ASC").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *Repository) UpdateNotification(tx *gorm.DB, n *NotificationSetting) error {
	return tx.Save(n).Error
}

func (r *Repository) DeleteNotification(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&NotificationSetting{}).Error
}
