package discovery

func NewMockDiscoveryInterface() DiscoveryInterface {
	return &MockDiscovery{}
}

type MockDiscovery struct{}

func (d *MockDiscovery) MigrateServer() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) EventServer() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) HostServer() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) ProcServer() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) TopoServer() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) DataCollect() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) AuditCtrl() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) HostCtrl() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) ObjectCtrl() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) ProcCtrl() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) GseProcServ() Interface {
	return &mockServer{}
}

type mockServer struct{}

func (*mockServer) GetServers() ([]string, error) {
	return []string{"http://127.0.0.1:8080"}, nil
}
