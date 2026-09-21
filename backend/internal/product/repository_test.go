package product_test

import (
	"backend-go/internal/entities"
	"backend-go/internal/product"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository product.RepositoryInterface
}

func (s *Suite) SetupSuite() {
	db, sqlMock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.mock = sqlMock
	s.repository = product.NewRepository()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestRepository_Create() {
	tenantID := uuid.New()
	id := uuid.New()
	now := time.Now()

	newProduct := &product.Product{
		TenantID:          tenantID,
		SKU:               "SKU-001",
		Name:              "Laptop Pro",
		Price:             15000000,
		Stock:             10,
		LowStockThreshold: 5,
		Status:            product.StatusActive,
		Tags:              entities.StringList{"laptop", "gaming"},
		CreatedBy:         "user-1",
		CreatedAt:         now,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectCommit()

	err := s.repository.Create(s.DB, newProduct)
	require.NoError(s.T(), err)
	require.Equal(s.T(), id, newProduct.ID)
}

func (s *Suite) TestRepository_GetByID_ScopedToTenant() {
	tenantID := uuid.New()
	expected := product.Product{
		ID:           uuid.New(),
		TenantID:     tenantID,
		SKU:          "SKU-001",
		Name:         "Laptop Pro",
		Price:        15000000,
		Stock:        10,
		Status:       product.StatusActive,
		CreatedBy:    "user-1",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	productRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "sku", "name", "price", "stock", "status",
		"created_by", "created_at", "updated_at",
	}).
		AddRow(expected.ID, expected.TenantID, expected.SKU, expected.Name,
			expected.Price, expected.Stock, expected.Status,
			expected.CreatedBy, expected.CreatedAt, expected.UpdatedAt)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND tenant_id = $2 ORDER BY "products"."id" LIMIT $3`)).
		WithArgs(expected.ID, tenantID, 1).
		WillReturnRows(productRows)

	// Preload images — empty
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product_images" WHERE "product_images"."product_id" = $1`)).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "alt", "sort_order", "is_primary", "created_at"}))

	result, err := s.repository.GetByID(s.DB, tenantID, expected.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected.ID, result.ID)
	require.Equal(s.T(), expected.SKU, result.SKU)
}

func (s *Suite) TestRepository_GetByID_CrossTenantIsolation() {
	tenantA := uuid.New()
	tenantB := uuid.New()
	id := uuid.New()

	// Product exists only for tenant A — tenant B must get not found
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND tenant_id = $2 ORDER BY "products"."id" LIMIT $3`)).
		WithArgs(id, tenantB, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := s.repository.GetByID(s.DB, tenantB, id)
	require.ErrorIs(s.T(), err, gorm.ErrRecordNotFound)
	_ = tenantA
}

func (s *Suite) TestRepository_List_WithSearchAndStatus() {
	tenantID := uuid.New()

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" WHERE tenant_id = $1 AND (name ILIKE $2 OR sku ILIKE $3) AND status = $4`)).
		WithArgs(tenantID, "%laptop%", "%laptop%", "active").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	id := uuid.New()
	productRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "sku", "name", "price", "stock", "status",
	}).
		AddRow(id, tenantID, "SKU-001", "Laptop Pro", 15000000, 10, "active")

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE tenant_id = $1 AND (name ILIKE $2 OR sku ILIKE $3) AND status = $4 ORDER BY created_at DESC LIMIT $5`)).
		WithArgs(tenantID, "%laptop%", "%laptop%", "active", 10).
		WillReturnRows(productRows)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product_images" WHERE "product_images"."product_id" = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "alt", "sort_order", "is_primary", "created_at"}))

	result, err := s.repository.List(s.DB, tenantID, &entities.ListParams{
		Page:   1,
		Limit:  10,
		Search: "laptop",
		Status: "active",
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, result.Pagination.Total)
	require.Len(s.T(), result.Data, 1)
	require.Equal(s.T(), "Laptop Pro", result.Data[0].Name)
}

func (s *Suite) TestRepository_Delete_ScopedToTenant() {
	tenantID := uuid.New()
	id := uuid.New()

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "products" WHERE id = $1 AND tenant_id = $2`)).
		WithArgs(id, tenantID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.Delete(s.DB, tenantID, id)
	require.NoError(s.T(), err)
}
