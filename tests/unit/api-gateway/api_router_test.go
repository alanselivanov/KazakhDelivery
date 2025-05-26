package apigateway_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceClient interface {
	MakeRequest(ctx context.Context, method, path string, body []byte, headers map[string]string) (int, []byte, error)
}

type MockServiceClient struct {
	mock.Mock
}

func (m *MockServiceClient) MakeRequest(ctx context.Context, method, path string, body []byte, headers map[string]string) (int, []byte, error) {
	args := m.Called(ctx, method, path, body, headers)
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

type ServiceRegistry struct {
	UserService      ServiceClient
	InventoryService ServiceClient
	OrderService     ServiceClient
}

type APIRouter struct {
	ServiceRegistry *ServiceRegistry
}

func NewAPIRouter(registry *ServiceRegistry) *APIRouter {
	return &APIRouter{
		ServiceRegistry: registry,
	}
}

func (r *APIRouter) RouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		headers := make(map[string]string)
		for key := range req.Header {
			headers[key] = req.Header.Get(key)
		}

		path := req.URL.Path
		var statusCode int
		var responseBody []byte
		var err error

		switch {
		case path == "/healthz":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
			return

		case path == "/api/auth/login" || path == "/api/auth/register" || path == "/api/users":
			statusCode, responseBody, err = r.ServiceRegistry.UserService.MakeRequest(
				ctx,
				req.Method,
				path,
				[]byte{},
				headers,
			)

		case path == "/api/products" || path == "/api/inventory":
			statusCode, responseBody, err = r.ServiceRegistry.InventoryService.MakeRequest(
				ctx,
				req.Method,
				path,
				[]byte{},
				headers,
			)

		case path == "/api/orders":
			statusCode, responseBody, err = r.ServiceRegistry.OrderService.MakeRequest(
				ctx,
				req.Method,
				path,
				[]byte{},
				headers,
			)

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Route not found"}`))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal service error"}`))
			return
		}

		w.WriteHeader(statusCode)
		w.Write(responseBody)
	}
}

func TestAPIRouter_RouteHandler_HealthCheck(t *testing.T) {
	registry := &ServiceRegistry{
		UserService:      &MockServiceClient{},
		InventoryService: &MockServiceClient{},
		OrderService:     &MockServiceClient{},
	}
	router := NewAPIRouter(registry)

	req := httptest.NewRequest("GET", "/healthz", nil)
	recorder := httptest.NewRecorder()

	handler := router.RouteHandler()
	handler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"status":"ok"}`, recorder.Body.String())
}

func TestAPIRouter_RouteHandler_UserService(t *testing.T) {
	mockUserService := new(MockServiceClient)
	registry := &ServiceRegistry{
		UserService:      mockUserService,
		InventoryService: &MockServiceClient{},
		OrderService:     &MockServiceClient{},
	}
	router := NewAPIRouter(registry)

	expectedResponse := []byte(`{"id":"123","username":"test_user"}`)
	mockUserService.On("MakeRequest", mock.Anything, "POST", "/api/auth/login", mock.Anything, mock.Anything).
		Return(http.StatusOK, expectedResponse, nil)

	req := httptest.NewRequest("POST", "/api/auth/login", nil)
	recorder := httptest.NewRecorder()

	handler := router.RouteHandler()
	handler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, string(expectedResponse), recorder.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestAPIRouter_RouteHandler_InventoryService(t *testing.T) {
	mockInventoryService := new(MockServiceClient)
	registry := &ServiceRegistry{
		UserService:      &MockServiceClient{},
		InventoryService: mockInventoryService,
		OrderService:     &MockServiceClient{},
	}
	router := NewAPIRouter(registry)

	expectedResponse := []byte(`[{"id":"prod1","name":"Product 1","price":9.99}]`)
	mockInventoryService.On("MakeRequest", mock.Anything, "GET", "/api/products", mock.Anything, mock.Anything).
		Return(http.StatusOK, expectedResponse, nil)

	req := httptest.NewRequest("GET", "/api/products", nil)
	recorder := httptest.NewRecorder()

	handler := router.RouteHandler()
	handler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, string(expectedResponse), recorder.Body.String())
	mockInventoryService.AssertExpectations(t)
}

func TestAPIRouter_RouteHandler_OrderService(t *testing.T) {
	mockOrderService := new(MockServiceClient)
	registry := &ServiceRegistry{
		UserService:      &MockServiceClient{},
		InventoryService: &MockServiceClient{},
		OrderService:     mockOrderService,
	}
	router := NewAPIRouter(registry)

	expectedResponse := []byte(`{"id":"order123","status":"processing"}`)
	mockOrderService.On("MakeRequest", mock.Anything, "GET", "/api/orders", mock.Anything, mock.Anything).
		Return(http.StatusOK, expectedResponse, nil)

	req := httptest.NewRequest("GET", "/api/orders", nil)
	recorder := httptest.NewRecorder()

	handler := router.RouteHandler()
	handler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, string(expectedResponse), recorder.Body.String())
	mockOrderService.AssertExpectations(t)
}

func TestAPIRouter_RouteHandler_NotFound(t *testing.T) {
	registry := &ServiceRegistry{
		UserService:      &MockServiceClient{},
		InventoryService: &MockServiceClient{},
		OrderService:     &MockServiceClient{},
	}
	router := NewAPIRouter(registry)

	req := httptest.NewRequest("GET", "/api/nonexistent", nil)
	recorder := httptest.NewRecorder()

	handler := router.RouteHandler()
	handler(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, `{"error":"Route not found"}`, recorder.Body.String())
}

func TestAPIRouter_RouteHandler_ServiceError(t *testing.T) {
	mockOrderService := new(MockServiceClient)
	registry := &ServiceRegistry{
		UserService:      &MockServiceClient{},
		InventoryService: &MockServiceClient{},
		OrderService:     mockOrderService,
	}
	router := NewAPIRouter(registry)

	mockOrderService.On("MakeRequest", mock.Anything, "GET", "/api/orders", mock.Anything, mock.Anything).
		Return(0, []byte{}, errors.New("service unavailable"))

	req := httptest.NewRequest("GET", "/api/orders", nil)
	recorder := httptest.NewRecorder()

	handler := router.RouteHandler()
	handler(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, `{"error":"Internal service error"}`, recorder.Body.String())
	mockOrderService.AssertExpectations(t)
}
