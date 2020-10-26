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

package extensions

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// collectObjectsByObjectIDs collect business object that belongs to business or global object, which both id must in objectIDs
func (am *AuthManager) collectObjectsByObjectIDs(ctx context.Context, header http.Header, businessID int64, objectIDs ...string) ([]metadata.Object, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	objectIDs = util.StrArrayUnique(objectIDs)

	// get model by objID
	cond := condition.CreateCondition().Field(common.BKObjIDField).In(objectIDs)
	fCond := cond.ToMapStr()
	util.AddModelBizIDConditon(fCond, businessID)
	fCond.Remove(metadata.BKMetadata)
	queryCond := &metadata.QueryCondition{Condition: fCond}

	resp, err := am.clientSet.CoreService().Model().ReadModel(ctx, header, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model by id: %+v failed, err: %+v", objectIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get model by id: %+v failed, not found", objectIDs)
	}
	if len(resp.Data.Info) != len(objectIDs) {
		return nil, fmt.Errorf("get model by id: %+v failed, get multiple model", objectIDs)
	}

	objects := make([]metadata.Object, 0)
	for _, item := range resp.Data.Info {
		objects = append(objects, item.Spec)
	}

	return objects, nil
}

// MakeResourcesByObjects make object resource with businessID and objects
func (am *AuthManager) MakeResourcesByObjects(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) ([]meta.ResourceAttribute, error) {
	// prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, object := range objects {
		// instance
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       object.ObjectName,
				InstanceID: object.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      0,
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

func (am *AuthManager) AuthorizeResourceCreate(ctx context.Context, header http.Header, businessID int64, resourceType meta.ResourceType) error {
	if !am.Enabled() {
		return nil
	}

	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:   resourceType,
			Action: meta.Create,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      businessID,
	}

	return am.batchAuthorize(ctx, header, resource)
}
