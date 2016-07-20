package osruntime

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// OSRuntime is used to provide the GOOS and GOARCH variables to the DownloadWithURLAndRuntime
// Generally, this should not be instantiated by hand. Instead, call Download() and it'll take
// care of this part for you
type OSRuntime struct {
	GOOS, GOARCH string
}

// New creates a new runtime instance with the values prepopulated from the
// currently running process
func New() OSRuntime {
	return OSRuntime{
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
	}
}

// BinPath returns the expected binPath for the given OS
func BinPath(osRuntime OSRuntime) (string, error) {
	switch osRuntime.GOOS {
	case "darwin", "linux":
		return filepath.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors", "bin"), nil
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "MeshbluConnectors", "bin"), nil
	}
	return "", fmt.Errorf("Unsupported OS: %v", osRuntime.GOOS)
}
