package category

import (
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *CreateCategoryRequest, userID string) (CategoryResponse, error)
	GetByID(ctx context.Context, tenantID uuid.UUID, id string) (CategoryResponse, error)
	List(ctx context.Context, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[CategoryResponse], error)
	Update(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateCategoryRequest, userID string) (CategoryResponse, error)
	Delete(ctx context.Context, tenantID uuid.UUID, id string) error

	// Cross-service
	GetByIDRaw(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Category, error)
}

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req *CreateCategoryRequest, userID string) (CategoryResponse, error) {
	var cat Category

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		slug := req.Slug
		if slug == "" {
			slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
		}

		if _, err := s.repo.GetBySlug(tx, tenantID, slug); err == nil {
			return apperror.BadRequest("slug kategori sudah digunakan", nil)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Internal("gagal memeriksa slug kategori", err)
		}

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		cat = Category{
			TenantID:  tenantID,
			ParentID:  req.ParentID,
			Name:      req.Name,
			Slug:      slug,
			Description: req.Description,
			Image:     req.Image,
			IsActive:  isActive,
			SortOrder: req.SortOrder,
			CreatedBy: userID,
		}

		if err := s.repo.Create(tx, &cat); err != nil {
			return apperror.Internal("gagal membuat kategori", err)
		}
		return nil
	})

	if err != nil {
		return CategoryResponse{}, err
	}

	return cat.ToResponse(), nil
}

func (s *Service) GetByID(ctx context.Context, tenantID uuid.UUID, id string) (CategoryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return CategoryResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	cat, err := s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CategoryResponse{}, apperror.NotFound("kategori tidak ditemukan", err)
		}
		return CategoryResponse{}, apperror.Internal("gagal mengambil kategori", err)
	}

	return cat.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[CategoryResponse], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), tenantID, search, page, limit)
	if err != nil {
		return pagination.Response[CategoryResponse]{}, apperror.Internal("gagal mengambil daftar kategori", err)
	}

	responses := make([]CategoryResponse, len(result.Data))
	for i, c := range result.Data {
		responses[i] = c.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) Update(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateCategoryRequest, userID string) (CategoryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return CategoryResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var cat Category

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cat, err = s.repo.GetByID(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("kategori tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil kategori", err)
		}

		if req.ParentID != nil {
			cat.ParentID = req.ParentID
		}
		if req.Name != nil {
			cat.Name = *req.Name
		}
		if req.Slug != nil {
			cat.Slug = *req.Slug
		}
		if req.Description != nil {
			cat.Description = *req.Description
		}
		if req.Image != nil {
			cat.Image = *req.Image
		}
		if req.SortOrder != nil {
			cat.SortOrder = *req.SortOrder
		}
		if req.IsActive != nil {
			cat.IsActive = *req.IsActive
		}
		cat.UpdatedBy = userID

		if err := s.repo.Update(tx, &cat); err != nil {
			return apperror.Internal("gagal memperbarui kategori", err)
		}
		return nil
	})

	if txErr != nil {
		return CategoryResponse{}, txErr
	}

	return cat.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("kategori tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil kategori", err)
	}

	if err := s.repo.Delete(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus kategori", err)
	}

	return nil
}

func (s *Service) GetByIDRaw(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Category, error) {
	return s.repo.GetByID(tx, tenantID, id)
}
