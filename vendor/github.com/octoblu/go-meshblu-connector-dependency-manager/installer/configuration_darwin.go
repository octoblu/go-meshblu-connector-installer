package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// finalDependencyFileName gets the final dependency filename
func finalDependencyFileName(depType, tag string) string {
	if depType == NodeType {
		return NodeType
	}
	return ""
}

// getResourceURI defines the uri to download
func getResourceURI(depType, tag string) string {
	if depType == NodeType {
		return getNodeURI(tag)
	}
	return ""
}

func getNodeURI(tag string) string {
	return strings.Replace("https://nodejs.org/dist/:tag:/node-:tag:-darwin-x64.tar.gz", ":tag:", tag, -1)
}

// getBinPath defines the target location
func getBinPath() string {
	return filepath.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors", "bin")
}

// extractBin allows you too extract the bin from the download
func extractBin(depType, target, tag string) error {
	if depType == NodeType {
		return extractNode(target, tag)
	}
	return fmt.Errorf("Unsupported platform")
}

// extractNode extracts the node dependencies
func extractNode(target, tag string) error {
	folderName := strings.Replace("node-:tag:-darwin-x64", ":tag:", tag, -1)
	nodePath := filepath.Join(target, folderName, "bin", "node")
	nodeSymPath := filepath.Join(target, "node")
	os.Remove(nodeSymPath)
	err := os.Symlink(nodePath, nodeSymPath)
	if err != nil {
		return err
	}

	npmPath := filepath.Join(target, folderName, "bin", "npm")
	npmSymPath := filepath.Join(target, "npm")
	os.Remove(npmSymPath)
	err = os.Symlink(npmPath, npmSymPath)
	if err != nil {
		return err
	}
	return nil
}
