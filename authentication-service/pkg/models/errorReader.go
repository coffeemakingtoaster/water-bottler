package models

import "errors"

type ErrorReader struct{}

func (ErrorReader) Read(p []byte) (int, error) {
	return 0, errors.New("custom error from reader for testing")
}
