package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/fx"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

func main() {
	ctx := context.Background()

	zeroTime := time.Now()
	idGen := &atomicx.Int64{}

	coll := &collector{}

	dnsLookupResults := fx.Parallel(ctx, fx.Parallelism(2),
		dslx.DNSLookupInput(
			dslx.DomainName("www.google.com"),
			dslx.DNSLookupOptionZeroTime(zeroTime),
			dslx.DNSLookupOptionLogger(log.Log),
			dslx.DNSLookupOptionIDGenerator(idGen),
		),
		dslx.DNSLookupGetaddrinfo(coll),
		dslx.DNSLookupUDP("8.8.8.8:53", coll),
	)

	coll.dump()

	endpoints := dslx.AddressSet(dnsLookupResults...).
		Add("142.250.184.100").
		RemoveBogons().
		ToEndpoints(
			dslx.EndpointNetwork("tcp"),
			dslx.EndpointPort(443),
			dslx.EndpointOptionDomain("www.google.com"),
			dslx.EndpointOptionIDGenerator(idGen),
			dslx.EndpointOptionLogger(log.Log),
			dslx.EndpointOptionZeroTime(zeroTime),
		)

	connpool := &dslx.ConnPool{}
	defer connpool.Close()

	tlsHandshakeErrors := &dslx.ErrorLogger{}

	_ = fx.Map(ctx, fx.Parallelism(2),
		fx.ComposeResult4(
			dslx.TCPConnect(connpool, coll),
			dslx.RecordErrors(
				tlsHandshakeErrors,
				dslx.TLSHandshake(connpool, coll),
			),
			dslx.HTTPTransportTLS(),
			dslx.HTTPRequest(coll),
		),
		endpoints...,
	)

	log.Infof("%+v", tlsHandshakeErrors.Errors())

	coll.dump()
}

type collector struct {
	odump []*dslx.Observations
	mu    *sync.Mutex
}

// MergeObservations implements ObservationCollector.MergeObservations.
func (c *collector) MergeObservations(obs ...*dslx.Observations) {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.odump = append(c.odump, obs...)
}

func (c *collector) dump() {
	data, err := json.Marshal(c.odump)
	runtimex.PanicOnError(err, "json.Marshal failed")
	fmt.Printf("%s\n", string(data))
}
