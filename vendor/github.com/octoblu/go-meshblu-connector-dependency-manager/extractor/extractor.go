package extractor

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor ungzips and untars the source to the target
type Extractor interface {
	Do(source, target string) error
	Ungzip(source, target string) (string, error)
	Unzip(source, target string) error
	Untar(tarball, target string) error
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
	var filesToRemove []string
	var err error
	if isGZ(source) {
		filesToRemove = append(filesToRemove, source)
		source, err = client.Ungzip(source, target)
		if err != nil {
			return err
		}
	}
	if isTar(source) {
		filesToRemove = append(filesToRemove, source)
		err = client.Untar(source, target)
		if err != nil {
			return err
		}
	}
	if isZip(source) {
		filesToRemove = append(filesToRemove, source)
		err = client.Unzip(source, target)
		if err != nil {
			return err
		}
	}
	for _, fileToRemove := range filesToRemove {
		err = os.Remove(fileToRemove)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ungzip the source to the target
func (client *Client) Ungzip(source, target string) (string, error) {
	reader, err := os.Open(source)
	if err != nil {
		return "", fmt.Errorf("Ungzip, os.Open: %v", err.Error())
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return "", fmt.Errorf("Ungzip, gzip.NewReader: %v", err.Error())
	}
	defer archive.Close()

	filename := filepath.Base(source)
	target = filepath.Join(target, strings.Replace(filename, ".gz", "", 1))
	writer, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("Ungzip, os.Create: %v", err.Error())
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)

	return target, err
}

// Unzip extracts a zip file
func (client *Client) Unzip(source, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(target, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(file *zip.File) error {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(target, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, file := range reader.File {
		err := extractAndWriteFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// Untar the source to the target
func (client *Client) Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return fmt.Errorf("Untar, os.Open: %v", err.Error())
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Untar, tarReader.Next: %v", err.Error())
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return fmt.Errorf("Untar, os.MkdirAll: %v", err.Error())
			}
			continue
		}

		if header.Typeflag == tar.TypeSymlink {
			os.Remove(path)
			err = os.Symlink(header.Linkname, path)
			if err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("Untar, os.OpenFile: %v", err.Error())
		}

		_, err = io.Copy(file, tarReader)
		if err != nil {
			file.Close()
			return fmt.Errorf("Untar, io.Copy: %v", err.Error())
		}
		file.Close()
	}
	return nil
}

func getFileName(source string) string {
	segments := strings.Split(source, "/")
	return segments[len(segments)-1]
}

func isGZ(source string) bool {
	fileName := getFileName(source)
	return strings.Index(fileName, ".gz") > -1
}

func isTar(source string) bool {
	fileName := getFileName(source)
	return strings.Index(fileName, ".tar") > -1
}

func isZip(source string) bool {
	fileName := getFileName(source)
	return strings.Index(fileName, ".zip") > -1
}
