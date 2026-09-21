package review

import (
	"backend-go/internal/product"
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *CreateReviewRequest) (ReviewResponse, error)
	List(ctx context.Context, tenantID uuid.UUID, productID string, onlyPublished bool, page, limit int) (pagination.Response[ReviewResponse], error)
	Moderate(ctx context.Context, tenantID uuid.UUID, id string, req *ModerateReviewRequest) (ReviewResponse, error)
	Delete(ctx context.Context, tenantID uuid.UUID, id string) error
}

type Service struct {
	db         *gorm.DB
	repo       RepositoryInterface
	productSvc product.ServiceInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface, productSvc product.ServiceInterface) ServiceInterface {
	return &Service{db: db, repo: repo, productSvc: productSvc}
}

// recalcRating refreshes the product's rating aggregates after a change.
func (s *Service) recalcRating(tx *gorm.DB, tenantID uuid.UUID, productID uuid.UUID) error {
	avg, count, err := s.repo.AggregateRating(tx, tenantID, productID)
	if err != nil {
		return err
	}
	avg = math.Round(avg*100) / 100
	return s.productSvc.UpdateRating(tx, tenantID, productID, avg, int(count))
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req *CreateReviewRequest) (ReviewResponse, error) {
	var review Review

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.productSvc.GetByIDRaw(tx, tenantID, req.ProductID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("produk tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil produk", err)
		}

		review = Review{
			TenantID:       tenantID,
			ProductID:      req.ProductID,
			CustomerID:     req.CustomerID,
			CustomerName:   req.CustomerName,
			Rating:         req.Rating,
			Title:          req.Title,
			Content:        req.Content,
			IsVerified:     req.CustomerID != nil,
			IsPublished:    req.CustomerID != nil,
		}

		if err := s.repo.Create(tx, &review); err != nil {
			return apperror.Internal("gagal membuat review", err)
		}

		if review.IsPublished {
			if err := s.recalcRating(tx, tenantID, req.ProductID); err != nil {
				return apperror.Internal("gagal memperbarui rating produk", err)
			}
		}
		return nil
	})

	if txErr != nil {
		return ReviewResponse{}, txErr
	}

	return review.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, productID string, onlyPublished bool, page, limit int) (pagination.Response[ReviewResponse], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), tenantID, productID, onlyPublished, page, limit)
	if err != nil {
		return pagination.Response[ReviewResponse]{}, apperror.Internal("gagal mengambil daftar review", err)
	}

	responses := make([]ReviewResponse, len(result.Data))
	for i, r := range result.Data {
		responses[i] = r.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) Moderate(ctx context.Context, tenantID uuid.UUID, id string, req *ModerateReviewRequest) (ReviewResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ReviewResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var review Review

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		review, err = s.repo.GetByID(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("review tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil review", err)
		}

		if req.IsPublished != nil {
			review.IsPublished = *req.IsPublished
		}
		if req.IsVerified != nil {
			review.IsVerified = *req.IsVerified
		}

		if err := s.repo.Update(tx, &review); err != nil {
			return apperror.Internal("gagal memperbarui review", err)
		}

		if err := s.recalcRating(tx, tenantID, review.ProductID); err != nil {
			return apperror.Internal("gagal memperbarui rating produk", err)
		}
		return nil
	})

	if txErr != nil {
		return ReviewResponse{}, txErr
	}

	return review.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	review, err := s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("review tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil review", err)
	}

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Delete(tx, tenantID, uid); err != nil {
			return apperror.Internal("gagal menghapus review", err)
		}

		if review.IsPublished {
			if err := s.recalcRating(tx, tenantID, review.ProductID); err != nil {
				return apperror.Internal("gagal memperbarui rating produk", err)
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}
