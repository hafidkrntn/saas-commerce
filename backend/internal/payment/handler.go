package payment

import (
	"backend-go/internal/entities"
	"backend-go/pkg/response"
	paymentclient "backend-go/third_party/payment"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListMethods(c *gin.Context) {
	data, err := h.svc.ListMethods(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateMethod(c *gin.Context) {
	var req CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateMethod(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) UpdateMethod(c *gin.Context) {
	var req UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateMethod(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) DeleteMethod(c *gin.Context) {
	data, err := h.svc.ListMethods(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}

	if err := h.svc.DeleteMethod(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}

	response.Deleted(c, data)
}

func (h *Handler) Charge(c *gin.Context) {
	var req ChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Charge(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) GetByOrder(c *gin.Context) {
	data, err := h.svc.GetPaymentByOrder(c.Request.Context(), entities.TenantID(c), c.Param("orderId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

// HandleWebhook is called by the payment provider (public route).
func (h *Handler) HandleWebhook(c *gin.Context) {
	var req paymentclient.MockWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.HandleWebhook(c.Request.Context(), entities.TenantID(c), req.TransactionID, req.Status)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}
