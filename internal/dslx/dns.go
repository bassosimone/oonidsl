package dslx

//
// DNS measurements
//

import (
	"context"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/bassosimone/oonidsl/internal/model"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// DomainName is a domain name to resolve.
type DomainName string

// DNSLookupOption is an option you can pass to DNSLookupInput.
type DNSLookupOption func(*DNSLookupInputState)

// DNSLookupOptionIDGenerator configures a specific ID generator.
// See DNSLookupInputState docs for additional information.
func DNSLookupOptionIDGenerator(value *atomicx.Int64) DNSLookupOption {
	return func(dis *DNSLookupInputState) {
		dis.IDGenerator = value
	}
}

// DNSLookupOptionLogger configures a specific logger.
// See DNSLookupInputState docs for additional information.
func DNSLookupOptionLogger(value model.Logger) DNSLookupOption {
	return func(dis *DNSLookupInputState) {
		dis.Logger = value
	}
}

// DNSLookupOptionZeroTime configures the measurement's zero time.
// See DNSLookupInputState docs for additional information.
func DNSLookupOptionZeroTime(value time.Time) DNSLookupOption {
	return func(dis *DNSLookupInputState) {
		dis.ZeroTime = value
	}
}

// DNSLookupInput creates state for resolving a domain name. The only mandatory
// argument is obviously the domain name to resolve. You can also supply optional
// values by passing options to this function.
func DNSLookupInput(domain DomainName, options ...DNSLookupOption) *DNSLookupInputState {
	state := &DNSLookupInputState{
		Domain:      string(domain),
		IDGenerator: &atomicx.Int64{},
		Logger:      model.DiscardLogger,
		ZeroTime:    time.Now(),
	}
	for _, option := range options {
		option(state)
	}
	return state
}

// DNSLookupInputState contains state for resolving a domain name.
//
// You should construct this type using the DNSLookupInput constructor
// as well as DNSLookupOption options to fill optional values. If you
// want to construct this type manually, please make sure you initialize
// all the variables marked as MANDATORY.
type DNSLookupInputState struct {
	// Domain is the MANDATORY domain name to lookup.
	Domain string

	// IDGenerator is the MANDATORY ID generator. We will use this field
	// to assign unique IDs to distinct sub-measurements. The default
	// construction implemented by DNSLookupInput creates a new generator
	// that starts counting from zero, leading to the first trace having
	// one as its index.
	IDGenerator *atomicx.Int64

	// Logger is the MANDATORY logger to use. The default construction
	// implemented by DNSLookupInput uses model.DiscardLogger.
	Logger model.Logger

	// ZeroTime is the MANDATORY zero time of the measurement. We will
	// use this field as the zero value to compute relative elapsed times
	// when generating measurements. The default construction by
	// DNSLookupInit initializes this field with the current time.
	ZeroTime time.Time
}

// DNSLookupResultState is the state returned by a DNS lookup. This struct
// will obviously contain the results of the DNS lookup as well as state
// inherited from the DNSLookupInputState. If you want to initialize this
// struct manually, make sure you follow specific instructions for each field.
type DNSLookupResultState struct {
	// Addresses contains the nonempty resolved addresses.
	Addresses []string

	// Domain is the domain we resolved. We inherit this field
	// from the value inside the DNSLookupInputState.
	Domain string

	// IDGenerator is the ID generator. We inherit this field
	// from the value inside the DNSLookupInputState.
	IDGenerator *atomicx.Int64

	// Logger is the logger to use. We inherit this field
	// from the value inside the DNSLookupInputState.
	Logger model.Logger

	// Trace is the trace we're currently using. This struct is
	// created by the various Apply functions using values inside
	// the DNSLookupInputState to initialize the Trace.
	Trace *measurexlite.Trace

	// ZeroTime is the zero time of the measurement. We inherit this field
	// from the value inside the DNSLookupInputState.
	ZeroTime time.Time
}

// DNSLookupGetaddrinfo returns a function that resolves a domain name to
// IP addresses using libc's getaddrinfo function.
func DNSLookupGetaddrinfo() Func[*DNSLookupInputState, *Result[*DNSLookupResultState]] {
	return &dnsLookupGetaddrinfoFunc{}
}

// dnsLookupGetaddrinfoFunc is the function returned by DNSLookupGetaddrinfo.
type dnsLookupGetaddrinfoFunc struct{}

// Apply implements Func.
func (f *dnsLookupGetaddrinfoFunc) Apply(
	ctx context.Context, input *DNSLookupInputState) *Result[*DNSLookupResultState] {

	// create trace
	trace := measurexlite.NewTrace(input.IDGenerator.Add(1), input.ZeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		input.Logger,
		"[#%d] DNSLookup[getaddrinfo] %s",
		trace.Index,
		input.Domain,
	)

	// setup
	const timeout = 4 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	resolver := trace.NewStdlibResolver(input.Logger)

	// lookup
	addrs, err := resolver.LookupHost(ctx, input.Domain)

	// stop the operation logger
	ol.Stop(err)

	return &Result[*DNSLookupResultState]{
		Error:        err,
		Observations: maybeTraceToObservations(trace),
		Skipped:      false,
		State: &DNSLookupResultState{
			Addresses:   addrs,
			Domain:      input.Domain,
			IDGenerator: input.IDGenerator,
			Logger:      input.Logger,
			Trace:       trace,
			ZeroTime:    input.ZeroTime,
		},
	}
}

// DNSLookupUDP returns a function that resolves a domain name to
// IP addresses using the given DNS-over-UDP resolver.
func DNSLookupUDP(resolver string) Func[*DNSLookupInputState, *Result[*DNSLookupResultState]] {
	return &dnsLookupUDPFunc{
		Resolver: resolver,
	}
}

// dnsLookupUDPFunc is the type returned by DNSLookupUDP. If you want
// to init this type manually, make sure you set the MANDATORY fields.
type dnsLookupUDPFunc struct {
	// Resolver is the MANDATORY resolver to use.
	Resolver string
}

// Apply implements Func.
func (f *dnsLookupUDPFunc) Apply(
	ctx context.Context, input *DNSLookupInputState) *Result[*DNSLookupResultState] {

	// create trace
	trace := measurexlite.NewTrace(input.IDGenerator.Add(1), input.ZeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		input.Logger,
		"[#%d] DNSLookup[%s/udp] %s",
		trace.Index,
		f.Resolver,
		input.Domain,
	)

	// setup
	const timeout = 4 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	resolver := trace.NewParallelUDPResolver(
		input.Logger,
		netxlite.NewDialerWithoutResolver(input.Logger),
		f.Resolver,
	)

	// lookup
	addrs, err := resolver.LookupHost(ctx, input.Domain)

	// stop the operation logger
	ol.Stop(err)

	return &Result[*DNSLookupResultState]{
		Error:        err,
		Observations: maybeTraceToObservations(trace),
		Skipped:      false,
		State: &DNSLookupResultState{
			Addresses:   addrs, // maybe empty
			Domain:      input.Domain,
			IDGenerator: input.IDGenerator,
			Logger:      input.Logger,
			Trace:       trace,
			ZeroTime:    input.ZeroTime,
		},
	}
}
