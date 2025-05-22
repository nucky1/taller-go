package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sales-api/internal/sale"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Fake client: usuario NO existente
type fakeClientNotFound struct{}

func (f *fakeClientNotFound) Get(url string) (*resty.Response, error) {
	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       http.NoBody,
		},
	}, nil
}

// Fake client: usuario existente
type fakeClientOK struct{}

func (f *fakeClientOK) Get(url string) (*resty.Response, error) {
	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       http.NoBody,
		},
	}, nil
}

// helper para crear handler
func newHandler(client HTTPClient, logger *zap.Logger) handler {
	return handler{
		saleService: sale.NewService(sale.NewLocalStorage()),
		httpClient:  client,
		logger:      logger,
	}
}

// ======================= CREATE =======================//
func TestCreateSale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()

	t.Run("Crear Venta: usuario inválido @400", func(t *testing.T) {
		router := gin.New()
		h := newHandler(&fakeClientNotFound{}, logger)
		router.POST("/sales", h.handleCreate)

		body := `{"user_id": "no-existe", "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "el usuario no existe")
	})

	t.Run("Crear Venta: usuario válido @201", func(t *testing.T) {
		router := gin.New()
		h := newHandler(&fakeClientOK{}, logger)
		router.POST("/sales", h.handleCreate)

		body := `{"user_id": "abc123", "amount": 200}`
		req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), `"user_id":"abc123"`)
		assert.Contains(t, rec.Body.String(), `"amount":200`)
		assert.Contains(t, rec.Body.String(), `"estado"`)
		assert.Contains(t, rec.Body.String(), `"id"`)
	})
}

//======================= CREATE =======================//

// ======================= UPDATE =======================//

func TestUpdateSale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()

	t.Run("actualiza correctamente", func(t *testing.T) {
		router := gin.New()
		h := newHandler(&fakeClientOK{}, logger)
		router.POST("/sales", h.handleCreate)
		router.PATCH("/sales/:id", h.handleUpdate)

		createBody := `{"user_id": "abc123", "amount": 150}`
		reqCreate := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(createBody))
		reqCreate.Header.Set("Content-Type", "application/json")
		recCreate := httptest.NewRecorder()
		router.ServeHTTP(recCreate, reqCreate)
		require.Equal(t, http.StatusCreated, recCreate.Code)

		bodyStr := recCreate.Body.String()
		idStart := strings.Index(bodyStr, `"id":"`) + len(`"id":"`)
		idEnd := strings.Index(bodyStr[idStart:], `"`) + idStart
		saleID := bodyStr[idStart:idEnd]

		updateBody := `{"estado": "approved"}`
		reqUpdate := httptest.NewRequest(http.MethodPatch, "/sales/"+saleID, bytes.NewBufferString(updateBody))
		reqUpdate.Header.Set("Content-Type", "application/json")
		recUpdate := httptest.NewRecorder()
		router.ServeHTTP(recUpdate, reqUpdate)

		require.Equal(t, http.StatusOK, recUpdate.Code)
		assert.Contains(t, recUpdate.Body.String(), `"estado":"approved"`)
	})
	t.Run("actualiza correctamente de pending a approved", func(t *testing.T) {
		router := gin.New()
		h := newHandler(&fakeClientOK{}, logger)
		router.POST("/sales", h.handleCreate)
		router.PATCH("/sales/:id", h.handleUpdate)

		createBody := `{"user_id": "abc123", "amount": 150}`
		reqCreate := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(createBody))
		reqCreate.Header.Set("Content-Type", "application/json")
		recCreate := httptest.NewRecorder()
		router.ServeHTTP(recCreate, reqCreate)
		require.Equal(t, http.StatusCreated, recCreate.Code)

		bodyStr := recCreate.Body.String()
		idStart := strings.Index(bodyStr, `"id":"`) + len(`"id":"`)
		idEnd := strings.Index(bodyStr[idStart:], `"`) + idStart
		saleID := bodyStr[idStart:idEnd]

		updateBody := `{"estado": "approved"}`
		reqUpdate := httptest.NewRequest(http.MethodPatch, "/sales/"+saleID, bytes.NewBufferString(updateBody))
		reqUpdate.Header.Set("Content-Type", "application/json")
		recUpdate := httptest.NewRecorder()
		router.ServeHTTP(recUpdate, reqUpdate)

		require.Equal(t, http.StatusOK, recUpdate.Code)
		assert.Contains(t, recUpdate.Body.String(), `"estado":"approved"`)
	})

	t.Run("error por estado inválido", func(t *testing.T) {
		router := gin.New()
		h := newHandler(&fakeClientOK{}, logger)
		router.POST("/sales", h.handleCreate)
		router.PATCH("/sales/:id", h.handleUpdate)

		createBody := `{"user_id": "abc123", "amount": 150}`
		reqCreate := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(createBody))
		reqCreate.Header.Set("Content-Type", "application/json")
		recCreate := httptest.NewRecorder()
		router.ServeHTTP(recCreate, reqCreate)
		require.Equal(t, http.StatusCreated, recCreate.Code)

		bodyStr := recCreate.Body.String()
		idStart := strings.Index(bodyStr, `"id":"`) + len(`"id":"`)
		idEnd := strings.Index(bodyStr[idStart:], `"`) + idStart
		saleID := bodyStr[idStart:idEnd]

		updateBody := `{"estado": "cancelled"}`
		reqUpdate := httptest.NewRequest(http.MethodPatch, "/sales/"+saleID, bytes.NewBufferString(updateBody))
		reqUpdate.Header.Set("Content-Type", "application/json")
		recUpdate := httptest.NewRecorder()
		router.ServeHTTP(recUpdate, reqUpdate)

		require.Equal(t, http.StatusBadRequest, recUpdate.Code)
		assert.Contains(t, recUpdate.Body.String(), "estado no válido para cambio")
	})

	t.Run("error por id inexistente", func(t *testing.T) {
		router := gin.New()
		h := newHandler(&fakeClientOK{}, logger)
		router.PATCH("/sales/:id", h.handleUpdate)

		updateBody := `{"estado": "approved"}`
		reqUpdate := httptest.NewRequest(http.MethodPatch, "/sales/no-existe", bytes.NewBufferString(updateBody))
		reqUpdate.Header.Set("Content-Type", "application/json")
		recUpdate := httptest.NewRecorder()
		router.ServeHTTP(recUpdate, reqUpdate)

		require.Equal(t, http.StatusNotFound, recUpdate.Code)
		assert.Contains(t, recUpdate.Body.String(), "sale not found")
	})
}

//======================= UPDATE =======================//
