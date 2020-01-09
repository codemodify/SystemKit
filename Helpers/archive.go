package Helpers

import (
	"archive/zip"
	"io"
	"os"
)

// ZipFiles - fully qualified path to where the new zip file
func ZipFiles(filename string, files []string, usePath bool) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file, usePath); err != nil {
			return err
		}
	}

	return nil
}

// AddFileToZip -
func AddFileToZip(zipWriter *zip.Writer, filename string, usePath bool) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	if usePath {
		header.Name = filename
	}
	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)

	return err
}
