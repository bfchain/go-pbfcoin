// Copyright 2014 The go-pbfcoin Authors
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

package runtime

import (
	"math/big"
	"time"

	"github.com/pbfcoin/go-pbfcoin/common"
	"github.com/pbfcoin/go-pbfcoin/core/state"
	"github.com/pbfcoin/go-pbfcoin/core/vm"
	"github.com/pbfcoin/go-pbfcoin/crypto"
	"github.com/pbfcoin/go-pbfcoin/pbfdb"
)

// Config is a basic type specifing certain configuration flags for running
// the EVM.
type Config struct {
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    *big.Int
	GasPrice    *big.Int
	Value       *big.Int
	DisableJit  bool // "disable" so it's enabled by default
	Debug       bool

	GpbfashFn func(n uint64) common.Hash
}

// sets defaults on the config
func setDefaults(cfg *Config) {
	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == nil {
		cfg.GasLimit = new(big.Int).Set(common.MaxBig)
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GpbfashFn == nil {
		cfg.GpbfashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Sha3([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
}

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Executes sets up a in memory, temporarily, environment for the execution of
// the given code. It enabled the JIT by default and make sure that it's restored
// to it's original state afterwards.
func Execute(code, input []byte, cfg *Config) ([]byte, *state.StateDB, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	// defer the call to setting back the original values
	defer func(debug, forceJit, enableJit bool) {
		vm.Debug = debug
		vm.ForceJit = forceJit
		vm.EnableJit = enableJit
	}(vm.Debug, vm.ForceJit, vm.EnableJit)

	vm.ForceJit = !cfg.DisableJit
	vm.EnableJit = !cfg.DisableJit
	vm.Debug = cfg.Debug

	var (
		db, _      = pbfdb.NewMemDatabase()
		statedb, _ = state.New(common.Hash{}, db)
		vmenv      = NewEnv(cfg, statedb)
		sender     = statedb.CreateAccount(cfg.Origin)
		receiver   = statedb.CreateAccount(common.StringToAddress("contract"))
	)
	// set the receiver's (the executing contract) code for execution.
	receiver.SetCode(code)

	// Call the code with the given configuration.
	ret, err := vmenv.Call(
		sender,
		receiver.Address(),
		input,
		cfg.GasLimit,
		cfg.GasPrice,
		cfg.Value,
	)

	if cfg.Debug {
		vm.StdErrFormat(vmenv.StructLogs())
	}
	return ret, statedb, err
}
