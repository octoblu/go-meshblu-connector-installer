package assembler

import (
	"os"

	"github.com/octoblu/go-meshblu-connector-assembler/foreverizer"
	"github.com/octoblu/go-meshblu-connector-assembler/ignition"
	"github.com/octoblu/go-meshblu-connector-assembler/meshbluconfig"
	"github.com/octoblu/go-meshblu-connector-assembler/serviceconfig"
)
import De "github.com/tj/go-debug"

var debug = De.Debug("meshblu-connector-assembler:doitaller")

// Assemble does all the things a connector assembler does
// including:
// [x]  creating directories
// [x]  writing the meshblu config
// [x]  writing the service config
// [x]  installing the ignition
// [x]  foreverize
func Assemble(opts Options) error {
	var err error

	debug("createDirectories")
	err = createDirectories(opts)
	if err != nil {
		return err
	}

	debug("writeMeshbluConfig")
	err = writeMeshbluConfig(opts)
	if err != nil {
		return err
	}

	debug("writeServiceConfig")
	err = writeServiceConfig(opts)
	if err != nil {
		return err
	}

	debug("installIgnition")
	err = installIgnition(opts)
	if err != nil {
		return err
	}

	debug("foreverize")
	err = foreverize(opts)
	if err != nil {
		return err
	}

	return nil
}

func createDirectories(opts Options) error {
	var err error

	outputDir := opts.OutputDirectory
	logDir := opts.LogDirectory
	binDir := opts.BinDirectory

	debug("creating directories: %v", outputDir)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	debug("creating log directory: %v", logDir)
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	debug("creating bin directory: %v", binDir)
	err = os.MkdirAll(binDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func foreverize(opts Options) error {
	return foreverizer.Foreverize(foreverizer.Options{
		ServiceName:     opts.ServiceName,
		ServiceType:     opts.ServiceType,
		ServiceUsername: opts.ServiceUsername,
		ServicePassword: opts.ServicePassword,
		DisplayName:     opts.DisplayName,
		Description:     opts.Description,
		IgnitionPath:    opts.IgnitionPath,
	})
}

func installIgnition(opts Options) error {
	return ignition.Install(ignition.InstallOptions{
		IgnitionURL:  opts.IgnitionURL,
		IgnitionPath: opts.IgnitionPath,
	})
}

func writeMeshbluConfig(opts Options) error {
	return meshbluconfig.Write(meshbluconfig.Options{
		DirPath:    opts.ConnectorDirectory,
		UUID:       opts.UUID,
		Token:      opts.Token,
		Domain:     opts.Domain,
		ResolveSrv: opts.ResolveSrv,
	})
}

func writeServiceConfig(opts Options) error {
	return serviceconfig.Write(serviceconfig.Options{
		ServiceName:   opts.ServiceName,
		ServiceType:   opts.ServiceType,
		DisplayName:   opts.DisplayName,
		Description:   opts.Description,
		ConnectorName: opts.ConnectorName,
		GithubSlug:    opts.GithubSlug,
		Tag:           opts.Tag,
		BinPath:       opts.BinDirectory,
		Dir:           opts.ConnectorDirectory,
		LogDir:        opts.LogDirectory,
	})
}
