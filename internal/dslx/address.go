package dslx

//
// Manipulate sets of IP addresses
//

import (
	"net"
	"strconv"

	"github.com/bassosimone/oonidsl/internal/netxlite"
)

// AddressSet transforms DNS lookup results into a set of IP addresses.
func AddressSet(dns ...*ErrorOr[*DNSLookupResultState]) *AddressSetState {
	uniq := make(map[string]bool)
	for _, e := range dns {
		if e.Error() != nil {
			continue
		}
		for _, a := range e.Unwrap().Addresses {
			uniq[a] = true
		}
	}
	return &AddressSetState{uniq}
}

// AddressSetState is the state created by AddressSet. The zero value
// struct is invalid, please use AddressSet to construct.
type AddressSetState struct {
	M map[string]bool
}

// Add adds a (possibly-new) address to the set.
func (as *AddressSetState) Add(addrs ...string) *AddressSetState {
	for _, addr := range addrs {
		as.M[addr] = true
	}
	return as
}

// RemoveBogons removes bogons from the set.
func (as *AddressSetState) RemoveBogons() *AddressSetState {
	zap := []string{}
	for addr := range as.M {
		if netxlite.IsBogon(addr) {
			zap = append(zap, addr)
		}
	}
	for _, addr := range zap {
		delete(as.M, addr)
	}
	return as
}

// ToEndpoints transforms this set of IP addresses to a list of endpoints. We will
// combine each IP address with the network and the port to construct an endpoint and
// we will also apply any additional option to each endpoint.
func (as *AddressSetState) ToEndpoints(
	network EndpointNetwork, port uint16, options ...EndpointOption) (v []*EndpointState) {
	for addr := range as.M {
		v = append(v, Endpoint(
			network,
			EndpointAddress(net.JoinHostPort(addr, strconv.Itoa(int(port)))),
			options...,
		))
	}
	return
}
