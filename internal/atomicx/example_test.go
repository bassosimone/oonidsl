package atomicx_test

import (
	"fmt"

	"github.com/bassosimone/oonidsl/internal/atomicx"
)

func Example_typicalUsage() {
	v := &atomicx.Int64{}
	v.Add(1)
	fmt.Printf("%d\n", v.Load())
	// Output: 1
}
