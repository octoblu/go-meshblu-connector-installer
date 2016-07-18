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
	return fmt.Sprintf("MeshbluConnector-%s", uuid)
}

func getDefaultServiceDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors")
}
