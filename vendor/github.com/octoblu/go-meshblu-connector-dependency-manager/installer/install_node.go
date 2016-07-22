package installer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/octoblu/go-meshblu-connector-installer/osruntime"
	"github.com/octoblu/unzipit"
	"github.com/spf13/afero"

	De "github.com/visionmedia/go-debug"
)

var de = De.Debug("meshblu-connector-installer:install_node")

// InstallNode installs the specified version of Node.JS
func InstallNode(tag, binPath string) error {
	return InstallNodeWithoutDefaults(tag, binPath, "https://nodejs.org", osruntime.New())
}

// InstallNodeWithoutDefaults installs the specified version of Node.JS
func InstallNodeWithoutDefaults(tag, binPath, baseURLStr string, osRuntime osruntime.OSRuntime) error {
	de("InstallNodeWithoutDefaults: %v | %v | %v | %v | %v", tag, binPath, baseURLStr, osRuntime.GOOS, osRuntime.GOARCH)
	if exists, err := nodeIsAlreadyInstalled(binPath, osRuntime); err != nil {
		return err
	} else if exists {
		de("node was already installed, skipping")
		return nil
	}

	packageURL, err := nodeURL(baseURLStr, tag, osRuntime)
	if err != nil {
		return err
	}
	de("resolved packageURL: %v", packageURL)

	response, err := http.Get(packageURL.String())
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("Expected HTTP status code 200, received: %v", response.StatusCode)
	}

	if osRuntime.GOOS == Windows {
		return installNodeOnWindowsFS(binPath, response.Body)
	}
	return installNodeAndNPMOnFS(binPath, response.Body)
}

func installNodeOnWindowsFS(binPath string, reader io.Reader) error {
	filePath := filepath.Join(binPath, "node.exe")
	return afero.WriteReader(afero.NewOsFs(), filePath, reader)
}

func installNodeAndNPMOnFS(binPath string, compressedReader io.Reader) error {
	archivePath, err := unzipit.UnpackStream(compressedReader, binPath)
	if err != nil {
		return err
	}

	actualNPMPath := filepath.Join(archivePath, "bin/npm")
	actualNodePath := filepath.Join(archivePath, "bin/node")
	pathToNodeSymlink := filepath.Join(binPath, "node")
	pathToNPMSymlink := filepath.Join(binPath, "npm")

	err = os.Symlink(actualNodePath, pathToNodeSymlink)
	if err != nil {
		return err
	}

	err = os.Symlink(actualNPMPath, pathToNPMSymlink)
	if err != nil {
		return err
	}

	return nil
}

func nodeIsAlreadyInstalled(binDir string, osRuntime osruntime.OSRuntime) (bool, error) {
	if osRuntime.GOOS == Windows {
		return afero.Exists(afero.NewOsFs(), filepath.Join(binDir, "node.exe"))
	}

	return afero.Exists(afero.NewOsFs(), filepath.Join(binDir, "node"))
}

func nodeURL(baseURLStr, tag string, osRuntime osruntime.OSRuntime) (*url.URL, error) {
	nodeURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, err
	}

	fileName, err := nodeFileName(tag, osRuntime)
	if err != nil {
		return nil, err
	}

	filePath, err := nodeFilePath(tag, osRuntime)
	if err != nil {
		return nil, err
	}

	nodeURL.Path = fmt.Sprintf("%v/%v", filePath, fileName)
	return nodeURL, nil
}

func nodeFileName(tag string, osRuntime osruntime.OSRuntime) (string, error) {
	if osRuntime.GOOS == "windows" {
		return "node.exe", nil
	}

	nodeArch, ok := ArchMap[osRuntime.GOARCH]
	if !ok {
		return "", fmt.Errorf("Unsupported architecture: %v", osRuntime.GOARCH)
	}

	return fmt.Sprintf("node-%v-%v-%v.tar.gz", tag, osRuntime.GOOS, nodeArch), nil
}

func nodeFilePath(tag string, osRuntime osruntime.OSRuntime) (string, error) {
	nodeArch, ok := ArchMap[osRuntime.GOARCH]
	if !ok {
		return "", fmt.Errorf("Unsupported architecture: %v", osRuntime.GOARCH)
	}

	if osRuntime.GOOS == "windows" {
		return fmt.Sprintf("/dist/%v/win-%v", tag, nodeArch), nil
	}

	return fmt.Sprintf("/dist/%v", tag), nil
}

// // Install is a convenience method for constructing an installer client
// // and calling client.Do
// func Install(depType, tag string) error {
// 	// client := New()
// 	// return client.Do(depType, tag)
// 	return nil
// }
