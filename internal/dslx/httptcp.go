package dslx

//
// TCP adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTCP converts a TCP connection into an HTTP transport.
func HTTPTransportTCP() Function[*TCPConnectResultState, ErrorOr[*HTTPTransportState]] {
	return &httpTransportTCPFunction{}
}

// httpTransportTCPFunction is the function returned by HTTPTransportTCP
type httpTransportTCPFunction struct{}

// Apply implements Function
func (f *httpTransportTCPFunction) Apply(
	ctx context.Context, input *TCPConnectResultState) ErrorOr[*HTTPTransportState] {
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
	return NewErrorOr(result, nil)
}
