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

const Miner_JS = `
web3._extend({
	property: 'miner',
	methods:
	[
		new web3._extend.method({
			name: 'start',
			call: 'miner_start',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.method({
			name: 'stop',
			call: 'miner_stop',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.method({
			name: 'setpbferbase',
			call: 'miner_setpbferbase',
			params: 1,
			inputFormatter: [web3._extend.formatters.formatInputInt],
			outputFormatter: web3._extend.formatters.formatOutputBool
		}),
		new web3._extend.method({
			name: 'setExtra',
			call: 'miner_setExtra',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.method({
			name: 'setGasPrice',
			call: 'miner_setGasPrice',
			params: 1,
			inputFormatter: [web3._extend.utils.fromDecial]
		}),
		new web3._extend.method({
			name: 'startAutoDAG',
			call: 'miner_startAutoDAG',
			params: 0,
			inputFormatter: []
		}),
		new web3._extend.method({
			name: 'stopAutoDAG',
			call: 'miner_stopAutoDAG',
			params: 0,
			inputFormatter: []
		}),
		new web3._extend.method({
			name: 'makeDAG',
			call: 'miner_makeDAG',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputDefaultBlockNumberFormatter]
		})
	],
	properties:
	[
		new web3._extend.Property({
			name: 'hashrate',
			getter: 'miner_hashrate',
			outputFormatter: web3._extend.utils.toDecimal
		})
	]
});
`
