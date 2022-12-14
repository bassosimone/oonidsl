package dslx

//
// Collecting observations
//

import (
	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/bassosimone/oonidsl/internal/model"
)

// Observations is the skeleton shared by most OONI measurements where
// we group observations by type using standard test keys.
type Observations struct {
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
}

// ObservationsProducer is anything from which we can extract observations.
type ObservationsProducer interface {
	// Observations exctracts and returns observations. This function
	// MUST have a once semantics: after the first call it MUST return
	// a nil or zero-length slice to the caller.
	Observations() []*Observations
}

// ExtractObservations extracts observations from a list of producers.
func ExtractObservations[T ObservationsProducer](producers ...fx.Result[T]) (out []*Observations) {
	for _, p := range producers {
		if p.IsErr() {
			continue
		}
		v := p.Unwrap()
		out = append(out, v.Observations()...)
	}
	return
}

// maybeTraceToObservations returns the observations inside the
// trace taking into account the case where trace is nil.
func maybeTraceToObservations(trace *measurexlite.Trace) (out []*Observations) {
	if trace != nil {
		out = append(out, &Observations{
			NetworkEvents: trace.NetworkEvents(),
			Queries:       trace.DNSLookupsFromRoundTrip(),
			Requests:      []*model.ArchivalHTTPRequestResult{}, // no extractor inside trace!
			TCPConnect:    trace.TCPConnects(),
			TLSHandshakes: trace.TLSHandshakes(),
		})
	}
	return
}

// ObservationsCollector is anything where we can store observations, e.g. TestKeys.
type ObservationsCollector interface {
	// MergeObservations merges collected observations into the test keys.
	// When implementing MergeObservations, the programmer is responsible to make this method
	// go routine safe.
	MergeObservations(obs ...*Observations)
}
