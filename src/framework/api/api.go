package api

import (
	"configcenter/src/framework/core/publisher"
	"context"
)

// CMDBFrameworkInit create a new framework
func CMDBFrameworkInit() (*Framework, error) {

	ctx, cancel := context.WithCancel(context.Background())

	// create Framework
	fr := &Framework{
		workers: MapWorker{},
		context: ctx,
		cancel:  cancel,
	}

	// set publisher manager
	publisherMgr := publisher.New()
	fr.Publisher = publisherMgr

	return fr, nil
}

// CMDBFrameworkUnInit destory the framework
func CMDBFrameworkUnInit(fr *Framework) error {
	return fr.stop()
}
