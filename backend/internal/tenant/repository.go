package tenant

import (
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(tx *gorm.DB, t *Tenant) error
	GetByID(tx *gorm.DB, id uuid.UUID) (Tenant, error)
	GetBySlug(tx *gorm.DB, slug string) (Tenant, error)
	List(tx *gorm.DB, search string, page, limit int) (pagination.Response[Tenant], error)
	Update(tx *gorm.DB, t *Tenant) error
	Delete(tx *gorm.DB, id uuid.UUID) error
	ListPlans(tx *gorm.DB) ([]Plan, error)
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) Create(tx *gorm.DB, t *Tenant) error {
	return tx.Create(t).Error
}

func (r *Repository) GetByID(tx *gorm.DB, id uuid.UUID) (Tenant, error) {
	var t Tenant
	if err := tx.Where("id = ?", id).First(&t).Error; err != nil {
		return Tenant{}, err
	}
	return t, nil
}

func (r *Repository) GetBySlug(tx *gorm.DB, slug string) (Tenant, error) {
	var t Tenant
	if err := tx.Where("slug = ?", slug).First(&t).Error; err != nil {
		return Tenant{}, err
	}
	return t, nil
}

func (r *Repository) List(tx *gorm.DB, search string, page, limit int) (pagination.Response[Tenant], error) {
	p := pagination.Resolve(page, limit)

	query := tx.Model(&Tenant{})
	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[Tenant]{}, err
	}

	var tenants []Tenant
	if err := query.Order("created_at DESC").Offset(p.Offset).Limit(p.Limit).Find(&tenants).Error; err != nil {
		return pagination.Response[Tenant]{}, err
	}

	return pagination.NewResponse(tenants, int(total), p.Page, p.Limit), nil
}

func (r *Repository) Update(tx *gorm.DB, t *Tenant) error {
	return tx.Save(t).Error
}

func (r *Repository) Delete(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("id = ?", id).Delete(&Tenant{}).Error
}

func (r *Repository) ListPlans(tx *gorm.DB) ([]Plan, error) {
	var plans []Plan
	if err := tx.Order("price ASC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}
