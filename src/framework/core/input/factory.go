package input

// New create a new Manager instance
func New() Manager {
	mgr := &manager{
		inputers:    MapInputer{},
		inputerChan: make(chan *wrapInputer, 1024),
	}
	return mgr
}
