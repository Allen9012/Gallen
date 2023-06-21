package codec

import (
	"github/Allen9012/tinyrpc/header"
	"github/Allen9012/tinyrpc/serializer"
	"io"
	"sync"
)

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/21
  @desc:
  @modified by:
**/

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	compressor compressor.CompressType // rpc compress type(raw,gzip,snappy,zlib)
	serializer serializer.Serializer
	response   header.ResponseHeader // rpc response header
	mutex      sync.Mutex            // protect pending map
	pending    map[uint64]string
}
