package installer

import (
	"io"
	"os"
)

// copyFile copies a file
func copyFile(source, target string) error {
	os.Remove(target)
	fileRead, err := os.Open(source)
	if err != nil {
		return err
	}
	defer fileRead.Close()

	fileWrite, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	_, err = io.Copy(fileWrite, fileRead)
	if err != nil {
		return err
	}
	return nil
}
