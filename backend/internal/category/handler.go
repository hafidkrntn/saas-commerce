package category

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
	var req CreateCategoryRequest
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
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	result, err := h.svc.List(c.Request.Context(), entities.TenantID(c), c.Query("search"), page, limit)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateCategoryRequest
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
