package category

import (
	"backend-go/pkg/pagination"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(tx *gorm.DB, c *Category) error
	GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Category, error)
	GetBySlug(tx *gorm.DB, tenantID uuid.UUID, slug string) (Category, error)
	List(tx *gorm.DB, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[Category], error)
	Update(tx *gorm.DB, c *Category) error
	Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
}

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) Create(tx *gorm.DB, c *Category) error {
	return tx.Create(c).Error
}

func (r *Repository) GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Category, error) {
	var c Category
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error
	if err != nil {
		return Category{}, err
	}
	return c, nil
}

func (r *Repository) GetBySlug(tx *gorm.DB, tenantID uuid.UUID, slug string) (Category, error) {
	var c Category
	err := tx.Where("slug = ? AND tenant_id = ?", slug, tenantID).First(&c).Error
	if err != nil {
		return Category{}, err
	}
	return c, nil
}

func (r *Repository) List(tx *gorm.DB, tenantID uuid.UUID, search string, page, limit int) (pagination.Response[Category], error) {
	p := pagination.Resolve(page, limit)

	query := tx.Model(&Category{}).Where("tenant_id = ?", tenantID)
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[Category]{}, err
	}

	var categories []Category
	if err := query.Order("sort_order ASC, name ASC").Offset(p.Offset).Limit(p.Limit).Find(&categories).Error; err != nil {
		return pagination.Response[Category]{}, err
	}

	return pagination.NewResponse(categories, int(total), p.Page, p.Limit), nil
}

func (r *Repository) Update(tx *gorm.DB, c *Category) error {
	return tx.Save(c).Error
}

func (r *Repository) Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&Category{}).Error
}
