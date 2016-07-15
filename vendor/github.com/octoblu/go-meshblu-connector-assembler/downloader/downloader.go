package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-assembler:downloader")

// Downloader interface with a way of downloading connector bundles
type Downloader interface {
	Download(downloadURI string) (string, error)
	GetBody(downloadURI string) (io.ReadCloser, error)
}

// Client interfaces with remote cdn
type Client struct {
	OutputDirectory string
}

// New constructs new Downloader instance
func New(OutputDirectory string) Downloader {
	return &Client{OutputDirectory}
}

// Download downloads the connector the local directory
func (client *Client) Download(downloadURI string) (string, error) {
	debug("downloading: %v", downloadURI)

	downloadFile := client.getDownloadFile(downloadURI)
	debug("to: %v", downloadFile)
	outputStream, err := os.Create(downloadFile)

	if err != nil {
		return "", fmt.Errorf("error opening output stream: %v", err.Error())
	}

	defer outputStream.Close()

	body, err := client.GetBody(downloadURI)
	if err != nil {
		return "", err
	}
	defer body.Close()

	_, err = io.Copy(outputStream, body)

	if err != nil {
		return "", fmt.Errorf("error downloading to file %v", err.Error())
	}

	debug("downloaded!")

	return downloadFile, nil
}

// GetBody gets the res.Body from download uri
func (client *Client) GetBody(downloadURI string) (io.ReadCloser, error) {
	debug("downloading: %v", downloadURI)
	response, err := http.Get(downloadURI)

	if err != nil {
		return nil, fmt.Errorf("http error downloading: %v", err.Error())
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("download returned invalid response code: %v", response.StatusCode)
	}

	debug("successful retrieved body")

	return response.Body, nil
}

func (client *Client) getDownloadFile(downloadURI string) string {
	fileName := getFileName(downloadURI)
	return filepath.Join(client.OutputDirectory, fileName)
}

func getFileName(source string) string {
	segments := strings.Split(source, "/")
	return segments[len(segments)-1]
}
