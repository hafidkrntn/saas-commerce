package product

import (
	"backend-go/internal/entities"
	"backend-go/pkg/response"
	"backend-go/pkg/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the product domain.
type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) tenantID(c *gin.Context) (uuid.UUID, bool) {
	tc := entities.GetTenantContext(c)
	id, err := tc.RequireTenantID()
	if err != nil {
		response.ErrForbidden(c, err)
		return uuid.Nil, false
	}
	return id, true
}

// Create handles POST /products
func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	tenantID, ok := h.tenantID(c)
	if !ok {
		return
	}

	userID := c.GetString("UserId")
	data, err := h.svc.Create(c.Request.Context(), tenantID, &req, userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, data, &req)
}

// GetByID handles GET /products/:id
func (h *Handler) GetByID(c *gin.Context) {
	tenantID, ok := h.tenantID(c)
	if !ok {
		return
	}

	result, err := h.svc.GetByID(c.Request.Context(), tenantID, c.Param("id"))
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// List handles GET /products
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := h.tenantID(c)
	if !ok {
		return
	}

	var params entities.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.svc.List(c.Request.Context(), tenantID, &params)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// Update handles PUT /products/:id
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := h.tenantID(c)
	if !ok {
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	userID := c.GetString("UserId")
	data, err := h.svc.Update(c.Request.Context(), tenantID, c.Param("id"), &req, userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Updated(c, data, &req)
}

// Delete handles DELETE /products/:id
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := h.tenantID(c)
	if !ok {
		return
	}

	id, err := utilities.ParseUUID(c.Param("id"))
	if err != nil {
		response.ErrBadRequest(c, err)
		return
	}

	data, err := h.svc.GetByID(c.Request.Context(), tenantID, id.String())
	if err != nil {
		response.Err(c, err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), tenantID, id.String()); err != nil {
		response.Err(c, err)
		return
	}

	response.Deleted(c, data)
}
