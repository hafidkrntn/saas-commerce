package customer

import (
	"backend-go/internal/entities"
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *CreateCustomerRequest, userID string) (CustomerResponse, error)
	GetByID(ctx context.Context, tenantID uuid.UUID, id string) (CustomerResponse, error)
	List(ctx context.Context, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[CustomerResponse], error)
	Update(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateCustomerRequest, userID string) (CustomerResponse, error)
	Delete(ctx context.Context, tenantID uuid.UUID, id string) error
	AddAddress(ctx context.Context, tenantID uuid.UUID, customerID string, req *AddressRequest) (AddressResponse, error)
	GetActivities(ctx context.Context, tenantID uuid.UUID, customerID string, limit int) ([]ActivityResponse, error)

	// Cross-service
	FindOrCreateByEmail(tx *gorm.DB, tenantID uuid.UUID, name, email string) (Customer, error)
	RecordOrderStats(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, amount float64) error
	AddActivity(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, activityType, description string) error
}

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req *CreateCustomerRequest, userID string) (CustomerResponse, error) {
	var customer Customer

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if req.Email != "" {
			if existing, err := s.repo.GetByEmail(tx, tenantID, req.Email); err == nil && existing.ID != uuid.Nil {
				return apperror.BadRequest("email customer sudah terdaftar", nil)
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.Internal("gagal memeriksa email customer", err)
			}
		}

		status := req.Status
		if status == "" {
			status = StatusActive
		}

		customer = Customer{
			TenantID:  tenantID,
			Name:      req.Name,
			Email:     req.Email,
			Phone:     req.Phone,
			Avatar:    req.Avatar,
			Status:    status,
			Tags:      req.Tags,
			CreatedBy: userID,
		}

		if err := s.repo.Create(tx, &customer); err != nil {
			return apperror.Internal("gagal membuat customer", err)
		}

		if req.Address != nil {
			address := &CustomerAddress{
				TenantID:   tenantID,
				CustomerID: customer.ID,
				Label:      req.Address.Label,
				Line1:      req.Address.Line1,
				Line2:      req.Address.Line2,
				City:       req.Address.City,
				State:      req.Address.State,
				Zip:        req.Address.Zip,
				Country:    req.Address.Country,
				IsDefault:  req.Address.IsDefault,
			}
			if err := s.repo.AddAddress(tx, address); err != nil {
				return apperror.Internal("gagal menyimpan alamat customer", err)
			}
		}

		return nil
	})

	if err != nil {
		return CustomerResponse{}, err
	}

	return customer.ToResponse(), nil
}

func (s *Service) GetByID(ctx context.Context, tenantID uuid.UUID, id string) (CustomerResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return CustomerResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	customer, err := s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CustomerResponse{}, apperror.NotFound("customer tidak ditemukan", err)
		}
		return CustomerResponse{}, apperror.Internal("gagal mengambil customer", err)
	}

	return customer.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[CustomerResponse], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), tenantID, params)
	if err != nil {
		return pagination.Response[CustomerResponse]{}, apperror.Internal("gagal mengambil daftar customer", err)
	}

	responses := make([]CustomerResponse, len(result.Data))
	for i, c := range result.Data {
		responses[i] = c.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) Update(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateCustomerRequest, userID string) (CustomerResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return CustomerResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var customer Customer

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		customer, err = s.repo.GetByID(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("customer tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil customer", err)
		}

		if req.Name != nil {
			customer.Name = *req.Name
		}
		if req.Email != nil {
			customer.Email = *req.Email
		}
		if req.Phone != nil {
			customer.Phone = *req.Phone
		}
		if req.Avatar != nil {
			customer.Avatar = *req.Avatar
		}
		if req.Status != nil {
			customer.Status = *req.Status
		}
		if req.Tags != nil {
			customer.Tags = req.Tags
		}
		customer.UpdatedBy = userID

		if err := s.repo.Update(tx, &customer); err != nil {
			return apperror.Internal("gagal memperbarui customer", err)
		}
		return nil
	})

	if txErr != nil {
		return CustomerResponse{}, txErr
	}

	return customer.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetByID(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("customer tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil customer", err)
	}

	if err := s.repo.Delete(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus customer", err)
	}

	return nil
}

func (s *Service) AddAddress(ctx context.Context, tenantID uuid.UUID, customerID string, req *AddressRequest) (AddressResponse, error) {
	cid, err := uuid.Parse(customerID)
	if err != nil {
		return AddressResponse{}, apperror.BadRequest("id customer tidak valid", err)
	}

	var address CustomerAddress

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.GetByID(tx, tenantID, cid); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("customer tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil customer", err)
		}

		country := req.Country
		if country == "" {
			country = "Indonesia"
		}

		address = CustomerAddress{
			TenantID:   tenantID,
			CustomerID: cid,
			Label:      req.Label,
			Line1:      req.Line1,
			Line2:      req.Line2,
			City:       req.City,
			State:      req.State,
			Zip:        req.Zip,
			Country:    country,
			IsDefault:  req.IsDefault,
		}

		if err := s.repo.AddAddress(tx, &address); err != nil {
			return apperror.Internal("gagal menyimpan alamat", err)
		}
		return nil
	})

	if txErr != nil {
		return AddressResponse{}, txErr
	}

	return AddressResponse{
		ID:        address.ID,
		Label:     address.Label,
		Line1:     address.Line1,
		Line2:     address.Line2,
		City:      address.City,
		State:     address.State,
		Zip:       address.Zip,
		Country:   address.Country,
		IsDefault: address.IsDefault,
	}, nil
}

func (s *Service) GetActivities(ctx context.Context, tenantID uuid.UUID, customerID string, limit int) ([]ActivityResponse, error) {
	cid, err := uuid.Parse(customerID)
	if err != nil {
		return nil, apperror.BadRequest("id customer tidak valid", err)
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	activities, err := s.repo.ListActivities(s.db.WithContext(ctx), tenantID, cid, limit)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil aktivitas customer", err)
	}

	responses := make([]ActivityResponse, len(activities))
	for i, a := range activities {
		responses[i] = ActivityResponse{
			ID:          a.ID,
			Type:        a.Type,
			Description: a.Description,
			Timestamp:   a.CreatedAt,
		}
	}

	return responses, nil
}

// =============================================================================
// Cross-service methods
// =============================================================================

func (s *Service) FindOrCreateByEmail(tx *gorm.DB, tenantID uuid.UUID, name, email string) (Customer, error) {
	if email != "" {
		customer, err := s.repo.GetByEmail(tx, tenantID, email)
		if err == nil {
			return customer, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return Customer{}, err
		}
	}

	customer := Customer{
		TenantID: tenantID,
		Name:     name,
		Email:    email,
		Status:   StatusActive,
	}
	if err := s.repo.Create(tx, &customer); err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (s *Service) RecordOrderStats(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, amount float64) error {
	return s.repo.RecordOrderStats(tx, tenantID, customerID, amount, time.Now())
}

func (s *Service) AddActivity(tx *gorm.DB, tenantID uuid.UUID, customerID uuid.UUID, activityType, description string) error {
	return s.repo.AddActivity(tx, &CustomerActivity{
		TenantID:    tenantID,
		CustomerID:  customerID,
		Type:        activityType,
		Description: description,
	})
}
