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

package model

import (
	"configcenter/framework/common/rest"
)

// Interface TODO
type Interface interface {
	Model() ModelInterface
	Attribute() AttributeInterface
	Module() ModuleInterface
	Set() SetInterface
	Instance() InstanceInterface
}

// NewModelClient TODO
func NewModelClient(client rest.ClientInterface) Interface {
	return &modelAPI{
		client: client,
	}
}

type modelAPI struct {
	client rest.ClientInterface
}

// Attribute TODO
func (m *modelAPI) Attribute() AttributeInterface {
	return &attribute{
		client: m.client,
	}
}

// Module TODO
func (m *modelAPI) Module() ModuleInterface {
	return &module{
		client: m.client,
	}
}

// Set TODO
func (m *modelAPI) Set() SetInterface {
	return &setClient{
		client: m.client,
	}
}

// Instance TODO
func (m *modelAPI) Instance() InstanceInterface {
	return &instClient{
		client: m.client,
	}
}

// Model TODO
func (m *modelAPI) Model() ModelInterface {
	return &modelClient{
		client: m.client,
	}
}
