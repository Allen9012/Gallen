/**
  @author: Allen
  @since: 2023/3/25
  @desc: //gob和Json的接口正好一样
**/
package codec

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
)

/* GobCodec */

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g *GobCodec) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

func (g *GobCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()
	// 1. 先写入header
	if err = g.enc.Encode(header); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	// 2. 再写入body
	if err = g.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

/*	*JsonCodec */

var _ Codec = (*JsonCodec)(nil)

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

func (j *JsonCodec) Close() error {
	return j.conn.Close()
}

func (j *JsonCodec) ReadHeader(header *Header) error {
	return j.dec.Decode(header)
}

func (j *JsonCodec) ReadBody(body interface{}) error {
	return j.dec.Decode(body)
}

func (j *JsonCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		// 该方法的作用是强制刷新缓冲区，确保所有的数据都被写入目标。
		_ = j.buf.Flush()
		if err != nil {
			_ = j.Close()
		}
	}()
	if err := j.enc.Encode(header); err != nil {
		log.Println("rpc: json error encoding header:", err)
		return err
	}
	if err := j.enc.Encode(body); err != nil {
		log.Println("rpc: json error encoding body:", err)
		return err
	}
	return nil
}
