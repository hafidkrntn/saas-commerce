package inventory

import (
	"backend-go/internal/entities"
	"backend-go/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListWarehouses(c *gin.Context) {
	data, err := h.svc.ListWarehouses(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateWarehouse(c *gin.Context) {
	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateWarehouse(c.Request.Context(), entities.TenantID(c), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) UpdateWarehouse(c *gin.Context) {
	var req UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateWarehouse(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) DeleteWarehouse(c *gin.Context) {
	if err := h.svc.DeleteWarehouse(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}
	response.Deleted(c, nil)
}

func (h *Handler) CreateMovement(c *gin.Context) {
	var req StockMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateMovement(c.Request.Context(), entities.TenantID(c), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) ListMovements(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	data, err := h.svc.ListMovements(c.Request.Context(), entities.TenantID(c), c.Query("product_id"), page, limit)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateShipment(c *gin.Context) {
	var req CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateShipment(c.Request.Context(), entities.TenantID(c), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) ListShipments(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	data, err := h.svc.ListShipments(c.Request.Context(), entities.TenantID(c), c.Query("status"), page, limit)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpdateShipment(c *gin.Context) {
	var req UpdateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateShipment(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}
