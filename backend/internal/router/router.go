package router

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
	"backend-go/internal/setting"
	"backend-go/internal/tenant"
	"backend-go/internal/user"
	"backend-go/middleware"
	"backend-go/pkg/response"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGinEngine(
	db *gorm.DB,
	authRouter *auth.Router,
	tenantRouter *tenant.Router,
	userRouter *user.Router,
	categoryRouter *category.Router,
	productRouter *product.Router,
	customerRouter *customer.Router,
	orderRouter *order.Router,
	paymentRouter *payment.Router,
	inventoryRouter *inventory.Router,
	settingRouter *setting.Router,
	analyticsRouter *analytics.Router,
	reviewRouter *review.Router,
) *gin.Engine {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "development" && appEnv != "staging" {
		gin.SetMode(gin.ReleaseMode)
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "saas-ecommerce"
	}

	app := gin.New()

	// Enable activity logging (used by response helpers)
	response.SetDB(db)

	// Middlewares
	app.Use(gin.Logger())
	app.Use(middleware.Recovery())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": appName})
	})

	// Readiness check
	app.GET("/readiness", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	api := app.Group("/api")
	v2 := api.Group("/v2")

	// Public routes (no auth)
	authRouter.RegisterRoutes(v2, v2.Group("", middleware.Auth()))

	// Protected tenant-scoped routes
	protected := v2.Group("", middleware.Auth(), middleware.RequireTenant())
	{
		categoryRouter.RegisterRoutes(protected)
		productRouter.RegisterRoutes(protected)
		customerRouter.RegisterRoutes(protected)
		orderRouter.RegisterRoutes(protected)
		paymentRouter.RegisterRoutes(protected)
		inventoryRouter.RegisterRoutes(protected)
		settingRouter.RegisterRoutes(protected)
		analyticsRouter.RegisterRoutes(protected)
		reviewRouter.RegisterRoutes(protected)
		userRouter.RegisterRoutes(protected)
	}

	// Platform admin routes
	platform := v2.Group("/platform", middleware.Auth(), middleware.RequirePlatform())
	{
		tenantRouter.RegisterRoutes(platform)
	}

	return app
}
