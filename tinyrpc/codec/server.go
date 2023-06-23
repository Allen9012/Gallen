package codec

import (
	"bufio"
	"github.com/Allen9012/tinyrpc/compressor"
	"github.com/Allen9012/tinyrpc/header"
	"github.com/Allen9012/tinyrpc/serializer"

	"hash/crc32"
	"io"
	"net/rpc"
	"sync"
)

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/22
  @desc:
  @modified by:
**/

// serverCodec is the implementation of rpc.ServerCodec
var _ rpc.ServerCodec = (*serverCodec)(nil)

type reqCtx struct {
	requestID    uint64
	compressType compressor.CompressType
}

type serverCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	request    header.RequestHeader
	serializer serializer.Serializer
	mutex      sync.Mutex // protects seq, pending
	seq        uint64
	pending    map[uint64]*reqCtx
}

// NewServerCodec Create a new server codec
func NewServerCodec(conn io.ReadWriteCloser, serializer serializer.Serializer) rpc.ServerCodec {
	return &serverCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		serializer: serializer,
		pending:    make(map[uint64]*reqCtx),
	}
}

// ReadRequestHeader read the rpc request header from the io stream
func (s *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	s.request.ResetHeader()     // 重置响应头
	data, err := recvFrame(s.r) // 读取请求头字符窜
	if err != nil {
		return err
	}
	err = s.request.UnMarshal(data) // 将字节串反序列化
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.seq++                                                                                        // 递增序列号
	s.pending[s.seq] = &reqCtx{requestID: s.request.ID, compressType: s.request.GetCompressType()} // 自增序号与请求头部的ID进行绑定
	r.ServiceMethod = s.request.Method
	r.Seq = s.seq
	s.mutex.Unlock()
	return nil
}

// ReadRequestBody read the rpc request body from the io stream
func (s *serverCodec) ReadRequestBody(param any) error {
	if param == nil {
		if s.request.RequestLen != 0 { // 废除多余部分
			if err := read(s.r, make([]byte, s.request.RequestLen)); err != nil {
				return err
			}
		}
		return nil
	}
	//	根据长度，读取该长度的字节串
	reqBody := make([]byte, s.request.RequestLen)
	if err := read(s.r, reqBody); err != nil {
		return err
	}
	//	校验
	if s.request.Checksum != 0 {
		if crc32.ChecksumIEEE(reqBody) != s.request.Checksum {
			return ErrorUnexpectedChecksum
		}
	}
	// 判断压缩器是否存在
	if _, ok := compressor.Compressors[s.request.GetCompressType()]; !ok {
		return ErrorNotFoundCompressor
	}
	// 解压
	req, err := compressor.Compressors[s.request.GetCompressType()].Unzip(reqBody)
	if err != nil {
		return err
	}
	// 反序列化
	return s.serializer.Unmarshal(req, param)
}

// WriteResponse Write the rpc response header and body to the io stream
func (s *serverCodec) WriteResponse(r *rpc.Response, param any) error {
	// 排除取不到序列号
	s.mutex.Lock()
	reqCtx, ok := s.pending[r.Seq]
	if !ok {
		s.mutex.Unlock()
		return ErrorInvalidSequence
	}
	delete(s.pending, r.Seq)
	s.mutex.Unlock()
	// 判断有没有Error
	if r.Error != "" {
		param = nil
	}
	if _, ok := compressor.
		Compressors[reqCtx.compressType]; !ok {
		return ErrorNotFoundCompressor
	}
	// 响应头返回
	var respBody []byte
	var err error
	if param != nil {
		respBody, err = s.serializer.Marshal(param)
		if err != nil {
			return err
		}
	}

	compressedRespBody, err := compressor.
		Compressors[reqCtx.compressType].Zip(respBody)
	if err != nil {
		return err
	}
	h := header.ResponsePool.Get().(*header.ResponseHeader)
	defer func() {
		h.ResetHeader()
		header.ResponsePool.Put(h)
	}()
	h.ID = reqCtx.requestID
	h.Error = r.Error
	h.ResponseLen = uint32(len(compressedRespBody))
	h.Checksum = crc32.ChecksumIEEE(compressedRespBody)
	h.CompressType = reqCtx.compressType

	if err = sendFrame(s.w, h.Marshal()); err != nil {
		return err
	}

	if err = write(s.w, compressedRespBody); err != nil {
		return err
	}
	s.w.(*bufio.Writer).Flush()
	return nil
}

func (s *serverCodec) Close() error {
	return s.c.Close()
}
