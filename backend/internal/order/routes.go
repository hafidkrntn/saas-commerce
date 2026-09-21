package order

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
	orders := rg.Group("/orders")
	{
		orders.GET("", r.handler.List)
		orders.POST("", r.handler.Create)
		orders.GET("/:id", r.handler.GetByID)
		orders.PATCH("/:id/status", r.handler.UpdateStatus)
	}
}
