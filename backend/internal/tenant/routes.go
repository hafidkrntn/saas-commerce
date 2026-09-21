package tenant

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

// RegisterRoutes mounts tenant management under /platform (platform admin only;
// the RequirePlatform middleware is applied by the router).
func (r *Router) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenants", r.handler.List)
	rg.POST("/tenants", r.handler.Create)
	rg.GET("/tenants/:id", r.handler.GetByID)
	rg.PUT("/tenants/:id", r.handler.Update)
	rg.DELETE("/tenants/:id", r.handler.Delete)
	rg.GET("/plans", r.handler.ListPlans)
}
