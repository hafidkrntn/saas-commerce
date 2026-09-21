package user

import (
	"backend-go/pkg/apperror"
	"backend-go/pkg/pagination"
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	CreateUser(ctx context.Context, tenantID uuid.UUID, req *CreateUserRequest, userID string) (UserResponse, error)
	GetUser(ctx context.Context, tenantID uuid.UUID, id string) (UserResponse, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[UserResponse], error)
	UpdateUser(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateUserRequest, userID string) (UserResponse, error)
	DeleteUser(ctx context.Context, tenantID uuid.UUID, id string) error

	CreateRole(ctx context.Context, tenantID uuid.UUID, req *CreateRoleRequest, userID string) (RoleResponse, error)
	GetRole(ctx context.Context, tenantID uuid.UUID, id string) (RoleResponse, error)
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]RoleResponse, error)
	UpdateRole(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateRoleRequest, userID string) (RoleResponse, error)
	DeleteRole(ctx context.Context, tenantID uuid.UUID, id string) error
}

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

// =============================================================================
// Users
// =============================================================================

func (s *Service) CreateUser(ctx context.Context, tenantID uuid.UUID, req *CreateUserRequest, userID string) (UserResponse, error) {
	var user User

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.GetUserByEmail(tx, tenantID, req.Email); err == nil {
			return apperror.BadRequest("email sudah digunakan", nil)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Internal("gagal memeriksa email", err)
		}

		roleID, err := uuid.Parse(req.RoleID)
		if err != nil {
			return apperror.BadRequest("role_id tidak valid", err)
		}

		role, err := s.repo.GetRole(tx, tenantID, roleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("role tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil role", err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return apperror.Internal("gagal mengenkripsi password", err)
		}

		user = User{
			TenantID:     &tenantID,
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: string(hash),
			Status:       StatusActive,
			CreatedBy:    userID,
		}

		if err := s.repo.CreateUser(tx, &user); err != nil {
			return apperror.Internal("gagal membuat user", err)
		}

		if err := s.repo.AssignRole(tx, user.ID, role.ID); err != nil {
			return apperror.Internal("gagal menetapkan role", err)
		}

		return nil
	})

	if txErr != nil {
		return UserResponse{}, txErr
	}

	return s.toUserResponse(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, tenantID uuid.UUID, id string) (UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UserResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	user, err := s.repo.GetUser(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserResponse{}, apperror.NotFound("user tidak ditemukan", err)
		}
		return UserResponse{}, apperror.Internal("gagal mengambil user", err)
	}

	return s.toUserResponse(ctx, user)
}

func (s *Service) ListUsers(ctx context.Context, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[UserResponse], error) {
	result, err := s.repo.ListUsers(s.db.WithContext(ctx), tenantID, search, page, limit)
	if err != nil {
		return pagination.Response[UserResponse]{}, apperror.Internal("gagal mengambil daftar user", err)
	}

	responses := make([]UserResponse, len(result.Data))
	for i, u := range result.Data {
		responses[i], err = s.toUserResponse(ctx, u)
		if err != nil {
			return pagination.Response[UserResponse]{}, err
		}
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) UpdateUser(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateUserRequest, userID string) (UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UserResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var user User

	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, err = s.repo.GetUser(tx, tenantID, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("user tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil user", err)
		}

		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.Status != nil {
			user.Status = *req.Status
		}
		if req.Password != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return apperror.Internal("gagal mengenkripsi password", err)
			}
			user.PasswordHash = string(hash)
		}
		user.UpdatedBy = userID

		if err := s.repo.UpdateUser(tx, &user); err != nil {
			return apperror.Internal("gagal memperbarui user", err)
		}
		return nil
	})

	if txErr != nil {
		return UserResponse{}, txErr
	}

	return s.toUserResponse(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	_, err = s.repo.GetUser(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil user", err)
	}

	if err := s.repo.DeleteUser(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus user", err)
	}

	return nil
}

func (s *Service) toUserResponse(ctx context.Context, user User) (UserResponse, error) {
	role, err := s.repo.GetUserRoleName(s.db.WithContext(ctx), user.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return UserResponse{}, apperror.Internal("gagal mengambil role user", err)
	}

	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Role:      role,
		Status:    user.Status,
		LastLogin: user.LastLoginAt,
		CreatedAt: user.CreatedAt,
	}, nil
}

// =============================================================================
// Roles
// =============================================================================

func (s *Service) CreateRole(ctx context.Context, tenantID uuid.UUID, req *CreateRoleRequest, userID string) (RoleResponse, error) {
	role := Role{
		TenantID:    &tenantID,
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		CreatedBy:   userID,
	}

	if err := s.repo.CreateRole(s.db.WithContext(ctx), &role); err != nil {
		return RoleResponse{}, apperror.Internal("gagal membuat role", err)
	}

	return s.toRoleResponse(ctx, role)
}

func (s *Service) GetRole(ctx context.Context, tenantID uuid.UUID, id string) (RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return RoleResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	role, err := s.repo.GetRole(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoleResponse{}, apperror.NotFound("role tidak ditemukan", err)
		}
		return RoleResponse{}, apperror.Internal("gagal mengambil role", err)
	}

	return s.toRoleResponse(ctx, role)
}

func (s *Service) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]RoleResponse, error) {
	roles, err := s.repo.ListRoles(s.db.WithContext(ctx), tenantID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil daftar role", err)
	}

	responses := make([]RoleResponse, len(roles))
	for i, r := range roles {
		responses[i], err = s.toRoleResponse(ctx, r)
		if err != nil {
			return nil, err
		}
	}
	return responses, nil
}

func (s *Service) UpdateRole(ctx context.Context, tenantID uuid.UUID, id string, req *UpdateRoleRequest, userID string) (RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return RoleResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	role, err := s.repo.GetRole(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoleResponse{}, apperror.NotFound("role tidak ditemukan", err)
		}
		return RoleResponse{}, apperror.Internal("gagal mengambil role", err)
	}

	if role.IsSystem {
		return RoleResponse{}, apperror.Forbidden("role sistem tidak dapat diubah", nil)
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Permissions != nil {
		role.Permissions = *req.Permissions
	}
	role.UpdatedBy = userID

	if err := s.repo.UpdateRole(s.db.WithContext(ctx), &role); err != nil {
		return RoleResponse{}, apperror.Internal("gagal memperbarui role", err)
	}

	return s.toRoleResponse(ctx, role)
}

func (s *Service) DeleteRole(ctx context.Context, tenantID uuid.UUID, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	role, err := s.repo.GetRole(s.db.WithContext(ctx), tenantID, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("role tidak ditemukan", err)
		}
		return apperror.Internal("gagal mengambil role", err)
	}

	if role.IsSystem {
		return apperror.Forbidden("role sistem tidak dapat dihapus", nil)
	}

	count, err := s.repo.CountRoleUsers(s.db.WithContext(ctx), uid)
	if err != nil {
		return apperror.Internal("gagal memeriksa role", err)
	}
	if count > 0 {
		return apperror.BadRequest("role masih digunakan oleh user", nil)
	}

	if err := s.repo.DeleteRole(s.db.WithContext(ctx), tenantID, uid); err != nil {
		return apperror.Internal("gagal menghapus role", err)
	}

	return nil
}

func (s *Service) toRoleResponse(ctx context.Context, role Role) (RoleResponse, error) {
	users, err := s.repo.CountRoleUsers(s.db.WithContext(ctx), role.ID)
	if err != nil {
		return RoleResponse{}, apperror.Internal("gagal menghitung user role", err)
	}

	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
		Users:       users,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
	}, nil
}
