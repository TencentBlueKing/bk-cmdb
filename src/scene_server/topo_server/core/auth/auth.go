/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package topoauth

import (
	"context"
	"net/http"

	"configcenter/src/auth"
	"configcenter/src/auth/meta"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type TopoAuth struct {
	authorizer auth.Authorize
}

func (ta *TopoAuth) RegisterObject(ctx context.Context, header http.Header, object metadata.Object) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.Model,
			Name:       object.ObjectID,
			InstanceID: object.ID,
		},
		SupplierAccount: util.GetOwnerID(header),
	}

	if err := ta.authorizer.RegisterResource(ctx, resource); err != nil {
		return err
	}

	return nil
}

func (ta *TopoAuth) RegisterObjectsBatch(header http.Header) error {

	return nil
}

func (ta *TopoAuth) DeregisterObject() error {

	return nil
}

func (ta *TopoAuth) RegisterInstance() error {

	return nil
}

func (ta *TopoAuth) DeregisterInstance() error {

	return nil
}
