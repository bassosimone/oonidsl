package dslx

//
// TLS measurements
//

import (
	"context"
	"crypto/tls"
	"io"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/bassosimone/oonidsl/internal/model"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// TLSHandshakeOption is an option you can pass to TLSHandshake.
type TLSHandshakeOption func(*tlsHandshakeFunction)

// TLSHandshakeOptionInsecureSkipVerify controls whether TLS verification is enabled.
func TLSHandshakeOptionInsecureSkipVerify(value bool) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunction) {
		thf.InsecureSkipVerify = value
	}
}

// TLSHandshakeOptionNextProto allows to configure the ALPN protocols.
func TLSHandshakeOptionNextProto(value []string) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunction) {
		thf.NextProto = value
	}
}

// TLSHandshakeOptionServerName allows to configure the SNI to use.
func TLSHandshakeOptionServerName(value string) TLSHandshakeOption {
	return func(thf *tlsHandshakeFunction) {
		thf.ServerName = value
	}
}

// TLSHandshake returns a function performing TSL handshakes.
func TLSHandshake(options ...TLSHandshakeOption) Function[
	*TCPConnectResultState, ErrorOr[*TLSHandshakeResultState]] {
	f := &tlsHandshakeFunction{
		InsecureSkipVerify: false,
		NextProto:          []string{},
		ServerName:         "",
	}
	for _, option := range options {
		option(f)
	}
	return f
}

// tlsHandshakeFunction performs TLS handshakes.
type tlsHandshakeFunction struct {
	// InsecureSkipVerify allows to skip TLS verification.
	InsecureSkipVerify bool

	// NextProto contains the ALPNs to negotiate.
	NextProto []string

	// ServerName is the ServerName to handshake for.
	ServerName string
}

// Apply implements Function.
func (f *tlsHandshakeFunction) Apply(
	ctx context.Context, input *TCPConnectResultState) ErrorOr[*TLSHandshakeResultState] {
	// keep using the same trace
	trace := input.Trace

	// use defaults or user-configured overrides
	serverName := f.serverName(input)
	nextProto := f.nextProto()

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		input.Logger,
		"[#%d] TLSHandshake with %s SNI=%s ALPN=%v",
		trace.Index,
		input.Address,
		serverName,
		nextProto,
	)

	// setup
	handshaker := trace.NewTLSHandshakerStdlib(input.Logger)
	config := &tls.Config{
		NextProtos:         nextProto,
		InsecureSkipVerify: f.InsecureSkipVerify,
		RootCAs:            netxlite.NewDefaultCertPool(),
		ServerName:         serverName,
	}
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	conn, state, err := handshaker.Handshake(ctx, input.Conn, config)

	// stop the operation logger
	ol.Stop(err)

	// start preparing the message to emit on the stdout
	result := &TLSHandshakeResultState{
		Address:     input.Address,
		Conn:        nil, // set later
		Domain:      input.Domain,
		IDGenerator: input.IDGenerator,
		Logger:      input.Logger,
		Network:     input.Network,
		TLSState:    state,
		Trace:       trace,
		ZeroTime:    input.ZeroTime,
	}

	// deal with the connections
	if err != nil {
		measurexlite.MaybeClose(input.Conn) // we own it
	} else {
		result.Conn = conn.(netxlite.TLSConn) // guaranteed to work
	}

	return NewErrorOr(result, err)
}

func (f *tlsHandshakeFunction) serverName(input *TCPConnectResultState) string {
	if f.ServerName != "" {
		return f.ServerName
	}
	return input.Domain
}

func (f *tlsHandshakeFunction) nextProto() []string {
	if len(f.NextProto) > 0 {
		return f.NextProto
	}
	return []string{"h2", "http/1.1"}
}

// TLSHandshakeResultState is the state generated by a TLS handshake. If you
// initialize manually, init at least the ones marked as MANDATORY.
type TLSHandshakeResultState struct {
	// Address is the MANDATORY address we tried to connect to.
	Address string

	// Conn is the established TLS conn.
	Conn netxlite.TLSConn

	// Domain is the OPTIONAL domain we resolved.
	Domain string

	// IDGenerator is the MANDATORY ID generator to use.
	IDGenerator *atomicx.Int64

	// Logger is the MANDATORY logger to use.
	Logger model.Logger

	// Network is the MANDATORY network we tried to use when connecting.
	Network string

	// TLSState is the possibly-empty TLS connection state.
	TLSState tls.ConnectionState

	// Trace is the MANDATORY trace we're using.
	Trace *measurexlite.Trace

	// ZeroTime is the MANDATORY zero time of the measurement.
	ZeroTime time.Time
}

var _ ObservationsProducer = &TLSHandshakeResultState{}

// Observations implements ObservationsProducer
func (s *TLSHandshakeResultState) Observations() []*Observations {
	return maybeTraceToObservations(s.Trace)
}

var _ io.Closer = &TLSHandshakeResultState{}

// Close implements io.Closer
func (s *TLSHandshakeResultState) Close() error {
	if s.Conn != nil {
		return s.Conn.Close()
	}
	return nil
}
