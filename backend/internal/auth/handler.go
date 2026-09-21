package auth

import (
	"backend-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, data)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.Refresh(c.Request.Context(), &req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, data)
}

func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	if err := h.svc.Logout(c.Request.Context(), &req); err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) Me(c *gin.Context) {
	data, err := h.svc.Me(c.Request.Context(), c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, data)
}
