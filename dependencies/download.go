package dependencies

import (
	"fmt"
	"net/http"
	"net/url"
)

// Download downloads dependencies
func Download() {

}

// DownloadWithURLAndRuntime can be used to override the URL and runtime parameters
// when downloading. You probably want the much easier Download function instead
func DownloadWithURLAndRuntime(assemblerVersion, dependencyManagerVersion, urlStr string, runtime Runtime) error {
	var err error

	err = download(urlStr, fmt.Sprintf("/octoblu/go-meshblu-connector-assembler/releases/download/%v/meshblu-connector-assembler-%v-%v", assemblerVersion, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}
	err = download(urlStr, fmt.Sprintf("/octoblu/go-meshblu-connector-dependency-manager/releases/download/%v/meshblu-connector-dependency-manager-%v-%v", dependencyManagerVersion, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}

	return nil
}

func download(urlStr, path string) error {
	requestURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return err
	}

	requestURL.Path = path
	_, err = http.Get(requestURL.String())
	return err
}
