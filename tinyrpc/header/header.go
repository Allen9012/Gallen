/**
  @author: Allen
  @since: 2023/6/18
  @desc: //header
**/
package header

import (
	"encoding/binary"
	"errors"
	"github.com/Allen9012/tinyrpc/compressor"
	"sync"
)

const (
	// MaxHeaderSize = 2 + 10 + 10 + 10 + 4 (10 refer to binary.MaxVarintLen64)
	MaxHeaderSize = 36

	Uint32Size = 4
	Uint16Size = 2
)

var UnmarshalError = errors.New("an error occurred in Unmarshal")

// RequestHeader request header structure looks like:
// +--------------+----------------+----------+------------+----------+
// | CompressType |      Method    |    ID    | RequestLen | Checksum |
// +--------------+----------------+----------+------------+----------+
// |    uint16    | uvarint+string |  uvarint |   uvarint  |  uint32  |
// +--------------+----------------+----------+------------+----------+
type RequestHeader struct {
	sync.RWMutex
	CompressType compressor.CompressType //它表示RPC的协议内容的压缩类型，TinyRPC支持四种压缩类型，Raw、Gzip、Snappy、Zlib
	Method       string                  // 方法名
	ID           uint64                  // 请求ID
	RequestLen   uint32                  // 请求体长度
	Checksum     uint32                  //请求体校验 使用CRC32摘要算法
}

func (r *RequestHeader) Marshal() []byte {
	r.RLock()
	defer r.RUnlock()

	idx := 0
	header := make([]byte, MaxHeaderSize)
	// 小端模式写入压缩类型
	binary.LittleEndian.PutUint16(header[idx:], uint16(r.CompressType))
	idx += Uint16Size
	idx += writeString(header[idx:], r.Method)
	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += binary.PutUvarint(header[idx:], uint64(r.RequestLen))

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header[:idx]
}

func (r *RequestHeader) UnMarshal(data []byte) (err error) {
	r.RLock()
	defer r.RUnlock()
	if len(data) == 0 {
		return UnmarshalError
	}
	defer func() {
		if r := recover(); r != nil {
			err = UnmarshalError
		}
	}()
	idx, size := 0, 0
	r.CompressType = compressor.CompressType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size
	r.Method, size = readString(data[idx:])
	idx += size
	r.ID, size = binary.Uvarint(data[idx:])
	idx += size
	length, size := binary.Uvarint(data[idx:])
	r.RequestLen = uint32(length)
	idx += size

	r.Checksum = binary.LittleEndian.Uint32(data[idx:]) // 读取uvarint类型的校验码
	return
}

// GetCompressType get compress type
func (r *RequestHeader) GetCompressType() compressor.CompressType {
	r.Lock()
	defer r.Unlock()
	return compressor.CompressType(r.CompressType)
}

// ResponseHeader request header structure looks like:
// +--------------+---------+----------------+-------------+----------+
// | CompressType |    ID   |      Error     | ResponseLen | Checksum |
// +--------------+---------+----------------+-------------+----------+
// |    uint16    | uvarint | uvarint+string |    uvarint  |  uint32  |
// +--------------+---------+----------------+-------------+----------+
type ResponseHeader struct {
	sync.RWMutex
	CompressType compressor.CompressType // 压缩类型
	ID           uint64                  // 请求ID
	Error        string                  // 错误信息
	ResponseLen  uint32                  // 响应体长度
	Checksum     uint32                  // 响应体校验
}

func (r *ResponseHeader) Marshal() []byte {
	r.RLock()
	defer r.RUnlock()
	idx := 0
	// 多出一个ErrorMessage的长度
	header := make([]byte, MaxHeaderSize+len(r.Error))
	// 小端模式写入压缩类型
	binary.LittleEndian.PutUint16(header[idx:], uint16(r.CompressType))
	idx += Uint16Size

	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += writeString(header[idx:], r.Error)
	idx += binary.PutUvarint(header[idx:], uint64(r.ResponseLen))

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header[:idx]
}

func (r *ResponseHeader) UnMarshal(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return UnmarshalError
	}

	defer func() {
		if r := recover(); r != nil {
			err = UnmarshalError
		}
	}()
	idx, size := 0, 0
	r.CompressType = compressor.CompressType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size

	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	r.Error, size = readString(data[idx:])
	idx += size

	length, size := binary.Uvarint(data[idx:])
	r.ResponseLen = uint32(length)
	idx += size

	r.Checksum = binary.LittleEndian.Uint32(data[idx:])
	return
}

// GetCompressType get compress type
func (r *ResponseHeader) GetCompressType() compressor.CompressType {
	r.Lock()
	defer r.Unlock()
	return compressor.CompressType(r.CompressType)
}

// 读取Method的string
func readString(data []byte) (string, int) {
	idx := 0
	length, size := binary.Uvarint(data) // 读取一个uvarint类型表示字符的长度
	idx += size
	str := string(data[idx : idx+int(length)])
	idx += int(length)
	return str, idx
}

// 把Method转化为合适的byte字节数组（uvarint+string）
func writeString(data []byte, str string) int {
	idx := 0
	idx += binary.PutUvarint(data, uint64(len(str)))
	copy(data[idx:], str)
	idx += len(str)
	return idx
}
