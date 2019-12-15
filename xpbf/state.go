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

package xpbf

import (
	"github.com/pbfcoin/go-pbfcoin/common"
	"github.com/pbfcoin/go-pbfcoin/core/state"
)

type State struct {
	xpbf  *Xpbf
	state *state.StateDB
}

func NewState(xpbf *Xpbf, statedb *state.StateDB) *State {
	return &State{xpbf, statedb}
}

func (self *State) State() *state.StateDB {
	return self.state
}

func (self *State) Get(addr string) *Object {
	return &Object{self.state.GetStateObject(common.HexToAddress(addr))}
}

func (self *State) SafeGet(addr string) *Object {
	return &Object{self.safeGet(addr)}
}

func (self *State) safeGet(addr string) *state.StateObject {
	object := self.state.GetStateObject(common.HexToAddress(addr))
	if object == nil {
		object = state.NewStateObject(common.HexToAddress(addr), self.xpbf.backend.ChainDb())
	}

	return object
}
