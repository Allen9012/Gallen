package json

import "errors"

/**
  Copyright © 2023 github.com/Allen9012 All rights reserved.
  @author: Allen
  @since: 2023/6/23
  @desc:
  @modified by:
**/

// Request .
type Request struct {
	A float64 `json:"a,omitempty"`
	B float64 `json:"b,omitempty"`
}

// Response .
type Response struct {
	C float64 `json:"c,omitempty"`
}

// TestService Defining Computational Digital Services
type TestService struct{}

// Add addition
func (this *TestService) Add(args *Request, reply *Response) error {
	reply.C = args.A + args.B
	return nil
}

// Sub subtraction
func (this *TestService) Sub(args *Request, reply *Response) error {
	reply.C = args.A - args.B
	return nil
}

// Mul multiplication
func (this *TestService) Mul(args *Request, reply *Response) error {
	reply.C = args.A * args.B
	return nil
}

// Div division
func (this *TestService) Div(args *Request, reply *Response) error {
	if args.B == 0 {
		return errors.New("divided is zero")
	}
	reply.C = args.A / args.B
	return nil
}
