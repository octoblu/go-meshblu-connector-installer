// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package unzipit allows you to easily unpack *.tar.gz, *.tar.bzip2, *.zip and *.tar files.
// There are not CGO involved nor hard dependencies of any type.
package unzipit

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	magicZIP  = []byte{0x50, 0x4b, 0x03, 0x04}
	magicGZ   = []byte{0x1f, 0x8b}
	magicBZIP = []byte{0x42, 0x5a}
	magicTAR  = []byte{0x75, 0x73, 0x74, 0x61, 0x72} // at offset 257
)

// Check whether a file has the magic number for tar, gzip, bzip2 or zip files
//
// Note that this function does not advance the Reader.
//
// 50 4b 03 04 for pkzip format
// 1f 8b for .gz
// 42 5a for bzip
// 75 73 74 61 72 at offset 257 for tar files
func magicNumber(reader *bufio.Reader, offset int) (string, error) {
	headerBytes, err := reader.Peek(offset + 5)
	if err != nil {
		return "", err
	}

	magic := headerBytes[offset : offset+5]

	if bytes.Equal(magicTAR, magic) {
		return "tar", nil
	}

	if bytes.Equal(magicZIP, magic[0:4]) {
		return "zip", nil
	}

	if bytes.Equal(magicGZ, magic[0:2]) {
		return "gzip", nil
	} else if bytes.Equal(magicBZIP, magic[0:2]) {
		return "bzip", nil
	}

	return "", nil
}

// Unpack unpacks a compressed and archived file and places result output in destination
// path.
//
// File formats supported are:
//   - .tar.gz
//   - .tar.bzip2
//   - .zip
//   - .tar
//
// If it cannot recognize the file format, it will save the file, as is, to the
// destination path.
func Unpack(file *os.File, destPath string) (string, error) {
	if file == nil {
		return "", errors.New("You must provide a valid file to unpack")
	}

	var err error
	if destPath == "" {
		destPath, err = ioutil.TempDir(os.TempDir(), "unpackit-")
		if err != nil {
			return "", err
		}
	}

	// Makes sure despPath exists
	if err := os.MkdirAll(destPath, 0740); err != nil {
		return "", err
	}

	r := bufio.NewReader(file)
	return UnpackStream(r, destPath)
}

// UnpackStream unpacks a compressed stream. Note that if the stream is a using ZIP
// compression (but only ZIP compression), it's going to get buffered in its entirety
// to memory prior to decompression.
func UnpackStream(reader io.Reader, destPath string) (string, error) {
	r := bufio.NewReader(reader)

	// Reads magic number from the stream so we can better determine how to proceed
	ftype, err := magicNumber(r, 0)
	if err != nil {
		return "", err
	}

	var decompressingReader *bufio.Reader
	switch ftype {
	case "gzip":
		decompressingReader, err = GunzipStream(r)
		if err != nil {
			return "", err
		}
	case "bzip":
		decompressingReader, err = Bunzip2Stream(r)
		if err != nil {
			return "", err
		}
	case "zip":
		// Like TAR, ZIP is also an archiving format, therefore we can just return
		// after it finishes
		return UnzipStream(r, destPath)
	default:
		decompressingReader = r
	}

	// Check magic number in offset 257 too see if this is also a TAR file
	ftype, err = magicNumber(decompressingReader, 257)
	if ftype == "tar" {
		return Untar(decompressingReader, destPath)
	}

	// If it's not a TAR archive then save it to disk as is.
	destRawFile := filepath.Join(destPath, sanitize(path.Base("tarstream")))

	// Creates destination file
	destFile, err := os.Create(destRawFile)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Copies data to destination file
	if _, err := io.Copy(destFile, decompressingReader); err != nil {
		return "", err
	}

	return destPath, nil
}

// Bunzip2 decompresses a bzip2 file and returns the decompressed stream
func Bunzip2(file *os.File) (*bufio.Reader, error) {
	freader := bufio.NewReader(file)
	bzip2Reader, err := Bunzip2Stream(freader)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(bzip2Reader), nil
}

// Bunzip2Stream unpacks a bzip2 stream
func Bunzip2Stream(reader io.Reader) (*bufio.Reader, error) {
	return bufio.NewReader(bzip2.NewReader(reader)), nil
}

// Gunzip decompresses a gzip file and returns the decompressed stream
func Gunzip(file *os.File) (*bufio.Reader, error) {
	freader := bufio.NewReader(file)
	gunzipReader, err := GunzipStream(freader)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(gunzipReader), nil
}

// GunzipStream unpacks a gzipped stream
func GunzipStream(reader io.Reader) (*bufio.Reader, error) {
	var decompressingReader *gzip.Reader
	var err error
	if decompressingReader, err = gzip.NewReader(reader); err != nil {
		return nil, err
	}

	return bufio.NewReader(decompressingReader), nil
}

// Unzip decompresses and unarchives a ZIP archive, returning the final path or an error
func Unzip(file *os.File, destPath string) (string, error) {
	fstat, err := file.Stat()
	if err != nil {
		return "", err
	}

	zr, err := zip.NewReader(file, fstat.Size())
	if err != nil {
		return "", err
	}

	return unpackZip(zr, destPath)
}

// UnzipStream unpacks a ZIP stream. Because of the nature of the ZIP format,
// the stream is copied to memory before decompression.
func UnzipStream(r io.Reader, destPath string) (string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	memReader := bytes.NewReader(data)
	zr, err := zip.NewReader(memReader, int64(len(data)))
	if err != nil {
		return "", err
	}

	return unpackZip(zr, destPath)
}

func unpackZip(zr *zip.Reader, destPath string) (string, error) {
	for _, f := range zr.File {
		err := unzipFile(f, destPath)
		if err != nil {
			return "", err
		}
	}
	return destPath, nil
}

func unzipFile(f *zip.File, destPath string) error {
	pathSeparator := string(os.PathSeparator)
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			log.Println(err)
		}
	}()

	filePath := sanitize(f.Name)

	// If the file is a subdirectory, it creates it before attempting to
	// create the actual file
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filepath.Join(destPath, filePath), f.Mode()); err != nil {
			return err
		}
		return nil
	}
	sepInd := strings.LastIndex(filePath, pathSeparator)
	if sepInd > -1 {
		directory := filePath[0:sepInd]
		if err := os.MkdirAll(filepath.Join(destPath, directory), 0740); err != nil {
			return err
		}
	}

	file, err := os.Create(filepath.Join(destPath, filePath))
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	file.Chmod(f.Mode())
	os.Chtimes(file.Name(), time.Now(), f.ModTime())

	if _, err := io.CopyN(file, rc, int64(f.UncompressedSize64)); err != nil {
		return err
	}

	return nil
}

// Untar unarchives a TAR archive and returns the final destination path or an error
func Untar(data io.Reader, destPath string) (string, error) {
	// Makes sure destPath exists
	if err := os.MkdirAll(destPath, 0740); err != nil {
		return "", err
	}

	tr := tar.NewReader(data)

	// Iterate through the files in the archive.
	rootdir := destPath
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}

		if err != nil {
			return rootdir, err
		}

		// Skip pax_global_header with the commit ID this archive was created from
		if hdr.Name == "pax_global_header" {
			continue
		}

		fp := filepath.Join(destPath, sanitize(hdr.Name))
		if hdr.FileInfo().IsDir() {
			if rootdir == destPath {
				rootdir = fp
			}
			if err := os.MkdirAll(fp, os.FileMode(hdr.Mode)); err != nil {
				return rootdir, err
			}
			continue
		}

		_, untarErr := untarFile(hdr, tr, fp, rootdir)
		if untarErr != nil {
			return rootdir, untarErr
		}
	}

	return rootdir, nil
}

func untarFile(hdr *tar.Header, tr *tar.Reader, fp, rootdir string) (string, error) {
	parentDir, _ := filepath.Split(fp)

	if err := os.MkdirAll(parentDir, 0740); err != nil {
		return rootdir, err
	}

	if hdr.Typeflag == tar.TypeSymlink {
		err := os.Symlink(hdr.Linkname, fp)
		return rootdir, err
	}

	file, err := os.Create(fp)
	if err != nil {
		return rootdir, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	file.Chmod(os.FileMode(hdr.Mode))
	os.Chtimes(file.Name(), time.Now(), hdr.ModTime)

	if _, err := io.Copy(file, tr); err != nil {
		return rootdir, err
	}

	return rootdir, nil
}

// Sanitizes name to avoid overwriting sensitive system files when unarchiving
func sanitize(name string) string {
	// Gets rid of volume drive label in Windows
	if len(name) > 1 && name[1] == ':' && runtime.GOOS == "windows" {
		name = name[2:]
	}

	name = filepath.Clean(name)
	name = filepath.ToSlash(name)
	for strings.HasPrefix(name, "../") {
		name = name[3:]
	}
	return name
}
