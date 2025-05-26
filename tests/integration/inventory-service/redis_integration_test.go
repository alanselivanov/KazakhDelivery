package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"advenced_alan_not_copy/tests/shared/utils"

	"github.com/go-redis/redis/v8"
)

func TestRedisIntegration(t *testing.T) {
	dockerAvailable, dockerMessage := utils.IsDockerAvailable()
	t.Logf("Docker status: %s", dockerMessage)

	if !dockerAvailable {
		t.Skip("Skipping integration test: Docker is not available")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Logf("Starting Redis container for integration test")
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)

	defer func() {
		t.Logf("Terminating Redis container")
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	host, err := redisContainer.Host(ctx)
	require.NoError(t, err)

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	pong, err := redisClient.Ping(ctx).Result()
	require.NoError(t, err)
	assert.Equal(t, "PONG", pong)

	t.Run("Set and Get", func(t *testing.T) {
		err := redisClient.Set(ctx, "test:key", "test-value", time.Minute).Err()
		require.NoError(t, err)

		value, err := redisClient.Get(ctx, "test:key").Result()
		require.NoError(t, err)
		assert.Equal(t, "test-value", value)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete-key"
		err := redisClient.Set(ctx, key, "delete-me", time.Minute).Err()
		require.NoError(t, err)

		err = redisClient.Del(ctx, key).Err()
		require.NoError(t, err)

		_, err = redisClient.Get(ctx, key).Result()
		assert.ErrorIs(t, err, redis.Nil)
	})

	t.Run("Pattern Deletion", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			key := fmt.Sprintf("test:pattern:%d", i)
			err := redisClient.Set(ctx, key, fmt.Sprintf("value-%d", i), time.Minute).Err()
			require.NoError(t, err)
		}

		keys, err := redisClient.Keys(ctx, "test:pattern:*").Result()
		require.NoError(t, err)
		assert.Len(t, keys, 5)

		for _, key := range keys {
			err := redisClient.Del(ctx, key).Err()
			require.NoError(t, err)
		}

		keys, err = redisClient.Keys(ctx, "test:pattern:*").Result()
		require.NoError(t, err)
		assert.Len(t, keys, 0)
	})
}
