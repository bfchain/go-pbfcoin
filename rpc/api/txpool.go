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

package api

import (
	"github.com/pbfcoin/go-pbfcoin/pbf"
	"github.com/pbfcoin/go-pbfcoin/rpc/codec"
	"github.com/pbfcoin/go-pbfcoin/rpc/shared"
	"github.com/pbfcoin/go-pbfcoin/xpbf"
)

const (
	TxPoolApiVersion = "1.0"
)

var (
	// mapping between methods and handlers
	txpoolMapping = map[string]txpoolhandler{
		"txpool_status": (*txPoolApi).Status,
	}
)

// net callback handler
type txpoolhandler func(*txPoolApi, *shared.Request) (interface{}, error)

// txpool api provider
type txPoolApi struct {
	xpbf    *xpbf.Xpbf
	pbfcoin *pbf.pbfcoin
	methods map[string]txpoolhandler
	codec   codec.ApiCoder
}

// create a new txpool api instance
func NewTxPoolApi(xpbf *xpbf.Xpbf, pbf *pbf.pbfcoin, coder codec.Codec) *txPoolApi {
	return &txPoolApi{
		xpbf:    xpbf,
		pbfcoin: pbf,
		methods: txpoolMapping,
		codec:   coder.New(nil),
	}
}

// collection with supported methods
func (self *txPoolApi) methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *txPoolApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.method)
}

func (self *txPoolApi) Name() string {
	return shared.TxPoolApiName
}

func (self *txPoolApi) ApiVersion() string {
	return TxPoolApiVersion
}

func (self *txPoolApi) Status(req *shared.Request) (interface{}, error) {
	pending, queue := self.pbfcoin.TxPool().Stats()
	return map[string]int{
		"pending": pending,
		"queued":  queue,
	}, nil
}
