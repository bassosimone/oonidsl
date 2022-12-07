package main

//
// Top-level measurement algorithm
//

import (
	"context"
	"sync"
	"time"

	"github.com/bassosimone/oonidsl/internal/atomicx"
)

// measure is the top-level measurement algorithm.
//
// Arguments:
//
// - ctx is the context;
//
// - idGen allows to assign unique IDs to submeasurements;
//
// - zeroTime is the "zero time" of the measurement;
//
// - tk contains the experiment results;
//
// - wg allows us to synchronize with our parent.
func measure(
	ctx context.Context,
	idGen *atomicx.Int64,
	zeroTime time.Time,
	tk *testKeys,
	wg *sync.WaitGroup,
) error {
	// run DCs measurements in background
	wg.Add(1)
	go measureDCs(ctx, idGen, zeroTime, tk, wg)

	// run web measurements in background
	wg.Add(1)
	go measureWeb(ctx, idGen, zeroTime, tk, wg)

	// wait for measurements to terminate
	wg.Wait()

	// make sure we fail the measurement if the main
	// context is cancelled (e.g., because the user
	// has hit ^C and has forced an early termination)
	return ctx.Err()
}
