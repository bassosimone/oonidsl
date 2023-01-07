package dslx

//
// QUIC measurements
//

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/bassosimone/oonidsl/internal/netxlite"
	"github.com/lucas-clemente/quic-go"
)

// QUICHandshakeOption is an option you can pass to QUICHandshake.
type QUICHandshakeOption func(*quicHandshakeFunc)

// QUICHandshakeOptionInsecureSkipVerify controls whether QUIC verification is enabled.
func QUICHandshakeOptionInsecureSkipVerify(value bool) QUICHandshakeOption {
	return func(thf *quicHandshakeFunc) {
		thf.InsecureSkipVerify = value
	}
}

// QUICHandshakeOptionRootCAs allows to configure custom root CAs.
func QUICHandshakeOptionRootCAs(value *x509.CertPool) QUICHandshakeOption {
	return func(thf *quicHandshakeFunc) {
		thf.RootCAs = value
	}
}

// QUICHandshakeOptionServerName allows to configure the SNI to use.
func QUICHandshakeOptionServerName(value string) QUICHandshakeOption {
	return func(thf *quicHandshakeFunc) {
		thf.ServerName = value
	}
}

// QUICHandshake returns a function performing QUIC handshakes.
func QUICHandshake(pool *QUICConnPool, options ...QUICHandshakeOption) Func[
	*Endpoint, *Maybe[*TLSConnection]] {
	f := &quicHandshakeFunc{
		InsecureSkipVerify: false,
		Pool:               pool,
		RootCAs:            netxlite.NewDefaultCertPool(),
		ServerName:         "",
	}
	for _, option := range options {
		option(f)
	}
	return f
}

// quicHandshakeFunc performs QUIC handshakes.
type quicHandshakeFunc struct {
	// InsecureSkipVerify allows to skip TLS verification.
	InsecureSkipVerify bool

	// Pool is the QUICConnPool that owns us.
	Pool *QUICConnPool

	// RootCAs contains the Root CAs to use.
	RootCAs *x509.CertPool

	// ServerName is the ServerName to handshake for.
	ServerName string
}

// Apply implements Func.
func (f *quicHandshakeFunc) Apply(
	ctx context.Context, input *Endpoint) *Maybe[*TLSConnection] {
	// create trace
	trace := measurexlite.NewTrace(input.IDGenerator.Add(1), input.ZeroTime)

	// use defaults or user-configured overrides
	serverName := f.serverName(input)

	// start the operation logger
	ol := measurexlite.NewOperationLogger(
		input.Logger,
		"[#%d] QUICHandshake with %s SNI=%s",
		trace.Index,
		input.Address,
		serverName,
	)

	// setup
	quicListener := netxlite.NewQUICListener()
	quicDialer := trace.NewQUICDialerWithoutResolver(quicListener, input.Logger)
	config := &tls.Config{
		NextProtos:         []string{"h3"},
		InsecureSkipVerify: f.InsecureSkipVerify,
		RootCAs:            f.RootCAs,
		ServerName:         serverName,
	}
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// handshake
	quicConn, err := quicDialer.DialContext(ctx, input.Address, config, &quic.Config{})

	// possibly register established conn for late close
	f.Pool.maybeRegister(quicConn)

	// stop the operation logger
	ol.Stop(err)

	// start preparing the message to emit on the stdout
	state := &TLSConnection{
		Address:     input.Address,
		QUICConn:    quicConn,
		Domain:      input.Domain,
		IDGenerator: input.IDGenerator,
		Logger:      input.Logger,
		Network:     input.Network,
		TLSConfig:   config,
		TLSState:    quicConn.ConnectionState().TLS.ConnectionState, // TODO unsafe?
		Trace:       trace,
		ZeroTime:    input.ZeroTime,
	}

	return &Maybe[*TLSConnection]{
		Error:        err,
		Observations: maybeTraceToObservations(trace),
		Skipped:      false,
		State:        state,
	}
}

func (f *quicHandshakeFunc) serverName(input *Endpoint) string {
	if f.ServerName != "" {
		return f.ServerName
	}
	return input.Domain
}
