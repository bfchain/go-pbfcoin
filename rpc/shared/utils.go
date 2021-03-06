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

package shared

import "strings"

const (
	AdminApiName    = "admin"
	pbfApiName      = "pbf"
	DbApiName       = "db"
	DebugApiName    = "debug"
	MergedApiName   = "merged"
	MinerApiName    = "miner"
	NetApiName      = "net"
	ShhApiName      = "shh"
	TxPoolApiName   = "txpool"
	PersonalApiName = "personal"
	Web3ApiName     = "web3"

	JsonRpcVersion = "2.0"
)

var (
	// All API's
	AllApis = strings.Join([]string{
		AdminApiName, DbApiName, pbfApiName, DebugApiName, MinerApiName, NetApiName,
		ShhApiName, TxPoolApiName, PersonalApiName, Web3ApiName,
	}, ",")
)
