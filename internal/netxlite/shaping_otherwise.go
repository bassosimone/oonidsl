//go:build !shaping

package netxlite

import (
	"github.com/bassosimone/oonidsl/internal/model"
)

func newMaybeShapingDialer(dialer model.Dialer) model.Dialer {
	return dialer
}
