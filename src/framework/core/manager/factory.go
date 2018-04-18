package manager

// New return a new  Manager instance
func New() *Manager {
	return &Manager{}
}

// Delete delete the framework instance
func Delete(mgr *Manager) error {

	if nil != mgr {
		return mgr.stop()
	}

	return nil
}
