package main

//
// Top-level measurement algorithm
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// measure is the top-level measurement algorithm.
func measure(ctx context.Context, state *measurementState) error {
	// run DCs measurements in background
	state.wg.Add(1)

	// dnslookup
	go measureTargets(ctx, state)

	// wait for measurements to terminate
	state.wg.Wait()

	// make sure we fail the measurement if the main
	// context is cancelled (e.g., because the user
	// has hit ^C and has forced an early termination)
	return ctx.Err()
}

func measureTargets(ctx context.Context, state *measurementState) {
	defer state.wg.Done()

	domains := []string{
		"textsecure-service.whispersystems.org",
		"storage.signal.org",
		"api.directory.signal.org",
		"cdn.signal.org",
		"cdn2.signal.org",
		"sfu.voip.signal.org",
	}

	// construct getaddrinfo resolver
	lookup := dslx.DNSLookupGetaddrinfo()

	// create the established connections pool
	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	certPool, err := newCertPool()
	if err != nil {
		// TODO
	}

	var successes *dslx.CounterState[*dslx.HTTPRequestResultState]
	var httpsResults []fx.Result[*dslx.HTTPRequestResultState]

	for _, d := range domains {
		// describe the DNS measurement input
		dnsInput := dslx.DNSLookupInput(
			dslx.DomainName(d),
			dslx.DNSLookupOptionIDGenerator(state.idGen), // do I have to increment this?
			dslx.DNSLookupOptionLogger(state.logger),
			dslx.DNSLookupOptionZeroTime(state.zeroTime),
		)
		// run the DNS Lookup
		dnsResults := fx.Map(
			ctx,
			fx.Parallelism(3),
			lookup,
			dnsInput,
		)
		// extract and merge observations with the test keys
		state.tk.mergeObservations(dslx.ExtractObservations(dnsResults...)...)

		// if the lookup has failed mark the whole web measurement as failed
		// TODO

		// obtain a unique set of IP addresses w/o bogons inside it
		ipAddrs := dslx.AddressSet(dnsResults...).RemoveBogons()

		// create the set of endpoints
		endpoints := ipAddrs.ToEndpoints(
			dslx.EndpointNetwork("tcp"),
			dslx.EndpointPort(443),
			dslx.EndpointOptionDomain(d),
			dslx.EndpointOptionIDGenerator(state.idGen),
			dslx.EndpointOptionLogger(state.logger),
			dslx.EndpointOptionZeroTime(state.zeroTime),
		)

		// count the number of successes
		successes = dslx.Counter[*dslx.HTTPRequestResultState]()

		// create function for the 443/tcp/tls/https measurement
		httpsFunction := fx.ComposeFlat6(
			dslx.TCPConnect(connpool),
			dslx.TLSHandshake(
				connpool,
				dslx.TLSHandshakeOptionRootCAs(certPool),
			),
			dslx.HTTPTransportTLS(),
			dslx.HTTPJustUseOneConn(), // stop subsequent connections
			dslx.HTTPRequest(),
			successes.Func(), // number of times we arrive here
		)

		// run 443/tcp/tls/https measurement
		httpsResults = fx.Map(
			ctx,
			fx.Parallelism(2),
			httpsFunction,
			endpoints...,
		)

		// extract and merge observations with the test keys
		state.tk.mergeObservations(dslx.ExtractObservations(httpsResults...)...)
	}

	// if we saw successes, then it's not blocked
	// TODO: re-define what success means here!
	if successes.Value() > 0 {
		state.tk.setResultSuccess()
		return
	}

	// attempt to set a meaningful error, if that's possible
	if err := dslx.FirstErrorExcludingBrokenIPv6Errors(httpsResults...); err != nil {
		state.tk.setResultFailure(err)
		return
	}

	// otherwise fallback to whatever is the first error
	if err := dslx.FirstError(httpsResults...); err != nil {
		state.tk.setResultFailure(err)
		return
	}

	// the last resort is to set an unknown failure error
	state.tk.setResultFailure(netxlite.ErrUnknown)

}

// setResultSuccess sets the result of the experiment in case of success
func (tk *testKeys) setResultSuccess() {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	tk.SignalBackendFailure = nil
	tk.SignalBackendStatus = "ok"
}

// setResultFailure sets the result of the experiment in case of failure
func (tk *testKeys) setResultFailure(err error) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	s := err.Error()
	tk.SignalBackendFailure = &s
	tk.SignalBackendStatus = "blocked"
}
