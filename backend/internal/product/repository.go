package product

import (
	"backend-go/internal/entities"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RepositoryInterface defines all data-access methods for this domain.
type RepositoryInterface interface {
	Create(tx *gorm.DB, product *Product) error
	GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Product, error)
	GetBySKU(tx *gorm.DB, tenantID uuid.UUID, sku string) (Product, error)
	List(tx *gorm.DB, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[Product], error)
	Update(tx *gorm.DB, product *Product) error
	Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error
	ReplaceImages(tx *gorm.DB, productID uuid.UUID, urls []string) error
	UpdateRating(tx *gorm.DB, tenantID uuid.UUID, productID uuid.UUID, rating float64, reviewCount int) error
}

// Repository implements RepositoryInterface.
type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

// --- CRUD Operations ---

func (r *Repository) Create(tx *gorm.DB, product *Product) error {
	return tx.Create(product).Error
}

func (r *Repository) GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Product, error) {
	var product Product
	err := tx.
		Preload("Images").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&product).Error
	if err != nil {
		return Product{}, err
	}

	r.enrich(tx, &product)
	return product, nil
}

// enrich fills transient fields (category name) via join.
func (r *Repository) enrich(tx *gorm.DB, product *Product) {
	if product.CategoryID != nil {
		tx.Table("categories").Select("name").
			Where("id = ?", *product.CategoryID).Pluck("name", &product.CategoryName)
	}
}

func (r *Repository) GetBySKU(tx *gorm.DB, tenantID uuid.UUID, sku string) (Product, error) {
	var product Product
	err := tx.
		Where("sku = ? AND tenant_id = ?", sku, tenantID).
		First(&product).Error
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (r *Repository) List(tx *gorm.DB, tenantID uuid.UUID, params *entities.ListParams) (pagination.Response[Product], error) {
	p := pagination.Resolve(params.Page, params.Limit)

	query := tx.Model(&Product{}).Where("tenant_id = ?", tenantID)

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR sku ILIKE ?", search, search)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.CategoryID != "" {
		if categoryID, err := uuid.Parse(params.CategoryID); err == nil {
			query = query.Where("category_id = ?", categoryID)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[Product]{}, err
	}

	var products []Product
	if err := query.
		Preload("Images").
		Order("created_at DESC").
		Offset(p.Offset).
		Limit(p.Limit).
		Find(&products).Error; err != nil {
		return pagination.Response[Product]{}, err
	}

	for i := range products {
		r.enrich(tx, &products[i])
	}

	return pagination.NewResponse(products, int(total), p.Page, p.Limit), nil
}

func (r *Repository) Update(tx *gorm.DB, product *Product) error {
	return tx.Save(product).Error
}

func (r *Repository) Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&Product{}).Error
}

// ReplaceImages deletes existing images and inserts the new list (ordered).
func (r *Repository) ReplaceImages(tx *gorm.DB, productID uuid.UUID, urls []string) error {
	if err := tx.Where("product_id = ?", productID).Delete(&ProductImage{}).Error; err != nil {
		return err
	}

	for i, url := range urls {
		if err := tx.Create(&ProductImage{
			ProductID: productID,
			URL:       url,
			SortOrder: i,
			IsPrimary: i == 0,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateRating recalculates the product's aggregate rating/review count.
func (r *Repository) UpdateRating(tx *gorm.DB, tenantID uuid.UUID, productID uuid.UUID, rating float64, reviewCount int) error {
	return tx.Model(&Product{}).
		Where("id = ? AND tenant_id = ?", productID, tenantID).
		Updates(map[string]any{"rating": rating, "review_count": reviewCount}).Error
}
