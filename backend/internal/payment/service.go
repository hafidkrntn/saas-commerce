package payment

import (
	"backend-go/internal/entities"
	"backend-go/internal/order"
	"backend-go/pkg/apperror"
	paymentclient "backend-go/third_party/payment"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	ListMethods(ctx context.Context, tenantID uuid.UUID) ([]PaymentMethodResponse, error)
	CreateMethod(ctx context.Context, tenantID uuid.UUID, req *CreatePaymentMethodRequest) (PaymentMethodResponse, error)
	UpdateMethod(ctx context.Context, tenantID uuid.UUID, id string, req *UpdatePaymentMethodRequest) (PaymentMethodResponse, error)
	DeleteMethod(ctx context.Context, tenantID uuid.UUID, id string) error

	Charge(ctx context.Context, tenantID uuid.UUID, req *ChargeRequest) (ChargeResponse, error)
	GetPaymentByOrder(ctx context.Context, tenantID uuid.UUID, orderID string) (PaymentResponse, error)
	HandleWebhook(ctx context.Context, tenantID uuid.UUID, transactionID, status string) (PaymentResponse, error)
}

type Service struct {
	db       *gorm.DB
	repo     RepositoryInterface
	orderSvc order.ServiceInterface
	provider paymentclient.ClientInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface, orderSvc order.ServiceInterface, provider paymentclient.ClientInterface) ServiceInterface {
	return &Service{
		db:       db,
		repo:     repo,
		orderSvc: orderSvc,
		provider: provider,
	}
}

// =============================================================================
// Payment methods
// =============================================================================

func (s *Service) ListMethods(ctx context.Context, tenantID uuid.UUID) ([]PaymentMethodResponse, error) {
	methods, err := s.repo.ListPaymentMethods(s.db.WithContext(ctx), tenantID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar metode pembayaran", err)
	}

	responses := make([]PaymentMethodResponse, len(methods))
	for i, m := range methods {
		responses[i] = m.ToResponse()
	}
	return responses, nil
}

func (s *Service) CreateMethod(ctx context.Context, tenantID uuid.UUID, req *CreatePaymentMethodRequest) (PaymentMethodResponse, error) {
	isEnabled := true
	if req.IsEnabled != nil {
		isEnabled = *req.IsEnabled
	}

	method := PaymentMethod{
		TenantID:    tenantID,
		Name:        req.Name,
		Type:        req.Type,
		IsEnabled:   isEnabled,
		Description: req.Description,
		Config:      entities.JSONB(req.Config),
	}

	if err := s.repo.CreatePaymentMethod(s.db.WithContext(ctx), &method); err != nil {
		return PaymentMethodResponse{}, apperror.Internal("gagal membuat metode pembayaran", err)
	}

	return method.ToResponse(), nil
}

func (s *Service) UpdateMethod(ctx context.Context, tenantID uuid.UUID, id string, req *UpdatePaymentMethodRequest) (PaymentMethodResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return PaymentMethodResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	method, err := s.repo.GetPaymentMethod(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PaymentMethodResponse{}, apperror.NotFound("metode pembayaran tidak ditemukan", err)
		}
		return PaymentMethodResponse{}, apperror.Internal("gagal mengambil metode pembayaran", err)
	}

	if req.Name != nil {
		method.Name = *req.Name
	}
	if req.Type != nil {
		method.Type = *req.Type
	}
	if req.IsEnabled != nil {
		method.IsEnabled = *req.IsEnabled
	}
	if req.Description != nil {
		method.Description = *req.Description
	}
	if req.Config != nil {
		method.Config = entities.JSONB(req.Config)
	}

	if err := s.repo.UpdatePaymentMethod(s.db.WithContext(ctx), &method); err != nil {
		return PaymentMethodResponse{}, apperror.Internal("gagal memperbarui metode pembayaran", err)
	}

	return method.ToResponse(), nil
}

func (s *Service) DeleteMethod(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetPaymentMethod(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("metode pembayaran tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil metode pembayaran", err)
	}

	if err := s.repo.DeletePaymentMethod(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus metode pembayaran", err)
	}

	return nil
}

// =============================================================================
// Charge & webhook
// =============================================================================

// Charge creates a payment for a pending order and calls the provider.
func (s *Service) Charge(ctx context.Context, tenantID uuid.UUID, req *ChargeRequest) (ChargeResponse, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return ChargeResponse{}, apperror.BadRequest("order_id tidak valid", err)
	}

	var result ChargeResponse

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Order must exist and still be awaiting payment
		order, err := s.repo.GetPaymentByOrder(tx, tenantID, orderID)
		if err == nil && order.Status == PayStatusPaid {
			return apperror.BadRequest("order sudah dibayar", nil)
		}

		ord, err := s.orderSvc.GetByIDRaw(tx, tenantID, orderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("order tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil order", err)
		}

		payment := Payment{
			TenantID:        tenantID,
			OrderID:         orderID,
			PaymentMethodID: ord.PaymentMethodID,
			Provider:        "mock",
			Status:          PayStatusPending,
			Amount:          ord.TotalAmount,
			Currency:        "IDR",
		}

		if err := s.repo.CreatePayment(tx, &payment); err != nil {
			return apperror.Internal("gagal membuat pembayaran", err)
		}

		charge, err := s.provider.CreateCharge(ctx, &paymentclient.CreateChargeRequest{
			OrderID:     orderID.String(),
			Amount:      ord.TotalAmount,
			Currency:    "IDR",
			Description: "Order " + ord.OrderNumber,
		})
		if err != nil {
			return apperror.Internal("gagal menghubungi payment provider", err)
		}

		if err := tx.Model(&Payment{}).Where("id = ?", payment.ID).
			Update("provider_transaction_id", charge.ID).Error; err != nil {
			return apperror.Internal("gagal memperbarui pembayaran", err)
		}
		payment.ProviderTransactionID = charge.ID

		result = ChargeResponse{
			Payment:     payment.ToResponse(),
			RedirectURL: charge.RedirectURL,
			Status:      charge.Status,
		}
		return nil
	})

	if txErr != nil {
		return ChargeResponse{}, txErr
	}

	return result, nil
}

func (s *Service) GetPaymentByOrder(ctx context.Context, tenantID uuid.UUID, orderID string) (PaymentResponse, error) {
	oid, err := uuid.Parse(orderID)
	if err != nil {
		return PaymentResponse{}, apperror.BadRequest("order_id tidak valid", err)
	}

	payment, err := s.repo.GetPaymentByOrder(s.db.WithContext(ctx), tenantID, oid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PaymentResponse{}, apperror.NotFound("pembayaran tidak ditemukan", err)
		}
		return PaymentResponse{}, apperror.Internal("gagal mengambil pembayaran", err)
	}

	return payment.ToResponse(), nil
}

// HandleWebhook processes provider callbacks (status updates).
func (s *Service) HandleWebhook(ctx context.Context, tenantID uuid.UUID, transactionID, status string) (PaymentResponse, error) {
	var result Payment

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment, err := s.repo.GetPaymentByTransaction(tx, tenantID, transactionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("transaksi tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil pembayaran", err)
		}

		now := time.Now()
		var paidAt *time.Time
		switch status {
		case PayStatusPaid:
			paidAt = &now
		case PayStatusFailed, PayStatusPending, PayStatusRefunded:
		default:
			return apperror.BadRequest("status tidak dikenal", nil)
		}

		if err := s.repo.UpdatePaymentStatus(tx, payment.ID, status, paidAt); err != nil {
			return apperror.Internal("gagal memperbarui status pembayaran", err)
		}
		payment.Status = status
		payment.PaidAt = paidAt

		// Mirror to the order
		switch status {
		case PayStatusPaid:
			if err := s.orderSvc.MarkPaid(tx, tenantID, payment.OrderID); err != nil {
				return apperror.Internal("gagal memperbarui status order", err)
			}
		case PayStatusFailed:
			if err := s.orderSvc.MarkFailed(tx, tenantID, payment.OrderID); err != nil {
				return apperror.Internal("gagal memperbarui status order", err)
			}
		}

		result = payment
		return nil
	})

	if txErr != nil {
		return PaymentResponse{}, txErr
	}

	return result.ToResponse(), nil
}
