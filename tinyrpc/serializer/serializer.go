package serializer

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/21
  @desc:
  @modified by:
**/

// Serializer is interface, each serializer has Marshal and Unmarshal functions
type Serializer interface {
	Marshal(message interface{}) ([]byte, error)
	Unmarshal(data []byte, message interface{}) error
}
