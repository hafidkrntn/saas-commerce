package payment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreatePaymentMethod(tx *gorm.DB, m *PaymentMethod) error
	GetPaymentMethod(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (PaymentMethod, error)
	ListPaymentMethods(tx *gorm.DB, tenantID uuid.UUID) ([]PaymentMethod, error)
	UpdatePaymentMethod(tx *gorm.DB, m *PaymentMethod) error
	DeletePaymentMethod(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error

	CreatePayment(tx *gorm.DB, p *Payment) error
	GetPaymentByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Payment, error)
	GetPaymentByTransaction(tx *gorm.DB, tenantID uuid.UUID, transactionID string) (Payment, error)
	GetPaymentByOrder(tx *gorm.DB, tenantID uuid.UUID, orderID uuid.UUID) (Payment, error)
	UpdatePaymentStatus(tx *gorm.DB, id uuid.UUID, status string, paidAt *time.Time) error
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) CreatePaymentMethod(tx *gorm.DB, m *PaymentMethod) error {
	return tx.Create(m).Error
}

func (r *Repository) GetPaymentMethod(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (PaymentMethod, error) {
	var m PaymentMethod
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&m).Error
	if err != nil {
		return PaymentMethod{}, err
	}
	return m, nil
}

func (r *Repository) ListPaymentMethods(tx *gorm.DB, tenantID uuid.UUID) ([]PaymentMethod, error) {
	var methods []PaymentMethod
	err := tx.Where("tenant_id = ?", tenantID).Order("name ASC").Find(&methods).Error
	if err != nil {
		return nil, err
	}
	return methods, nil
}

func (r *Repository) UpdatePaymentMethod(tx *gorm.DB, m *PaymentMethod) error {
	return tx.Save(m).Error
}

func (r *Repository) DeletePaymentMethod(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&PaymentMethod{}).Error
}

func (r *Repository) CreatePayment(tx *gorm.DB, p *Payment) error {
	return tx.Create(p).Error
}

func (r *Repository) GetPaymentByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Payment, error) {
	var p Payment
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&p).Error
	if err != nil {
		return Payment{}, err
	}
	return p, nil
}

func (r *Repository) GetPaymentByTransaction(tx *gorm.DB, tenantID uuid.UUID, transactionID string) (Payment, error) {
	var p Payment
	err := tx.Where("provider_transaction_id = ? AND tenant_id = ?", transactionID, tenantID).First(&p).Error
	if err != nil {
		return Payment{}, err
	}
	return p, nil
}

func (r *Repository) GetPaymentByOrder(tx *gorm.DB, tenantID uuid.UUID, orderID uuid.UUID) (Payment, error) {
	var p Payment
	err := tx.Where("order_id = ? AND tenant_id = ?", orderID, tenantID).Order("created_at DESC").First(&p).Error
	if err != nil {
		return Payment{}, err
	}
	return p, nil
}

func (r *Repository) UpdatePaymentStatus(tx *gorm.DB, id uuid.UUID, status string, paidAt *time.Time) error {
	return tx.Model(&Payment{}).
		Where("id = ?", id).
		Updates(map[string]any{"status": status, "paid_at": paidAt}).Error
}
