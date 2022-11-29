package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/dslx"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

func dump(v any) {
	data, err := json.Marshal(v)
	runtimex.PanicOnError(err, "json.Marshal failed")
	fmt.Printf("%s\n", string(data))
}

func main() {
	ctx := context.Background()

	zeroTime := time.Now()
	idGen := &atomicx.Int64{}

	dnsLookupResults := dslx.Parallel(ctx, dslx.Parallelism(2),
		dslx.DNSLookupInput(
			dslx.DomainName("www.google.com"),
			dslx.DNSLookupOptionZeroTime(zeroTime),
			dslx.DNSLookupOptionLogger(log.Log),
			dslx.DNSLookupOptionIDGenerator(idGen),
		),
		dslx.DNSLookupGetaddrinfo(),
		dslx.DNSLookupUDP("8.8.8.8:53"),
	)

	dnsObservations := dslx.ExtractObservations(dnsLookupResults...)
	dump(dnsObservations)

	endpoints := dslx.AddressSet(dnsLookupResults...).
		Add("142.250.184.100").
		RemoveBogons().
		ToEndpoints(
			dslx.EndpointNetwork("tcp"),
			443,
			dslx.EndpointOptionDomain("www.google.com"),
			dslx.EndpointOptionIDGenerator(idGen),
			dslx.EndpointOptionLogger(log.Log),
			dslx.EndpointOptionZerotime(zeroTime),
		)

	endpointsResults := dslx.Map(ctx, dslx.Parallelism(2),
		dslx.Compose5(
			dslx.TCPConnect(),
			dslx.TLSHandshake(),
			dslx.HTTPTransportTLS(),
			dslx.HTTPRequest(),
			dslx.Close[*dslx.HTTPRequestResultState](),
		),
		endpoints...,
	)

	endpointsObservations := dslx.ExtractObservations(endpointsResults...)
	dump(endpointsObservations)
}
