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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/pbfcoin/go-pbfcoin/metrics"
)

var (
	hashInMeter      = metrics.NewMeter("pbf/downloader/hashes/in")
	hashReqTimer     = metrics.NewTimer("pbf/downloader/hashes/req")
	hashDropMeter    = metrics.NewMeter("pbf/downloader/hashes/drop")
	hashTimeoutMeter = metrics.NewMeter("pbf/downloader/hashes/timeout")

	blockInMeter      = metrics.NewMeter("pbf/downloader/blocks/in")
	blockReqTimer     = metrics.NewTimer("pbf/downloader/blocks/req")
	blockDropMeter    = metrics.NewMeter("pbf/downloader/blocks/drop")
	blockTimeoutMeter = metrics.NewMeter("pbf/downloader/blocks/timeout")

	headerInMeter      = metrics.NewMeter("pbf/downloader/headers/in")
	headerReqTimer     = metrics.NewTimer("pbf/downloader/headers/req")
	headerDropMeter    = metrics.NewMeter("pbf/downloader/headers/drop")
	headerTimeoutMeter = metrics.NewMeter("pbf/downloader/headers/timeout")

	bodyInMeter      = metrics.NewMeter("pbf/downloader/bodies/in")
	bodyReqTimer     = metrics.NewTimer("pbf/downloader/bodies/req")
	bodyDropMeter    = metrics.NewMeter("pbf/downloader/bodies/drop")
	bodyTimeoutMeter = metrics.NewMeter("pbf/downloader/bodies/timeout")

	receiptInMeter      = metrics.NewMeter("pbf/downloader/receipts/in")
	receiptReqTimer     = metrics.NewTimer("pbf/downloader/receipts/req")
	receiptDropMeter    = metrics.NewMeter("pbf/downloader/receipts/drop")
	receiptTimeoutMeter = metrics.NewMeter("pbf/downloader/receipts/timeout")

	stateInMeter      = metrics.NewMeter("pbf/downloader/states/in")
	stateReqTimer     = metrics.NewTimer("pbf/downloader/states/req")
	stateDropMeter    = metrics.NewMeter("pbf/downloader/states/drop")
	stateTimeoutMeter = metrics.NewMeter("pbf/downloader/states/timeout")
)
