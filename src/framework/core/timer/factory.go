package timer

import (
	"configcenter/src/framework/core/timer/regular"
	"context"
)

// New create a new timer instance
func New(ctx context.Context) Timer {

	target := &timer{}

	target.regularTimer = regular.New(ctx)

	return target
}
