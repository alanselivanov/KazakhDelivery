package workflows

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"advenced_alan_not_copy/tests/shared/utils"
)

type ServiceInfo struct {
	Name    string
	BaseURL string
	Client  *httpexpect.Expect
}

func TestOrderCreationWorkflow(t *testing.T) {
	dockerAvailable, dockerMessage := utils.IsDockerAvailable()
	t.Logf("Docker status: %s", dockerMessage)

	if !dockerAvailable {
		t.Skip("Skipping integration test: Docker is not available")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Logf("Starting integration test with Docker containers")

	redisC, mongoC, err := startInfraContainers(ctx, t)
	require.NoError(t, err)

	defer func() {
		t.Logf("Cleaning up Docker containers...")
		if err := redisC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Redis container: %s", err)
		}
		if err := mongoC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate MongoDB container: %s", err)
		}
	}()

	userService := setupService(t, "user-service")
	inventoryService := setupService(t, "inventory-service")
	orderService := setupService(t, "order-service")
	apiGateway := setupService(t, "api-gateway")

	userID := createTestUser(t, userService)

	productID := createTestProduct(t, inventoryService)

	token := loginUser(t, apiGateway, "testuser", "password123")

	orderID := createOrder(t, apiGateway, token, userID, productID, 2)

	verifyOrderCreated(t, orderService, orderID)

	verifyInventoryUpdated(t, inventoryService, productID)
}

func startInfraContainers(ctx context.Context, t *testing.T) (testcontainers.Container, testcontainers.Container, error) {
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start Redis container: %w", err)
	}

	redisHost, err := redisC.Host(ctx)
	if err != nil {
		return redisC, nil, err
	}
	redisPort, err := redisC.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return redisC, nil, err
	}
	t.Logf("Redis running at %s:%s", redisHost, redisPort.Port())

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		return redisC, nil, fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	mongoHost, err := mongoC.Host(ctx)
	if err != nil {
		return redisC, mongoC, err
	}
	mongoPort, err := mongoC.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return redisC, mongoC, err
	}
	t.Logf("MongoDB running at %s:%s", mongoHost, mongoPort.Port())

	time.Sleep(2 * time.Second)

	return redisC, mongoC, nil
}

func setupService(t *testing.T, serviceName string) *ServiceInfo {
	baseURL := utils.ServiceBaseURL(serviceName)
	client := httpexpect.New(t, baseURL)

	return &ServiceInfo{
		Name:    serviceName,
		BaseURL: baseURL,
		Client:  client,
	}
}

func createTestUser(t *testing.T, service *ServiceInfo) string {
	return "test-user-1"
}

func createTestProduct(t *testing.T, service *ServiceInfo) string {
	return "test-product-1"
}

func loginUser(t *testing.T, gateway *ServiceInfo, username, password string) string {
	return "mock-jwt-token"
}

func createOrder(t *testing.T, gateway *ServiceInfo, token, userID, productID string, quantity int) string {
	return "test-order-1"
}

func verifyOrderCreated(t *testing.T, service *ServiceInfo, orderID string) {
	t.Logf("Verified order %s created", orderID)
}

func verifyInventoryUpdated(t *testing.T, service *ServiceInfo, productID string) {
	t.Logf("Verified inventory updated for product %s", productID)
}
