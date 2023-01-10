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

	// NetworkEvents contains network events.
	NetworkEvents []*model.ArchivalNetworkEvent `json:"network_events"`

	// Queries contains DNS queries.
	Queries []*model.ArchivalDNSLookupResult `json:"queries"`

	// Requests contains HTTP results.
	Requests []*model.ArchivalHTTPRequestResult `json:"requests"`

	// QUICHandshakes contains QUIC handshakes results.
	QUICHandshakes []*model.ArchivalTLSOrQUICHandshakeResult `json:"quic_handshakes"`

	// Failure contains the failure of the experiment.
	Failure *string `json:"failure"`

	// mu provides mutual exclusion for accessing the test keys.
	mu sync.Mutex
}

// mergeObservations merges collected observations into the test keys.
func (tk *testKeys) mergeObservations(obs ...*dslx.Observations) {
	defer tk.mu.Unlock()
	tk.mu.Lock()
	for _, o := range obs {
		tk.NetworkEvents = append(tk.NetworkEvents, o.NetworkEvents...)
		tk.Queries = append(tk.Queries, o.Queries...)
		tk.Requests = append(tk.Requests, o.Requests...)
		tk.QUICHandshakes = append(tk.QUICHandshakes, o.QUICHandshakes...)
	}
}
