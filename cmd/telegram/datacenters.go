package main

//
// Measuring Data Centers (DCs)
//

import (
	"context"
	"net"
	"strconv"

	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/fx"
)

// measureDCs measures telegram data centers.
func measureDCs(ctx context.Context, state *measurementState) {
	// tell the parent we terminated
	defer state.wg.Done()

	// ipAddrs contains the DCs IP addresses
	var ipAddrs = dslx.AddressSet().Add(
		"149.154.175.50",
		"149.154.167.51",
		"149.154.175.100",
		"149.154.167.91",
		"149.154.171.5",
	)

	// construct the list of endpoints to measure: we need to
	// measure each IP address with port 80 and 443
	var (
		endpoints []*dslx.EndpointState
		ports     = []int{80, 443}
	)
	for addr := range ipAddrs.M {
		for _, port := range ports {
			endpoints = append(endpoints, dslx.Endpoint(
				dslx.EndpointNetwork("tcp"),
				dslx.EndpointAddress(net.JoinHostPort(addr, strconv.Itoa(port))),
				dslx.EndpointOptionIDGenerator(state.idGen),
				dslx.EndpointOptionLogger(state.logger),
				dslx.EndpointOptionZerotime(state.zeroTime),
			))
		}
	}

	var (
		// tcpConnectSuccessCounter counts the number of TCP successes
		tcpConnectSuccessCounter = dslx.Counter[*dslx.TCPConnectResultState]()

		// httpRoundTripSuccessCounter counts the number of HTTP successes
		httpRoundTripSuccessCounter = dslx.Counter[*dslx.HTTPRequestResultState]()
	)

	// create the established connections pool
	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	// construct the function to measure the endpoints
	function := fx.ComposeFlat5(
		dslx.TCPConnect(connpool),
		tcpConnectSuccessCounter.Func(), // count number of successful TCP connects
		dslx.HTTPTransportTCP(),
		dslx.HTTPRequest(
			dslx.HTTPRequestOptionMethod("POST"),
		),
		httpRoundTripSuccessCounter.Func(), // count number of successful HTTP round trips
	)

	// measure all the endpoints in parallel and collect the results
	results := fx.Map(
		ctx,
		fx.Parallelism(3),
		function,
		endpoints...,
	)

	// extract and merge observations with the test keys
	state.tk.mergeObservations(dslx.ExtractObservations(results...)...)

	// set top-level keys indicating DCs blocking
	state.tk.setDCBlocking(tcpConnectSuccessCounter, httpRoundTripSuccessCounter)
}

// setDCBlocking sets the blocking status of data centers based on
// the number of times we completed TCP and HTTP operations.
//
// Arguments:
//
// - tcpSuccessCount is the number of times a TCP connect succeded;
//
// - httpSuccessCount is like tcpSuccessCount but for HTTP.
//
// We say there is TCP blocking if no TCP connect succeded. Likewise, we
// say there is HTTP blocking when no HTTP round trip succeded.
func (tk *testKeys) setDCBlocking(
	tcpSuccessCounter *dslx.CounterState[*dslx.TCPConnectResultState],
	httpSuccessCounter *dslx.CounterState[*dslx.HTTPRequestResultState],
) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	tk.TelegramTCPBlocking = tcpSuccessCounter.Value() <= 0
	tk.TelegramHTTPBlocking = httpSuccessCounter.Value() <= 0
}
