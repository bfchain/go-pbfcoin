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
	NetApiVersion = "1.0"
)

var (
	// mapping between methods and handlers
	netMapping = map[string]npbfandler{
		"net_peerCount": (*netApi).PeerCount,
		"net_listening": (*netApi).IsListening,
		"net_version":   (*netApi).Version,
	}
)

// net callback handler
type npbfandler func(*netApi, *shared.Request) (interface{}, error)

// net api provider
type netApi struct {
	xpbf    *xpbf.Xpbf
	pbfcoin *pbf.pbfcoin
	methods map[string]npbfandler
	codec   codec.ApiCoder
}

// create a new net api instance
func NewNetApi(xpbf *xpbf.Xpbf, pbf *pbf.pbfcoin, coder codec.Codec) *netApi {
	return &netApi{
		xpbf:    xpbf,
		pbfcoin: pbf,
		methods: netMapping,
		codec:   coder.New(nil),
	}
}

// collection with supported methods
func (self *netApi) methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *netApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.method)
}

func (self *netApi) Name() string {
	return shared.NetApiName
}

func (self *netApi) ApiVersion() string {
	return NetApiVersion
}

// Number of connected peers
func (self *netApi) PeerCount(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xpbf.PeerCount()), nil
}

func (self *netApi) IsListening(req *shared.Request) (interface{}, error) {
	return self.xpbf.IsListening(), nil
}

func (self *netApi) Version(req *shared.Request) (interface{}, error) {
	return self.xpbf.NetworkVersion(), nil
}
