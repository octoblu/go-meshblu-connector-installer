package serviceconfig

import (
	"encoding/json"
	"path/filepath"

	"github.com/spf13/afero"
)

// Options define the options passed to Write when
// saving a service file
type Options struct {
	ServiceName   string
	ServiceType   string
	DisplayName   string
	Description   string
	ConnectorName string
	GithubSlug    string
	Tag           string
	BinPath       string
	Dir           string
	LogDir        string
}

type serviceFile struct {
	ServiceName   string
	ServiceType   string
	DisplayName   string
	Description   string
	ConnectorName string
	GithubSlug    string
	Tag           string
	BinPath       string
	Dir           string

	Stderr, Stdout string
}

// Write a ServiceConfig JSON to the file system
func Write(options Options) error {
	return WriteWithFS(options, afero.NewOsFs())
}

// WriteWithFS does everything Write does, only you get to supply
// your own FileSystem!
func WriteWithFS(options Options, fs afero.Fs) error {
	serviceFile := optionsToServiceFile(options)

	data, err := json.MarshalIndent(serviceFile, "", "  ")
	if err != nil {
		return err
	}

	serviceConfigPath := filepath.Join(options.Dir, "service.json")
	err = afero.WriteFile(fs, serviceConfigPath, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func optionsToServiceFile(options Options) serviceFile {
	return serviceFile{
		Stdout: filepath.Join(options.LogDir, "connector.log"),
		Stderr: filepath.Join(options.LogDir, "connector-error.log"),

		ServiceName:   options.ServiceName,
		ServiceType:   options.ServiceType,
		DisplayName:   options.DisplayName,
		Description:   options.Description,
		ConnectorName: options.ConnectorName,
		GithubSlug:    options.GithubSlug,
		Tag:           options.Tag,
		BinPath:       options.BinPath,
		Dir:           options.Dir,
	}
}
