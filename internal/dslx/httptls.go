package dslx

//
// TLS adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTLS converts a TLS connection into an HTTP transport.
func HTTPTransportTLS() fx.Func[*TLSHandshakeResultState, fx.Result[*HTTPTransportState]] {
	return &httpTransportTLSFunc{}
}

// httpTransportTLSFunc is the function returned by HTTPTransportTLS.
type httpTransportTLSFunc struct{}

// Apply implements Func.
func (f *httpTransportTLSFunc) Apply(
	ctx context.Context, input *TLSHandshakeResultState) fx.Result[*HTTPTransportState] {
	// create transport
	httpTransport := netxlite.NewHTTPTransport(
		input.Logger,
		netxlite.NewNullDialer(),
		netxlite.NewSingleUseTLSDialer(input.Conn),
	)

	result := &HTTPTransportState{
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
	}
	return fx.Ok(result)
}
