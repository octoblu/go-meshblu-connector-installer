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

func getServiceDirectory() string {
	programFilesDir := os.Getenv("PROGRAMFILESX86")
	if programFilesDir == "" {
		programFilesDir = os.Getenv("PROGRAMFILES")
	}
	return filepath.Join(programFilesDir, "MeshbluConnectors")
}

func getUserServiceDirectory() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "MeshbluConnectors")
}
