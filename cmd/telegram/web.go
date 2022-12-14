package main

//
// Measuring web.telegram.org
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// measureWeb measures telegram web.
func measureWeb(ctx context.Context, state *measurementState) {
	const webDomain = "web.telegram.org"

	// tell the parent we terminated
	defer state.wg.Done()

	// describe the DNS measurement input
	dnsInput := dslx.DNSLookupInput(
		dslx.DomainName(webDomain),
		dslx.DNSLookupOptionIDGenerator(state.idGen),
		dslx.DNSLookupOptionLogger(state.logger),
		dslx.DNSLookupOptionZeroTime(state.zeroTime),
	)

	// construct getaddrinfo resolver
	getaddrinfoResolver := dslx.DNSLookupGetaddrinfo(state.tk)

	// perform the DNS lookup
	dnsResults := getaddrinfoResolver.Apply(ctx, dnsInput)

	// if the lookup has failed mark the whole web measurement as failed
	if dnsResults.IsErr() {
		state.tk.setWebResultFailure(dnsResults.UnwrapErr())
		return
	}

	// obtain a unique set of IP addresses w/o bogons inside it
	ipAddrs := dslx.AddressSet(dnsResults).RemoveBogons()

	// if the set is empty we only got bogons
	if len(ipAddrs.M) <= 0 {
		state.tk.setWebResultFailure(netxlite.ErrDNSBogon)
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
		dslx.EndpointOptionIDGenerator(state.idGen),
		dslx.EndpointOptionLogger(state.logger),
		dslx.EndpointOptionZeroTime(state.zeroTime),
	)

	// count the number of successes
	successes := dslx.Counter[*dslx.HTTPRequestResultState]()

	// create function for the 443/tcp measurement
	httpsFunction := fx.ComposeResult6(
		dslx.TCPConnect(connpool, state.tk),
		dslx.TLSHandshake(connpool, state.tk),
		dslx.HTTPTransportTLS(),
		dslx.HTTPJustUseOneConn(), // stop subsequent connections
		dslx.HTTPRequest(state.tk),
		successes.Func(), // number of times we arrive here
	)

	// run 443/tcp measurement
	httpsResults := fx.Map(
		ctx,
		fx.Parallelism(2),
		httpsFunction,
		httpsEndpoints...,
	)

	// if we saw successes, then it's not blocked
	if successes.Value() > 0 {
		state.tk.setWebResultSuccess()
		return
	}

	// attempt to set a meaningful error, if that's possible
	if err := dslx.FirstErrorExcludingBrokenIPv6Errors(httpsResults...); err != nil {
		state.tk.setWebResultFailure(err)
		return
	}

	// otherwise fallback to whatever is the first error
	if err := dslx.FirstError(httpsResults...); err != nil {
		state.tk.setWebResultFailure(err)
		return
	}

	// the last resort is to set an unknown failure error
	state.tk.setWebResultFailure(netxlite.ErrUnknown)
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
