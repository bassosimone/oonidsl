package main

//
// Top-level measurement algorithm
//

import "context"

// measure is the top-level measurement algorithm.
func measure(ctx context.Context, state *measurementState) error {
	// run DCs measurements in background
	state.wg.Add(1)
	go measureDCs(ctx, state)

	// run web measurements in background
	state.wg.Add(1)
	go measureWeb(ctx, state)

	// wait for measurements to terminate
	state.wg.Wait()

	// make sure we fail the measurement if the main
	// context is cancelled (e.g., because the user
	// has hit ^C and has forced an early termination)
	return ctx.Err()
}
