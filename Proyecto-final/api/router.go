package api

import (
	"net/http"
	"sales-api/internal/sale"

	"github.com/gin-gonic/gin"
)

// InitRoutes registers all sale CRUD endpoints on the given Gin engine.
// It initializes the storage, service, and handler, then binds each HTTP
// method and path to the appropriate handler function.
func InitRoutes(e *gin.Engine) {
	storage := sale.NewLocalStorage()
	service := sale.NewService(storage)

	h := handler{
		saleService: service,
	}

	e.POST("/sales", h.handleCreate)
	//e.GET("/sales/:id", h.handleRead)
	e.PATCH("/sales/:id", h.handleUpdate)
	//e.DELETE("/sales/:id", h.handleDelete)

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
