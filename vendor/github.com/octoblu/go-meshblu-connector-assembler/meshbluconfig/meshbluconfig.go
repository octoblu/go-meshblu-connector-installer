package meshbluconfig

import (
	"encoding/json"
	"path/filepath"

	"github.com/spf13/afero"
)

type meshbluConfig struct {
	UUID       string `json:"uuid"`
	Token      string `json:"token"`
	Domain     string `json:"domain,omitempty"`
	ResolveSRV bool   `json:"resolveSrv"`
}

// Options are the meshblu options that will be written to the file system
type Options struct {
	DirPath    string
	UUID       string
	Token      string
	Domain     string
	ResolveSRV bool
}

// Write a meshblu JSON to the file system
func Write(options Options) error {
	return WriteWithFS(options, afero.NewOsFs())
}

// WriteWithFS does everything Write does, only you get to supply
// your own FileSystem!
func WriteWithFS(options Options, fs afero.Fs) error {
	config := meshbluConfig{
		UUID:       options.UUID,
		Token:      options.Token,
		Domain:     options.Domain,
		ResolveSRV: options.ResolveSRV,
	}
	data, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		return err
	}

	meshbluConfigPath := filepath.Join(options.DirPath, "meshblu.json")
	err = afero.WriteFile(fs, meshbluConfigPath, data, 0644)

	if err != nil {
		return err
	}

	return nil
}
