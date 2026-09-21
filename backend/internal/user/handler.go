package user

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

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateUser(c.Request.Context(), entities.TenantID(c), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) GetUser(c *gin.Context) {
	data, err := h.svc.GetUser(c.Request.Context(), entities.TenantID(c), c.Param("id"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	data, err := h.svc.ListUsers(c.Request.Context(), entities.TenantID(c), c.Query("search"), page, limit)
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateUser(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	if err := h.svc.DeleteUser(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}
	response.Deleted(c, nil)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.CreateRole(c.Request.Context(), entities.TenantID(c), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Created(c, data, &req)
}

func (h *Handler) GetRole(c *gin.Context) {
	data, err := h.svc.GetRole(c.Request.Context(), entities.TenantID(c), c.Param("id"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) ListRoles(c *gin.Context) {
	data, err := h.svc.ListRoles(c.Request.Context(), entities.TenantID(c))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	data, err := h.svc.UpdateRole(c.Request.Context(), entities.TenantID(c), c.Param("id"), &req, c.GetString("UserId"))
	if err != nil {
		response.Err(c, err)
		return
	}
	response.Updated(c, data, &req)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	if err := h.svc.DeleteRole(c.Request.Context(), entities.TenantID(c), c.Param("id")); err != nil {
		response.Err(c, err)
		return
	}
	response.Deleted(c, nil)
}
