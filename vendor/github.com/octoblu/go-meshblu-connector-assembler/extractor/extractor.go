package extractor

import (
	"io"
	"os"

	"github.com/octoblu/go-meshblu-connector-assembler/downloader"
	"github.com/octoblu/unzipit"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-assembler:extractor")

// Extractor ungzips and untars the source to the target
type Extractor interface {
	Do(source, target string) error
	DoWithBody(body io.ReadCloser, target string) error
	DoWithURI(downloadURI, target string) error
}

// Client interfaces with the Extractor
type Client struct {
}

// New constructs a new Extractor
func New() Extractor {
	return &Client{}
}

// Do extracts the a .gz, .tar, and/or .zip file
func (client *Client) Do(source, target string) error {
	debug("extracting file to target")
	file, err := os.Open(source)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = unzipit.Unpack(file, target)
	if err != nil {
		return err
	}
	return nil
}

// DoWithBody extracts the a .gz, .tar, and/or .zip from io.ReadCloser
func (client *Client) DoWithBody(body io.ReadCloser, target string) error {
	debug("extracting body of http response")
	defer body.Close()
	_, err := unzipit.UnpackStream(body, target)
	if err != nil {
		return err
	}
	return nil
}

// DoWithURI extracts the a .gz, .tar, and/or .zip from a url
func (client *Client) DoWithURI(downloadURI, target string) error {
	debug("downloading from uri and extracting")
	d := downloader.New(target)
	body, err := d.GetBody(downloadURI)
	if err != nil {
		return err
	}
	defer body.Close()
	_, err = unzipit.UnpackStream(body, target)
	if err != nil {
		return err
	}
	return nil
}
