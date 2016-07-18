package ignition

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/spf13/afero"
)

import De "github.com/tj/go-debug"

var debug = De.Debug("meshblu-connector-assembler:install")

// InstallOptions are options for the Install function
type InstallOptions struct {
	IgnitionURL, IgnitionPath string
}

// Install downloads the ignition script and
// installs it into the correct place
func Install(options InstallOptions) error {
	return InstallWithoutDefaults(options, afero.NewOsFs())
}

// InstallWithoutDefaults downloads the ignition script and
// installs it into the correct place on the file
// system specified
func InstallWithoutDefaults(options InstallOptions, fs afero.Fs) error {
	debug("Downloading ignition")
	response, err := http.Get(options.IgnitionURL)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("Ignition download expected 200, Received: %v", response.StatusCode)
	}

	debug("Reading response.Body")
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	debug("Making the directory")
	err = fs.MkdirAll(filepath.Base(options.IgnitionPath), 0755)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, options.IgnitionPath, data, 0755)
}
