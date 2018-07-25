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

package objcontroller

import (
	"fmt"

	"configcenter/src/apimachinery/objcontroller/identifier"
	"configcenter/src/apimachinery/objcontroller/inst"
	"configcenter/src/apimachinery/objcontroller/meta"
	"configcenter/src/apimachinery/objcontroller/openapi"
	"configcenter/src/apimachinery/objcontroller/privilege"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

type ObjControllerClientInterface interface {
	Instance() inst.InstanceInterface
	Meta() meta.MetaInterface
	Identifier() identifier.IdentifierInterface
	OpenAPI() openapi.OpenApiInterface
	Privilege() privilege.PrivilegeInterface
}

func NewObjectControllerInterface(c *util.Capability, version string) ObjControllerClientInterface {
	base := fmt.Sprintf("/object/%s", version)
	return &objectctrl{
		client: rest.NewRESTClient(c, base),
	}
}

type objectctrl struct {
	client rest.ClientInterface
}

func (o *objectctrl) Instance() inst.InstanceInterface {
	return inst.NewInstanceInterface(o.client)
}

func (o *objectctrl) Meta() meta.MetaInterface {
	return meta.NewmetaInterface(o.client)
}

func (o *objectctrl) OpenAPI() openapi.OpenApiInterface {
	return openapi.NewOpenApiInterface(o.client)
}

func (o *objectctrl) Privilege() privilege.PrivilegeInterface {
	return privilege.NewPrivilegeInterface(o.client)
}

func (o *objectctrl) Identifier() identifier.IdentifierInterface {
	return identifier.NewIdentifierInterface(o.client)
}
