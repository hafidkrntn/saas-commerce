package customer

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

func (h *Handler) Create(c *gin.Context) {
	var req CreateCustomerRequest
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

func (h *Handler) GetByID(c *gin.Context) {
	result, err := h.svc.GetByID(c.Request.Context(), entities.TenantID(c), c.Param("id"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

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

func (h *Handler) Update(c *gin.Context) {
	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Update(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Updated(c, data, &req)
}

func (h *Handler) Delete(c *gin.Context) {
	data, err := h.svc.GetByID(c.Request.Context(), entities.TenantID(c), c.Param("id"))
	if err != nil {
		response.Err(c, err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}

	response.Deleted(c, data)
}

func (h *Handler) AddAddress(c *gin.Context) {
	var req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.AddAddress(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, data, &req)
}

func (h *Handler) GetActivities(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))

	data, err := h.svc.GetActivities(c.Request.Context(), entities.TenantID(c), c.Param("id"), limit)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, data)
}
