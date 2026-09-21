package auth

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

// RegisterRoutes registers auth routes. Public endpoints (login, refresh)
// are registered on the bare group before the auth middleware is applied.
func (r *Router) RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.POST("/auth/login", r.handler.Login)
	public.POST("/auth/refresh", r.handler.Refresh)

	protected.POST("/auth/logout", r.handler.Logout)
	protected.GET("/auth/me", r.handler.Me)
}
