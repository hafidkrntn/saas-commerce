package order

import (
	"backend-go/internal/entities"
	"backend-go/pkg/response"
	"backend-go/pkg/utilities"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /orders
func (h *Handler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Create(c.Request.Context(), entities.TenantID(c), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, data, &req)
}

// GetByID handles GET /orders/:id
func (h *Handler) GetByID(c *gin.Context) {
	result, err := h.svc.GetByID(c.Request.Context(), entities.TenantID(c), c.Param("id"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// List handles GET /orders
func (h *Handler) List(c *gin.Context) {
	var params entities.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.svc.List(c.Request.Context(), entities.TenantID(c), &params)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// UpdateStatus handles PATCH /orders/:id/status
func (h *Handler) UpdateStatus(c *gin.Context) {
	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateStatus(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Updated(c, data, &req)
}

// Delete handles DELETE /orders/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := utilities.ParseUUID(c.Param("id"))
	if err != nil {
		response.ErrBadRequest(c, err)
		return
	}

	data, err := h.svc.GetByID(c.Request.Context(), entities.TenantID(c), id.String())
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Deleted(c, data)
}
