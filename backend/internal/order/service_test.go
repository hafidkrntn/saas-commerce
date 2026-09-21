package order_test

import (
	"backend-go/internal/customer"
	"backend-go/internal/order"
	"backend-go/internal/product"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	customermocks "backend-go/internal/customer/mocks"
	ordermocks "backend-go/internal/order/mocks"
	productmocks "backend-go/internal/product/mocks"
)

type ServiceSuite struct {
	suite.Suite
	DB          *gorm.DB
	mock        sqlmock.Sqlmock
	repo        *ordermocks.MockRepositoryInterface
	productSvc  *productmocks.MockServiceInterface
	customerSvc *customermocks.MockServiceInterface
	service     order.ServiceInterface
}

func (s *ServiceSuite) SetupTest() {
	db, sqlMock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.mock = sqlMock
	ctrl := gomock.NewController(s.T())

	s.repo = ordermocks.NewMockRepositoryInterface(ctrl)
	s.productSvc = productmocks.NewMockServiceInterface(ctrl)
	s.customerSvc = customermocks.NewMockServiceInterface(ctrl)

	s.service = order.NewService(s.DB, s.repo, s.productSvc, s.customerSvc)
}

func (s *ServiceSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestCreate_Success() {
	tenantID := uuid.New()
	productID := uuid.New()
	customerID := uuid.New()
	userID := "user-1"

	req := &order.CreateOrderRequest{
		CustomerName:  "Budi",
		CustomerEmail: "budi@example.com",
		Items: []order.OrderItemRequest{
			{ProductID: productID.String(), Quantity: 2},
		},
		ShippingCost: 20000,
	}

	// Transaction scaffolding + order number lookup
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`SELECT .* FROM "orders" WHERE order_number LIKE`).
		WithArgs(sqlmock.AnyArg(), tenantID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"order_number"}))

	// Cross-module calls
	s.productSvc.EXPECT().GetByIDRaw(gomock.Any(), tenantID, productID).
		Return(makeProduct(productID, "Laptop", 1500000, 10), nil)
	s.productSvc.EXPECT().ReserveStock(gomock.Any(), tenantID, productID, 2).Return(nil)

	s.customerSvc.EXPECT().FindOrCreateByEmail(gomock.Any(), tenantID, "Budi", "budi@example.com").
		Return(customerRow(customerID), nil)

	// Repository calls
	s.repo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ *gorm.DB, o *order.Order) error {
			o.ID = uuid.New()
			o.OrderNumber = "ORD/2608/00001"
			o.CreatedAt = s.DB.NowFunc()
			return nil
		})
	s.repo.EXPECT().AddTimeline(gomock.Any(), gomock.Any()).Return(nil)
	s.customerSvc.EXPECT().RecordOrderStats(gomock.Any(), tenantID, customerID, gomock.Any()).Return(nil)
	s.customerSvc.EXPECT().AddActivity(gomock.Any(), tenantID, customerID, "order", gomock.Any()).Return(nil)

	s.mock.ExpectCommit()

	result, err := s.service.Create(context.Background(), tenantID, req, userID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), order.StatusPending, result.Status)
	require.Equal(s.T(), float64(3020000), result.TotalAmount)
	require.Len(s.T(), result.Items, 1)
}

func (s *ServiceSuite) TestCreate_ProductNotFound() {
	tenantID := uuid.New()
	productID := uuid.New()

	req := &order.CreateOrderRequest{
		CustomerName: "Budi",
		Items: []order.OrderItemRequest{
			{ProductID: productID.String(), Quantity: 1},
		},
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`SELECT .* FROM "orders" WHERE order_number LIKE`).
		WithArgs(sqlmock.AnyArg(), tenantID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"order_number"}))

	s.productSvc.EXPECT().GetByIDRaw(gomock.Any(), tenantID, productID).
		Return(makeProduct(productID, "Laptop", 1500000, 10), gorm.ErrRecordNotFound)

	s.mock.ExpectRollback()

	_, err := s.service.Create(context.Background(), tenantID, req, "user-1")
	require.Error(s.T(), err)
}

func (s *ServiceSuite) TestCreate_InsufficientStock() {
	tenantID := uuid.New()
	productID := uuid.New()

	req := &order.CreateOrderRequest{
		CustomerName: "Budi",
		Items: []order.OrderItemRequest{
			{ProductID: productID.String(), Quantity: 99},
		},
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`SELECT .* FROM "orders" WHERE order_number LIKE`).
		WithArgs(sqlmock.AnyArg(), tenantID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"order_number"}))

	s.productSvc.EXPECT().GetByIDRaw(gomock.Any(), tenantID, productID).
		Return(makeProduct(productID, "Laptop", 1500000, 10), nil)

	s.mock.ExpectRollback()

	_, err := s.service.Create(context.Background(), tenantID, req, "user-1")
	require.Error(s.T(), err)
}

func (s *ServiceSuite) TestUpdateStatus_CancelRestoresStock() {
	tenantID := uuid.New()
	orderID := uuid.New()
	productID := uuid.New()

	existing := order.Order{
		ID:           orderID,
		TenantID:     tenantID,
		OrderNumber:  "ORD/2608/00001",
		CustomerName: "Budi",
		Status:       order.StatusPending,
		Items: []order.OrderItem{
			{ProductID: productID, Quantity: 2, UnitPrice: 1000, Subtotal: 2000},
		},
	}

	s.mock.ExpectBegin()

	s.repo.EXPECT().GetByID(gomock.Any(), tenantID, orderID).Return(existing, nil)
	s.productSvc.EXPECT().RestoreStock(gomock.Any(), tenantID, productID, 2).Return(nil)
	s.repo.EXPECT().UpdateStatus(gomock.Any(), tenantID, orderID, order.StatusCancelled, "user-1").Return(nil)
	s.repo.EXPECT().AddTimeline(gomock.Any(), gomock.Any()).Return(nil)

	s.mock.ExpectCommit()

	result, err := s.service.UpdateStatus(context.Background(), tenantID, orderID.String(),
		&order.UpdateOrderStatusRequest{Status: order.StatusCancelled}, "user-1")
	require.NoError(s.T(), err)
	require.Equal(s.T(), order.StatusCancelled, result.Status)
}

func (s *ServiceSuite) TestUpdateStatus_InvalidTransition() {
	tenantID := uuid.New()
	orderID := uuid.New()

	existing := order.Order{
		ID:           orderID,
		TenantID:     tenantID,
		OrderNumber:  "ORD/2608/00001",
		CustomerName: "Budi",
		Status:       order.StatusCancelled,
	}

	s.mock.ExpectBegin()
	s.repo.EXPECT().GetByID(gomock.Any(), tenantID, orderID).Return(existing, nil)
	s.mock.ExpectRollback()

	_, err := s.service.UpdateStatus(context.Background(), tenantID, orderID.String(),
		&order.UpdateOrderStatusRequest{Status: order.StatusShipped}, "user-1")
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "tidak dapat berubah")
}

func (s *ServiceSuite) TestGetByID_NotFound() {
	tenantID := uuid.New()
	orderID := uuid.New()

	s.repo.EXPECT().GetByID(gomock.Any(), tenantID, orderID).Return(order.Order{}, gorm.ErrRecordNotFound)

	_, err := s.service.GetByID(context.Background(), tenantID, orderID.String())
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "tidak ditemukan")
}

func (s *ServiceSuite) TestGetByID_InvalidID() {
	_, err := s.service.GetByID(context.Background(), uuid.New(), "not-a-uuid")
	require.Error(s.T(), err)
}

// --- test helpers ---

type productModel = product.Product
type customerModel = customer.Customer

func makeProduct(id uuid.UUID, name string, price float64, stock int) (p productModel) {
	p.ID = id
	p.Name = name
	p.Price = price
	p.Stock = stock
	p.Status = "active"
	return p
}

func customerRow(id uuid.UUID) (c customerModel) {
	c.ID = id
	return c
}
