package dslx

//
// TLS adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTLS converts a TLS connection into an HTTP transport.
func HTTPTransportTLS() Function[*TLSHandshakeResultState, *HTTPTransportState] {
	return &httpTransportTLSFunction{}
}

// httpTransportTLSFunction is the function returned by HTTPTransportTLS.
type httpTransportTLSFunction struct{}

// Apply implements Function.
func (f *httpTransportTLSFunction) Apply(
	ctx context.Context, input *TLSHandshakeResultState) *HTTPTransportState {

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
		netxlite.NewNullDialer(),
		netxlite.NewSingleUseTLSDialer(input.Conn),
	)

	return &HTTPTransportState{
		Address:               input.Address,
		Domain:                input.Domain,
		Err:                   nil,
		IDGenerator:           input.IDGenerator,
		Logger:                input.Logger,
		Network:               input.Network,
		Scheme:                "https",
		TLSNegotiatedProtocol: input.TLSState.NegotiatedProtocol,
		Trace:                 input.Trace,
		Transport:             httpTransport,
		UnderlyingCloser:      input.Conn,
		ZeroTime:              input.ZeroTime,
	}
}
