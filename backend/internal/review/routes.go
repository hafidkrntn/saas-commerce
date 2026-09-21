package review

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
	reviews := rg.Group("/reviews")
	{
		reviews.GET("", r.handler.List)
		reviews.POST("", r.handler.Create)
		reviews.PATCH("/:id/moderate", r.handler.Moderate)
		reviews.DELETE("/:id", r.handler.Delete)
	}
}
