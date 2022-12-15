package main

//
// Top-level measurement algorithm
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

// measure is the top-level measurement algorithm.
func measure(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
) error {
	domains := []string{
		"textsecure-service.whispersystems.org",
		"storage.signal.org",
		"api.directory.signal.org",
		"cdn.signal.org",
		"cdn2.signal.org",
		"sfu.voip.signal.org",
	}
	errch := make(chan error, len(domains))
	wg := &sync.WaitGroup{}

	for _, domain := range domains {
		wg.Add(1)
		go measureTarget(ctx, logger, idGen, zeroTime, tk, wg, domain, errch)
	}

	// wait for measurements to terminate
	wg.Wait()

	for {
		select {
		case e := <-errch:
			if e != nil {
				tk.setResultFailure(e)
				return ctx.Err()
			}
		default:
			tk.setResultSuccess()
			return ctx.Err()
		}
	}
}

func measureTarget(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	wg *sync.WaitGroup,
	domain string,
	errch chan error,
) {
	defer wg.Done()
	errch <- doMeasureTarget(ctx, logger, idGen, zeroTime, tk, wg, domain)
}

func doMeasureTarget(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	wg *sync.WaitGroup,
	domain string,
) error {
	// describe the DNS measurement input
	dnsInput := dslx.DNSLookupInput(
		dslx.DomainName(domain),
		dslx.DNSLookupOptionIDGenerator(idGen),
		dslx.DNSLookupOptionLogger(logger),
		dslx.DNSLookupOptionZeroTime(zeroTime),
	)
	// construct getaddrinfo resolver
	lookup := dslx.DNSLookupGetaddrinfo()
	// run the DNS Lookup
	dnsResult := lookup.Apply(ctx, dnsInput)

	// extract and merge observations with the test keys
	tk.mergeObservations(dslx.ExtractObservations(dnsResult)...)

	// if the lookup has failed we return
	if dnsResult.Error != nil {
		return dnsResult.Error
	}

	// obtain a unique set of IP addresses w/o bogons inside it
	ipAddrs := dslx.AddressSet(dnsResult).RemoveBogons()

	// create the set of endpoints
	endpoints := ipAddrs.ToEndpoints(
		dslx.EndpointNetwork("tcp"),
		dslx.EndpointPort(443),
		dslx.EndpointOptionDomain(domain),
		dslx.EndpointOptionIDGenerator(idGen),
		dslx.EndpointOptionLogger(logger),
		dslx.EndpointOptionZeroTime(zeroTime),
	)

	// count the number of successes
	successes := dslx.Counter[*dslx.HTTPRequestResultState]()

	// create the established connections pool
	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	// create the certificate pool
	certPool, err := newCertPool()
	if err != nil {
		// TODO
	}

	// create function for the 443/tcp/tls/https measurement
	httpsFunction := dslx.Compose6(
		dslx.TCPConnect(connpool),
		dslx.TLSHandshake(
			connpool,
			dslx.TLSHandshakeOptionRootCAs(certPool),
		),
		dslx.HTTPTransportTLS(),
		dslx.HTTPJustUseOneConn(), // TODO: do we want this?
		dslx.HTTPRequest(),
		successes.Func(), // number of times we arrive here
	)

	// run 443/tcp/tls/https measurement
	httpsResults := dslx.Map(
		ctx,
		dslx.Parallelism(2),
		httpsFunction,
		endpoints...,
	)

	// extract and merge observations with the test keys
	tk.mergeObservations(dslx.ExtractObservations(httpsResults...)...)

	// if we saw successes, then this domain is not blocked
	if successes.Value() > 0 {
		return nil
	}

	// attempt to set a meaningful error, if that's possible
	if err := dslx.FirstErrorExcludingBrokenIPv6Errors(httpsResults...); err != nil {
		return err
	}

	// otherwise fallback to whatever is the first error
	if err := dslx.FirstError(httpsResults...); err != nil {
		return err
	}

	// the last resort is to set an unknown failure error
	return netxlite.ErrUnknown
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
