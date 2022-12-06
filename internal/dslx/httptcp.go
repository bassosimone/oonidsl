package dslx

//
// TCP adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTCP converts a TCP connection into an HTTP transport.
func HTTPTransportTCP() fx.Func[*TCPConnectResultState, fx.Result[*HTTPTransportState]] {
	return &httpTransportTCPFunc{}
}

// httpTransportTCPFunc is the function returned by HTTPTransportTCP
type httpTransportTCPFunc struct{}

// Apply implements Func
func (f *httpTransportTCPFunc) Apply(
	ctx context.Context, input *TCPConnectResultState) fx.Result[*HTTPTransportState] {
	// create transport
	httpTransport := netxlite.NewHTTPTransport(
		input.Logger,
		netxlite.NewSingleUseDialer(input.Conn),
		netxlite.NewNullTLSDialer(),
	)

	result := &HTTPTransportState{
		Address:               input.Address,
		Domain:                input.Domain,
		IDGenerator:           input.IDGenerator,
		Logger:                input.Logger,
		Network:               input.Network,
		Scheme:                "http",
		TLSNegotiatedProtocol: "",
		Trace:                 input.Trace,
		Transport:             httpTransport,
		ZeroTime:              input.ZeroTime,
	}
	return fx.Ok(result)
}
