package utils

import (
	"fmt"
	"os"
)

func CheckDockerAvailability() int {
	fmt.Println("Checking Docker availability...")

	available, message := IsDockerAvailable()

	fmt.Printf("Docker check result: %s\n", message)

	if available {
		return 0
	} else {
		return 1
	}
}

func init() {
	if len(os.Args) > 0 && os.Args[0] == "check_docker.go" {
		exitCode := CheckDockerAvailability()
		os.Exit(exitCode)
	}
}
