package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	idGen := &atomicx.Int64{}
	tk := &testKeys{}
	zeroTime := time.Now()

	err := measure(ctx, idGen, zeroTime, tk, wg)
	runtimex.PanicOnError(err, "measure failed unexpectedly")

	data, err := json.Marshal(tk)
	runtimex.PanicOnError(err, "json.Marshal failed unexpectedly")
	fmt.Printf("%s\n", string(data))
}
