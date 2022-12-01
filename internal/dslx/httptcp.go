package dslx

//
// TCP adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTCP converts a TCP connection into an HTTP transport.
func HTTPTransportTCP() Function[*ErrorOr[*TCPConnectResultState], *ErrorOr[*HTTPTransportState]] {
	return &httpTransportTCPFunction{}
}

// httpTransportTCPFunction is the function returned by HTTPTransportTCP
type httpTransportTCPFunction struct{}

// Apply implements Function
func (f *httpTransportTCPFunction) Apply(ctx context.Context,
	maybeInput *ErrorOr[*TCPConnectResultState]) *ErrorOr[*HTTPTransportState] {

	// if the previous stage failed, forward the error
	if maybeInput.err != nil {
		return NewErrorOr[*HTTPTransportState](nil, maybeInput.err)
	}
	input := maybeInput.Unwrap()

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
		UnderlyingCloser:      input.Conn,
		ZeroTime:              input.ZeroTime,
	}
	return NewErrorOr(result, nil)
}
