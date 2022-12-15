package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/bassosimone/oonidsl/internal/atomicx"
	"github.com/bassosimone/oonidsl/internal/runtimex"
)

func main() {
	tk := &testKeys{}
	ctx := context.Background()

	err := measure(ctx, log.Log, &atomicx.Int64{}, time.Now(), tk)
	runtimex.PanicOnError(err, "measure failed unexpectedly")

	data, err := json.Marshal(tk)
	runtimex.PanicOnError(err, "json.Marshal failed unexpectedly")
	fmt.Printf("%s\n", string(data))
}
