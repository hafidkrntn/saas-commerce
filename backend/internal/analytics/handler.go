package analytics

import (
	"backend-go/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Summary(c *gin.Context) {
	from, to := periodFromQuery(c)
	data, err := h.svc.Summary(c.Request.Context(), c.GetString("TenantId"), from, to)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Revenue(c *gin.Context) {
	months, _ := strconv.Atoi(c.Query("months"))
	data, err := h.svc.Revenue(c.Request.Context(), c.GetString("TenantId"), months)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Sales(c *gin.Context) {
	days, _ := strconv.Atoi(c.Query("days"))
	data, err := h.svc.Sales(c.Request.Context(), c.GetString("TenantId"), days)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) TopProducts(c *gin.Context) {
	from, to := periodFromQuery(c)
	limit, _ := strconv.Atoi(c.Query("limit"))
	data, err := h.svc.TopProducts(c.Request.Context(), c.GetString("TenantId"), from, to, limit)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) TopCategories(c *gin.Context) {
	from, to := periodFromQuery(c)
	limit, _ := strconv.Atoi(c.Query("limit"))
	data, err := h.svc.TopCategories(c.Request.Context(), c.GetString("TenantId"), from, to, limit)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func periodFromQuery(c *gin.Context) (time.Time, time.Time) {
	from, _ := time.Parse(time.RFC3339, c.Query("from"))
	to, _ := time.Parse(time.RFC3339, c.Query("to"))
	return from, to
}
