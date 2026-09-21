package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	GetUserByEmail(tx *gorm.DB, email string) (UserRow, error)
	GetUserByID(tx *gorm.DB, id uuid.UUID) (UserRow, error)
	GetUserRoles(tx *gorm.DB, userID uuid.UUID) ([]RoleRow, error)
	GetTenant(tx *gorm.DB, id uuid.UUID) (TenantBrief, error)
	UpdateLastLogin(tx *gorm.DB, userID uuid.UUID, at time.Time) error
	CreateRefreshToken(tx *gorm.DB, token *RefreshToken) error
	GetRefreshToken(tx *gorm.DB, tokenHash string) (RefreshToken, error)
	RevokeRefreshToken(tx *gorm.DB, id uuid.UUID) error
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) GetUserByEmail(tx *gorm.DB, email string) (UserRow, error) {
	var user UserRow
	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		return UserRow{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(tx *gorm.DB, id uuid.UUID) (UserRow, error) {
	var user UserRow
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		return UserRow{}, err
	}
	return user, nil
}

func (r *Repository) GetUserRoles(tx *gorm.DB, userID uuid.UUID) ([]RoleRow, error) {
	var roles []RoleRow
	err := tx.Table("roles").
		Select("roles.id, roles.tenant_id, roles.name, roles.permissions").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) GetTenant(tx *gorm.DB, id uuid.UUID) (TenantBrief, error) {
	var tenant TenantBrief
	err := tx.Table("tenants").
		Select("tenants.id, tenants.name, tenants.slug, COALESCE(plans.code, 'free') AS plan_code, tenants.subscription_status").
		Joins("LEFT JOIN plans ON plans.id = tenants.plan_id").
		Where("tenants.id = ?", id).
		Scan(&tenant).Error
	if err != nil {
		return TenantBrief{}, err
	}
	return tenant, nil
}

func (r *Repository) UpdateLastLogin(tx *gorm.DB, userID uuid.UUID, at time.Time) error {
	return tx.Model(&UserRow{}).Where("id = ?", userID).Update("last_login_at", at).Error
}

func (r *Repository) CreateRefreshToken(tx *gorm.DB, token *RefreshToken) error {
	return tx.Create(token).Error
}

func (r *Repository) GetRefreshToken(tx *gorm.DB, tokenHash string) (RefreshToken, error) {
	var token RefreshToken
	if err := tx.Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		return RefreshToken{}, err
	}
	return token, nil
}

func (r *Repository) RevokeRefreshToken(tx *gorm.DB, id uuid.UUID) error {
	now := time.Now()
	return tx.Model(&RefreshToken{}).Where("id = ?", id).Update("revoked_at", now).Error
}
