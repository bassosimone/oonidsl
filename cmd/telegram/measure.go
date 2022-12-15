package main

//
// Top-level measurement algorithm
//

import (
	"context"
	"sync"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/model"
)

// measure is the top-level measurement algorithm.
func measure(
	ctx context.Context,
	logger model.Logger,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
) error {
	wg := &sync.WaitGroup{}

	// run DCs measurements in background
	wg.Add(1)
	go measureDCs(ctx, logger, idGen, zeroTime, tk, wg)

	// run web measurements in background
	wg.Add(1)
	go measureWeb(ctx, logger, idGen, zeroTime, tk, wg)

	// wait for measurements to terminate
	wg.Wait()

	// make sure we fail the measurement if the main
	// context is cancelled (e.g., because the user
	// has hit ^C and has forced an early termination)
	return ctx.Err()
}
