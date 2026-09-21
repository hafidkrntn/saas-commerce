package order

import (
	"backend-go/internal/entities"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(tx *gorm.DB, order *Order) error
	GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Order, error)
	List(tx *gorm.DB, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[Order], error)
	UpdateStatus(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, status OrderStatus, userID string) error
	UpdatePaymentStatus(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, status PaymentStatus) error
	AddTimeline(tx *gorm.DB, event *OrderTimeline) error
	FindPendingPaymentOrder(tx *gorm.DB, tenantID uuid.UUID, orderID uuid.UUID) (Order, error)
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) Create(tx *gorm.DB, order *Order) error {
	return tx.Create(order).Error
}

func (r *Repository) GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Order, error) {
	var order Order
	err := tx.
		Preload("Items").
		Preload("Timeline", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&order).Error
	if err != nil {
		return Order{}, err
	}

	r.enrich(tx, &order)
	return order, nil
}

// enrich fills transient fields (customer email, payment method name) via joins.
func (r *Repository) enrich(tx *gorm.DB, order *Order) {
	if order.CustomerID != nil {
		tx.Table("customers").Select("email").
			Where("id = ?", *order.CustomerID).Pluck("email", &order.CustomerEmail)
	}
	if order.PaymentMethodID != nil {
		tx.Table("payment_methods").Select("name").
			Where("id = ?", *order.PaymentMethodID).Pluck("name", &order.PaymentMethodName)
	}
}

func (r *Repository) List(tx *gorm.DB, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[Order], error) {
	p := pagination.Resolve(params.Page, params.Limit)

	query := tx.Model(&Order{}).Where("tenant_id = ?", tenantID)

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("order_number ILIKE ? OR customer_name ILIKE ?", search, search)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.DateFrom != "" {
		query = query.Where("created_at >= ?", params.DateFrom)
	}
	if params.DateTo != "" {
		query = query.Where("created_at <= ?", params.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[Order]{}, err
	}

	var orders []Order
	if err := query.
		Preload("Items").
		Order("created_at DESC").
		Offset(p.Offset).
		Limit(p.Limit).
		Find(&orders).Error; err != nil {
		return pagination.Response[Order]{}, err
	}

	for i := range orders {
		r.enrich(tx, &orders[i])
	}

	return pagination.NewResponse(orders, int(total), p.Page, p.Limit), nil
}

func (r *Repository) UpdateStatus(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, status OrderStatus, userID string) error {
	return tx.Model(&Order{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(map[string]any{"status": status, "updated_by": userID}).Error
}

func (r *Repository) UpdatePaymentStatus(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, status PaymentStatus) error {
	return tx.Model(&Order{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("payment_status", status).Error
}

func (r *Repository) AddTimeline(tx *gorm.DB, event *OrderTimeline) error {
	return tx.Create(event).Error
}

func (r *Repository) FindPendingPaymentOrder(tx *gorm.DB, tenantID uuid.UUID, orderID uuid.UUID) (Order, error) {
	var order Order
	err := tx.Where("id = ? AND tenant_id = ? AND payment_status = ?", orderID, tenantID, PaymentPending).
		First(&order).Error
	if err != nil {
		return Order{}, err
	}
	return order, nil
}
