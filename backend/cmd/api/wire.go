//go:build wireinject

package main

import (
	"backend-go/internal/analytics"
	"backend-go/internal/auth"
	"backend-go/internal/category"
	"backend-go/internal/customer"
	"backend-go/internal/inventory"
	"backend-go/internal/order"
	"backend-go/internal/payment"
	"backend-go/internal/product"
	"backend-go/internal/review"
	"backend-go/internal/router"
	"backend-go/internal/setting"
	"backend-go/internal/tenant"
	"backend-go/internal/user"
	"backend-go/pkg/database"
	paymentclient "backend-go/third_party/payment"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var authSet = wire.NewSet(
	auth.NewRepository,
	auth.NewService,
	auth.NewHandler,
	auth.NewRouter,
)

var tenantSet = wire.NewSet(
	tenant.NewRepository,
	tenant.NewService,
	tenant.NewHandler,
	tenant.NewRouter,
)

var userSet = wire.NewSet(
	user.NewRepository,
	user.NewService,
	user.NewHandler,
	user.NewRouter,
)

var categorySet = wire.NewSet(
	category.NewRepository,
	category.NewService,
	category.NewHandler,
	category.NewRouter,
)

var productSet = wire.NewSet(
	product.NewRepository,
	product.NewService,
	product.NewHandler,
	product.NewRouter,
)

var customerSet = wire.NewSet(
	customer.NewRepository,
	customer.NewService,
	customer.NewHandler,
	customer.NewRouter,
)

var orderSet = wire.NewSet(
	order.NewRepository,
	order.NewService,
	order.NewHandler,
	order.NewRouter,
)

var paymentSet = wire.NewSet(
	MockProviderBaseURL,
	paymentclient.NewMockClient,
	payment.NewRepository,
	payment.NewService,
	payment.NewHandler,
	payment.NewRouter,
	wire.Bind(new(paymentclient.ClientInterface), new(*paymentclient.MockClient)),
)

var inventorySet = wire.NewSet(
	inventory.NewRepository,
	inventory.NewService,
	inventory.NewHandler,
	inventory.NewRouter,
)

var settingSet = wire.NewSet(
	setting.NewRepository,
	setting.NewService,
	setting.NewHandler,
	setting.NewRouter,
)

var analyticsSet = wire.NewSet(
	analytics.NewService,
	analytics.NewHandler,
	analytics.NewRouter,
)

var reviewSet = wire.NewSet(
	review.NewRepository,
	review.NewService,
	review.NewHandler,
	review.NewRouter,
)

func Init(log *logrus.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		database.NewGormDB,
		router.NewGinEngine,
		authSet,
		tenantSet,
		userSet,
		categorySet,
		productSet,
		customerSet,
		orderSet,
		paymentSet,
		inventorySet,
		settingSet,
		analyticsSet,
		reviewSet,
	))
}
