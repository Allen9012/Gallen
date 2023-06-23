package compressor

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/23
  @desc:
  @modified by:
**/

// RawCompressor implements the Compressor interface
type RawCompressor struct {
}

// Zip .
func (_ RawCompressor) Zip(data []byte) ([]byte, error) {
	return data, nil
}

// Unzip .
func (_ RawCompressor) Unzip(data []byte) ([]byte, error) {
	return data, nil
}
