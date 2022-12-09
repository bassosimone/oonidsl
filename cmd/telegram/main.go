package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/model"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

// measurementState contains state variables used when measuring.
type measurementState struct {
	// dGen allows to assign unique IDs to submeasurements.
	idGen *atomicx.Int64

	// logger contains the logger.
	logger model.Logger

	// tk contains the experiment results.
	tk *testKeys

	// zeroTime is the "zero time" of the measurement.
	zeroTime time.Time

	// wg allows waiting for background goroutines.
	wg *sync.WaitGroup
}

func main() {
	state := &measurementState{
		idGen:    &atomicx.Int64{},
		logger:   log.Log,
		tk:       &testKeys{},
		zeroTime: time.Now(),
		wg:       &sync.WaitGroup{},
	}
	ctx := context.Background()

	err := measure(ctx, state)
	runtimex.PanicOnError(err, "measure failed unexpectedly")

	data, err := json.Marshal(state.tk)
	runtimex.PanicOnError(err, "json.Marshal failed unexpectedly")
	fmt.Printf("%s\n", string(data))
}
