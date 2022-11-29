package dslx

//
// Manipulate endpoints
//

import (
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/model"
)

type (
	// EndpointNetwork is the network of the endpoint
	EndpointNetwork string

	// EndpointAddress is the endpoint address.
	EndpointAddress string
)

// EndpointState is the state necessary for measuring a single endpoint. You
// should construct from an AddressSetState. Otherwise, in case of manual initialization,
// make sure you initialize all the fields marked as MANDATORY.
type EndpointState struct {
	// Address is the MANDATORY endpoint address.
	Address string

	// Domain is the OPTIONAL domain used to resolve the endpoints' IP address.
	Domain string

	// IDGenerator is MANDATORY the ID generator to use.
	IDGenerator *atomicx.Int64

	// Logger is the MANDATORY logger to use.
	Logger model.Logger

	// Network is the MANDATORY endpoint network.
	Network string

	// ZeroTime is the MANDATORY zero time of the measurement.
	ZeroTime time.Time
}

// EndpointOption is an option you can use to construct EndpointState.
type EndpointOption func(*EndpointState)

// EndpointOptionDomain allows to set the domain.
func EndpointOptionDomain(value string) EndpointOption {
	return func(es *EndpointState) {
		es.Domain = value
	}
}

// EndpointOptionIDGenerator allows to set the ID generator.
func EndpointOptionIDGenerator(value *atomicx.Int64) EndpointOption {
	return func(es *EndpointState) {
		es.IDGenerator = value
	}
}

// EndpointOptionLogger allows to set the logger.
func EndpointOptionLogger(value model.Logger) EndpointOption {
	return func(es *EndpointState) {
		es.Logger = value
	}
}

// EndpointOptionZerotime allows to set the zero time.
func EndpointOptionZerotime(value time.Time) EndpointOption {
	return func(es *EndpointState) {
		es.ZeroTime = value
	}
}

// Endpoint creates state for measuring a network Endpoint (i.e., a three
// tuple composed of a network protocol, an IP address, and a port).
//
// Arguments:
//
// - network is either "tcp" or "udp";
//
// - address is the Endpoint address represented as an IP address followed by ":"
// followed by a port. IPv6 addresses must be quoted (e.g., "[::1]:80");
//
// - options contains additional options.
func Endpoint(
	network EndpointNetwork, address EndpointAddress, options ...EndpointOption) *EndpointState {
	epnt := &EndpointState{
		Address:     string(address),
		Domain:      "",
		IDGenerator: &atomicx.Int64{},
		Logger:      model.DiscardLogger,
		Network:     string(network),
		ZeroTime:    time.Now(),
	}
	for _, option := range options {
		option(epnt)
	}
	return epnt
}
