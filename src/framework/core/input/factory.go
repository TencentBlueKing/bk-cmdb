package input

// New create a new Manager instance
func New() Manager {
	mgr := &manager{inputers: MapInputer{}}
	return mgr
}
