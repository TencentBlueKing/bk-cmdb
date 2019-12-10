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
	"strconv"
	"strings"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * plat represent cloud plat here
 */

func (am *AuthManager) CollectAllPlats(ctx context.Context, header http.Header) ([]PlatSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{}),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDPlat, &cond)
	if err != nil {
		blog.V(3).Infof("get all plats, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get all plats, err: %+v", err)
	}
	plats := make([]PlatSimplify, 0)
	for _, cls := range result.Data.Info {
		plat := PlatSimplify{}
		_, err = plat.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get all plat failed, err: %+v", err)
		}
		plats = append(plats, plat)
	}
	return plats, nil
}

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

func (am *AuthManager) MakeResourcesByPlatID(header http.Header, action meta.Action, platIDs ...int64) ([]meta.ResourceAttribute, error) {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		blog.Errorf("MakeResourcesByPlatID failed, collectPlatByIDs failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("collectPlatByIDs failed, err: %+v", err)
	}
	return am.MakeResourcesByPlat(header, action, plats...)
}

// be careful: plat is registered as a common instance in iam
func (am *AuthManager) MakeResourcesByPlat(header http.Header, action meta.Action, plats ...PlatSimplify) ([]meta.ResourceAttribute, error) {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)

	platModels, err := am.collectObjectsByObjectIDs(ctx, header, 0, common.BKInnerObjIDPlat)
	if err != nil {
		blog.Errorf("get plat model failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get plat model failed, err: %+v", err)
	}
	if len(platModels) == 0 {
		blog.Errorf("get plat model failed, not found, rid: %s", rid)
		return nil, fmt.Errorf("get plat model failed, not found")
	}
	platModel := platModels[0]

	resources := make([]meta.ResourceAttribute, 0)
	for _, plat := range plats {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Plat,
				Name:       plat.BKCloudNameField,
				InstanceID: plat.BKCloudIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
			Layers: []meta.Item{
				{
					Type:       meta.Model,
					Name:       platModel.ObjectName,
					InstanceID: platModel.ID,
				},
			},
		}

		resources = append(resources, resource)
	}
	return resources, nil
}

func (am *AuthManager) AuthorizeByPlat(ctx context.Context, header http.Header, action meta.Action, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	rid := util.GetHTTPCCRequestID(header)

	// make auth resources
	resources, err := am.MakeResourcesByPlat(header, action, plats...)
	if err != nil {
		blog.Errorf("AuthorizeByPlat failed, MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}

	return am.authorize(ctx, header, 0, resources...)
}

func (am *AuthManager) AuthorizeByPlatIDs(ctx context.Context, header http.Header, action meta.Action, platIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plat by id failed, err: %+d", err)
	}
	return am.AuthorizeByPlat(ctx, header, action, plats...)
}

func (am *AuthManager) UpdateRegisteredPlat(ctx context.Context, header http.Header, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(plats) == 0 {
		return nil
	}

	rid := util.GetHTTPCCRequestID(header)

	// make auth resources
	resources, err := am.MakeResourcesByPlat(header, meta.EmptyAction, plats...)
	if err != nil {
		blog.Errorf("UpdateRegisteredPlat failed, MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredPlatByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredPlat(ctx, header, plats...)
}

func (am *AuthManager) UpdateRegisteredPlatByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredPlat(ctx, header, plats...)
}

func (am *AuthManager) DeregisterPlatByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister plats failed, get plats by id failed, err: %+v", err)
	}
	return am.DeregisterPlat(ctx, header, plats...)
}

func (am *AuthManager) RegisterPlat(ctx context.Context, header http.Header, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(plats) == 0 {
		return nil
	}

	rid := util.GetHTTPCCRequestID(header)

	// make auth resources
	resources, err := am.MakeResourcesByPlat(header, meta.EmptyAction, plats...)
	if err != nil {
		blog.Errorf("RegisterPlat failed, MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterPlatByID(ctx context.Context, header http.Header, platIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(platIDs) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plats by id failed, err: %+v", err)
	}
	return am.RegisterPlat(ctx, header, plats...)
}

func (am *AuthManager) DeregisterPlat(ctx context.Context, header http.Header, plats ...PlatSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(plats) == 0 {
		return nil
	}

	rid := util.GetHTTPCCRequestID(header)

	// make auth resources
	resources, err := am.MakeResourcesByPlat(header, meta.EmptyAction, plats...)
	if err != nil {
		blog.Errorf("DeregisterPlat failed, MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterPlatByID(ctx context.Context, header http.Header, platIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(platIDs) == 0 {
		return nil
	}

	plats, err := am.collectPlatByIDs(ctx, header, platIDs...)
	if err != nil {
		return fmt.Errorf("get plats by id failed, err: %+v", err)
	}
	return am.DeregisterPlat(ctx, header, plats...)
}

func (am *AuthManager) ListAuthorizedPlatIDs(ctx context.Context, username string) ([]int64, error) {
	authorizedResources, err := am.Authorize.ListAuthorizedResources(ctx, username, 0, meta.Plat, meta.FindMany)
	if err != nil {
		return nil, err
	}

	authorizedPlatIDs := make([]int64, 0)
	for _, iamResource := range authorizedResources {
		if len(iamResource) == 0 {
			continue
		}
		resource := iamResource[len(iamResource)-1]
		if strings.HasPrefix(resource.ResourceID, "plat:") {
			parts := strings.Split(resource.ResourceID, ":")
			if len(parts) < 2 {
				return nil, fmt.Errorf("parse platID from iamResource failed,  iamResourceID: %s, format error", resource.ResourceID)
			}
			platID, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse platID from iamResource failed, iamResourceID: %s, err: %+v", resource.ResourceID, err)
			}
			authorizedPlatIDs = append(authorizedPlatIDs, platID)
		}
	}
	return authorizedPlatIDs, nil
}
