package dslx

//
// TLS adapters for HTTP
//

import (
	"context"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportTLS converts a TLS connection into an HTTP transport.
func HTTPTransportTLS() Function[
	*ErrorOr[*TLSHandshakeResultState], *ErrorOr[*HTTPTransportState]] {
	return &httpTransportTLSFunction{}
}

// httpTransportTLSFunction is the function returned by HTTPTransportTLS.
type httpTransportTLSFunction struct{}

// Apply implements Function.
func (f *httpTransportTLSFunction) Apply(ctx context.Context,
	maybeInput *ErrorOr[*TLSHandshakeResultState]) *ErrorOr[*HTTPTransportState] {

	// if the previous stage failed, forward the error
	if maybeInput.err != nil {
		return NewErrorOr[*HTTPTransportState](nil, maybeInput.err)
	}
	input := maybeInput.Unwrap()

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
		UnderlyingCloser:      input.Conn,
		ZeroTime:              input.ZeroTime,
	}
	return NewErrorOr(result, nil)
}
