package output

// New create a new Manager instance
func New() Manager {

	mgr := &manager{
		outputers: MapOutputer{},
	}

	return mgr
}
