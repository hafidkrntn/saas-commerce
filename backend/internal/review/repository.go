package review

import (
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) Create(tx *gorm.DB, review *Review) error {
	return tx.Create(review).Error
}

func (r *Repository) GetByID(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) (Review, error) {
	var review Review
	err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).First(&review).Error
	if err != nil {
		return Review{}, err
	}
	return review, nil
}

func (r *Repository) List(tx *gorm.DB, tenantID uuid.UUID, productID string, onlyPublished bool, page, limit int) (pagination.Response[Review], error) {
	p := pagination.Resolve(page, limit)

	query := tx.Model(&Review{}).Where("tenant_id = ?", tenantID)
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if onlyPublished {
		query = query.Where("is_published = ?", true)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[Review]{}, err
	}

	var reviews []Review
	if err := query.Order("created_at DESC").Offset(p.Offset).Limit(p.Limit).Find(&reviews).Error; err != nil {
		return pagination.Response[Review]{}, err
	}

	return pagination.NewResponse(reviews, int(total), p.Page, p.Limit), nil
}

func (r *Repository) Update(tx *gorm.DB, review *Review) error {
	return tx.Save(review).Error
}

func (r *Repository) Delete(tx *gorm.DB, tenantID uuid.UUID, id uuid.UUID) error {
	return tx.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&Review{}).Error
}

func (r *Repository) AggregateRating(tx *gorm.DB, tenantID uuid.UUID, productID uuid.UUID) (float64, int64, error) {
	var avg float64
	var count int64

	err := tx.Model(&Review{}).
		Where("tenant_id = ? AND product_id = ? AND is_published = ?", tenantID, productID, true).
		Select("COALESCE(AVG(rating), 0), COUNT(*)").Row().Scan(&avg, &count)
	if err != nil {
		return 0, 0, err
	}

	return avg, count, nil
}
