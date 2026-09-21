package user

import (
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateUser(tx *gorm.DB, u *User) error
	GetUser(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (User, error)
	GetUserByEmail(tx *gorm.DB, tenantID uuid.UUID, email string) (User, error)
	ListUsers(tx *gorm.DB, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[User], error)
	UpdateUser(tx *gorm.DB, u *User) error
	DeleteUser(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	GetUserRoleName(tx *gorm.DB, userID uuid.UUID) (string, error)
	AssignRole(tx *gorm.DB, userID uuid.UUID, roleID uuid.UUID) error
	UnassignRole(tx *gorm.DB, userID uuid.UUID, roleID uuid.UUID) error
	DeleteUserRoles(tx *gorm.DB, userID uuid.UUID) error

	CreateRole(tx *gorm.DB, r *Role) error
	GetRole(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Role, error)
	ListRoles(tx *gorm.DB, tenantID uuid.UUID) ([]Role, error)
	UpdateRole(tx *gorm.DB, r *Role) error
	DeleteRole(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	CountRoleUsers(tx *gorm.DB, roleID uuid.UUID) (int64, error)
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

// --- Users ---

func (r *Repository) CreateUser(tx *gorm.DB, u *User) error {
	return tx.Create(u).Error
}

func (r *Repository) GetUser(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (User, error) {
	var u User
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *Repository) GetUserByEmail(tx *gorm.DB, tenantID uuid.UUID, email string) (User, error) {
	var u User
	err := tx.Where("email = ? AND tenant_id = ?", email, tenantID).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *Repository) ListUsers(tx *gorm.DB, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[User], error) {
	p := pagination.Resolve(page, limit)

	query := tx.Model(&User{}).Where("tenant_id = ?", tenantID)
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[User]{}, err
	}

	var users []User
	if err := query.Order("created_at DESC").Offset(p.Offset).Limit(p.Limit).Find(&users).Error; err != nil {
		return pagination.Response[User]{}, err
	}

	return pagination.NewResponse(users, int(total), p.Page, p.Limit), nil
}

func (r *Repository) UpdateUser(tx *gorm.DB, u *User) error {
	return tx.Save(u).Error
}

func (r *Repository) DeleteUser(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&User{}).Error
}

func (r *Repository) GetUserRoleName(tx *gorm.DB, userID uuid.UUID) (string, error) {
	var name string
	err := tx.Table("roles").
		Select("roles.name").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Order("roles.permissions DESC").
		Limit(1).
		Pluck("roles.name", &name).Error
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *Repository) AssignRole(tx *gorm.DB, userID uuid.UUID, roleID uuid.UUID) error {
	return tx.Create(&UserRole{UserID: userID, RoleID: roleID}).Error
}

func (r *Repository) UnassignRole(tx *gorm.DB, userID uuid.UUID, roleID uuid.UUID) error {
	return tx.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
}

func (r *Repository) DeleteUserRoles(tx *gorm.DB, userID uuid.UUID) error {
	return tx.Where("user_id = ?", userID).Delete(&UserRole{}).Error
}

// --- Roles ---

func (r *Repository) CreateRole(tx *gorm.DB, role *Role) error {
	return tx.Create(role).Error
}

func (r *Repository) GetRole(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Role, error) {
	var role Role
	err := tx.Where("id = ? AND (tenant_id = ? OR tenant_id IS NULL)", id, tenantID).First(&role).Error
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

func (r *Repository) ListRoles(tx *gorm.DB, tenantID uuid.UUID) ([]Role, error) {
	var roles []Role
	err := tx.Where("tenant_id = ? OR tenant_id IS NULL", tenantID).Order("permissions DESC").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) UpdateRole(tx *gorm.DB, role *Role) error {
	return tx.Save(role).Error
}

func (r *Repository) DeleteRole(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&Role{}).Error
}

func (r *Repository) CountRoleUsers(tx *gorm.DB, roleID uuid.UUID) (int64, error) {
	var count int64
	err := tx.Model(&UserRole{}).Where("role_id = ?", roleID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
