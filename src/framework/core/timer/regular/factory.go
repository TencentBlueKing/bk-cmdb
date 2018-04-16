package regular

import (
	"context"
)

// New create a regular timer
func New(ctx context.Context) Regular {

	// create a new regular instance
	target := &regular{
		tasks: MapRegular{},
	}

	//  start the timer main loop
	go target.run(ctx)

	return target
}
