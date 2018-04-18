package manager

import (
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/publisher"
	"configcenter/src/framework/core/timer"
	"context"
)

// Manager contains the basic framework data and the publisher client used to publis events.
type Manager struct {
	InputerMgr input.Manager
	Timer      timer.Timer
	Publisher  publisher.Publisher
}

// stop used to stop the business cycles.
func (cli *Manager) stop() error {

	return cli.InputerMgr.Stop()
}

// Run start the business cycle until the stop method is called.
func (cli *Manager) Run(ctx context.Context, cancel context.CancelFunc) {

	cli.InputerMgr.Run(ctx, cancel)
}
