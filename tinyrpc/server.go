package tinyrpc

import (
	"github.com/Allen9012/tinyrpc/codec"
	"github.com/Allen9012/tinyrpc/serializer"
	"log"
	"net"
	"net/rpc"
)

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/23
  @desc:
  @modified by:
**/

// Server rpc server based on net/rpc implementation
type Server struct {
	*rpc.Server
	serializer.Serializer
}

// NewServer Create a new rpc server
func NewServer(opts ...Option) *Server {
	options := options{
		serializer: serializer.Proto,
	}
	return &Server{
		Server:     rpc.NewServer(),
		Serializer: options.serializer,
	}
}

//Register register rpc function
func (s *Server) Register(rcvr interface{}) error {
	return s.Server.Register(rcvr)
}

//RegisterName register the rpc function with the specified name
func (s *Server) RegisterName(name string, rcvr interface{}) error {
	return s.Server.RegisterName(name, rcvr)
}

//Serve start service
func (s *Server) Serve(lis net.Listener) {
	log.Printf("tinyrpc server start on: %s\n", lis.Addr().String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		go s.Server.ServeCodec(codec.NewServerCodec(conn, s.Serializer))
	}
}
