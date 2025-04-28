package api

import (
	"errors"
	"net/http"
	"sales-api/internal/sale"

	"github.com/gin-gonic/gin"
)

// "github.com/go-playground/validator/v10"
// handler holds the sale service and implements HTTP handlers for sale CRUD.
type handler struct {
	saleService *sale.Service
}

// handleCreate handles POST /sales
func (h *handler) handleCreate(ctx *gin.Context) {
	// request payload
	var req struct {
		user_id string  `json:"user_id" binding:"required"`
		amount  float32 `json:"amount" binding:"required,gt=0"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := http.Get("localhost:8080/user/" + req.user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if resp.StatusCode != 200 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error en la consulta de usuarios."})
		return
	}
	u := &sale.Sale{
		user_id: req.user_id,
		amount:  req.amount,
	}
	if err := h.saleService.Create(u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

// handleRead handles GET /sales/:id
func (h *handler) handleRead(ctx *gin.Context) {
	/*
		id := ctx.Param("id")

		u, err := h.saleService.Get(id)
		if err != nil {
			if errors.Is(err, sale.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, u)
	*/
}

// handleUpdate handles PUT /sales/:id
func (h *handler) handleUpdate(ctx *gin.Context) {
	id := ctx.Param("id")

	// bind partial update fields
	var fields *sale.UpdateFields
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.saleService.Update(id, fields)
	if err != nil {
		if errors.Is(err, sale.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// handleDelete handles DELETE /sales/:id
func (h *handler) handleDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.saleService.Delete(id); err != nil {
		if errors.Is(err, sale.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
