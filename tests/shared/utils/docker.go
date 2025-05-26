package utils

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func IsDockerAvailable() (bool, string) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:      "alpine:latest",
		WaitingFor: wait.ForLog("ready"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})

	if err == nil && container != nil {
		_ = container.Terminate(ctx)
		return true, "Docker is available and working properly"
	}

	var message string
	if err != nil {
		message = fmt.Sprintf("Docker is not available: %s", err.Error())
	} else {
		message = "Docker is not available: container could not be created"
	}

	return false, message
}
