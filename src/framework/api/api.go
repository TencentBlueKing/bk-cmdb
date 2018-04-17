package api

import (
	"configcenter/src/framework/core/publisher"
	"configcenter/src/framework/core/timer"
	"context"
)

// CreateFramework create a new framework
func CreateFramework() (*Framework, error) {

	ctx, cancel := context.WithCancel(context.Background())

	// create Framework
	fr := &Framework{
		workers: MapWorker{},
		context: ctx,
		cancel:  cancel,
	}

	/** initialize the default configuration */

	// set the timer
	fr.Timer = timer.New(ctx)

	// set publisher manager
	fr.Publisher = publisher.New()

	return fr, nil
}

// DestoryFramework destory the framework
func DestoryFramework(fr *Framework) error {
	return fr.stop()
}
