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

// Package toposerver TODO
package toposerver

import (
	"fmt"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/toposerver/association"
	"configcenter/src/apimachinery/toposerver/inst"
	"configcenter/src/apimachinery/toposerver/kube"
	"configcenter/src/apimachinery/toposerver/object"
	"configcenter/src/apimachinery/toposerver/resourcedir"
	"configcenter/src/apimachinery/toposerver/settemplate"
	"configcenter/src/apimachinery/util"
)

// TopoServerClientInterface TODO
type TopoServerClientInterface interface {
	Instance() inst.InstanceInterface
	Object() object.ObjectInterface
	Association() association.AssociationInterface
	SetTemplate() settemplate.SetTemplateInterface
	ResourceDirectory() resourcedir.ResourceDirectoryInterface
	Kube() kube.KubeOperationInterface
}

// NewTopoServerClient TODO
func NewTopoServerClient(c *util.Capability, version string) TopoServerClientInterface {
	base := fmt.Sprintf("/topo/%s", version)
	return &topoServer{
		restCli: rest.NewRESTClient(c, base),
	}
}

type topoServer struct {
	restCli rest.ClientInterface
}

// Instance TODO
func (t *topoServer) Instance() inst.InstanceInterface {
	return inst.NewInstanceClient(t.restCli)
}

// Kube container data related interface initialization.
func (t *topoServer) Kube() kube.KubeOperationInterface {
	return kube.NewKubeOperationInterface(t.restCli)
}

// Object TODO
func (t *topoServer) Object() object.ObjectInterface {
	return object.NewObjectInterface(t.restCli)
}

// Association TODO
func (t *topoServer) Association() association.AssociationInterface {
	return association.NewAssociationInterface(t.restCli)
}

// SetTemplate TODO
func (t *topoServer) SetTemplate() settemplate.SetTemplateInterface {
	return settemplate.NewSetTemplateInterface(t.restCli)
}

// ResourceDirectory TODO
func (t *topoServer) ResourceDirectory() resourcedir.ResourceDirectoryInterface {
	return resourcedir.NewResourceDirectoryInterface(t.restCli)
}
