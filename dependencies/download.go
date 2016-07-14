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
func DownloadWithURLAndRuntime(version, urlStr string, runtime Runtime) error {
	requestURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return err
	}

	requestURL.Path = fmt.Sprintf("/octoblu/go-fu/releases/download/%v/fu-%v-%v", version, runtime.GOOS, runtime.GOARCH)
	_, err = http.Get(requestURL.String())
	return err
}
