package inventory

import (
	"backend-go/internal/product"
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	ListWarehouses(ctx context.Context, tenantID uuid.UUID) ([]WarehouseResponse, error)
	CreateWarehouse(ctx context.Context, tenantID uuid.UUID, req *CreateWarehouseRequest, userID string) (WarehouseResponse, error)
	UpdateWarehouse(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateWarehouseRequest) (WarehouseResponse, error)
	DeleteWarehouse(ctx context.Context, tenantID uuid.UUID, id string) error

	CreateMovement(ctx context.Context, tenantID uuid.UUID, req *StockMovementRequest, userID string) (StockMovementResponse, error)
	ListMovements(ctx context.Context, tenantID uuid.UUID, productID string, page, limit int) (pagination.Response[StockMovementResponse], error)

	CreateShipment(ctx context.Context, tenantID uuid.UUID, req *CreateShipmentRequest, userID string) (ShipmentResponse, error)
	ListShipments(ctx context.Context, tenantID uuid.UUID, status string, page, limit int) (pagination.Response[ShipmentResponse], error)
	UpdateShipment(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateShipmentRequest, userID string) (ShipmentResponse, error)
}

type Service struct {
	db         *gorm.DB
	repo       RepositoryInterface
	productSvc product.ServiceInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface, productSvc product.ServiceInterface) ServiceInterface {
	return &Service{db: db, repo: repo, productSvc: productSvc}
}

// =============================================================================
// Warehouses
// =============================================================================

func (s *Service) ListWarehouses(ctx context.Context, tenantID uuid.UUID) ([]WarehouseResponse, error) {
	warehouses, err := s.repo.ListWarehouses(s.db.WithContext(ctx), tenantID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar gudang", err)
	}

	responses := make([]WarehouseResponse, len(warehouses))
	for i, w := range warehouses {
		responses[i] = w.ToResponse()
	}
	return responses, nil
}

func (s *Service) CreateWarehouse(ctx context.Context, tenantID uuid.UUID, req *CreateWarehouseRequest, userID string) (WarehouseResponse, error) {
	warehouse := Warehouse{
		TenantID:  tenantID,
		Name:      req.Name,
		Address:   req.Address,
		IsDefault: req.IsDefault,
		CreatedBy: userID,
	}

	if err := s.repo.CreateWarehouse(s.db.WithContext(ctx), &warehouse); err != nil {
		return WarehouseResponse{}, apperror.Internal("gagal membuat gudang", err)
	}

	return warehouse.ToResponse(), nil
}

func (s *Service) UpdateWarehouse(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateWarehouseRequest) (WarehouseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return WarehouseResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	warehouse, err := s.repo.GetWarehouse(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return WarehouseResponse{}, apperror.NotFound("gudang tidak ditemukan", err)
		}
		return WarehouseResponse{}, apperror.Internal("gagal mengambil gudang", err)
	}

	if req.Name != nil {
		warehouse.Name = *req.Name
	}
	if req.Address != nil {
		warehouse.Address = *req.Address
	}
	if req.IsDefault != nil {
		warehouse.IsDefault = *req.IsDefault
	}

	if err := s.repo.UpdateWarehouse(s.db.WithContext(ctx), &warehouse); err != nil {
		return WarehouseResponse{}, apperror.Internal("gagal memperbarui gudang", err)
	}

	return warehouse.ToResponse(), nil
}

func (s *Service) DeleteWarehouse(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetWarehouse(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("gudang tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil gudang", err)
	}

	if err := s.repo.DeleteWarehouse(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus gudang", err)
	}

	return nil
}

// =============================================================================
// Stock movements
// =============================================================================

// CreateMovement records a movement and adjusts the product stock.
func (s *Service) CreateMovement(ctx context.Context, tenantID uuid.UUID, req *StockMovementRequest, userID string) (StockMovementResponse, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return StockMovementResponse{}, apperror.BadRequest("product_id tidak valid", err)
	}

	var movement StockMovement

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.productSvc.GetByIDRaw(tx, tenantID, productID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("produk tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil produk", err)
		}

		switch req.Type {
		case "in", "return":
			if err := s.productSvc.AddStock(tx, tenantID, productID, req.Quantity, false); err != nil {
				return err
			}
		case "out":
			if err := s.productSvc.ReserveStock(tx, tenantID, productID, req.Quantity); err != nil {
				return err
			}
			if err := s.productSvc.ReleaseReserved(tx, tenantID, productID, req.Quantity); err != nil {
				return err
			}
		case "adjustment":
			if req.Quantity > 0 {
				if err := s.productSvc.AddStock(tx, tenantID, productID, req.Quantity, false); err != nil {
					return err
				}
			}
		}

		movement = StockMovement{
			TenantID:  tenantID,
			ProductID: productID,
			Type:      req.Type,
			Quantity:  req.Quantity,
			Reference: req.Reference,
			Note:      req.Note,
			UserName:  userID,
		}

		if err := s.repo.CreateMovement(tx, &movement); err != nil {
			return apperror.Internal("gagal mencatat pergerakan stok", err)
		}
		return nil
	})

	if txErr != nil {
		return StockMovementResponse{}, txErr
	}

	return movement.ToResponse(), nil
}

func (s *Service) ListMovements(ctx context.Context, tenantID uuid.UUID, productID string, page, limit int) (pagination.Response[StockMovementResponse], error) {
	result, err := s.repo.ListMovements(s.db.WithContext(ctx), tenantID, productID, page, limit)
	if err != nil {
		return pagination.Response[StockMovementResponse]{}, apperror.Internal("gagal mengambil daftar pergerakan stok", err)
	}

	responses := make([]StockMovementResponse, len(result.Data))
	for i, m := range result.Data {
		responses[i] = m.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

// =============================================================================
// Incoming shipments
// =============================================================================

func (s *Service) CreateShipment(ctx context.Context, tenantID uuid.UUID, req *CreateShipmentRequest, userID string) (ShipmentResponse, error) {
	var shipment IncomingShipment

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.productSvc.GetByIDRaw(tx, tenantID, req.ProductID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("produk tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil produk", err)
		}

		shipment = IncomingShipment{
			TenantID:     tenantID,
			ProductID:    req.ProductID,
			Quantity:     req.Quantity,
			ExpectedDate: req.ExpectedDate,
			Supplier:     req.Supplier,
			Status:       ShipmentScheduled,
			Notes:        req.Notes,
			CreatedBy:    userID,
		}

		if err := s.repo.CreateShipment(tx, &shipment); err != nil {
			return apperror.Internal("gagal membuat shipment", err)
		}
		return nil
	})

	if txErr != nil {
		return ShipmentResponse{}, txErr
	}

	return shipment.ToResponse(), nil
}

func (s *Service) ListShipments(ctx context.Context, tenantID uuid.UUID, status string, page, limit int) (pagination.Response[ShipmentResponse], error) {
	result, err := s.repo.ListShipments(s.db.WithContext(ctx), tenantID, status, page, limit)
	if err != nil {
		return pagination.Response[ShipmentResponse]{}, apperror.Internal("gagal mengambil daftar shipment", err)
	}

	responses := make([]ShipmentResponse, len(result.Data))
	for i, sh := range result.Data {
		responses[i] = sh.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) UpdateShipment(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateShipmentRequest, userID string) (ShipmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ShipmentResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var shipment IncomingShipment

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shipment, err = s.repo.GetShipment(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("shipment tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil shipment", err)
		}

		previousStatus := shipment.Status

		if req.Quantity != nil {
			shipment.Quantity = *req.Quantity
		}
		if req.ExpectedDate != nil {
			shipment.ExpectedDate = req.ExpectedDate
		}
		if req.Supplier != nil {
			shipment.Supplier = *req.Supplier
		}
		if req.Status != nil {
			shipment.Status = *req.Status
		}
		if req.Notes != nil {
			shipment.Notes = *req.Notes
		}
		shipment.UpdatedBy = userID

		if err := s.repo.UpdateShipment(tx, &shipment); err != nil {
			return apperror.Internal("gagal memperbarui shipment", err)
		}

		// On delivery, stock the product and record a movement
		if shipment.Status == ShipmentDelivered && previousStatus != ShipmentDelivered {
			if err := s.productSvc.AddStock(tx, tenantID, shipment.ProductID, shipment.Quantity, true); err != nil {
				return err
			}
			if err := s.repo.CreateMovement(tx, &StockMovement{
				TenantID:  tenantID,
				ProductID: shipment.ProductID,
				Type:      "in",
				Quantity:  shipment.Quantity,
				Reference: fmt.Sprintf("SHIP/%s", shipment.ID.String()[:8]),
				Note:      "Barang masuk dari shipment",
				UserName:  userID,
			}); err != nil {
				return apperror.Internal("gagal mencatat pergerakan stok", err)
			}
		}

		return nil
	})

	if txErr != nil {
		return ShipmentResponse{}, txErr
	}

	return shipment.ToResponse(), nil
}
