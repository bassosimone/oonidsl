package main

//
// Measuring web.telegram.org
//

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// measureWeb measures telegram web.
//
// Arguments:
//
// - ctx is the context;
//
// - idGen allows to assign unique IDs to submeasurements;
//
// - zeroTime is the "zero time" of the measurement;
//
// - tk contains the experiment results;
//
// - wg allows us to synchronize with our parent.
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

	// extract and merge observations with the test keys
	tk.mergeObservations(dslx.ExtractObservations(dnsResults)...)

	// if the lookup has failed mark the whole web measurement as failed
	if dnsResults.IsErr() {
		setWebResultFailure(tk, dnsResults.UnwrapErr())
		return
	}

	// obtain a unique set of IP addresses w/o bogons inside it
	ipAddrs := dslx.AddressSet(dnsResults).RemoveBogons()

	// if the set is empty we only got bogons
	if len(ipAddrs.M) <= 0 {
		setWebResultFailure(tk, netxlite.ErrDNSBogon)
		return
	}

	// create the established connections pool
	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	// create endpoints for the 443/tcp measurement
	httpsEndpoints := ipAddrs.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(443),
		dslx.EndpointOptionDomain(webDomain),
		dslx.EndpointOptionIDGenerator(idGen),
		dslx.EndpointOptionLogger(log.Log),
		dslx.EndpointOptionZerotime(zeroTime),
	)

	// create function for the 443/tcp measurement
	httpsFunction := fx.ComposeFlat5(
		dslx.TCPConnect(connpool),
		dslx.TLSHandshake(connpool),
		dslx.HTTPTransportTLS(),
		dslx.HTTPJustUseOneConn(), // stop subsequent connections
		dslx.HTTPRequest(),
	)

	// start 443/tcp measurement in async fashion
	httpsResults := fx.Map(
		ctx,
		fx.Parallelism(2),
		httpsFunction,
		httpsEndpoints...,
	)

	// extract and merge observations with the test keys
	tk.mergeObservations(dslx.ExtractObservations(httpsResults...)...)

	// TODO(bassosimone): here we should set the web failure
	// TODO(bassosimone): we should filter failed TCP
	// connect attempts caused by missing IPv6
}

// setWebResultFailure results the result of the web experiment in case of failure
func setWebResultFailure(tk *testKeys, err error) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	s := err.Error()
	tk.TelegramWebFailure = &s
	tk.TelegramWebStatus = "blocked"
}
