package header

import "sync"

/**
  @author: Allen
  @since: 2023/6/19
  @desc: pool
**/

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{
		New: func() any {
			return &RequestHeader{}
		}}

	ResponsePool = sync.Pool{
		New: func() any {
			return &ResponseHeader{}
		}}
}

// ResetHeader reset request header
func (r *RequestHeader) ResetHeader() {
	r.CompressType = 0
	r.Method = ""
	r.ID = 0
	r.RequestLen = 0
	r.Checksum = 0
}

// ResetHeader reset response header
func (r *ResponseHeader) ResetHeader() {
	r.CompressType = 0
	r.ID = 0
	r.ResponseLen = 0
	r.Error = ""
	r.Checksum = 0
}
