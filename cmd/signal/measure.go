package main

//
// Top-level measurement algorithm
//

import (
	"context"
	"crypto/x509"
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
	certPool, err := newCertPool()
	if err != nil {
		return err // fundamental error, let's not submit
	}

	domains := []string{
		"textsecure-service.whispersystems.org",
		"storage.signal.org",
		"api.directory.signal.org",
		"cdn.signal.org",
		"cdn2.signal.org",
		"sfu.voip.signal.org",
	}

	// run measurements in parallel
	errch := make(chan error)
	for _, domain := range domains {
		go measureTarget(ctx, logger, idGen, zeroTime, tk, domain, certPool, errch)
	}

	// collect the result of each measurement
	var errors []error
	for range domains {
		errors = append(errors, <-errch)
	}

	// set the final result
	for _, err := range errors {
		if err != nil {
			tk.setResultFailure(err)
			return ctx.Err()
		}
	}
	tk.setResultSuccess()
	return ctx.Err()
}

func measureTarget(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	domain string,
	certPool *x509.CertPool,
	errch chan error,
) {
	// Note: this pattern ensures we write the output channel exactly once
	errch <- doMeasureTarget(ctx, logger, idGen, zeroTime, tk, domain, certPool)
}

func doMeasureTarget(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	domain string,
	certPool *x509.CertPool,
) error {
	// describe the DNS measurement input
	dnsInput := dslx.NewDomainToResolve(
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
	ipAddrs := dslx.NewAddressSet(dnsResult).RemoveBogons()

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
	successes := dslx.Counter[*dslx.HTTPResponse]()

	// create the established connections pool
	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	// create function for the 443/tcp/tls/https measurement
	httpsFunction := dslx.Compose6(
		dslx.TCPConnect(connpool),
		dslx.TLSHandshake(
			connpool,
			dslx.TLSHandshakeOptionRootCAs(certPool),
		),
		dslx.HTTPTransportTLS(),
		dslx.HTTPJustUseOneConn(), // TODO(bassosimone): do we want this?
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
