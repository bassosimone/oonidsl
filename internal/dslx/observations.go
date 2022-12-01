package dslx

//
// Collecting observations
//

import (
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
func ExtractObservations[T ObservationsProducer](producers ...ErrorOr[T]) (out []*Observations) {
	for _, p := range producers {
		if p.Error() != nil {
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
