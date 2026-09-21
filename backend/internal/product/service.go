package product

import (
	"backend-go/internal/entities"
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================================================================
// Service Interface
// =============================================================================

type ServiceInterface interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *CreateProductRequest, userID string) (ProductResponse, error)
	GetByID(ctx context.Context, tenantID uuid.UUID, id string) (ProductResponse, error)
	List(ctx context.Context, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[ProductResponse], error)
	Update(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateProductRequest, userID string) (ProductResponse, error)
	Delete(ctx context.Context, tenantID uuid.UUID, id string) error

	// Cross-service methods (accept tx so callers can share a transaction)
	GetBySKU(tx *gorm.DB, tenantID uuid.UUID, sku string) (Product, error)
	GetByIDRaw(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Product, error)
	ReserveStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int) error
	ReleaseReserved(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int) error
	RestoreStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int) error
	AddStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int, incoming bool) error
	UpdateRating(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, rating float64, reviewCount int) error
}

// =============================================================================
// Service Implementation
// =============================================================================

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

// =============================================================================
// CRUD Implementations
// =============================================================================

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req *CreateProductRequest, userID string) (ProductResponse, error) {
	var product Product

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existing, err := s.repo.GetBySKU(tx, tenantID, strings.TrimSpace(req.SKU))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Internal("gagal memeriksa sku produk", err)
		}
		if existing.ID != uuid.Nil {
			return apperror.BadRequest("sku produk sudah digunakan", nil)
		}

		status := req.Status
		if status == "" {
			status = StatusActive
		}

		product = Product{
			TenantID:          tenantID,
			SKU:               strings.TrimSpace(req.SKU),
			Name:              req.Name,
			Description:       req.Description,
			Price:             req.Price,
			CompareAtPrice:    req.CompareAtPrice,
			CostPrice:         req.CostPrice,
			Stock:             req.Stock,
			LowStockThreshold: req.LowStockThreshold,
			Status:            status,
			CategoryID:        req.CategoryID,
			Tags:              req.Tags,
			IsActive:          true,
			CreatedBy:         userID,
		}

		if err := s.repo.Create(tx, &product); err != nil {
			return apperror.Internal("gagal membuat produk", err)
		}

		if len(req.Images) > 0 {
			if err := s.repo.ReplaceImages(tx, product.ID, req.Images); err != nil {
				return apperror.Internal("gagal menyimpan gambar produk", err)
			}
			// Reload so the response carries the fresh image list
			product, err = s.repo.GetByID(tx, tenantID, product.ID)
			if err != nil {
				return apperror.Internal("gagal memuat ulang produk", err)
			}
		}

		return nil
	})

	if err != nil {
		return ProductResponse{}, err
	}

	return product.ToResponse(), nil
}

func (s *Service) GetByID(ctx context.Context, tenantID uuid.UUID, id string) (ProductResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ProductResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	product, err := s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProductResponse{}, apperror.NotFound("produk tidak ditemukan", err)
		}
		return ProductResponse{}, apperror.Internal("gagal mengambil produk", err)
	}

	return product.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[ProductResponse], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), tenantID, params)
	if err != nil {
		return pagination.Response[ProductResponse]{}, apperror.Internal("gagal mengambil daftar produk", err)
	}

	responses := make([]ProductResponse, len(result.Data))
	for i, p := range result.Data {
		responses[i] = p.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) Update(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateProductRequest, userID string) (ProductResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ProductResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var product Product

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product, err = s.repo.GetByID(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("produk tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil produk", err)
		}

		if req.SKU != nil {
			product.SKU = strings.TrimSpace(*req.SKU)
		}
		if req.Name != nil {
			product.Name = *req.Name
		}
		if req.Description != nil {
			product.Description = *req.Description
		}
		if req.Price != nil {
			product.Price = *req.Price
		}
		if req.CompareAtPrice != nil {
			product.CompareAtPrice = *req.CompareAtPrice
		}
		if req.CostPrice != nil {
			product.CostPrice = *req.CostPrice
		}
		if req.LowStockThreshold != nil {
			product.LowStockThreshold = *req.LowStockThreshold
		}
		if req.Status != nil {
			product.Status = *req.Status
		}
		if req.CategoryID != nil {
			product.CategoryID = req.CategoryID
		}
		if req.Tags != nil {
			product.Tags = req.Tags
		}
		if req.IsActive != nil {
			product.IsActive = *req.IsActive
		}
		product.UpdatedBy = userID

		if err := s.repo.Update(tx, &product); err != nil {
			return apperror.Internal("gagal memperbarui produk", err)
		}

		if req.Images != nil {
			if err := s.repo.ReplaceImages(tx, product.ID, req.Images); err != nil {
				return apperror.Internal("gagal menyimpan gambar produk", err)
			}
			// Reload so the response carries the fresh image list
			product, err = s.repo.GetByID(tx, tenantID, product.ID)
			if err != nil {
				return apperror.Internal("gagal memuat ulang produk", err)
			}
		}

		return nil
	})

	if txErr != nil {
		return ProductResponse{}, txErr
	}

	return product.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("produk tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil produk", err)
	}

	if err := s.repo.Delete(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus produk", err)
	}

	return nil
}

// =============================================================================
// Cross-service methods
// =============================================================================

func (s *Service) GetBySKU(tx *gorm.DB, tenantID uuid.UUID, sku string) (Product, error) {
	return s.repo.GetBySKU(tx, tenantID, sku)
}

func (s *Service) GetByIDRaw(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Product, error) {
	return s.repo.GetByID(tx, tenantID, id)
}

func (s *Service) adjustStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int, adjust func(p *Product) bool) error {
	product, err := s.repo.GetByID(tx, tenantID, id)
	if err != nil {
		return err
	}

	if !adjust(&product) {
		return apperror.BadRequest("stok tidak mencukupi", nil,
			apperror.F("available", product.Stock),
			apperror.F("requested", qty),
		)
	}

	return s.repo.Update(tx, &product)
}

// ReserveStock deducts available stock and adds to reserved (order placed).
func (s *Service) ReserveStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int) error {
	return s.adjustStock(tx, tenantID, id, qty, func(p *Product) bool {
		if p.Stock < qty {
			return false
		}
		p.Stock -= qty
		p.Reserved += qty
		return true
	})
}

// ReleaseReserved removes items from reserved (order shipped/delivered).
func (s *Service) ReleaseReserved(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int) error {
	return s.adjustStock(tx, tenantID, id, qty, func(p *Product) bool {
		if p.Reserved < qty {
			p.Reserved = 0
			return true
		}
		p.Reserved -= qty
		return true
	})
}

// RestoreStock returns reserved stock to available (order cancelled).
func (s *Service) RestoreStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int) error {
	return s.adjustStock(tx, tenantID, id, qty, func(p *Product) bool {
		if p.Reserved < qty {
			p.Reserved = 0
		} else {
			p.Reserved -= qty
		}
		p.Stock += qty
		return true
	})
}

// AddStock increases stock; when incoming is true the incoming counter is
// decremented (shipment received) instead of tracked.
func (s *Service) AddStock(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, qty int, incoming bool) error {
	return s.adjustStock(tx, tenantID, id, qty, func(p *Product) bool {
		p.Stock += qty
		if incoming && p.Incoming >= qty {
			p.Incoming -= qty
		}
		return true
	})
}

// UpdateRating updates the product's rating aggregates (from review module).
func (s *Service) UpdateRating(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID, rating float64, reviewCount int) error {
	return s.repo.UpdateRating(tx, tenantID, id, rating, reviewCount)
}
