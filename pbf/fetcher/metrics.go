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

// Contains the metrics collected by the fetcher.

package fetcher

import (
	"github.com/pbfcoin/go-pbfcoin/metrics"
)

var (
	propAnnounceInMeter   = metrics.NewMeter("pbf/fetcher/prop/announces/in")
	propAnnounceOutTimer  = metrics.NewTimer("pbf/fetcher/prop/announces/out")
	propAnnounceDropMeter = metrics.NewMeter("pbf/fetcher/prop/announces/drop")
	propAnnounceDOSMeter  = metrics.NewMeter("pbf/fetcher/prop/announces/dos")

	propBroadcastInMeter   = metrics.NewMeter("pbf/fetcher/prop/broadcasts/in")
	propBroadcastOutTimer  = metrics.NewTimer("pbf/fetcher/prop/broadcasts/out")
	propBroadcastDropMeter = metrics.NewMeter("pbf/fetcher/prop/broadcasts/drop")
	propBroadcastDOSMeter  = metrics.NewMeter("pbf/fetcher/prop/broadcasts/dos")

	blockFetchMeter  = metrics.NewMeter("pbf/fetcher/fetch/blocks")
	headerFetchMeter = metrics.NewMeter("pbf/fetcher/fetch/headers")
	bodyFetchMeter   = metrics.NewMeter("pbf/fetcher/fetch/bodies")

	blockFilterInMeter   = metrics.NewMeter("pbf/fetcher/filter/blocks/in")
	blockFilterOutMeter  = metrics.NewMeter("pbf/fetcher/filter/blocks/out")
	headerFilterInMeter  = metrics.NewMeter("pbf/fetcher/filter/headers/in")
	headerFilterOutMeter = metrics.NewMeter("pbf/fetcher/filter/headers/out")
	bodyFilterInMeter    = metrics.NewMeter("pbf/fetcher/filter/bodies/in")
	bodyFilterOutMeter   = metrics.NewMeter("pbf/fetcher/filter/bodies/out")
)
