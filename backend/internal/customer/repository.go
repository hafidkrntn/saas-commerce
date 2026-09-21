package customer

import (
	"backend-go/internal/entities"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(tx *gorm.DB, c *Customer) error
	GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Customer, error)
	GetByEmail(tx *gorm.DB, tenantID uuid.UUID, email string) (Customer, error)
	List(tx *gorm.DB, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[Customer], error)
	Update(tx *gorm.DB, c *Customer) error
	Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	RecordOrderStats(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, amount float64, orderAt any) error
	AddAddress(tx *gorm.DB, a *CustomerAddress) error
	UpdateAddress(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, a *CustomerAddress) error
	DeleteAddress(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, id uuid.UUID) error
	AddActivity(tx *gorm.DB, a *CustomerActivity) error
	ListActivities(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, limit int) ([]CustomerActivity, error)
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) Create(tx *gorm.DB, c *Customer) error {
	return tx.Create(c).Error
}

func (r *Repository) GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Customer, error) {
	var c Customer
	err := tx.Preload("Addresses").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&c).Error
	if err != nil {
		return Customer{}, err
	}
	return c, nil
}

func (r *Repository) GetByEmail(tx *gorm.DB, tenantID uuid.UUID, email string) (Customer, error) {
	var c Customer
	err := tx.Where("email = ? AND tenant_id = ?", email, tenantID).First(&c).Error
	if err != nil {
		return Customer{}, err
	}
	return c, nil
}

func (r *Repository) List(tx *gorm.DB, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[Customer], error) {
	p := pagination.Resolve(params.Page, params.Limit)

	query := tx.Model(&Customer{}).Where("tenant_id = ?", tenantID)

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", search, search, search)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[Customer]{}, err
	}

	var customers []Customer
	if err := query.Order("created_at DESC").Offset(p.Offset).Limit(p.Limit).Find(&customers).Error; err != nil {
		return pagination.Response[Customer]{}, err
	}

	return pagination.NewResponse(customers, int(total), p.Page, p.Limit), nil
}

func (r *Repository) Update(tx *gorm.DB, c *Customer) error {
	return tx.Save(c).Error
}

func (r *Repository) Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&Customer{}).Error
}

// RecordOrderStats increments order counters after an order is placed.
func (r *Repository) RecordOrderStats(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, amount float64, orderAt any) error {
	return tx.Model(&Customer{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(map[string]any{
			"total_orders":       gorm.Expr("total_orders + 1"),
			"total_spent":        gorm.Expr("total_spent + ?", amount),
			"lifetime_value":     gorm.Expr("lifetime_value + ?", amount),
			"average_order_value": gorm.Expr("CASE WHEN total_orders > 0 THEN (total_spent + ?) / (total_orders + 1) ELSE ? END", amount, amount),
			"last_order_at":      orderAt,
		}).Error
}

func (r *Repository) AddAddress(tx *gorm.DB, a *CustomerAddress) error {
	return tx.Create(a).Error
}

func (r *Repository) UpdateAddress(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, a *CustomerAddress) error {
	return tx.Model(&CustomerAddress{}).
		Where("id = ? AND customer_id = ? AND tenant_id = ?", a.ID, customerID, tenantID).
		Save(a).Error
}

func (r *Repository) DeleteAddress(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND customer_id = ? AND tenant_id = ?", id, customerID, tenantID).
		Delete(&CustomerAddress{}).Error
}

func (r *Repository) AddActivity(tx *gorm.DB, a *CustomerActivity) error {
	return tx.Create(a).Error
}

func (r *Repository) ListActivities(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, limit int) ([]CustomerActivity, error) {
	var activities []CustomerActivity
	err := tx.
		Where("customer_id = ? AND tenant_id = ?", customerID, tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}
