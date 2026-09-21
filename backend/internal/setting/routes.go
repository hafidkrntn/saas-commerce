package setting

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) RegisterRoutes(rg *gin.RouterGroup) {
	settings := rg.Group("/settings")
	{
		settings.GET("/store", r.handler.GetStore)
		settings.PUT("/store", r.handler.UpdateStore)
		settings.GET("/company", r.handler.GetCompany)
		settings.PUT("/company", r.handler.UpdateCompany)

		settings.GET("/tax-rates", r.handler.ListTaxRates)
		settings.POST("/tax-rates", r.handler.CreateTaxRate)
		settings.PUT("/tax-rates/:id", r.handler.UpdateTaxRate)
		settings.DELETE("/tax-rates/:id", r.handler.DeleteTaxRate)

		settings.GET("/shipping-zones", r.handler.ListShippingZones)
		settings.POST("/shipping-zones", r.handler.CreateShippingZone)
		settings.PUT("/shipping-zones/:id", r.handler.UpdateShippingZone)
		settings.DELETE("/shipping-zones/:id", r.handler.DeleteShippingZone)

		settings.GET("/notifications", r.handler.ListNotifications)
		settings.PUT("/notifications/:id", r.handler.UpdateNotification)
	}
}
