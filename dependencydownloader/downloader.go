package dependencydownloader

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/octoblu/go-meshblu-connector-installer/osruntime"
)

// GithubURL is the default url used for the Downloader.
// the use of this constant can be overriden by using the
// DownloadWithURLAndRuntime funtion
const GithubURL = "https://github.com"

// Downloader downloads dependencies
type Downloader interface {
	Download() error
}

type httpDownloader struct {
	assemblerVersion, dependencyManagerVersion string
	baseURL                                    *url.URL
	runtime                                    osruntime.OSRuntime
}

// Download constructs a new downloader and downloads in one go
func Download(assemblerVersion, dependencyManagerVersion string) error {
	// downloader := New(assemblerVersion, dependencyManagerVersion)
	// return downloader.Download()
	return nil
}

// New constructs a new Downloader
func New(assemblerVersion, dependencyManagerVersion string) Downloader {
	downloader, err := NewWithoutDefaults(assemblerVersion, dependencyManagerVersion, GithubURL, osruntime.New())
	if err != nil {
		log.Fatalln("This URL should never be invalid, but it is: ", err.Error())
	}

	return downloader
}

// NewWithoutDefaults constructs a new Downloader with provided url, runtime, and filesystem
func NewWithoutDefaults(assemblerVersion, dependencyManagerVersion, urlStr string, runtime osruntime.OSRuntime) (Downloader, error) {
	baseURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	downloader := &httpDownloader{
		assemblerVersion:         assemblerVersion,
		dependencyManagerVersion: dependencyManagerVersion,
		baseURL:                  baseURL,
		runtime:                  runtime,
	}
	return downloader, nil
}

func (downloader *httpDownloader) Download() error {
	var err error

	assemblerVersion := downloader.assemblerVersion
	dependencyManagerVersion := downloader.dependencyManagerVersion
	runtime := downloader.runtime

	err = downloader.download(fmt.Sprintf("/octoblu/go-meshblu-connector-assembler/releases/download/%v/meshblu-connector-assembler-%v-%v", assemblerVersion, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}
	err = downloader.download(fmt.Sprintf("/octoblu/go-meshblu-connector-dependency-manager/releases/download/%v/meshblu-connector-dependency-manager-%v-%v", dependencyManagerVersion, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}

	return nil
}

// download retrieves the file
func (downloader *httpDownloader) download(path string) error {
	requestURL := *downloader.baseURL
	requestURL.Path = path

	_, err := http.Get(requestURL.String())
	return err
}
