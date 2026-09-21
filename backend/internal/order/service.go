package order

import (
	"backend-go/internal/customer"
	"backend-go/internal/entities"
	"backend-go/internal/product"
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"backend-go/pkg/utilities"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================================================================
// Service Interface
// =============================================================================

type ServiceInterface interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *CreateOrderRequest, userID string) (OrderResponse, error)
	GetByID(ctx context.Context, tenantID uuid.UUID, id string) (OrderResponse, error)
	List(ctx context.Context, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[OrderResponse], error)
	UpdateStatus(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateOrderStatusRequest, userID string) (OrderResponse, error)

	// Cross-service (used by payment module)
	GetByIDRaw(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Order, error)
	MarkPaid(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	MarkFailed(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	MarkRefunded(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, userID string) error
}

// =============================================================================
// Service Implementation
// =============================================================================

type Service struct {
	db          *gorm.DB
	repo        RepositoryInterface
	productSvc  product.ServiceInterface
	customerSvc customer.ServiceInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface, productSvc product.ServiceInterface, customerSvc customer.ServiceInterface) ServiceInterface {
	return &Service{
		db:          db,
		repo:        repo,
		productSvc:  productSvc,
		customerSvc: customerSvc,
	}
}

// =============================================================================
// CRUD Implementations
// =============================================================================

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req *CreateOrderRequest, userID string) (OrderResponse, error) {
	var result Order

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderNumber := utilities.GenerateNumberTransactionsScoped(tx, tenantID, "orders", "order_number", "ORD")

		var subtotal float64
		items := make([]OrderItem, 0, len(req.Items))

		for _, itemReq := range req.Items {
			productID, err := uuid.Parse(itemReq.ProductID)
			if err != nil {
				return apperror.BadRequest("product_id tidak valid", err)
			}

			product, err := s.productSvc.GetByIDRaw(tx, tenantID, productID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperror.BadRequest(fmt.Sprintf("produk '%s' tidak ditemukan", itemReq.ProductID), err)
				}
				return apperror.Internal("gagal mengambil data produk", err)
			}

			if product.Status != "active" {
				return apperror.BadRequest(fmt.Sprintf("produk '%s' tidak aktif", product.Name), nil)
			}

			if product.Stock < itemReq.Quantity {
				return apperror.BadRequest(
					fmt.Sprintf("stok produk '%s' tidak mencukupi (tersedia: %d)", product.Name, product.Stock), nil,
				)
			}

			subtotalItem := product.Price * float64(itemReq.Quantity)
			subtotal += subtotalItem

			items = append(items, OrderItem{
				TenantID:     tenantID,
				ProductID:    product.ID,
				ProductName:  product.Name,
				ProductSKU:   product.SKU,
				ProductImage: firstImage(product.Images),
				Quantity:     itemReq.Quantity,
				UnitPrice:    product.Price,
				Subtotal:     subtotalItem,
			})

			if err := s.productSvc.ReserveStock(tx, tenantID, product.ID, itemReq.Quantity); err != nil {
				return err
			}
		}

		// Resolve customer (use provided id or find/create by email)
		var customerID *uuid.UUID
		if req.CustomerID != nil {
			customerID = req.CustomerID
		} else if req.CustomerEmail != "" {
			c, err := s.customerSvc.FindOrCreateByEmail(tx, tenantID, req.CustomerName, req.CustomerEmail)
			if err != nil {
				return apperror.Internal("gagal menyiapkan data customer", err)
			}
			customerID = &c.ID
		}

		total := subtotal + req.ShippingCost - req.Discount

		order := Order{
			TenantID:        tenantID,
			OrderNumber:     orderNumber,
			CustomerID:      customerID,
			CustomerName:    req.CustomerName,
			Status:          StatusPending,
			PaymentStatus:   PaymentPending,
			PaymentMethodID: req.PaymentMethodID,
			ShippingMethod:  req.ShippingMethod,
			ShippingAddress: req.ShippingAddress,
			Subtotal:        subtotal,
			ShippingCost:    req.ShippingCost,
			Discount:        req.Discount,
			TotalAmount:     total,
			Notes:           req.Notes,
			Items:           items,
			CreatedBy:       userID,
		}

		if err := s.repo.Create(tx, &order); err != nil {
			return apperror.Internal("gagal membuat order", err)
		}

		// Timeline: created
		if err := s.repo.AddTimeline(tx, &OrderTimeline{
			TenantID:  tenantID,
			OrderID:   order.ID,
			Type:      "created",
			Title:     "Order dibuat",
			UserName:  userID,
			CreatedAt: order.CreatedAt,
		}); err != nil {
			return apperror.Internal("gagal mencatat timeline", err)
		}

		// Customer stats + activity
		if customerID != nil {
			if err := s.customerSvc.RecordOrderStats(tx, tenantID, *customerID, total); err != nil {
				return apperror.Internal("gagal memperbarui statistik customer", err)
			}
			if err := s.customerSvc.AddActivity(tx, tenantID, *customerID, "order",
				fmt.Sprintf("Membuat order %s senilai Rp %s", orderNumber, utilities.FormatNumber(total))); err != nil {
				return apperror.Internal("gagal mencatat aktivitas customer", err)
			}
		}

		result = order
		return nil
	})

	if txErr != nil {
		return OrderResponse{}, txErr
	}

	return result.ToResponse(), nil
}

func (s *Service) GetByID(ctx context.Context, tenantID uuid.UUID, id string) (OrderResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return OrderResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	order, err := s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderResponse{}, apperror.NotFound("order tidak ditemukan", err)
		}
		return OrderResponse{}, apperror.Internal("gagal mengambil order", err)
	}

	return order.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[OrderResponse], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), tenantID, params)
	if err != nil {
		return pagination.Response[OrderResponse]{}, apperror.Internal("gagal mengambil daftar order", err)
	}

	responses := make([]OrderResponse, len(result.Data))
	for i, o := range result.Data {
		responses[i] = o.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) UpdateStatus(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateOrderStatusRequest, userID string) (OrderResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return OrderResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var result Order

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.GetByID(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("order tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil order", err)
		}

		if !CanTransition(order.Status, req.Status) {
			return apperror.BadRequest(fmt.Sprintf("status tidak dapat berubah dari '%s' ke '%s'", order.Status, req.Status), nil)
		}

		// Stock bookkeeping
		for _, item := range order.Items {
			switch req.Status {
			case StatusShipped:
				if err := s.productSvc.ReleaseReserved(tx, tenantID, item.ProductID, item.Quantity); err != nil {
					return err
				}
			case StatusCancelled:
				if err := s.productSvc.RestoreStock(tx, tenantID, item.ProductID, item.Quantity); err != nil {
					return err
				}
			}
		}

		if err := s.repo.UpdateStatus(tx, tenantID, uid, req.Status, userID); err != nil {
			return apperror.Internal("gagal memperbarui status order", err)
		}

		eventType := string(req.Status)
		titles := map[OrderStatus]string{
			StatusConfirmed:  "Order dikonfirmasi",
			StatusProcessing: "Order diproses",
			StatusShipped:    "Order dikirim",
			StatusDelivered:  "Order selesai",
			StatusCancelled:  "Order dibatalkan",
			StatusRefunded:   "Order direfund",
		}

		if err := s.repo.AddTimeline(tx, &OrderTimeline{
			TenantID: tenantID,
			OrderID:  order.ID,
			Type:     eventType,
			Title:    titles[req.Status],
			UserName: userID,
		}); err != nil {
			return apperror.Internal("gagal mencatat timeline", err)
		}

		result = order
		result.Status = req.Status
		return nil
	})

	if txErr != nil {
		return OrderResponse{}, txErr
	}

	return result.ToResponse(), nil
}

// =============================================================================
// Cross-service methods
// =============================================================================

func (s *Service) GetByIDRaw(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Order, error) {
	return s.repo.GetByID(tx, tenantID, id)
}

func (s *Service) MarkPaid(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.UpdatePaymentStatus(tx, tenantID, id, PaymentPaid); err != nil {
		return err
	}
	return s.repo.AddTimeline(tx, &OrderTimeline{
		TenantID: tenantID,
		OrderID:  id,
		Type:     "payment",
		Title:    "Pembayaran diterima",
	})
}

func (s *Service) MarkFailed(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return s.repo.UpdatePaymentStatus(tx, tenantID, id, PaymentFailed)
}

func (s *Service) MarkRefunded(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, userID string) error {
	if err := s.repo.UpdatePaymentStatus(tx, tenantID, id, PaymentRefunded); err != nil {
		return err
	}
	return s.repo.AddTimeline(tx, &OrderTimeline{
		TenantID: tenantID,
		OrderID:  id,
		Type:     "refunded",
		Title:    "Dana dikembalikan",
		UserName: userID,
	})
}

// firstImage returns the primary image url of a product image list.
func firstImage(images []product.ProductImage) string {
	if len(images) == 0 {
		return ""
	}
	for _, img := range images {
		if img.IsPrimary {
			return img.URL
		}
	}
	return images[0].URL
}
