package api

import (
	"errors"
	"net/http"
	"sales-api/internal/sale"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Define una interfaz para poder inyectar mock en testing
type HTTPClient interface {
	Get(url string) (*resty.Response, error)
}

// realHTTPClient implementa HTTPClient usando resty.Client
type realHTTPClient struct {
	client *resty.Client
}

func (r *realHTTPClient) Get(url string) (*resty.Response, error) {
	return r.client.R().Get(url)
}

// handler holds the sale service and implements HTTP handlers for sale CRUD.
type handler struct {
	saleService *sale.Service
	httpClient  HTTPClient
	logger      *zap.Logger
}

// handleCreate handles POST /sales
func (h *handler) handleCreate(ctx *gin.Context) {
	// request payload
	var req struct {
		UserID string  `json:"user_id" binding:"required"`
		Amount float32 `json:"amount" binding:"required,gt=0"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validar que el usuario exista
	resp, err := h.httpClient.Get("http://localhost:8080/users/" + req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al contactar servicio de usuarios"})
		return
	}
	if resp.StatusCode() != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el usuario no existe"})
		return
	}
	u := &sale.Sale{
		UserID: req.UserID,
		Amount: req.Amount,
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
		h.logger.Warn("binding error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.saleService.Update(id, fields)

	if err != nil {
		h.logger.Warn("update failed", zap.String("id", id), zap.Error(err))
		if errors.Is(err, sale.ErrSaleNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, sale.ErrInvalidStateChange) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, sale.ErrInvalidNewState) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// Crear endpoint GET /sales con filtros por user_id y status.
func (h *handler) handleList(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.Query("status")
	validStates := map[string]bool{"approved": true, "rejected": true, "pending": true}
	// no se pide esta validación, pero la coloco ya que no puede venir el id del user vacio!
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}
	if status != "" && !validStates[status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estado inválido"})
		return
	}

	sales, err := h.saleService.ListByUserAndStatus(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metadata := map[string]interface{}{
		"quantity":     len(sales),
		"approved":     0,
		"rejected":     0,
		"pending":      0,
		"total_amount": 0.0,
	}
	for _, s := range sales {
		metadata[s.Estado] = metadata[s.Estado].(int) + 1
		metadata["total_amount"] = metadata["total_amount"].(float64) + float64(s.Amount)
	}

	c.JSON(http.StatusOK, gin.H{
		"metadata": metadata,
		"results":  sales,
	})
}
