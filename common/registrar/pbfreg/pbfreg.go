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

package pbfreg

import (
	"math/big"

	"github.com/pbfcoin/go-pbfcoin/common/registrar"
	"github.com/pbfcoin/go-pbfcoin/xpbf"
)

// implements a versioned Registrar on an archiving full node
type pbfReg struct {
	backend  *xpbf.Xpbf
	registry *registrar.Registrar
}

func New(xe *xpbf.Xpbf) (self *pbfReg) {
	self = &pbfReg{backend: xe}
	self.registry = registrar.New(xe)
	return
}

func (self *pbfReg) Registry() *registrar.Registrar {
	return self.registry
}

func (self *pbfReg) Resolver(n *big.Int) *registrar.Registrar {
	xe := self.backend
	if n != nil {
		xe = self.backend.AtStateNum(n.Int64())
	}
	return registrar.New(xe)
}
