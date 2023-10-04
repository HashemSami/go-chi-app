package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrNotFound   = errors.New("models: resource could not be found")
	ErrEmailTaken = errors.New("models: email address is already in use")
)

// type to produce the error if the file extension doesn't match with the
// the extensions required by our app
type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

// this function will look at the first few bytes out of the file
// to determine the file type, and return error if the type is not required
func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	// ReadSeeker allows us to seek back to the starting point of the file
	// so when we pass back to the main function we can perform a full read
	// of that file
	testBytes := make([]byte, 512)
	// read up to the 512 bytes
	_, err := r.Read(testBytes)
	if err != nil {
		fmt.Errorf("checking content type: %w", err)
	}

	// return the file pointer back to the initial place
	_, err = r.Seek(0, 0)
	if err != nil {
		fmt.Errorf("checking content type: %w", err)
	}

	contentType := http.DetectContentType(testBytes)

	for _, t := range allowedTypes {
		if contentType == t {
			return nil
		}
	}

	// if the loop finishes and the content type has no match
	return FileError{
		Issue: fmt.Sprintf("invalid content type: %v", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if hasExtension(filename, allowedExtensions) {
		return nil
	}
	return FileError{
		Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename)),
	}
}
