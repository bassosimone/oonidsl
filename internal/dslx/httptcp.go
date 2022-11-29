package dslx

//
// TCP adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTCP converts a TCP connection into an HTTP transport.
func HTTPTransportTCP() Function[*TCPConnectResultState, *HTTPTransportState] {
	return &httpTransportTCPFunction{}
}

// httpTransportTCPFunction is the function returned by HTTPTransportTCP
type httpTransportTCPFunction struct{}

// Apply implements Function
func (f *httpTransportTCPFunction) Apply(
	ctx context.Context, input *TCPConnectResultState) *HTTPTransportState {

	// if the previous stage failed, forward the error
	if input.Err != nil {
		return &HTTPTransportState{
			Address:               input.Address,
			Domain:                input.Domain,
			Err:                   input.Err,
			IDGenerator:           input.IDGenerator,
			Logger:                input.Logger,
			Network:               input.Network,
			Scheme:                "",
			TLSNegotiatedProtocol: "",
			Trace:                 input.Trace,
			Transport:             nil,
			UnderlyingCloser:      nil,
			ZeroTime:              input.ZeroTime,
		}
	}

	// create transport
	httpTransport := netxlite.NewHTTPTransport(
		input.Logger,
		netxlite.NewSingleUseDialer(input.Conn),
		netxlite.NewNullTLSDialer(),
	)

	return &HTTPTransportState{
		Address:               input.Address,
		Domain:                input.Domain,
		Err:                   nil,
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
}
