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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * biz set related auth interface
 */

func (am *AuthManager) collectBizSetByIDs(ctx context.Context, header http.Header, rid string, bizSetIDs ...int64) (
	[]BizSetSimplify, error) {

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	bizSetIDs = util.IntArrayUnique(bizSetIDs)

	cond := metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKBizSetIDField: mapstr.MapStr{
				common.BKDBIN: bizSetIDs,
			},
		},
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDBizSet, &cond)
	if err != nil {
		blog.Errorf("get biz set by id failed, err: %v, rid: %s", err, rid)
		return nil, fmt.Errorf("get biz set by id failed, err: %v", err)
	}
	blog.V(5).Infof("get biz set by id result: %+v", result)

	instances := make([]BizSetSimplify, 0)
	for _, cls := range result.Info {
		instance := BizSetSimplify{}
		if _, err = instance.Parse(cls); err != nil {
			blog.Errorf("parse biz set from db data failed, err: %v, rid: %s", err, rid)
			return nil, fmt.Errorf("parse biz set from db data failed, err: %v", err)
		}

		instances = append(instances, instance)
	}
	return instances, nil
}

// makeResourcesByBizSet TODO
func (am *AuthManager) makeResourcesByBizSet(header http.Header, action meta.Action, rid string,
	bizSets ...BizSetSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, bizSet := range bizSets {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.BizSet,
				Name:       bizSet.BKBizSetNameField,
				InstanceID: bizSet.BKBizSetIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
		}

		resources = append(resources, resource)
	}
	return resources
}

// authorizeByBizSet authorize by biz set
func (am *AuthManager) authorizeByBizSet(ctx context.Context, header http.Header, action meta.Action, rid string,
	bizSets ...BizSetSimplify) error {

	if !am.Enabled() {
		return nil
	}

	if len(bizSets) == 0 {
		return nil
	}

	// make auth resources
	resources := am.makeResourcesByBizSet(header, action, rid, bizSets...)

	return am.batchAuthorize(ctx, header, resources...)
}

// AuthorizeByBizSetID TODO
func (am *AuthManager) AuthorizeByBizSetID(ctx context.Context, header http.Header, action meta.Action,
	bizSetIDs ...int64) error {
	if !am.Enabled() {
		return nil
	}

	rid := util.ExtractRequestIDFromContext(ctx)
	bizSets, err := am.collectBizSetByIDs(ctx, header, rid, bizSetIDs...)
	if err != nil {
		blog.Errorf("get biz set data by id failed, err: %v, rid: %s", err, rid)
		return fmt.Errorf("authorize biz set failed, get biz set by id failed, err: %v", err)
	}

	return am.authorizeByBizSet(ctx, header, action, rid, bizSets...)
}
