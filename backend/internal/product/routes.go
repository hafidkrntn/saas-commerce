package product

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
	products := rg.Group("/products")
	{
		products.GET("", r.handler.List)
		products.POST("", r.handler.Create)
		products.GET("/:id", r.handler.GetByID)
		products.PUT("/:id", r.handler.Update)
		products.DELETE("/:id", r.handler.Delete)
	}
}
