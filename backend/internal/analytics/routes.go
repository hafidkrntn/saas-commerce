package analytics

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
	analytics := rg.Group("/analytics")
	{
		analytics.GET("/summary", r.handler.Summary)
		analytics.GET("/revenue", r.handler.Revenue)
		analytics.GET("/sales", r.handler.Sales)
		analytics.GET("/top-products", r.handler.TopProducts)
		analytics.GET("/top-categories", r.handler.TopCategories)
	}
}
