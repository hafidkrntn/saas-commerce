package tenant

import (
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	Create(ctx context.Context, req *CreateTenantRequest, userID string) (TenantResponse, error)
	GetByID(ctx context.Context, id string) (TenantResponse, error)
	List(ctx context.Context, search string, page, limit int) (pagination.Response[TenantResponse], error)
	Update(ctx context.Context, id string, req *UpdateTenantRequest, userID string) (TenantResponse, error)
	Delete(ctx context.Context, id string) error
	ListPlans(ctx context.Context) ([]PlanResponse, error)
}

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateTenantRequest, userID string) (TenantResponse, error) {
	var tenant Tenant

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.GetBySlug(tx, req.Slug); err == nil {
			return apperror.BadRequest("slug tenant sudah digunakan", nil)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Internal("gagal memeriksa slug", err)
		}

		status := req.SubscriptionStatus
		if status == "" {
			status = "trial"
		}

		tenant = Tenant{
			Name:               req.Name,
			Slug:               req.Slug,
			Logo:               req.Logo,
			PlanID:             req.PlanID,
			SubscriptionStatus: status,
			TrialEndsAt:        req.TrialEndsAt,
			CreatedBy:          userID,
		}

		if err := s.repo.Create(tx, &tenant); err != nil {
			return apperror.Internal("gagal membuat tenant", err)
		}
		return nil
	})

	if err != nil {
		return TenantResponse{}, err
	}

	return tenant.ToResponse(), nil
}

func (s *Service) GetByID(ctx context.Context, id string) (TenantResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return TenantResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	tenant, err := s.repo.GetByID(s.db.WithContext(ctx), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TenantResponse{}, apperror.NotFound("tenant tidak ditemukan", err)
		}
		return TenantResponse{}, apperror.Internal("gagal mengambil tenant", err)
	}

	return tenant.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, search string, page, limit int) (pagination.Response[TenantResponse], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), search, page, limit)
	if err != nil {
		return pagination.Response[TenantResponse]{}, apperror.Internal("gagal mengambil daftar tenant", err)
	}

	responses := make([]TenantResponse, len(result.Data))
	for i, t := range result.Data {
		responses[i] = t.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) Update(ctx context.Context, id string, req *UpdateTenantRequest, userID string) (TenantResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return TenantResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var tenant Tenant

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tenant, err = s.repo.GetByID(tx, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("tenant tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil tenant", err)
		}

		if req.Name != nil {
			tenant.Name = *req.Name
		}
		if req.Slug != nil {
			tenant.Slug = *req.Slug
		}
		if req.Logo != nil {
			tenant.Logo = *req.Logo
		}
		if req.PlanID != nil {
			tenant.PlanID = req.PlanID
		}
		if req.SubscriptionStatus != nil {
			tenant.SubscriptionStatus = *req.SubscriptionStatus
		}
		if req.TrialEndsAt != nil {
			tenant.TrialEndsAt = req.TrialEndsAt
		}
		tenant.UpdatedBy = userID

		if err := s.repo.Update(tx, &tenant); err != nil {
			return apperror.Internal("gagal memperbarui tenant", err)
		}
		return nil
	})

	if txErr != nil {
		return TenantResponse{}, txErr
	}

	return tenant.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetByID(s.db.WithContext(ctx), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("tenant tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil tenant", err)
	}

	if err := s.repo.Delete(s.db.WithContext(ctx), uid); err != nil {
		return apperror.Internal("gagal menghapus tenant", err)
	}

	return nil
}

func (s *Service) ListPlans(ctx context.Context) ([]PlanResponse, error) {
	plans, err := s.repo.ListPlans(s.db.WithContext(ctx))
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar plan", err)
	}

	responses := make([]PlanResponse, len(plans))
	for i, p := range plans {
		responses[i] = p.ToResponse()
	}
	return responses, nil
}
