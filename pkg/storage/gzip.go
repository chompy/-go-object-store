package storage

import (
	"bytes"
	"compress/gzip"

	"github.com/pkg/errors"
)

func compressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := zw.Close(); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func uncompressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := zw.Close(); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}
