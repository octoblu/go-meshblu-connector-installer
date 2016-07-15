package dependencydownloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/spf13/afero"

	"github.com/octoblu/go-meshblu-connector-installer/osruntime"
)

// GithubURL is the default url used for the Downloader.
// the use of this constant can be overriden by using the
// DownloadWithURLAndRuntime funtion
const GithubURL = "https://github.com"

// Downloader downloads dependencies
type Downloader interface {
	// Downloader downloads the assembler and connector manager and stores them on the file system
	Download() error
}

type httpDownloader struct {
	assemblerVersion, dependencyManagerVersion string
	baseURL                                    *url.URL
	runtime                                    osruntime.OSRuntime
	fs                                         afero.Fs
}

// Download constructs a new downloader and downloads in one go
func Download(assemblerVersion, dependencyManagerVersion string) error {
	downloader := New(assemblerVersion, dependencyManagerVersion)
	return downloader.Download()
}

// New constructs a new Downloader
func New(assemblerVersion, dependencyManagerVersion string) Downloader {
	downloader, err := NewWithoutDefaults(assemblerVersion, dependencyManagerVersion, GithubURL, osruntime.New(), afero.NewOsFs())
	if err != nil {
		log.Fatalln("This URL should never be invalid, but it is: ", err.Error())
	}

	return downloader
}

// NewWithoutDefaults constructs a new Downloader with provided url, runtime, and filesystem
func NewWithoutDefaults(assemblerVersion, dependencyManagerVersion, urlStr string, runtime osruntime.OSRuntime, fs afero.Fs) (Downloader, error) {
	baseURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	downloader := &httpDownloader{
		assemblerVersion:         assemblerVersion,
		dependencyManagerVersion: dependencyManagerVersion,
		baseURL:                  baseURL,
		runtime:                  runtime,
		fs:                       fs,
	}
	return downloader, nil
}

// Downloader downloads the assembler and connector manager and stores them on the file system
func (downloader *httpDownloader) Download() error {
	var err error

	err = downloader.download("assembler", "assembler-installer", downloader.assemblerVersion)
	if err != nil {
		return err
	}
	err = downloader.download("dependency-manager", "dependency-manager", downloader.dependencyManagerVersion)
	if err != nil {
		return err
	}

	return nil
}

// download retrieves the file and stores it on the file system
func (downloader *httpDownloader) download(remoteName, fileName, assemblerVersion string) error {
	runtime := downloader.runtime
	remotePath := fmt.Sprintf("/octoblu/go-meshblu-connector-%v/releases/download/%v/meshblu-connector-%v-%v-%v", remoteName, assemblerVersion, remoteName, runtime.GOOS, runtime.GOARCH)

	requestURL := *downloader.baseURL
	requestURL.Path = remotePath

	response, err := http.Get(requestURL.String())
	if err != nil {
		return err
	}

	localPath := path.Join("/Users/octoblu/.octoblu/MeshbluConnectors/bin/", fileName)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return afero.WriteFile(downloader.fs, localPath, body, 0755)
}
