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
	DbApiversion = "1.0"
)

var (
	// mapping between methods and handlers
	DbMapping = map[string]dbhandler{
		"db_getString": (*dbApi).GetString,
		"db_putString": (*dbApi).PutString,
		"db_gpbfex":    (*dbApi).Gpbfex,
		"db_putHex":    (*dbApi).PutHex,
	}
)

// db callback handler
type dbhandler func(*dbApi, *shared.Request) (interface{}, error)

// db api provider
type dbApi struct {
	xpbf    *xpbf.Xpbf
	pbfcoin *pbf.pbfcoin
	methods map[string]dbhandler
	codec   codec.ApiCoder
}

// create a new db api instance
func NewDbApi(xpbf *xpbf.Xpbf, pbfcoin *pbf.pbfcoin, coder codec.Codec) *dbApi {
	return &dbApi{
		xpbf:    xpbf,
		pbfcoin: pbfcoin,
		methods: DbMapping,
		codec:   coder.New(nil),
	}
}

// collection with supported methods
func (self *dbApi) methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *dbApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.method]; ok {
		return callback(self, req)
	}

	return nil, &shared.NotImplementedError{req.method}
}

func (self *dbApi) Name() string {
	return shared.DbApiName
}

func (self *dbApi) ApiVersion() string {
	return DbApiversion
}

func (self *dbApi) GetString(req *shared.Request) (interface{}, error) {
	args := new(DbArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	if err := args.requirements(); err != nil {
		return nil, err
	}

	ret, err := self.xpbf.DbGet([]byte(args.Database + args.Key))
	return string(ret), err
}

func (self *dbApi) PutString(req *shared.Request) (interface{}, error) {
	args := new(DbArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	if err := args.requirements(); err != nil {
		return nil, err
	}

	return self.xpbf.DbPut([]byte(args.Database+args.Key), args.Value), nil
}

func (self *dbApi) Gpbfex(req *shared.Request) (interface{}, error) {
	args := new(DbHexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	if err := args.requirements(); err != nil {
		return nil, err
	}

	if res, err := self.xpbf.DbGet([]byte(args.Database + args.Key)); err == nil {
		return newHexData(res), nil
	} else {
		return nil, err
	}
}

func (self *dbApi) PutHex(req *shared.Request) (interface{}, error) {
	args := new(DbHexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	if err := args.requirements(); err != nil {
		return nil, err
	}

	return self.xpbf.DbPut([]byte(args.Database+args.Key), args.Value), nil
}
