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

package coreservice

import (
	"fmt"

	"configcenter/src/apimachinery/coreservice/association"
	"configcenter/src/apimachinery/coreservice/instance"
	"configcenter/src/apimachinery/coreservice/model"
	"configcenter/src/apimachinery/coreservice/synchronize"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

type CoreServiceClientInterface interface {
	Instance() instance.InstanceClientInterface
	Model() model.ModelClientInterface
	Association() association.AssociationClientInterface
	Synchronize() synchronize.SynchronizeClientInterface
}

func NewCoreServiceClient(c *util.Capability, version string) CoreServiceClientInterface {
	base := fmt.Sprintf("/api/%s", version)
	return &coreService{
		restCli: rest.NewRESTClient(c, base),
	}
}

type coreService struct {
	restCli rest.ClientInterface
}

func (c *coreService) Instance() instance.InstanceClientInterface {
	return instance.NewInstanceClientInterface(c.restCli)
}

func (c *coreService) Model() model.ModelClientInterface {
	return model.NewModelClientInterface(c.restCli)
}

func (c *coreService) Association() association.AssociationClientInterface {
	return association.NewAssociationClientInterface(c.restCli)
}

func (c *coreService) Synchronize() synchronize.SynchronizeClientInterface {
	return synchronize.NewSynchronizeClientInterface(c.restCli)
}
