// Copyright 2015 The go-pbfcoin Authors
// This file is part of the go-pbfcoin library.
//
// The go-pbfcoin library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pbfcoin library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pbfcoin library. If not, see <http://www.gnu.org/licenses/>.

// Package rpc implements the pbfcoin JSON-RPC API.
package rpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/pbfcoin/go-pbfcoin/rpc/comms"
	"github.com/pbfcoin/go-pbfcoin/rpc/shared"
)

// Xpbf is a native API interface to a remote node.
type Xpbf struct {
	client comms.pbfcoinClient
	reqId  uint32
}

// NewXpbf constructs a new native API interface to a remote node.
func NewXpbf(client comms.pbfcoinClient) *Xpbf {
	return &Xpbf{
		client: client,
	}
}

// Call invokes a method with the given parameters are the remote node.
func (self *Xpbf) Call(method string, params []interface{}) (map[string]interface{}, error) {
	// Assemble the json RPC request
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req := &shared.Request{
		Id:      atomic.AddUint32(&self.reqId, 1),
		Jsonrpc: "2.0",
		method:  method,
		Params:  data,
	}
	// Send the request over and retrieve the response
	if err := self.client.Send(req); err != nil {
		return nil, err
	}
	res, err := self.client.Recv()
	if err != nil {
		return nil, err
	}
	// Ensure the response is valid, and extract the results
	success, isSuccessResponse := res.(*shared.SuccessResponse)
	failure, isFailureResponse := res.(*shared.ErrorResponse)
	switch {
	case isFailureResponse:
		return nil, fmt.Errorf("method invocation failed: %v", failure.Error)

	case isSuccessResponse:
		return success.Result.(map[string]interface{}), nil

	default:
		return nil, fmt.Errorf("Invalid response type: %v", reflect.TypeOf(res))
	}
}
