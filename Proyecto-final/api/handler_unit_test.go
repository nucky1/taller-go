package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sales-api/internal/sale"
	"strings"
	"testing"
	"time"

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

// createTestSale crea una venta con estado inicial forzado (por ej: "pending")
func createTestSale(svc *sale.Service, userID string, amount float32, estado string) *sale.Sale {
	s := &sale.Sale{
		UserID:    userID,
		Amount:    amount,
		Estado:    estado,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
	_ = svc.Create(s)
	return s
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

	t.Run("actualiza correctamente de pending a approved", func(t *testing.T) {
		router := gin.New()
		storage := sale.NewLocalStorage()
		service := sale.NewService(storage)

		s := createTestSale(service, "abc123", 150, "pending")

		h := handler{
			saleService: service,
			httpClient:  &fakeClientOK{},
			logger:      logger,
		}
		router.PATCH("/sales/:id", h.handleUpdate)

		updateBody := `{"estado": "approved"}`
		reqUpdate := httptest.NewRequest(http.MethodPatch, "/sales/"+s.ID, bytes.NewBufferString(updateBody))
		reqUpdate.Header.Set("Content-Type", "application/json")
		recUpdate := httptest.NewRecorder()
		router.ServeHTTP(recUpdate, reqUpdate)

		require.Equal(t, http.StatusOK, recUpdate.Code)
		assert.Contains(t, recUpdate.Body.String(), `"estado":"approved"`)
	})

	t.Run("error por estado inválido", func(t *testing.T) {
		router := gin.New()
		storage := sale.NewLocalStorage()
		service := sale.NewService(storage)

		s := createTestSale(service, "abc123", 150, "pending")

		h := handler{
			saleService: service,
			httpClient:  &fakeClientOK{},
			logger:      logger,
		}
		router.PATCH("/sales/:id", h.handleUpdate)

		updateBody := `{"estado": "cancelled"}`
		reqUpdate := httptest.NewRequest(http.MethodPatch, "/sales/"+s.ID, bytes.NewBufferString(updateBody))
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

//======================= flujo completo POST → PATCH → GET (happy path) =======================//

func TestIntegration_FlujoCompleto(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()

	router := gin.New()
	storage := sale.NewLocalStorage()
	service := sale.NewService(storage)

	h := handler{
		saleService: service,
		httpClient:  &fakeClientOK{},
		logger:      logger,
	}

	router.POST("/sales", h.handleCreate)
	router.PATCH("/sales/:id", h.handleUpdate)
	router.GET("/sales", h.handleList)

	// 1. POST /sales
	body := `{"user_id": "abc123", "amount": 150}`
	req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var created sale.Sale
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	require.Equal(t, "abc123", created.UserID)
	require.Equal(t, float32(150), created.Amount)
	require.Equal(t, "pending", created.Estado)

	// 2. PATCH /sales/:id
	updateBody := `{"estado": "approved"}`
	reqPatch := httptest.NewRequest(http.MethodPatch, "/sales/"+created.ID, strings.NewReader(updateBody))
	reqPatch.Header.Set("Content-Type", "application/json")
	recPatch := httptest.NewRecorder()
	router.ServeHTTP(recPatch, reqPatch)

	require.Equal(t, http.StatusOK, recPatch.Code)
	var updated sale.Sale
	err = json.Unmarshal(recPatch.Body.Bytes(), &updated)
	require.NoError(t, err)
	require.Equal(t, "approved", updated.Estado)

	// 3. GET /sales?user_id=abc123&status=approved
	reqGet := httptest.NewRequest(http.MethodGet, "/sales?user_id=abc123&status=approved", nil)
	recGet := httptest.NewRecorder()
	router.ServeHTTP(recGet, reqGet)

	require.Equal(t, http.StatusOK, recGet.Code)
	bodyResp := recGet.Body.String()
	assert.Contains(t, bodyResp, `"approved":1`)
	assert.Contains(t, bodyResp, `"results"`)
	assert.Contains(t, bodyResp, `"user_id":"abc123"`)
	assert.Contains(t, bodyResp, `"estado":"approved"`)
}

//======================= flujo completo POST → PATCH → GET (happy path) =======================//
