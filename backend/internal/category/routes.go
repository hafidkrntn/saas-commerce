package category

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
	categories := rg.Group("/categories")
	{
		categories.GET("", r.handler.List)
		categories.POST("", r.handler.Create)
		categories.GET("/:id", r.handler.GetByID)
		categories.PUT("/:id", r.handler.Update)
		categories.DELETE("/:id", r.handler.Delete)
	}
}
