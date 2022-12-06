package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/model"
	"github.com/bassosimone/oonidsl/internal/netxlite"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

// testKeys contains the experiment testKeys
type testKeys struct {
	// NetworkEvents contains I/O events.
	NetworkEvents []*model.ArchivalNetworkEvent `json:"network_events"`

	// Queries contains the DNS queries results.
	Queries []*model.ArchivalDNSLookupResult `json:"queries"`

	// Requests contains HTTP request results.
	Requests []*model.ArchivalHTTPRequestResult `json:"requests"`

	// TCPConnect contains the TCP connect results.
	TCPConnect []*model.ArchivalTCPConnectResult `json:"tcp_connect"`

	// TLSHandshakes contains the TLS handshakes results.
	TLSHandshakes []*model.ArchivalTLSOrQUICHandshakeResult `json:"tls_handshakes"`

	// TelegramTCPBlocking indicates whether DCs are
	// blocked using TCP/IP interference.
	TelegramTCPBlocking bool `json:"telegram_tcp_blocking"`

	// TelegramHTTPBlocking indicates whether DCs are
	// blocked using HTTP interference.
	TelegramHTTPBlocking bool `json:"telegram_http_blocking"`

	// TelegramWebFailure is the failure in accessing telegram web.
	TelegramWebFailure *string `json:"telegram_web_failure"`

	// TelegramWebStatus is the status of telegram web.
	TelegramWebStatus string `json:"telegram_web_status"`

	// mu provides mutual exclusion
	mu sync.Mutex
}

// mergeObservations merges collected observations into the test keys
func (tk *testKeys) mergeObservations(obs ...*dslx.Observations) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	for _, o := range obs {
		tk.NetworkEvents = append(tk.NetworkEvents, o.NetworkEvents...)
		tk.Queries = append(tk.Queries, o.Queries...)
		tk.Requests = append(tk.Requests, o.Requests...)
		tk.TCPConnect = append(tk.TCPConnect, o.TCPConnect...)
		tk.TLSHandshakes = append(tk.TLSHandshakes, o.TLSHandshakes...)
	}
}

// setDCBlocking sets blocking rules depending on DC results
func (tk *testKeys) setDCBlocking(
	tcpSuccess *dslx.CounterState[*dslx.TCPConnectResultState],
	httpSuccess *dslx.CounterState[*dslx.HTTPRequestResultState],
) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	tk.TelegramTCPBlocking = tcpSuccess.Value() <= 0
	tk.TelegramHTTPBlocking = httpSuccess.Value() <= 0
}

// setWebResultFailure results the result of the web experiment in case of failure
func (tk *testKeys) setWebResultFailure(err error) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	s := err.Error()
	tk.TelegramWebFailure = &s
	tk.TelegramWebStatus = "blocked"
}

// measureDCs measures telegram access points
func measureDCs(
	ctx context.Context,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	wg *sync.WaitGroup,
) {
	// tell the parent we terminated
	defer wg.Done()

	// ipAddrs contains the access points IP addresses
	var ipAddrs = dslx.AddressSet().Add(
		"149.154.175.50",
		"149.154.167.51",
		"149.154.175.100",
		"149.154.167.91",
		"149.154.171.5",
	)

	// construct the list of endpoints to measure
	var (
		endpoints []*dslx.EndpointState
		ports     = []int{80, 443}
	)
	for addr := range ipAddrs.M {
		for _, port := range ports {
			endpoints = append(endpoints, dslx.Endpoint(
				dslx.EndpointNetwork("tcp"),
				dslx.EndpointAddress(net.JoinHostPort(addr, strconv.Itoa(port))),
				dslx.EndpointOptionIDGenerator(idGen),
				dslx.EndpointOptionLogger(log.Log),
				dslx.EndpointOptionZerotime(zeroTime),
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

	// construct the function to measure endpoints
	function := fx.ComposeFlat5(
		dslx.TCPConnect(connpool),
		tcpConnectSuccessCounter.Func(), // count number of times we reach this point
		dslx.HTTPTransportTCP(),
		dslx.HTTPRequest(
			dslx.HTTPRequestOptionMethod("POST"),
		),
		httpRoundTripSuccessCounter.Func(), // count number of times we reach this point
	)

	// measure all the endpoints in parallel and collect the results
	results := fx.Map(
		ctx,
		fx.Parallelism(4),
		function,
		endpoints...,
	)

	// extract observations from the above measurement
	observations := dslx.ExtractObservations(results...)

	// merge observations with the test keys
	tk.mergeObservations(observations...)

	// set top-level keys indicating DCs blocking
	tk.setDCBlocking(tcpConnectSuccessCounter, httpRoundTripSuccessCounter)
}

// measureWeb measures telegram web
func measureWeb(
	ctx context.Context,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	wg *sync.WaitGroup,
) {
	const webDomain = "web.telegram.org"

	// tell the parent we terminated
	defer wg.Done()

	// describe the DNS measurement input
	dnsInput := dslx.DNSLookupInput(
		dslx.DomainName(webDomain),
		dslx.DNSLookupOptionIDGenerator(idGen),
		dslx.DNSLookupOptionLogger(log.Log),
		dslx.DNSLookupOptionZeroTime(zeroTime),
	)

	// construct getaddrinfo resolver
	getaddrinfoResolver := dslx.DNSLookupGetaddrinfo()

	// perform the DNS lookup
	dnsResults := getaddrinfoResolver.Apply(ctx, dnsInput)

	// if the lookup has failed mark the whole web measurement as failed
	if dnsResults.IsErr() {
		tk.setWebResultFailure(dnsResults.UnwrapErr())
		return
	}

	// obtain a unique set of IP addresses w/o bogons inside it
	ipAddrs := dslx.AddressSet(dnsResults).RemoveBogons()

	// if the set is empty we only got bogons
	if len(ipAddrs.M) <= 0 {
		tk.setWebResultFailure(netxlite.ErrDNSBogon)
		return
	}

	// create endpoints for the 80/tcp measurement
	httpEndpoints := ipAddrs.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(80),
		dslx.EndpointOptionDomain(webDomain),
		dslx.EndpointOptionIDGenerator(idGen),
		dslx.EndpointOptionLogger(log.Log),
		dslx.EndpointOptionZerotime(zeroTime),
	)

	// create the established connections pool
	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	// create function for the 80/tcp measurement
	httpFunction := fx.ComposeFlat4(
		dslx.TCPConnect(connpool),
		dslx.HTTPTransportTCP(),
		dslx.HTTPRequest(),
		fx.Lambda(func(ctx context.Context, state *dslx.HTTPRequestResultState) fx.Result[*dslx.HTTPRequestResultState] {
			// TODO(bassosimone): analyze the HTTP response here.
			return fx.Ok(state)
		}),
	)

	// start 80/tcp measurement in async fashion
	httpResultsAsync := fx.MapAsync(
		ctx,
		fx.Parallelism(2),
		httpFunction,
		fx.Stream(httpEndpoints...),
	)

	// create endpoints for the 443/tcp measurement
	httpsEndpoints := ipAddrs.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(443),
		dslx.EndpointOptionDomain(webDomain),
		dslx.EndpointOptionIDGenerator(idGen),
		dslx.EndpointOptionLogger(log.Log),
		dslx.EndpointOptionZerotime(zeroTime),
	)

	// track TLS handshake errors in particular which allows us
	// potentially to perform follow up experiments
	tlsHandshakeErrs := &dslx.ErrorLogger{}

	// create function for the 443/tcp measurement
	httpsFunction := fx.ComposeFlat5(
		dslx.TCPConnect(connpool),
		dslx.RecordErrors(
			tlsHandshakeErrs,
			dslx.TLSHandshake(connpool),
		),
		dslx.HTTPTransportTLS(),
		dslx.HTTPRequest(),
		fx.Lambda(func(ctx context.Context, state *dslx.HTTPRequestResultState) fx.Result[*dslx.HTTPRequestResultState] {
			// TODO(bassosimone): analyze the HTTP response here.
			return fx.Ok(state)
		}),
	)

	// start 443/tcp measurement in async fashion
	httpsResultsAsync := fx.MapAsync(
		ctx,
		fx.Parallelism(2),
		httpsFunction,
		fx.Stream(httpsEndpoints...),
	)

	// await for completion of HTTP and HTTPS measurements
	_ = fx.ZipAndCollect(httpResultsAsync, httpsResultsAsync)
}

func main() {
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	idGen := &atomicx.Int64{}
	tk := &testKeys{}
	zeroTime := time.Now()

	wg.Add(1)
	go measureDCs(ctx, idGen, zeroTime, tk, wg)

	wg.Add(1)
	go measureWeb(ctx, idGen, zeroTime, tk, wg)

	wg.Wait()

	data, err := json.Marshal(tk)
	runtimex.PanicOnError(err, "json.Marshal failed unexpectedly")
	fmt.Printf("%s\n", string(data))
}
