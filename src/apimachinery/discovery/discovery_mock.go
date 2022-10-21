/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package discovery

// NewMockDiscoveryInterface TODO
func NewMockDiscoveryInterface() DiscoveryInterface {
	return &MockDiscovery{}
}

// MockDiscovery TODO
type MockDiscovery struct{}

// MigrateServer TODO
func (d *MockDiscovery) MigrateServer() Interface {
	return &mockServer{}
}

// ApiServer TODO
func (d *MockDiscovery) ApiServer() Interface {
	return &mockServer{}
}

// EventServer TODO
func (d *MockDiscovery) EventServer() Interface {
	return &mockServer{}
}

// HostServer TODO
func (d *MockDiscovery) HostServer() Interface {
	return &mockServer{}
}

// ProcServer TODO
func (d *MockDiscovery) ProcServer() Interface {
	return &mockServer{}
}

// TopoServer TODO
func (d *MockDiscovery) TopoServer() Interface {
	return &mockServer{}
}

// DataCollect TODO
func (d *MockDiscovery) DataCollect() Interface {
	return &mockServer{}
}

// AuditCtrl TODO
func (d *MockDiscovery) AuditCtrl() Interface {
	return &mockServer{}
}

// HostCtrl TODO
func (d *MockDiscovery) HostCtrl() Interface {
	return &mockServer{}
}

// ObjectCtrl TODO
func (d *MockDiscovery) ObjectCtrl() Interface {
	return &mockServer{}
}

// ProcCtrl TODO
func (d *MockDiscovery) ProcCtrl() Interface {
	return &mockServer{}
}

// GseProcServer TODO
func (d *MockDiscovery) GseProcServer() Interface {
	return &mockServer{}
}

// OperationServer TODO
func (d *MockDiscovery) OperationServer() Interface {
	return &mockServer{}
}

// CoreService TODO
func (d *MockDiscovery) CoreService() Interface {
	return &mockServer{}
}

// TaskServer TODO
func (d *MockDiscovery) TaskServer() Interface {
	return &mockServer{}
}

// CloudServer TODO
func (d *MockDiscovery) CloudServer() Interface {
	return &mockServer{}
}

// AuthServer TODO
func (d *MockDiscovery) AuthServer() Interface {
	return &mockServer{}
}

// CacheService TODO
func (d *MockDiscovery) CacheService() Interface {
	return &mockServer{}
}

// IsMaster TODO
func (d *MockDiscovery) IsMaster() bool {
	return true
}

// Server TODO
func (d *MockDiscovery) Server(name string) Interface {
	return emptyServerInst
}

type mockServer struct{}

// GetServers TODO
func (*mockServer) GetServers() ([]string, error) {
	return []string{"http://127.0.0.1:8080"}, nil
}

// IsMaster TODO
func (*mockServer) IsMaster(string) bool {
	return true
}

// GetServersChan TODO
func (s *mockServer) GetServersChan() chan []string {
	return nil
}
