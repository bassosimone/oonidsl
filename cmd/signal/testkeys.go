package main

//
// Experiment results (aka "test keys")
//

import (
	"sync"

	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/model"
)

// testKeys contains the experiment test keys.
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

	SignalBackendStatus string `json:"signal_backend_status"`

	SignalBackendFailure *string `json:"signal_backend_failure"`

	// mu provides mutual exclusion.
	mu sync.Mutex
}

// MergeObservations implements ObservationsCollector.MergeObservations.
func (tk *testKeys) MergeObservations(obs ...*dslx.Observations) {
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
