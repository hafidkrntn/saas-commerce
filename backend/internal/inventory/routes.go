package inventory

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
	inv := rg.Group("/inventory")
	{
		inv.GET("/warehouses", r.handler.ListWarehouses)
		inv.POST("/warehouses", r.handler.CreateWarehouse)
		inv.PUT("/warehouses/:id", r.handler.UpdateWarehouse)
		inv.DELETE("/warehouses/:id", r.handler.DeleteWarehouse)

		inv.GET("/stock-movements", r.handler.ListMovements)
		inv.POST("/stock-movements", r.handler.CreateMovement)

		inv.GET("/shipments", r.handler.ListShipments)
		inv.POST("/shipments", r.handler.CreateShipment)
		inv.PUT("/shipments/:id", r.handler.UpdateShipment)
	}
}
