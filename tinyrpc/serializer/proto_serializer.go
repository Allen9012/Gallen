package serializer

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/23
  @desc:
  @modified by:
**/

// ErrorNotImplementProtoMessage refers to param not implemented by proto.Message
var ErrorNotImplementProtoMessage = errors.New("param does not implement proto.Message")

var _ Serializer = (*ProtoSerializer)(nil)

var Proto = ProtoSerializer{}

// ProtoSerializer implements the Serializer interface
type ProtoSerializer struct {
}

// Marshal returns the protobuf encoding of message
func (_ ProtoSerializer) Marshal(message interface{}) ([]byte, error) {
	var body proto.Message
	if message == nil {
		return []byte{}, nil
	}
	var ok bool
	if body, ok = message.(proto.Message); !ok {
		return nil, ErrorNotImplementProtoMessage
	}
	return proto.Marshal(body)
}

// Unmarshal parses the protobuf-encoded data and stores the result
func (_ ProtoSerializer) Unmarshal(data []byte, message interface{}) error {
	var body proto.Message
	if message == nil {
		return nil
	}
	var ok bool
	if body, ok = message.(proto.Message); !ok {
		return ErrorNotImplementProtoMessage
	}
	return proto.Unmarshal(data, body)
}
