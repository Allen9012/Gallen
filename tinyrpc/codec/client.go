package codec

import (
	"bufio"
	"github/Allen9012/tinyrpc/compressor"
	"github/Allen9012/tinyrpc/header"
	"github/Allen9012/tinyrpc/serializer"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"
)

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/21
  @desc:
  @modified by:
**/

// clientCodec is a rpc client codec
var _ rpc.ClientCodec = (*clientCodec)(nil)

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

func NewClientCodec(conn io.ReadWriteCloser, compressor compressor.CompressType, serializer serializer.Serializer) rpc.ClientCodec {
	return &clientCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		compressor: compressor,
		serializer: serializer,
		pending:    make(map[uint64]string),
	}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//	判断压缩器是否存在
	if _, ok := compressor.Compressors[c.compressor]; !ok {
		return ErrorNotFoundCompressor
	}
	// 序列化器编码
	reqBody, err := c.serializer.Marshal(param)
	if err != nil {
		return err
	}
	//	压缩
	compressedReqBody, err := compressor.Compressors[c.compressor].Zip(reqBody)
	if err != nil {
		return err
	}
	// 从请求头部对象池取出请求头
	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
	//	设置请求头
	h.ID = r.Seq
	h.Method = r.ServiceMethod
	h.RequestLen = uint32(len(compressedReqBody))
	h.CompressType = c.compressor
	h.Checksum = crc32.ChecksumIEEE(compressedReqBody)
	//	发送请求头
	if err := sendFrame(c.w, h.Marshal()); err != nil {
		return err
	}
	// 发送请求体
	if err := write(c.w, compressedReqBody); err != nil {
		return err
	}
	//	刷新缓冲区
	c.w.(*bufio.Writer).Flush()
	return nil
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.response.ResetHeader()    // 重置响应头
	data, err := recvFrame(c.r) // 读取请求头字符窜
	if err != nil {
		return err
	}
	err = c.response.UnMarshal(data) // 序列化器继续解码
	if err != nil {
		return err
	}
	c.mutex.Lock()
	r.Seq = c.response.ID
	r.ServiceMethod = c.pending[r.Seq] // 从pending map中取出对应的方法名
	r.Error = c.response.Error         //	设置错误信息
	delete(c.pending, r.Seq)
	c.mutex.Unlock()
	return nil
}

func (c *clientCodec) ReadResponseBody(param any) error {
	if param == nil {
		if c.response.ResponseLen != 0 { // 废除多余部分
			if err := read(c.r, make([]byte, c.response.ResponseLen)); err != nil {
				return err
			}
		}
		return nil
	}
	//	根据长度，读取该长度的字节串
	respBody := make([]byte, c.response.ResponseLen)
	if err := read(c.r, respBody); err != nil {
		return err
	}
	//	校验
	if c.response.Checksum != 0 {
		if crc32.ChecksumIEEE(respBody) != c.response.Checksum {
			return ErrorUnexpectedChecksum
		}
	}
	// 判断压缩器是否存在
	if _, ok := compressor.Compressors[c.response.GetCompressType()]; !ok {
		return ErrorNotFoundCompressor
	}
	// 解压
	resp, err := compressor.Compressors[c.response.GetCompressType()].Unzip(respBody)
	if err != nil {
		return err
	}
	// 反序列化
	return c.serializer.Unmarshal(resp, param)
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}
