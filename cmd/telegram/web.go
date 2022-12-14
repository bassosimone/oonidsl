package main

//
// Measuring web.telegram.org
//

import (
	"context"
	"sync"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/model"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// measureWeb measures telegram web.
func measureWeb(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	wg *sync.WaitGroup,
) {
	const webDomain = "web.telegram.org"

	// tell the parent we terminated
	defer wg.Done()

	// describe the DNS measurement input
	dnsInput := dslx.NewDomainToResolve(
		dslx.DomainName(webDomain),
		dslx.DNSLookupOptionIDGenerator(idGen),
		dslx.DNSLookupOptionLogger(logger),
		dslx.DNSLookupOptionZeroTime(zeroTime),
	)

	// construct getaddrinfo resolver
	getaddrinfoResolver := dslx.DNSLookupGetaddrinfo()

	// perform the DNS lookup
	dnsResults := getaddrinfoResolver.Apply(ctx, dnsInput)

	// extract and merge observations with the test keys
	tk.mergeObservations(dslx.ExtractObservations(dnsResults)...)

	// if the lookup has failed mark the whole web measurement as failed
	if err := dnsResults.Error; err != nil {
		tk.setWebResultFailure(err)
		return
	}

	// obtain a unique set of IP addresses w/o bogons inside it
	ipAddrs := dslx.NewAddressSet(dnsResults).RemoveBogons()

	// if the set is empty we only got bogons
	if len(ipAddrs.M) <= 0 {
		tk.setWebResultFailure(netxlite.ErrDNSBogon)
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
		dslx.EndpointOptionLogger(logger),
		dslx.EndpointOptionZeroTime(zeroTime),
	)

	// count the number of successes
	successes := dslx.Counter[*dslx.HTTPResponse]()

	// create function for the 443/tcp measurement
	httpsFunction := dslx.Compose6(
		dslx.TCPConnect(connpool),
		dslx.TLSHandshake(connpool),
		dslx.HTTPTransportTLS(),
		dslx.HTTPJustUseOneConn(), // stop subsequent connections
		dslx.HTTPRequest(),
		successes.Func(), // number of times we arrive here
	)

	// run 443/tcp measurement
	httpsResults := dslx.Map(
		ctx,
		dslx.Parallelism(2),
		httpsFunction,
		httpsEndpoints...,
	)

	// extract and merge observations with the test keys
	tk.mergeObservations(dslx.ExtractObservations(httpsResults...)...)

	// if we saw successes, then it's not blocked
	if successes.Value() > 0 {
		tk.setWebResultSuccess()
		return
	}

	// attempt to set a meaningful error, if that's possible
	if err := dslx.FirstErrorExcludingBrokenIPv6Errors(httpsResults...); err != nil {
		tk.setWebResultFailure(err)
		return
	}

	// otherwise fallback to whatever is the first error
	if err := dslx.FirstError(httpsResults...); err != nil {
		tk.setWebResultFailure(err)
		return
	}

	// the last resort is to set an unknown failure error
	tk.setWebResultFailure(netxlite.ErrUnknown)
}

// setWebResultSuccess sets the result of the web experiment in case of success
func (tk *testKeys) setWebResultSuccess() {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	tk.TelegramWebFailure = nil
	tk.TelegramWebStatus = "ok"
}

// setWebResultFailure sets the result of the web experiment in case of failure
func (tk *testKeys) setWebResultFailure(err error) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	s := err.Error()
	tk.TelegramWebFailure = &s
	tk.TelegramWebStatus = "blocked"
}
