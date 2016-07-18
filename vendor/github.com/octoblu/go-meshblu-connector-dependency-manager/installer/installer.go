package installer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/octoblu/go-meshblu-connector-dependency-manager/extractor"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-dependency-manager:installer")

// Installer interfaces with the a remote download server
type Installer interface {
	Do(depType, tag string) error
}

// Client interfaces with the a remote download server
type Client struct {
}

// New constructs a new Installer instance
func New() *Client {
	return &Client{}
}

// Do download and install
func (client *Client) Do(depType, tag string) error {
	debug("installing dependency", depType, tag)
	uri := getResourceURI(depType, tag)
	if uri == "" {
		return fmt.Errorf("unsupported platform")
	}

	target := getBinPath()
	exists, err := filePathExists(filepath.Join(target, finalDependencyFileName(depType, tag)))
	if err != nil {
		return err
	}
	if exists {
		debug("dependency already exists")
		return nil
	}
	err = os.MkdirAll(target, 0777)
	if err != nil {
		return err
	}

	debug("downloading %s", uri)
	downloadFile, err := download(uri, target)
	if err != nil {
		return err
	}

	debug("extracting...", downloadFile, target)
	extractorClient := extractor.New()
	err = extractorClient.Do(downloadFile, target)
	if err != nil {
		return err
	}

	debug("extracting bin...")
	err = extractBin(depType, target, tag)
	if err != nil {
		return err
	}

	debug("done!")
	return nil
}

func getFileName(source string) (string, error) {
	uri, err := url.Parse(source)
	if err != nil {
		return "", err
	}
	segments := strings.Split(uri.Path, "/")
	return segments[len(segments)-1], nil
}

func download(uri, target string) (string, error) {
	fileName, err := getFileName(uri)
	if err != nil {
		return "", err
	}
	downloadFile := filepath.Join(target, fileName)
	outputStream, err := os.Create(downloadFile)

	if err != nil {
		debug("Error on os.Create", err.Error())
		return "", err
	}

	defer outputStream.Close()

	response, err := http.Get(uri)

	if err != nil {
		debug("Error on http.Get", err.Error())
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Download invalid status code: %v", response.StatusCode)
	}

	_, err = io.Copy(outputStream, response.Body)

	if err != nil {
		debug("Error on io.Copy", err.Error())
		return "", err
	}
	return downloadFile, nil
}

// filePathExists check if a file exists
func filePathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
