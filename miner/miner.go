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

// Package miner implements pbfcoin block creation and mining.
package miner

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/pbfcoin/go-pbfcoin/common"
	"github.com/pbfcoin/go-pbfcoin/core"
	"github.com/pbfcoin/go-pbfcoin/core/state"
	"github.com/pbfcoin/go-pbfcoin/core/types"
	"github.com/pbfcoin/go-pbfcoin/event"
	"github.com/pbfcoin/go-pbfcoin/logger"
	"github.com/pbfcoin/go-pbfcoin/logger/glog"
	"github.com/pbfcoin/go-pbfcoin/params"
	"github.com/pbfcoin/go-pbfcoin/pbf/downloader"
	"github.com/pbfcoin/go-pbfcoin/pow"
)

type Miner struct {
	mux *event.TypeMux

	worker *worker

	MinAcceptedGasPrice *big.Int

	threads  int
	coinbase common.Address
	mining   int32
	pbf      core.Backend
	pow      pow.PoW

	canStart    int32 // can start indicates whpbfer we can start the mining operation
	shouldStart int32 // should start indicates whpbfer we should start after sync
}

func New(pbf core.Backend, mux *event.TypeMux, pow pow.PoW) *Miner {
	miner := &Miner{pbf: pbf, mux: mux, pow: pow, worker: newWorker(common.Address{}, pbf), canStart: 1}
	go miner.update()

	return miner
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `Done` or `Failed` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your mining operation for as long as the DOS continues.
func (self *Miner) update() {
	events := self.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
out:
	for ev := range events.Chan() {
		switch ev.Data.(type) {
		case downloader.StartEvent:
			atomic.StoreInt32(&self.canStart, 0)
			if self.Mining() {
				self.Stop()
				atomic.StoreInt32(&self.shouldStart, 1)
				glog.V(logger.Info).Infoln("Mining operation aborted due to sync operation")
			}
		case downloader.DoneEvent, downloader.FailedEvent:
			shouldStart := atomic.LoadInt32(&self.shouldStart) == 1

			atomic.StoreInt32(&self.canStart, 1)
			atomic.StoreInt32(&self.shouldStart, 0)
			if shouldStart {
				self.Start(self.coinbase, self.threads)
			}
			// unsubscribe. we're only interested in this event once
			events.Unsubscribe()
			// stop immediately and ignore all further pending events
			break out
		}
	}
}

func (m *Miner) SetGasPrice(price *big.Int) {
	// FIXME block tests set a nil gas price. Quick dirty fix
	if price == nil {
		return
	}

	m.worker.setGasPrice(price)
}

func (self *Miner) Start(coinbase common.Address, threads int) {
	atomic.StoreInt32(&self.shouldStart, 1)
	self.threads = threads
	self.worker.coinbase = coinbase
	self.coinbase = coinbase

	if atomic.LoadInt32(&self.canStart) == 0 {
		glog.V(logger.Info).Infoln("Can not start mining operation due to network sync (starts when finished)")
		return
	}

	atomic.StoreInt32(&self.mining, 1)

	for i := 0; i < threads; i++ {
		self.worker.register(NewCpuAgent(i, self.pow))
	}

	glog.V(logger.Info).Infof("Starting mining operation (CPU=%d TOT=%d)\n", threads, len(self.worker.agents))

	self.worker.start()

	self.worker.commitNewWork()
}

func (self *Miner) Stop() {
	self.worker.stop()
	atomic.StoreInt32(&self.mining, 0)
	atomic.StoreInt32(&self.shouldStart, 0)
}

func (self *Miner) Register(agent Agent) {
	if self.Mining() {
		agent.Start()
	}
	self.worker.register(agent)
}

func (self *Miner) Unregister(agent Agent) {
	self.worker.unregister(agent)
}

func (self *Miner) Mining() bool {
	return atomic.LoadInt32(&self.mining) > 0
}

func (self *Miner) HashRate() (tot int64) {
	tot += self.pow.Gpbfashrate()
	// do we care this might race? is it worth we're rewriting some
	// aspects of the worker/locking up agents so we can get an accurate
	// hashrate?
	for agent := range self.worker.agents {
		tot += agent.GpbfashRate()
	}
	return
}

func (self *Miner) SetExtra(extra []byte) error {
	if uint64(len(extra)) > params.MaximumExtraDataSize.Uint64() {
		return fmt.Errorf("Extra exceeds max length. %d > %v", len(extra), params.MaximumExtraDataSize)
	}

	self.worker.extra = extra
	return nil
}

func (self *Miner) PendingState() *state.StateDB {
	return self.worker.pendingState()
}

func (self *Miner) PendingBlock() *types.Block {
	return self.worker.pendingBlock()
}

func (self *Miner) Setpbferbase(addr common.Address) {
	self.coinbase = addr
	self.worker.setpbferbase(addr)
}
