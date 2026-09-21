package order_test

import (
	"backend-go/internal/entities"
	"backend-go/internal/order"
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
	repository order.RepositoryInterface
}

func (s *Suite) SetupSuite() {
	db, sqlMock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.mock = sqlMock
	s.repository = order.NewRepository()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestRepository_Create() {
	now := time.Now()
	id := uuid.New()
	tenantID := uuid.New()

	newOrder := &order.Order{
		TenantID:      tenantID,
		OrderNumber:   "ORD/2608/00001",
		CustomerName:  "Customer A",
		Status:        order.StatusPending,
		PaymentStatus: order.PaymentPending,
		TotalAmount:   100.0,
		Notes:         "Test order",
		CreatedBy:     "user-1",
		CreatedAt:     now,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectCommit()

	err := s.repository.Create(s.DB, newOrder)
	require.NoError(s.T(), err)
	require.Equal(s.T(), id, newOrder.ID)
}

func (s *Suite) TestRepository_GetByID() {
	tenantID := uuid.New()
	expected := order.Order{
		ID:           uuid.New(),
		TenantID:     tenantID,
		OrderNumber:  "ORD/2608/00001",
		CustomerName: "Customer A",
		Status:       order.StatusPending,
		TotalAmount:  100.0,
		CreatedBy:    "user-1",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	orderRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "order_number", "customer_name", "status",
		"payment_status", "total_amount", "created_by", "created_at", "updated_at",
	}).
		AddRow(expected.ID, expected.TenantID, expected.OrderNumber, expected.CustomerName,
			expected.Status, expected.PaymentStatus, expected.TotalAmount,
			expected.CreatedBy, expected.CreatedAt, expected.UpdatedAt)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 AND tenant_id = $2 ORDER BY "orders"."id" LIMIT $3`)).
		WithArgs(expected.ID, tenantID, 1).
		WillReturnRows(orderRows)

	// Preload items — empty
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_items" WHERE "order_items"."order_id" = $1`)).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "unit_price", "subtotal"}))

	// Preload timeline — empty
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_timeline" WHERE "order_timeline"."order_id" = $1`)).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "type", "title", "description", "user_name", "created_at"}))

	result, err := s.repository.GetByID(s.DB, tenantID, expected.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected.ID, result.ID)
	require.Equal(s.T(), expected.OrderNumber, result.OrderNumber)
}

func (s *Suite) TestRepository_GetByID_NotFound() {
	tenantID := uuid.New()
	id := uuid.New()

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 AND tenant_id = $2 ORDER BY "orders"."id" LIMIT $3`)).
		WithArgs(id, tenantID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := s.repository.GetByID(s.DB, tenantID, id)
	require.ErrorIs(s.T(), err, gorm.ErrRecordNotFound)
}

func (s *Suite) TestRepository_List() {
	tenantID := uuid.New()

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "orders" WHERE tenant_id = $1`)).
		WithArgs(tenantID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	id := uuid.New()
	orderRows := sqlmock.NewRows([]string{
		"id", "tenant_id", "order_number", "customer_name", "status", "total_amount",
	}).
		AddRow(id, tenantID, "ORD/2608/00001", "Customer A", order.StatusPending, 50.0)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "orders" WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2`)).
		WithArgs(tenantID, 10).
		WillReturnRows(orderRows)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "order_items" WHERE "order_items"."order_id" = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "unit_price", "subtotal"}))

	result, err := s.repository.List(s.DB, tenantID, &entities.ListParams{Page: 1, Limit: 10})
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, result.Pagination.Total)
	require.Len(s.T(), result.Data, 1)
}

func (s *Suite) TestRepository_UpdateStatus() {
	tenantID := uuid.New()
	id := uuid.New()
	status := order.StatusConfirmed
	updatedBy := "user-2"

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET "status"=$1,"updated_by"=$2,"updated_at"=$3 WHERE id = $4 AND tenant_id = $5`)).
		WithArgs(status, updatedBy, sqlmock.AnyArg(), id, tenantID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.UpdateStatus(s.DB, tenantID, id, status, updatedBy)
	require.NoError(s.T(), err)
}

func (s *Suite) TestRepository_UpdatePaymentStatus() {
	tenantID := uuid.New()
	id := uuid.New()

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET "payment_status"=$1,"updated_at"=$2 WHERE id = $3 AND tenant_id = $4`)).
		WithArgs(order.PaymentPaid, sqlmock.AnyArg(), id, tenantID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.UpdatePaymentStatus(s.DB, tenantID, id, order.PaymentPaid)
	require.NoError(s.T(), err)
}
