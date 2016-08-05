package assembler

import (
	"fmt"
	"os"
	"path/filepath"
)

func getIgnitionPath(connectorDirectory string) string {
	return filepath.Join(connectorDirectory, "start")
}

func getServiceName(uuid string) string {
	return fmt.Sprintf("com.octoblu.%s", uuid)
}

func getServiceDirectory() string {
	return filepath.Join("/Library", "MeshbluConnectors")
}

func getUserServiceDirectory() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "MeshbluConnectors")
}
