package dslx

//
// TLS adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTLS converts a TLS connection into an HTTP transport.
func HTTPTransportTLS() Func[*TLSHandshakeResultState, *Result[*HTTPTransportState]] {
	return &httpTransportTLSFunc{}
}

// httpTransportTLSFunc is the function returned by HTTPTransportTLS.
type httpTransportTLSFunc struct{}

// Apply implements Func.
func (f *httpTransportTLSFunc) Apply(
	ctx context.Context, input *TLSHandshakeResultState) *Result[*HTTPTransportState] {
	httpTransport := netxlite.NewHTTPTransport(
		input.Logger,
		netxlite.NewNullDialer(),
		netxlite.NewSingleUseTLSDialer(input.Conn),
	)
	return &Result[*HTTPTransportState]{
		Error:        nil,
		Observations: nil,
		Skipped:      false,
		State: &HTTPTransportState{
			Address:               input.Address,
			Domain:                input.Domain,
			IDGenerator:           input.IDGenerator,
			Logger:                input.Logger,
			Network:               input.Network,
			Scheme:                "https",
			TLSNegotiatedProtocol: input.TLSState.NegotiatedProtocol,
			Trace:                 input.Trace,
			Transport:             httpTransport,
			ZeroTime:              input.ZeroTime,
		},
	}
}
