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

package logger

import (
	"math/big"
	"time"
)

type utctime8601 struct{}

func (utctime8601) MarshalJSON() ([]byte, error) {
	timestr := time.Now().UTC().Format(time.RFC3339Nano)
	// Bounds check
	if len(timestr) > 26 {
		timestr = timestr[:26]
	}
	return []byte(`"` + timestr + `Z"`), nil
}

type JsonLog interface {
	EventName() string
}

type LogEvent struct {
	// Guid string      `json:"guid"`
	Ts utctime8601 `json:"ts"`
	// Level string      `json:"level"`
}

type LogStarting struct {
	ClientString    string `json:"client_impl"`
	ProtocolVersion int    `json:"eth_version"`
	LogEvent
}

func (l *LogStarting) EventName() string {
	return "starting"
}

type P2PConnected struct {
	RemoteId            string `json:"remote_id"`
	RemoteAddress       string `json:"remote_addr"`
	RemoteVersionString string `json:"remote_version_string"`
	NumConnections      int    `json:"num_connections"`
	LogEvent
}

func (l *P2PConnected) EventName() string {
	return "p2p.connected"
}

type P2PDisconnected struct {
	NumConnections int    `json:"num_connections"`
	RemoteId       string `json:"remote_id"`
	LogEvent
}

func (l *P2PDisconnected) EventName() string {
	return "p2p.disconnected"
}

type pbfMinerNewBlock struct {
	BlockHash     string   `json:"block_hash"`
	BlockNumber   *big.Int `json:"block_number"`
	ChainHeadHash string   `json:"chain_head_hash"`
	BlockPrevHash string   `json:"block_prev_hash"`
	LogEvent
}

func (l *pbfMinerNewBlock) EventName() string {
	return "pbf.miner.new_block"
}

type pbfChainReceivedNewBlock struct {
	BlockHash     string   `json:"block_hash"`
	BlockNumber   *big.Int `json:"block_number"`
	ChainHeadHash string   `json:"chain_head_hash"`
	BlockPrevHash string   `json:"block_prev_hash"`
	RemoteId      string   `json:"remote_id"`
	LogEvent
}

func (l *pbfChainReceivedNewBlock) EventName() string {
	return "pbf.chain.received.new_block"
}

type pbfChainNewHead struct {
	BlockHash     string   `json:"block_hash"`
	BlockNumber   *big.Int `json:"block_number"`
	ChainHeadHash string   `json:"chain_head_hash"`
	BlockPrevHash string   `json:"block_prev_hash"`
	LogEvent
}

func (l *pbfChainNewHead) EventName() string {
	return "pbf.chain.new_head"
}

type pbfTxReceived struct {
	TxHash   string `json:"tx_hash"`
	RemoteId string `json:"remote_id"`
	LogEvent
}

func (l *pbfTxReceived) EventName() string {
	return "pbf.tx.received"
}

//
//
// The types below are legacy and need to be converted to new format or deleted
//
//

// type P2PConnecting struct {
// 	RemoteId       string `json:"remote_id"`
// 	RemoteEndpoint string `json:"remote_endpoint"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PConnecting) EventName() string {
// 	return "p2p.connecting"
// }

// type P2PHandshaked struct {
// 	RemoteCapabilities []string `json:"remote_capabilities"`
// 	RemoteId           string   `json:"remote_id"`
// 	NumConnections     int      `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PHandshaked) EventName() string {
// 	return "p2p.handshaked"
// }

// type P2PDisconnecting struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PDisconnecting) EventName() string {
// 	return "p2p.disconnecting"
// }

// type P2PDisconnectingBadHandshake struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PDisconnectingBadHandshake) EventName() string {
// 	return "p2p.disconnecting.bad_handshake"
// }

// type P2PDisconnectingBadProtocol struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PDisconnectingBadProtocol) EventName() string {
// 	return "p2p.disconnecting.bad_protocol"
// }

// type P2PDisconnectingReputation struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PDisconnectingReputation) EventName() string {
// 	return "p2p.disconnecting.reputation"
// }

// type P2PDisconnectingDHT struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PDisconnectingDHT) EventName() string {
// 	return "p2p.disconnecting.dht"
// }

// type P2PpbfDisconnectingBadBlock struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PpbfDisconnectingBadBlock) EventName() string {
// 	return "p2p.pbf.disconnecting.bad_block"
// }

// type P2PpbfDisconnectingBadTx struct {
// 	Reason         string `json:"reason"`
// 	RemoteId       string `json:"remote_id"`
// 	NumConnections int    `json:"num_connections"`
// 	LogEvent
// }

// func (l *P2PpbfDisconnectingBadTx) EventName() string {
// 	return "p2p.pbf.disconnecting.bad_tx"
// }

// type pbfNewBlockBroadcasted struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockBroadcasted) EventName() string {
// 	return "pbf.newblock.broadcasted"
// }

// type pbfNewBlockIsKnown struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockIsKnown) EventName() string {
// 	return "pbf.newblock.is_known"
// }

// type pbfNewBlockIsNew struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockIsNew) EventName() string {
// 	return "pbf.newblock.is_new"
// }

// type pbfNewBlockMissingParent struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockMissingParent) EventName() string {
// 	return "pbf.newblock.missing_parent"
// }

// type pbfNewBlockIsInvalid struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockIsInvalid) EventName() string {
// 	return "pbf.newblock.is_invalid"
// }

// type pbfNewBlockChainIsOlder struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockChainIsOlder) EventName() string {
// 	return "pbf.newblock.chain.is_older"
// }

// type pbfNewBlockChainIsCanonical struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockChainIsCanonical) EventName() string {
// 	return "pbf.newblock.chain.is_cannonical"
// }

// type pbfNewBlockChainNotCanonical struct {
// 	BlockNumber     int    `json:"block_number"`
// 	HeadHash        string `json:"head_hash"`
// 	BlockHash       string `json:"block_hash"`
// 	BlockDifficulty int    `json:"block_difficulty"`
// 	BlockPrevHash   string `json:"block_prev_hash"`
// 	LogEvent
// }

// func (l *pbfNewBlockChainNotCanonical) EventName() string {
// 	return "pbf.newblock.chain.not_cannonical"
// }

// type pbfTxCreated struct {
// 	TxHash    string `json:"tx_hash"`
// 	TxSender  string `json:"tx_sender"`
// 	TxAddress string `json:"tx_address"`
// 	TxHexRLP  string `json:"tx_hexrlp"`
// 	TxNonce   int    `json:"tx_nonce"`
// 	LogEvent
// }

// func (l *pbfTxCreated) EventName() string {
// 	return "pbf.tx.created"
// }

// type pbfTxBroadcasted struct {
// 	TxHash    string `json:"tx_hash"`
// 	TxSender  string `json:"tx_sender"`
// 	TxAddress string `json:"tx_address"`
// 	TxNonce   int    `json:"tx_nonce"`
// 	LogEvent
// }

// func (l *pbfTxBroadcasted) EventName() string {
// 	return "pbf.tx.broadcasted"
// }

// type pbfTxValidated struct {
// 	TxHash    string `json:"tx_hash"`
// 	TxSender  string `json:"tx_sender"`
// 	TxAddress string `json:"tx_address"`
// 	TxNonce   int    `json:"tx_nonce"`
// 	LogEvent
// }

// func (l *pbfTxValidated) EventName() string {
// 	return "pbf.tx.validated"
// }

// type pbfTxIsInvalid struct {
// 	TxHash    string `json:"tx_hash"`
// 	TxSender  string `json:"tx_sender"`
// 	TxAddress string `json:"tx_address"`
// 	Reason    string `json:"reason"`
// 	TxNonce   int    `json:"tx_nonce"`
// 	LogEvent
// }

// func (l *pbfTxIsInvalid) EventName() string {
// 	return "pbf.tx.is_invalid"
// }
