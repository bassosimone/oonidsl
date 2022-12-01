package dslx

//
// HTTP measurements core
//

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/measurexlite"
	"github.com/bassosimone/oonidsl/internal/model"
	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// HTTPTransportState is an HTTP transport bound to a TCP or TLS connection
// that would use such a connection only and for any input URL. You generally
// use HTTPTransportTCP or HTTPTransportTLS to create a new instance; if you
// want to initialize manually, make sure you init the MANDATORY fields.
type HTTPTransportState struct {
	// Address is the MANDATORY address we're connected to.
	Address string

	// Domain is the OPTIONAL domain from which the address was resolved.
	Domain string

	// IDGenerator is the MANDATORY ID generator.
	IDGenerator *atomicx.Int64

	// Logger is the MANDATORY logger to use.
	Logger model.Logger

	// Network is the MANDATORY network used by the underlying conn.
	Network string

	// Scheme is the MANDATORY URL scheme to use.
	Scheme string

	// TLSNegotiatedProtocol is the OPTIONAL negotiated protocol.
	TLSNegotiatedProtocol string

	// Trace is the MANDATORY trace we're using.
	Trace *measurexlite.Trace

	// Transport is the MANDATORY HTTP transport we're using.
	Transport model.HTTPTransport

	// UnderlyingCloser is the MANDATORY closer to close the underlying conn.
	UnderlyingCloser net.Conn

	// ZeroTime is the MANDATORY zero time of the measurement.
	ZeroTime time.Time
}

var _ io.Closer = &HTTPTransportState{}

// Close implements io.Closer
func (s *HTTPTransportState) Close() error {
	if s.UnderlyingCloser != nil {
		return s.UnderlyingCloser.Close()
	}
	return nil
}

// HTTPRequestOption is an option you can pass to HTTPRequest.
type HTTPRequestOption func(*httpRequestFunction)

// HTTPRequestOptionAccept sets the Accept header.
func HTTPRequestOptionAccept(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.Accept = value
	}
}

// HTTPRequestOptionAcceptLanguage sets the Accept header.
func HTTPRequestOptionAcceptLanguage(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.AcceptLanguage = value
	}
}

// HTTPRequestOptionHost sets the Host header.
func HTTPRequestOptionHost(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.Host = value
	}
}

// HTTPRequestOptionHost sets the request method.
func HTTPRequestOptionMethod(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.Method = value
	}
}

// HTTPRequestOptionReferer sets the Referer header.
func HTTPRequestOptionReferer(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.Referer = value
	}
}

// HTTPRequestOptionURLPath sets the URL path.
func HTTPRequestOptionURLPath(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.URLPath = value
	}
}

// HTTPRequestOptionUserAgent sets the UserAgent header.
func HTTPRequestOptionUserAgent(value string) HTTPRequestOption {
	return func(hrf *httpRequestFunction) {
		hrf.UserAgent = value
	}
}

// HTTPRequest issues an HTTP request using a transport and returns a response.
func HTTPRequest(options ...HTTPRequestOption) Function[*HTTPTransportState, *HTTPRequestResultState] {
	f := &httpRequestFunction{}
	for _, option := range options {
		option(f)
	}
	return f
}

// httpRequestFunction is the Function returned by HTTPRequest.
type httpRequestFunction struct {
	// Accept is the OPTIONAL accept header.
	Accept string

	// AcceptLanguage is the OPTIONAL accept-language header.
	AcceptLanguage string

	// Host is the OPTIONAL host header.
	Host string

	// Method is the OPTIONAL method.
	Method string

	// Referer is the OPTIONAL referer header.
	Referer string

	// URLPath is the OPTIONAL URL path.
	URLPath string

	// UserAgent is the OPTIONAL user-agent header.
	UserAgent string
}

// Apply implements Function.
func (f *httpRequestFunction) Apply(
	ctx context.Context, input *HTTPTransportState) *ErrorOr[*HTTPRequestResultState] {
	// create HTTP request
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var (
		body         []byte
		observations []*Observations
		resp         *http.Response
	)

	req, err := f.newHTTPRequest(ctx, input)
	if err == nil {

		// start the operation logger
		ol := measurexlite.NewOperationLogger(
			input.Logger,
			"[#%d] HTTPRequest %s with %s/%s",
			input.Trace.Index,
			req.URL.String(),
			input.Address,
			input.Network,
		)

		// perform HTTP transaction and collect the related observations
		resp, body, observations, err = f.do(ctx, input, req)

		// stop the operation logger
		ol.Stop(err)
	}

	result := &HTTPRequestResultState{
		Address:                  input.Address,
		Domain:                   input.Domain,
		HTTPObservationsOnce:     sync.Once{},
		HTTPObservations:         observations, // possibly nil
		HTTPRequest:              req,          // possibly nil
		HTTPResponse:             resp,         // possibly nil
		HTTPResponseBodySnapshot: body,         // possibly nil
		IDGenerator:              input.IDGenerator,
		Logger:                   input.Logger,
		Network:                  input.Network,
		Trace:                    input.Trace,
		UnderlyingCloser:         input.UnderlyingCloser,
		ZeroTime:                 input.ZeroTime,
	}

	return NewErrorOr(result, err)
}

func (f *httpRequestFunction) newHTTPRequest(
	ctx context.Context, input *HTTPTransportState) (*http.Request, error) {
	URL := &url.URL{
		Scheme:      input.Scheme,
		Opaque:      "",
		User:        nil,
		Host:        f.urlHost(input),
		Path:        f.urlPath(),
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	method := "GET"
	if f.Method != "" {
		method = f.Method
	}

	req, err := http.NewRequestWithContext(ctx, method, URL.String(), nil)
	if err != nil {
		return nil, err
	}

	if v := f.Host; v != "" { // note: if req.Host is empty, Go uses URL.Hostname
		req.Header.Set("Host", v) // ignored by Go but we want it into the measurement
		req.Host = v
	}

	if v := f.Accept; v != "" {
		req.Header.Set("Accept", v)
	}

	if v := f.AcceptLanguage; v != "" {
		req.Header.Set("Accept-Language", v)
	}

	if v := f.Referer; v != "" {
		req.Header.Set("Referer", v)
	}

	if v := f.UserAgent; v != "" { // not setting means using Go's default
		req.Header.Set("User-Agent", v)
	}

	return req, nil
}

func (f *httpRequestFunction) urlHost(input *HTTPTransportState) string {
	if input.Domain != "" {
		return input.Domain
	}
	addr, port, err := net.SplitHostPort(input.Address)
	if err != nil {
		input.Logger.Warnf("httpRequestFunction: cannot SplitHostPort for input.Address")
		return input.Address
	}
	switch {
	case port == "80" && input.Scheme == "http":
		return addr
	case port == "443" && input.Scheme == "https":
		return addr
	default:
		return input.Address // with port only if port is nonstandard
	}
}

func (f *httpRequestFunction) urlPath() string {
	if f.URLPath != "" {
		return f.URLPath
	}
	return "/"
}

func (f *httpRequestFunction) do(
	ctx context.Context,
	input *HTTPTransportState,
	req *http.Request,
) (*http.Response, []byte, []*Observations, error) {
	const maxbody = 1 << 19 // TODO(bassosimone): allow to configure this value?
	started := input.Trace.TimeSince(input.Trace.ZeroTime)
	observations := []*Observations{{}} // one entry!

	observations[0].NetworkEvents = append(observations[0].NetworkEvents,
		measurexlite.NewAnnotationArchivalNetworkEvent(
			input.Trace.Index,
			started,
			"http_transaction_start",
		))

	resp, err := input.Transport.RoundTrip(req)
	var body []byte
	if err == nil {
		defer resp.Body.Close()
		reader := io.LimitReader(resp.Body, maxbody)
		body, err = netxlite.ReadAllContext(ctx, reader) // TODO: enable streaming and measure speed
	}
	finished := input.Trace.TimeSince(input.Trace.ZeroTime)

	observations[0].NetworkEvents = append(observations[0].NetworkEvents,
		measurexlite.NewAnnotationArchivalNetworkEvent(
			input.Trace.Index,
			finished,
			"http_transaction_done",
		))

	observations[0].Requests = append(observations[0].Requests,
		measurexlite.NewArchivalHTTPRequestResult(
			input.Trace.Index,
			started,
			input.Network,
			input.Address,
			input.TLSNegotiatedProtocol,
			input.Transport.Network(),
			req,
			resp,
			maxbody,
			body,
			err,
			finished,
		))

	return resp, body, observations, err
}

// HTTPRequestResultState is the state generated by HTTP requests. Generally
// obtained by HTTPRequest().Apply. To init manually, init at least MANDATORY fields.
type HTTPRequestResultState struct {
	// Address is the MANDATORY address we're connected to.
	Address string

	// Domain is the OPTIONAL domain from which we determined Address.
	Domain string

	// HTTPObservationsOnce ensures we drain HTTPObservations just once.
	HTTPObservationsOnce sync.Once

	// HTTPObservations contains zero or more HTTP observations. These are
	// returned when you call the Observations method.
	HTTPObservations []*Observations

	// HTTPRequest is the possibly-nil HTTP request.
	HTTPRequest *http.Request

	// HTTPResponse is the HTTP response or nil if Err != nil.
	HTTPResponse *http.Response

	// HTTPResponseBodySnapshot is the response body or nil if Err != nil.
	HTTPResponseBodySnapshot []byte

	// IDGenerator is the MANDATORY ID generator.
	IDGenerator *atomicx.Int64

	// Logger is the MANDATORY logger to use.
	Logger model.Logger

	// Network is the MANDATORY network we're connected to.
	Network string

	// Trace is the MANDATORY trace we're using. The trace is drained
	// when you call the Observations method.
	Trace *measurexlite.Trace

	// UnderlyingCloser is the MANDATORY closer to close the underlying conn.
	UnderlyingCloser net.Conn

	// ZeroTime is the MANDATORY zero time of the measurement.
	ZeroTime time.Time
}

var _ ObservationsProducer = &HTTPRequestResultState{}

// Observations implements ObservationsProducer
func (s *HTTPRequestResultState) Observations() (out []*Observations) {
	if s.Trace != nil {
		out = append(out, maybeTraceToObservations(s.Trace)...)
	}
	// Note: we cannot modify the array because that may be a data race so we
	// use a sync.Once to ensure we have "once" semantics.
	s.HTTPObservationsOnce.Do(func() {
		out = append(out, s.HTTPObservations...)
	})
	return
}

var _ io.Closer = &HTTPRequestResultState{}

// Close implements io.Closer
func (s *HTTPRequestResultState) Close() error {
	if s.UnderlyingCloser != nil {
		return s.UnderlyingCloser.Close()
	}
	return nil
}
