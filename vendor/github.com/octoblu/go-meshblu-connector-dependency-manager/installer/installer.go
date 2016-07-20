package installer

import (
	"log"
	"net/url"
	"runtime/debug"

	"github.com/octoblu/go-meshblu-connector-installer/osruntime"
	"github.com/spf13/afero"
	De "github.com/tj/go-debug"
)

var de = De.Debug("meshblu-connector-dependency-manager:installer")

// Installer interfaces with the a remote download server
type Installer interface {
	Install(depType, tag string) error
	// Do(depType, tag string) error
}

// Client interfaces with the a remote download server
type Client struct {
	baseURL   *url.URL
	fs        afero.Fs
	osruntime osruntime.OSRuntime
}

// New constructs a new Installer instance
func New() Installer {
	client, err := NewWithoutDefaults("https://nodejs.org", afero.NewOsFs(), osruntime.New())
	if err != nil {
		debug.PrintStack()
		log.Fatalln(err.Error())
	}

	return client
}

// NewWithoutDefaults constructs a new Installer instance
func NewWithoutDefaults(baseURLStr string, fs afero.Fs, osruntime osruntime.OSRuntime) (Installer, error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, err
	}
	return &Client{baseURL: baseURL, fs: fs, osruntime: osruntime}, nil
}

// Install the package and tagged version
func (client *Client) Install(depType, tag string) error {
	return nil
}

// Do download and install
// func (client *Client) Do(depType, tag string) error {
// 	de("installing dependency", depType, tag)
//
// 	target := getBinPath()
// 	exists, err := filePathExists(filepath.Join(target, finalDependencyFileName(depType, tag)))
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		de("dependency already exists")
// 		return nil
// 	}
// 	err = os.MkdirAll(target, 0777)
// 	if err != nil {
// 		return err
// 	}
//
// 	de("downloading %s", uri)
// 	downloadFile, err := download(uri, target)
// 	if err != nil {
// 		return err
// 	}
//
// 	de("extracting...", downloadFile, target)
// 	extractorClient := extractor.New()
// 	err = extractorClient.Do(downloadFile, target)
// 	if err != nil {
// 		return err
// 	}
//
// 	de("extracting bin...")
// 	err = extractBin(depType, target, tag)
// 	if err != nil {
// 		return err
// 	}
//
// 	de("done!")
// 	return nil
// }

// func getFileName(source string) (string, error) {
// 	uri, err := url.Parse(source)
// 	if err != nil {
// 		return "", err
// 	}
// 	segments := strings.Split(uri.Path, "/")
// 	return segments[len(segments)-1], nil
// }

// func download(uri, target string) (string, error) {
// 	fileExists, err := fileExists()
//
// 	fileName, err := getFileName(uri)
// 	if err != nil {
// 		return "", err
// 	}
// 	downloadFile := filepath.Join(target, fileName)
// 	outputStream, err := os.Create(downloadFile)
//
// 	if err != nil {
// 		de("Error on os.Create", err.Error())
// 		return "", err
// 	}
//
// 	defer outputStream.Close()
//
// 	response, err := http.Get(uri)
//
// 	if err != nil {
// 		de("Error on http.Get", err.Error())
// 		return "", err
// 	}
// 	defer response.Body.Close()
//
// 	if response.StatusCode != 200 {
// 		return "", fmt.Errorf("Download invalid status code: %v", response.StatusCode)
// 	}
//
// 	_, err = io.Copy(outputStream, response.Body)
//
// 	if err != nil {
// 		de("Error on io.Copy", err.Error())
// 		return "", err
// 	}
// 	return downloadFile, nil
// }
//
// // filePathExists check if a file exists
// func filePathExists(path string) (bool, error) {
// 	_, err := os.Stat(path)
//
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
//
// 	return true, nil
// }
