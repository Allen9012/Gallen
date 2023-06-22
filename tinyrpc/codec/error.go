package codec

import "errors"

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/22
  @desc:
  @modified by:
**/

var (
	ErrorInvalidSequence        = errors.New("invalid sequence number in response")
	ErrorUnexpectedChecksum     = errors.New("unexpected checksum")
	ErrorNotFoundCompressor     = errors.New("not found compressor")
	ErrorCompressorTypeMismatch = errors.New("request and response Compressor type mismatch")
)
