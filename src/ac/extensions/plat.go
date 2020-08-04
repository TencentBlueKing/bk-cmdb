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
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * plat represent cloud plat here
 */

func (am *AuthManager) collectPlatByIDs(ctx context.Context, header http.Header, platIDs ...int64) ([]PlatSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	platIDs = util.IntArrayUnique(platIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKSubAreaField).In(platIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDPlat, &cond)
	if err != nil {
		blog.V(3).Infof("get plats by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get plats by id failed, err: %+v", err)
	}
	plats := make([]PlatSimplify, 0)
	for _, cls := range result.Data.Info {
		plat := PlatSimplify{}
		_, err = plat.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get plat by id failed, err: %+v", err)
		}
		plats = append(plats, plat)
	}
	return plats, nil
}

// be careful: plat is registered as a common instance in iam
func (am *AuthManager) MakeResourcesByPlat(header http.Header, action meta.Action, plats ...PlatSimplify) ([]meta.ResourceAttribute, error) {

	resources := make([]meta.ResourceAttribute, 0)
	for _, plat := range plats {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.CloudAreaInstance,
				Name:       plat.BKCloudNameField,
				InstanceID: plat.BKCloudIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
		}

		resources = append(resources, resource)
	}
	return resources, nil
}

func (am *AuthManager) AuthorizeByPlat(ctx context.Context, header http.Header, action meta.Action, plats ...PlatSimplify) error {
	if !am.Enabled() {
		return nil
	}

	rid := util.GetHTTPCCRequestID(header)

	// make auth resources
	resources, err := am.MakeResourcesByPlat(header, action, plats...)
	if err != nil {
		blog.Errorf("AuthorizeByPlat failed, MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}

	return am.batchAuthorize(ctx, header, resources...)
}

func (am *AuthManager) AuthorizeByPlatIDs(ctx context.Context, header http.Header, action meta.Action, platIDs ...int64) error {
	if !am.Enabled() {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plat by id failed, err: %+d", err)
	}
	return am.AuthorizeByPlat(ctx, header, action, plats...)
}
