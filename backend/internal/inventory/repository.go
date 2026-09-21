package inventory

import (
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateWarehouse(tx *gorm.DB, w *Warehouse) error
	GetWarehouse(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Warehouse, error)
	ListWarehouses(tx *gorm.DB, tenantID uuid.UUID) ([]Warehouse, error)
	UpdateWarehouse(tx *gorm.DB, w *Warehouse) error
	DeleteWarehouse(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error

	CreateMovement(tx *gorm.DB, m *StockMovement) error
	ListMovements(tx *gorm.DB, tenantID uuid.UUID, productID string, page, limit int) (pagination.Response[StockMovement], error)

	CreateShipment(tx *gorm.DB, s *IncomingShipment) error
	GetShipment(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (IncomingShipment, error)
	ListShipments(tx *gorm.DB, tenantID uuid.UUID, status string, page, limit int) (pagination.Response[IncomingShipment], error)
	UpdateShipment(tx *gorm.DB, s *IncomingShipment) error
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) CreateWarehouse(tx *gorm.DB, w *Warehouse) error {
	return tx.Create(w).Error
}

func (r *Repository) GetWarehouse(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Warehouse, error) {
	var w Warehouse
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&w).Error
	if err != nil {
		return Warehouse{}, err
	}
	return w, nil
}

func (r *Repository) ListWarehouses(tx *gorm.DB, tenantID uuid.UUID) ([]Warehouse, error) {
	var warehouses []Warehouse
	err := tx.Where("tenant_id = ?", tenantID).Order("is_default DESC, name ASC").Find(&warehouses).Error
	if err != nil {
		return nil, err
	}
	return warehouses, nil
}

func (r *Repository) UpdateWarehouse(tx *gorm.DB, w *Warehouse) error {
	return tx.Save(w).Error
}

func (r *Repository) DeleteWarehouse(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&Warehouse{}).Error
}

func (r *Repository) CreateMovement(tx *gorm.DB, m *StockMovement) error {
	return tx.Create(m).Error
}

func (r *Repository) ListMovements(tx *gorm.DB, tenantID uuid.UUID, productID string, page, limit int) (pagination.Response[StockMovement], error) {
	p := pagination.Resolve(page, limit)

	query := tx.Model(&StockMovement{}).Where("tenant_id = ?", tenantID)
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[StockMovement]{}, err
	}

	var movements []StockMovement
	if err := query.Order("created_at DESC").Offset(p.Offset).Limit(p.Limit).Find(&movements).Error; err != nil {
		return pagination.Response[StockMovement]{}, err
	}

	return pagination.NewResponse(movements, int(total), p.Page, p.Limit), nil
}

func (r *Repository) CreateShipment(tx *gorm.DB, s *IncomingShipment) error {
	return tx.Create(s).Error
}

func (r *Repository) GetShipment(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (IncomingShipment, error) {
	var s IncomingShipment
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&s).Error
	if err != nil {
		return IncomingShipment{}, err
	}
	return s, nil
}

func (r *Repository) ListShipments(tx *gorm.DB, tenantID uuid.UUID, status string, page, limit int) (pagination.Response[IncomingShipment], error) {
	p := pagination.Resolve(page, limit)

	query := tx.Model(&IncomingShipment{}).Where("tenant_id = ?", tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[IncomingShipment]{}, err
	}

	var shipments []IncomingShipment
	if err := query.Order("expected_date ASC, created_at DESC").Offset(p.Offset).Limit(p.Limit).Find(&shipments).Error; err != nil {
		return pagination.Response[IncomingShipment]{}, err
	}

	return pagination.NewResponse(shipments, int(total), p.Page, p.Limit), nil
}

func (r *Repository) UpdateShipment(tx *gorm.DB, s *IncomingShipment) error {
	return tx.Save(s).Error
}
