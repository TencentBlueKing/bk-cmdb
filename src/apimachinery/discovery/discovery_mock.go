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

func NewMockDiscoveryInterface() DiscoveryInterface {
	return &MockDiscovery{}
}

type MockDiscovery struct{}

func (d *MockDiscovery) MigrateServer() Interface {
	return &mockServer{}
}

func (d *MockDiscovery) ApiServer() Interface {
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
