package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	testsDir := filepath.Join(filepath.Dir(filename), "..", "..")
	projectRoot := filepath.Join(testsDir, "..")
	return projectRoot
}

func LoadTestEnv(serviceName string) error {
	envFile := filepath.Join(GetProjectRoot(), serviceName, ".env.test")

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf("test env file not found for service: %s", serviceName)
	}

	content, err := os.ReadFile(envFile)
	if err != nil {
		return err
	}

	lines := filepath.SplitList(string(content))
	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		parts := filepath.SplitList(line)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}

	return nil
}

func ServiceBaseURL(serviceName string) string {
	switch serviceName {
	case "api-gateway":
		return "http://localhost:8080"
	case "inventory-service":
		return "http://localhost:8081"
	case "order-service":
		return "http://localhost:8082"
	case "user-service":
		return "http://localhost:8083"
	default:
		return ""
	}
}
