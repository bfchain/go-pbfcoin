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

package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/pbfcoin/go-pbfcoin/cmd/utils"
	"github.com/pbfcoin/go-pbfcoin/jsre"
	"github.com/pbfcoin/go-pbfcoin/logger"
	"github.com/pbfcoin/go-pbfcoin/logger/glog"
	"github.com/pbfcoin/go-pbfcoin/rpc/comms"
	"github.com/pbfcoin/go-pbfcoin/rpc/shared"
	"github.com/pbfcoin/go-pbfcoin/rpc/useragent"
	"github.com/pbfcoin/go-pbfcoin/xpbf"

	"github.com/robertkrimen/otto"
)

type Jpbf struct {
	pbfApi shared.pbfcoinApi
	re     *jsre.JSRE
	client comms.pbfcoinClient
	fe     xpbf.Frontend
}

func NewJpbf(pbfApi shared.pbfcoinApi, re *jsre.JSRE, client comms.pbfcoinClient, fe xpbf.Frontend) *Jpbf {
	return &Jpbf{pbfApi, re, client, fe}
}

func (self *Jpbf) err(call otto.FunctionCall, code int, msg string, id interface{}) (response otto.Value) {
	m := shared.NewRpcErrorResponse(id, shared.JsonRpcVersion, code, fmt.Errorf(msg))
	errObj, _ := json.Marshal(m.Error)
	errRes, _ := json.Marshal(m)

	call.Otto.Run("ret_error = " + string(errObj))
	res, _ := call.Otto.Run("ret_response = " + string(errRes))

	return res
}

func (self *Jpbf) Send(call otto.FunctionCall) (response otto.Value) {
	reqif, err := call.Argument(0).Export()
	if err != nil {
		return self.err(call, -32700, err.Error(), nil)
	}

	jsonreq, err := json.Marshal(reqif)
	var reqs []shared.Request
	batch := true
	err = json.Unmarshal(jsonreq, &reqs)
	if err != nil {
		reqs = make([]shared.Request, 1)
		err = json.Unmarshal(jsonreq, &reqs[0])
		batch = false
	}

	call.Otto.Set("response_len", len(reqs))
	call.Otto.Run("var ret_response = new Array(response_len);")

	for i, req := range reqs {
		var respif interface{}
		err := self.client.Send(&req)
		if err != nil {
			return self.err(call, -32603, err.Error(), req.Id)
		}

	recv:
		respif, err = self.client.Recv()
		if err != nil {
			return self.err(call, -32603, err.Error(), req.Id)
		}

		agentreq, isRequest := respif.(*shared.Request)
		if isRequest {
			self.handleRequest(agentreq)
			goto recv // receive response after agent interaction
		}

		sucres, isSuccessResponse := respif.(*shared.SuccessResponse)
		errres, isErrorResponse := respif.(*shared.ErrorResponse)
		if !isSuccessResponse && !isErrorResponse {
			return self.err(call, -32603, fmt.Sprintf("Invalid response type (%T)", respif), req.Id)
		}

		call.Otto.Set("ret_jsonrpc", shared.JsonRpcVersion)
		call.Otto.Set("ret_id", req.Id)

		var res []byte
		if isSuccessResponse {
			res, err = json.Marshal(sucres.Result)
		} else if isErrorResponse {
			res, err = json.Marshal(errres.Error)
		}

		call.Otto.Set("ret_result", string(res))
		call.Otto.Set("response_idx", i)
		response, err = call.Otto.Run(`
		ret_response[response_idx] = { jsonrpc: ret_jsonrpc, id: ret_id, result: JSON.parse(ret_result) };
		`)
	}

	if !batch {
		call.Otto.Run("ret_response = ret_response[0];")
	}

	if call.Argument(1).IsObject() {
		call.Otto.Set("callback", call.Argument(1))
		call.Otto.Run(`
	    if (Object.prototype.toString.call(callback) == '[object Function]') {
			callback(null, ret_response);
		}
		`)
	}

	return
}

// handleRequest will handle user agent requests by interacting with the user and sending
// the user response back to the gpbf service
func (self *Jpbf) handleRequest(req *shared.Request) bool {
	var err error
	var args []interface{}
	if err = json.Unmarshal(req.Params, &args); err != nil {
		glog.V(logger.Info).Infof("Unable to parse agent request - %v\n", err)
		return false
	}

	switch req.method {
	case useragent.AskPasswordmethod:
		return self.askPassword(req.Id, req.Jsonrpc, args)
	case useragent.ConfirmTransactionmethod:
		return self.confirmTransaction(req.Id, req.Jsonrpc, args)
	}

	return false
}

// askPassword will ask the user to supply the password for a given account
func (self *Jpbf) askPassword(id interface{}, jsonrpc string, args []interface{}) bool {
	var err error
	var passwd string
	if len(args) >= 1 {
		if account, ok := args[0].(string); ok {
			fmt.Printf("Unlock account %s\n", account)
		} else {
			return false
		}
	}
	passwd, err = utils.PromptPassword("Passphrase: ", true)

	if err = self.client.Send(shared.NewRpcResponse(id, jsonrpc, passwd, err)); err != nil {
		glog.V(logger.Info).Infof("Unable to send user agent ask password response - %v\n", err)
	}

	return err == nil
}

func (self *Jpbf) confirmTransaction(id interface{}, jsonrpc string, args []interface{}) bool {
	// Accept all tx which are send from this console
	return self.client.Send(shared.NewRpcResponse(id, jsonrpc, true, nil)) == nil
}
