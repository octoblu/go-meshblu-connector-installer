package assembler

import (
	"fmt"
	"os"
	"path/filepath"
)

func getIgnitionPath(connectorDirectory string) string {
	return filepath.Join(connectorDirectory, "start.exe")
}

func getServiceName(uuid string) string {
	return fmt.Sprintf("MeshbluConnector-%s", uuid)
}

// getDefaultServiceDirectory gets the OS specific install directory
func getDefaultServiceDirectory() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "MeshbluConnectors")
}
