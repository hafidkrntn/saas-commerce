package review

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
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Create(c.Request.Context(), entities.TenantID(c), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	onlyPublished, _ := strconv.ParseBool(c.Query("published"))

	data, err := h.svc.List(c.Request.Context(), entities.TenantID(c), c.Query("product_id"), onlyPublished, page, limit)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Moderate(c *gin.Context) {
	var req ModerateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Moderate(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}
	response.Deleted(c, nil)
}
