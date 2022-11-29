package dslx

//
// TCP measurements
//

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/bassosimone/oonidsl/internal/model"
)

// TCPConnect returns a function that establishes TCP connections.
func TCPConnect() Function[*EndpointState, *TCPConnectResultState] {
	f := &tcpConnectFunction{}
	return f
}

// tcpConnectFunction is a function that establishes TCP connections.
type tcpConnectFunction struct{}

// Apply applies the function to its arguments.
func (f *tcpConnectFunction) Apply(
	ctx context.Context, input *EndpointState) *TCPConnectResultState {

	// create trace
	trace := measurexlite.NewTrace(input.IDGenerator.Add(1), input.ZeroTime)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		input.Logger,
		"[#%d] TCPConnect %s",
		trace.Index,
		input.Address,
	)

	// setup
	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	dialer := trace.NewDialerWithoutResolver(input.Logger)

	// connect
	conn, err := dialer.DialContext(ctx, "tcp", input.Address)

	// stop the operation logger
	ol.Stop(err)

	return &TCPConnectResultState{
		Address:     input.Address,
		Conn:        conn,
		Domain:      input.Domain,
		Err:         err,
		IDGenerator: input.IDGenerator,
		Logger:      input.Logger,
		Network:     input.Network,
		Trace:       trace,
		ZeroTime:    input.ZeroTime,
	}
}

// TCPConnectResultState is the state generated by a TCP connect. If you
// initialize manually, init at least the ones marked as MANDATORY.
type TCPConnectResultState struct {
	// Address is the MANDATORY address we tried to connect to.
	Address string

	// Conn is the possibly-nil TCP connection.
	Conn net.Conn

	// Domain is the OPTIONAL domain from which we resolved the Address.
	Domain string

	// Err is the error that occurred when connecting or nil.
	Err error

	// IDGenerator is the MANDATORY ID generator.
	IDGenerator *atomicx.Int64

	// Logger is the MANDATORY logger to use.
	Logger model.Logger

	// Network is the MANDATORY network we tried to use when connecting.
	Network string

	// Trace is the MANDATORY trace we're using.
	Trace *measurexlite.Trace

	// ZeroTime is the MANDATORY zero time of the measurement.
	ZeroTime time.Time
}

var _ ObservationsProducer = &TCPConnectResultState{}

// Observations implements ObservationsProducer
func (s *TCPConnectResultState) Observations() []*Observations {
	return maybeTraceToObservations(s.Trace)
}

var _ io.Closer = &TCPConnectResultState{}

// Close implements io.Closer
func (s *TCPConnectResultState) Close() error {
	if s.Conn != nil {
		return s.Conn.Close()
	}
	return nil
}
