package user

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
	users := rg.Group("/users")
	{
		users.GET("", r.handler.ListUsers)
		users.POST("", r.handler.CreateUser)
		users.GET("/:id", r.handler.GetUser)
		users.PUT("/:id", r.handler.UpdateUser)
		users.DELETE("/:id", r.handler.DeleteUser)
	}

	roles := rg.Group("/roles")
	{
		roles.GET("", r.handler.ListRoles)
		roles.POST("", r.handler.CreateRole)
		roles.GET("/:id", r.handler.GetRole)
		roles.PUT("/:id", r.handler.UpdateRole)
		roles.DELETE("/:id", r.handler.DeleteRole)
	}
}
