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
	"bytes"
	"encoding/json"
	"math/big"

	"fmt"

	"github.com/pbfcoin/go-pbfcoin/common"
	"github.com/pbfcoin/go-pbfcoin/common/natspec"
	"github.com/pbfcoin/go-pbfcoin/pbf"
	"github.com/pbfcoin/go-pbfcoin/rlp"
	"github.com/pbfcoin/go-pbfcoin/rpc/codec"
	"github.com/pbfcoin/go-pbfcoin/rpc/shared"
	"github.com/pbfcoin/go-pbfcoin/xpbf"
	"gopkg.in/fatih/set.v0"
)

const (
	pbfApiVersion = "1.0"
)

// pbf api provider
// See https://github.com/pbfcoin/wiki/wiki/JSON-RPC
type pbfApi struct {
	xpbf    *xpbf.Xpbf
	pbfcoin *pbf.pbfcoin
	methods map[string]pbfhandler
	codec   codec.ApiCoder
}

// pbf callback handler
type pbfhandler func(*pbfApi, *shared.Request) (interface{}, error)

var (
	pbfMapping = map[string]pbfhandler{
		"eth_accounts":                            (*pbfApi).Accounts,
		"eth_blockNumber":                         (*pbfApi).BlockNumber,
		"eth_getBalance":                          (*pbfApi).GetBalance,
		"eth_protocolVersion":                     (*pbfApi).ProtocolVersion,
		"eth_coinbase":                            (*pbfApi).Coinbase,
		"eth_mining":                              (*pbfApi).IsMining,
		"eth_syncing":                             (*pbfApi).IsSyncing,
		"eth_gasPrice":                            (*pbfApi).GasPrice,
		"eth_getStorage":                          (*pbfApi).GetStorage,
		"eth_storageAt":                           (*pbfApi).GetStorage,
		"eth_getStorageAt":                        (*pbfApi).GetStorageAt,
		"eth_getTransactionCount":                 (*pbfApi).GetTransactionCount,
		"eth_getBlockTransactionCountByHash":      (*pbfApi).GetBlockTransactionCountByHash,
		"eth_getBlockTransactionCountByNumber":    (*pbfApi).GetBlockTransactionCountByNumber,
		"eth_getUncleCountByBlockHash":            (*pbfApi).GetUncleCountByBlockHash,
		"eth_getUncleCountByBlockNumber":          (*pbfApi).GetUncleCountByBlockNumber,
		"eth_getData":                             (*pbfApi).GetData,
		"eth_getCode":                             (*pbfApi).GetData,
		"eth_getNatSpec":                          (*pbfApi).GetNatSpec,
		"eth_sign":                                (*pbfApi).Sign,
		"eth_sendRawTransaction":                  (*pbfApi).SubmitTransaction,
		"eth_submitTransaction":                   (*pbfApi).SubmitTransaction,
		"eth_sendTransaction":                     (*pbfApi).SendTransaction,
		"eth_signTransaction":                     (*pbfApi).SignTransaction,
		"eth_transact":                            (*pbfApi).SendTransaction,
		"eth_estimateGas":                         (*pbfApi).EstimateGas,
		"eth_call":                                (*pbfApi).Call,
		"eth_flush":                               (*pbfApi).Flush,
		"eth_getBlockByHash":                      (*pbfApi).GetBlockByHash,
		"eth_getBlockByNumber":                    (*pbfApi).GetBlockByNumber,
		"eth_getTransactionByHash":                (*pbfApi).GetTransactionByHash,
		"eth_getTransactionByBlockNumberAndIndex": (*pbfApi).GetTransactionByBlockNumberAndIndex,
		"eth_getTransactionByBlockHashAndIndex":   (*pbfApi).GetTransactionByBlockHashAndIndex,
		"eth_getUncleByBlockHashAndIndex":         (*pbfApi).GetUncleByBlockHashAndIndex,
		"eth_getUncleByBlockNumberAndIndex":       (*pbfApi).GetUncleByBlockNumberAndIndex,
		"eth_getCompilers":                        (*pbfApi).GetCompilers,
		"eth_compileSolidity":                     (*pbfApi).CompileSolidity,
		"eth_newFilter":                           (*pbfApi).NewFilter,
		"eth_newBlockFilter":                      (*pbfApi).NewBlockFilter,
		"eth_newPendingTransactionFilter":         (*pbfApi).NewPendingTransactionFilter,
		"eth_uninstallFilter":                     (*pbfApi).UninstallFilter,
		"eth_getFilterChanges":                    (*pbfApi).GetFilterChanges,
		"eth_getFilterLogs":                       (*pbfApi).GetFilterLogs,
		"eth_getLogs":                             (*pbfApi).GetLogs,
		"eth_hashrate":                            (*pbfApi).Hashrate,
		"eth_getWork":                             (*pbfApi).GetWork,
		"eth_submitWork":                          (*pbfApi).SubmitWork,
		"eth_submitHashrate":                      (*pbfApi).SubmitHashrate,
		"eth_resend":                              (*pbfApi).Resend,
		"eth_pendingTransactions":                 (*pbfApi).PendingTransactions,
		"eth_getTransactionReceipt":               (*pbfApi).GetTransactionReceipt,
	}
)

// create new pbfApi instance
func NewethApi(xpbf *xpbf.Xpbf, pbf *pbf.pbfcoin, codec codec.Codec) *pbfApi {
	return &pbfApi{xpbf, pbf, pbfMapping, codec.New(nil)}
}

// collection with supported methods
func (self *pbfApi) methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *pbfApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.method)
}

func (self *pbfApi) Name() string {
	return shared.pbfApiName
}

func (self *pbfApi) ApiVersion() string {
	return pbfApiVersion
}

func (self *pbfApi) Accounts(req *shared.Request) (interface{}, error) {
	return self.xpbf.Accounts(), nil
}

func (self *pbfApi) Hashrate(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xpbf.HashRate()), nil
}

func (self *pbfApi) BlockNumber(req *shared.Request) (interface{}, error) {
	num := self.xpbf.CurrentBlock().Number()
	return newHexNum(num.Bytes()), nil
}

func (self *pbfApi) GetBalance(req *shared.Request) (interface{}, error) {
	args := new(GetBalanceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xpbf.AtStateNum(args.BlockNumber).BalanceAt(args.Address), nil
}

func (self *pbfApi) ProtocolVersion(req *shared.Request) (interface{}, error) {
	return self.xpbf.pbfVersion(), nil
}

func (self *pbfApi) Coinbase(req *shared.Request) (interface{}, error) {
	return newHexData(self.xpbf.Coinbase()), nil
}

func (self *pbfApi) IsMining(req *shared.Request) (interface{}, error) {
	return self.xpbf.IsMining(), nil
}

func (self *pbfApi) IsSyncing(req *shared.Request) (interface{}, error) {
	origin, current, height := self.pbfcoin.Downloader().Progress()
	if current < height {
		return map[string]interface{}{
			"startingBlock": newHexNum(big.NewInt(int64(origin)).Bytes()),
			"currentBlock":  newHexNum(big.NewInt(int64(current)).Bytes()),
			"highestBlock":  newHexNum(big.NewInt(int64(height)).Bytes()),
		}, nil
	}
	return false, nil
}

func (self *pbfApi) GasPrice(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xpbf.DefaultGasPrice().Bytes()), nil
}

func (self *pbfApi) GetStorage(req *shared.Request) (interface{}, error) {
	args := new(GetStorageArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xpbf.AtStateNum(args.BlockNumber).State().SafeGet(args.Address).Storage(), nil
}

func (self *pbfApi) GetStorageAt(req *shared.Request) (interface{}, error) {
	args := new(GetStorageAtArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xpbf.AtStateNum(args.BlockNumber).StorageAt(args.Address, args.Key), nil
}

func (self *pbfApi) GetTransactionCount(req *shared.Request) (interface{}, error) {
	args := new(GetTxCountArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	count := self.xpbf.AtStateNum(args.BlockNumber).TxCountAt(args.Address)
	return fmt.Sprintf("%#x", count), nil
}

func (self *pbfApi) GetBlockTransactionCountByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xpbf.pbfBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *pbfApi) GetBlockTransactionCountByNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xpbf.pbfBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *pbfApi) GetUncleCountByBlockHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xpbf.pbfBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *pbfApi) GetUncleCountByBlockNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xpbf.pbfBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *pbfApi) GetData(req *shared.Request) (interface{}, error) {
	args := new(GetDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v := self.xpbf.AtStateNum(args.BlockNumber).CodeAtBytes(args.Address)
	return newHexData(v), nil
}

func (self *pbfApi) Sign(req *shared.Request) (interface{}, error) {
	args := new(NewSigArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v, err := self.xpbf.Sign(args.From, args.Data, false)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *pbfApi) SubmitTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	v, err := self.xpbf.PushTx(args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// JsonTransaction is returned as response by the JSON RPC. It contains the
// signed RLP encoded transaction as Raw and the signed transaction object as Tx.
type JsonTransaction struct {
	Raw string `json:"raw"`
	Tx  *tx    `json:"tx"`
}

func (self *pbfApi) SignTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	tx, err := self.xpbf.SignTransaction(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return JsonTransaction{"0x" + common.Bytes2Hex(data), newTx(tx)}, nil
}

func (self *pbfApi) SendTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	v, err := self.xpbf.Transact(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *pbfApi) GetNatSpec(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	var jsontx = fmt.Sprintf(`{"params":[{"to":"%s","data": "%s"}]}`, args.To, args.Data)
	notice := natspec.GetNotice(self.xpbf, jsontx, self.pbfcoin.HTTPClient())

	return notice, nil
}

func (self *pbfApi) EstimateGas(req *shared.Request) (interface{}, error) {
	_, gas, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if len(gas) == 0 {
		return newHexNum(0), nil
	} else {
		return newHexNum(common.String2Big(gas)), err
	}
}

func (self *pbfApi) Call(req *shared.Request) (interface{}, error) {
	v, _, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if v == "0x0" {
		return newHexData([]byte{}), nil
	} else {
		return newHexData(common.FromHex(v)), nil
	}
}

func (self *pbfApi) Flush(req *shared.Request) (interface{}, error) {
	return nil, shared.NewNotImplementedError(req.method)
}

func (self *pbfApi) doCall(params json.RawMessage) (string, string, error) {
	args := new(CallArgs)
	if err := self.codec.Decode(params, &args); err != nil {
		return "", "", err
	}

	return self.xpbf.AtStateNum(args.BlockNumber).Call(args.From, args.To, args.Value.String(), args.Gas.String(), args.GasPrice.String(), args.Data)
}

func (self *pbfApi) GetBlockByHash(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByHashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xpbf.pbfBlockByHash(args.BlockHash)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xpbf.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *pbfApi) GetBlockByNumber(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByNumberArgs)
	if err := json.Unmarshal(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xpbf.pbfBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xpbf.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *pbfApi) GetTransactionByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	tx, bhash, bnum, txi := self.xpbf.pbfTransactionByHash(args.Hash)
	if tx != nil {
		v := NewTransactionRes(tx)
		// if the blockhash is 0, assume this is a pending transaction
		if bytes.Compare(bhash.Bytes(), bytes.Repeat([]byte{0}, 32)) != 0 {
			v.BlockHash = newHexData(bhash)
			v.BlockNumber = newHexNum(bnum)
			v.TxIndex = newHexNum(txi)
		}
		return v, nil
	}
	return nil, nil
}

func (self *pbfApi) GetTransactionByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xpbf.pbfBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xpbf.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Transactions[args.Index], nil
	}
}

func (self *pbfApi) GetTransactionByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xpbf.pbfBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xpbf.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Transactions[args.Index], nil
}

func (self *pbfApi) GetUncleByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xpbf.pbfBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xpbf.Td(raw.Hash()), false)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Uncles[args.Index], nil
}

func (self *pbfApi) GetUncleByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xpbf.pbfBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xpbf.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Uncles[args.Index], nil
	}
}

func (self *pbfApi) GetCompilers(req *shared.Request) (interface{}, error) {
	var lang string
	if solc, _ := self.xpbf.Solc(); solc != nil {
		lang = "Solidity"
	}
	c := []string{lang}
	return c, nil
}

func (self *pbfApi) CompileSolidity(req *shared.Request) (interface{}, error) {
	solc, _ := self.xpbf.Solc()
	if solc == nil {
		return nil, shared.NewNotAvailableError(req.method, "solc (solidity compiler) not found")
	}

	args := new(SourceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	contracts, err := solc.Compile(args.Source)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (self *pbfApi) NewFilter(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	id := self.xpbf.NewLogFilter(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)
	return newHexNum(big.NewInt(int64(id)).Bytes()), nil
}

func (self *pbfApi) NewBlockFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xpbf.NewBlockFilter()), nil
}

func (self *pbfApi) NewPendingTransactionFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xpbf.NewTransactionFilter()), nil
}

func (self *pbfApi) UninstallFilter(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xpbf.UninstallFilter(args.Id), nil
}

func (self *pbfApi) GetFilterChanges(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	switch self.xpbf.GetFilterType(args.Id) {
	case xpbf.BlockFilterTy:
		return NewHashesRes(self.xpbf.BlockFilterChanged(args.Id)), nil
	case xpbf.TransactionFilterTy:
		return NewHashesRes(self.xpbf.TransactionFilterChanged(args.Id)), nil
	case xpbf.LogFilterTy:
		return NewLogsRes(self.xpbf.LogFilterChanged(args.Id)), nil
	default:
		return []string{}, nil // reply empty string slice
	}
}

func (self *pbfApi) GetFilterLogs(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return NewLogsRes(self.xpbf.Logs(args.Id)), nil
}

func (self *pbfApi) GetLogs(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return NewLogsRes(self.xpbf.AllLogs(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)), nil
}

func (self *pbfApi) GetWork(req *shared.Request) (interface{}, error) {
	self.xpbf.SetMining(true, 0)
	ret, err := self.xpbf.RemoteMining().GetWork()
	if err != nil {
		return nil, shared.NewNotReadyError("mining work")
	} else {
		return ret, nil
	}
}

func (self *pbfApi) SubmitWork(req *shared.Request) (interface{}, error) {
	args := new(SubmitWorkArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xpbf.RemoteMining().SubmitWork(args.Nonce, common.HexToHash(args.Digest), common.HexToHash(args.Header)), nil
}

func (self *pbfApi) SubmitHashrate(req *shared.Request) (interface{}, error) {
	args := new(SubmitHashRateArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return false, shared.NewDecodeParamError(err.Error())
	}
	self.xpbf.RemoteMining().SubmitHashrate(common.HexToHash(args.Id), args.Rate)
	return true, nil
}

func (self *pbfApi) Resend(req *shared.Request) (interface{}, error) {
	args := new(ResendArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	from := common.HexToAddress(args.Tx.From)

	pending := self.pbfcoin.TxPool().GetTransactions()
	for _, p := range pending {
		if pFrom, err := p.From(); err == nil && pFrom == from && p.SigHash() == args.Tx.tx.SigHash() {
			self.pbfcoin.TxPool().RemoveTx(common.HexToHash(args.Tx.Hash))
			return self.xpbf.Transact(args.Tx.From, args.Tx.To, args.Tx.Nonce, args.Tx.Value, args.GasLimit, args.GasPrice, args.Tx.Data)
		}
	}

	return nil, fmt.Errorf("Transaction %s not found", args.Tx.Hash)
}

func (self *pbfApi) PendingTransactions(req *shared.Request) (interface{}, error) {
	txs := self.pbfcoin.TxPool().GetTransactions()

	// grab the accounts from the account manager. This will help with determining which
	// transactions should be returned.
	accounts, err := self.pbfcoin.AccountManager().Accounts()
	if err != nil {
		return nil, err
	}

	// Add the accouns to a new set
	accountSet := set.New()
	for _, account := range accounts {
		accountSet.Add(account.Address)
	}

	var ltxs []*tx
	for _, tx := range txs {
		if from, _ := tx.From(); accountSet.Has(from) {
			ltxs = append(ltxs, newTx(tx))
		}
	}

	return ltxs, nil
}

func (self *pbfApi) GetTransactionReceipt(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	txhash := common.BytesToHash(common.FromHex(args.Hash))
	tx, bhash, bnum, txi := self.xpbf.pbfTransactionByHash(args.Hash)
	rec := self.xpbf.GetTxReceipt(txhash)
	// We could have an error of "not found". Should disambiguate
	// if err != nil {
	// 	return err, nil
	// }
	if rec != nil && tx != nil {
		v := NewReceiptRes(rec)
		v.BlockHash = newHexData(bhash)
		v.BlockNumber = newHexNum(bnum)
		v.TransactionIndex = newHexNum(txi)
		return v, nil
	}

	return nil, nil
}
