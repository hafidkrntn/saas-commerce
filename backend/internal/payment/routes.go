package payment

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
	payments := rg.Group("/payments")
	{
		payments.GET("/methods", r.handler.ListMethods)
		payments.POST("/methods", r.handler.CreateMethod)
		payments.PUT("/methods/:id", r.handler.UpdateMethod)
		payments.DELETE("/methods/:id", r.handler.DeleteMethod)
		payments.POST("/charge", r.handler.Charge)
		payments.GET("/order/:orderId", r.handler.GetByOrder)
		payments.POST("/webhook", r.handler.HandleWebhook)
	}
}
