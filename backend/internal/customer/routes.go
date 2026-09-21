package customer

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
	customers := rg.Group("/customers")
	{
		customers.GET("", r.handler.List)
		customers.POST("", r.handler.Create)
		customers.GET("/:id", r.handler.GetByID)
		customers.PUT("/:id", r.handler.Update)
		customers.DELETE("/:id", r.handler.Delete)
		customers.POST("/:id/addresses", r.handler.AddAddress)
		customers.GET("/:id/activities", r.handler.GetActivities)
	}
}
