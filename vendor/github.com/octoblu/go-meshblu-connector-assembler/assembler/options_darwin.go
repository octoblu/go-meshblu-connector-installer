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

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors")
}

func getUserName() (string, error) {
	return "", nil
}
