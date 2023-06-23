package compressor

import (
	"bytes"
	"github.com/golang/snappy"
	"io"
	"io/ioutil"
)

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/23
  @desc:
  @modified by:
**/

// SnappyCompressor implements the Compressor interface
type SnappyCompressor struct {
}

// Zip .
func (_ SnappyCompressor) Zip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := snappy.NewBufferedWriter(buf)
	defer func() {
		w.Close()
	}()
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Unzip .
func (_ SnappyCompressor) Unzip(data []byte) ([]byte, error) {
	r := snappy.NewReader(bytes.NewBuffer(data))
	data, err := ioutil.ReadAll(r)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
