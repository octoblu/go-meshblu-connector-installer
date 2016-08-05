package assembler

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// OptionsOptions defines the service configurations
type OptionsOptions struct {
	IgnitionTag     string
	ConnectorName   string
	GithubSlug      string
	Tag             string
	OutputDirectory string
	ServiceName     string
	Hostname        string
	Port            int
	UUID, Token     string
	Debug           bool
	ServiceType     string
	ServiceUsername string
	ServicePassword string
}

// Options define the options that the Assemble takes
type Options struct {
	// IgnitionTag     string
	ConnectorName string
	GithubSlug    string
	Tag           string

	IgnitionURL  string
	IgnitionPath string

	OutputDirectory    string
	ConnectorDirectory string
	BinDirectory       string
	LogDirectory       string

	ServiceName string
	DisplayName string
	Description string

	Hostname    string
	Port        int
	UUID, Token string
	// Debug           bool

	ServiceType     string
	ServiceUsername string
	ServicePassword string
}

// NewOptions constructs a new Options object using OptionsOptinos
func NewOptions(opts OptionsOptions) (*Options, error) {
	if opts.ConnectorName == "" {
		return nil, fmt.Errorf("Missing required opt: opts.ConnectorName")
	}
	if opts.GithubSlug == "" {
		return nil, fmt.Errorf("Missing required opt: opts.GithubSlug")
	}
	if opts.Tag == "" {
		return nil, fmt.Errorf("Missing required opt: opts.Tag")
	}
	if opts.UUID == "" {
		return nil, fmt.Errorf("Missing required opt: opts.UUID")
	}
	if opts.Token == "" {
		return nil, fmt.Errorf("Missing required opt: opts.Token")
	}
	if opts.IgnitionTag == "" {
		return nil, fmt.Errorf("Missing required opt: opts.IgnitionTag")
	}
	if opts.ServiceType == "" {
		return nil, fmt.Errorf("Missing required opt: opts.ServiceType")
	}
	if opts.ServiceType == "UserService" {
		if opts.ServiceUsername == "" {
			return nil, fmt.Errorf("Missing required opt: opts.ServiceUsername")
		}
	}

	OutputDirectory := opts.OutputDirectory
	if OutputDirectory == "" {
		OutputDirectory = getUserServiceDirectory()
		if opts.ServiceType == "Service" {
			OutputDirectory = getServiceDirectory()
		}
	}

	OutputDirectory, err := filepath.Abs(OutputDirectory)
	if err != nil {
		return nil, err
	}
	ConnectorDirectory := getConnectorDirectory(OutputDirectory, opts.UUID)

	return &Options{
		ConnectorDirectory: ConnectorDirectory,
		OutputDirectory:    OutputDirectory,
		Hostname:           getHostname(opts.Hostname),
		Port:               getPort(opts.Port),
		IgnitionURL:        getIgnitionURI(opts.IgnitionTag, runtime.GOOS, runtime.GOARCH),
		IgnitionPath:       getIgnitionPath(ConnectorDirectory),
		LogDirectory:       getLogDirectory(ConnectorDirectory),
		BinDirectory:       getBinDirectory(OutputDirectory),

		ServiceName: getServiceName(opts.UUID),
		DisplayName: getDisplayName(opts.UUID),
		Description: getDescription(opts.ConnectorName, opts.UUID),

		ConnectorName: opts.ConnectorName,
		GithubSlug:    opts.GithubSlug,
		Tag:           opts.Tag,
		UUID:          opts.UUID,
		Token:         opts.Token,

		ServiceType:     opts.ServiceType,
		ServiceUsername: opts.ServiceUsername,
		ServicePassword: opts.ServicePassword,
	}, nil
}

func getBinDirectory(outputDirectory string) string {
	return filepath.Join(outputDirectory, "bin")
}

func getConnectorDirectory(outputDirectory, uuid string) string {
	return filepath.Join(outputDirectory, uuid)
}

func getDescription(connector, uuid string) string {
	return fmt.Sprintf("MeshbluConnector (%v) %v", connector, uuid)
}

func getDisplayName(uuid string) string {
	return fmt.Sprintf("MeshbluConnector %s", uuid)
}

func getDownloadURL(connector, tag, githubSlug, goos, goarch string) string {
	baseURI := fmt.Sprintf("https://github.com/%s/releases/download", githubSlug)

	extension := "tar.gz"
	if runtime.GOOS == "windows" {
		extension = "zip"
	}

	fileName := fmt.Sprintf("%s-%s-%s.%s", connector, goos, goarch, extension)
	return fmt.Sprintf("%s/%s/%s", baseURI, tag, fileName)
}

func getHostname(hostname string) string {
	if hostname != "" {
		return hostname
	}

	return "meshblu.octoblu.com"
}

func getIgnitionURI(tag, goos, goarch string) string {
	baseURI := "https://github.com/octoblu/go-meshblu-connector-ignition/releases/download"
	return fmt.Sprintf("%s/%s/meshblu-connector-ignition-%s-%s", baseURI, tag, goos, goarch)
}

func getLogDirectory(connectorDirectory string) string {
	return filepath.Join(connectorDirectory, "log")
}

func getPort(port int) int {
	if port != 0 {
		return port
	}

	return 443
}
