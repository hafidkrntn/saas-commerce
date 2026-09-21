package setting

import (
	"backend-go/internal/entities"
	"backend-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetStore(c *gin.Context) {
	data, err := h.svc.GetStore(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpdateStore(c *gin.Context) {
	var req UpdateStoreSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateStore(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) GetCompany(c *gin.Context) {
	data, err := h.svc.GetCompany(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpdateCompany(c *gin.Context) {
	var req UpdateCompanySettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateCompany(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) ListTaxRates(c *gin.Context) {
	data, err := h.svc.ListTaxRates(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateTaxRate(c *gin.Context) {
	var req CreateTaxRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateTaxRate(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) UpdateTaxRate(c *gin.Context) {
	var req UpdateTaxRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateTaxRate(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) DeleteTaxRate(c *gin.Context) {
	if err := h.svc.DeleteTaxRate(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}
	response.Deleted(c, nil)
}

func (h *Handler) ListShippingZones(c *gin.Context) {
	data, err := h.svc.ListShippingZones(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) CreateShippingZone(c *gin.Context) {
	var req CreateShippingZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateShippingZone(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) UpdateShippingZone(c *gin.Context) {
	var req UpdateShippingZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateShippingZone(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) DeleteShippingZone(c *gin.Context) {
	if err := h.svc.DeleteShippingZone(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}
	response.Deleted(c, nil)
}

func (h *Handler) ListNotifications(c *gin.Context) {
	data, err := h.svc.ListNotifications(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpdateNotification(c *gin.Context) {
	var req UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateNotification(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}
