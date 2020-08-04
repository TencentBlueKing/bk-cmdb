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
 * set instance
 */

func (am *AuthManager) collectSetBySetIDs(ctx context.Context, header http.Header, setIDs ...int64) ([]SetSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKSetIDField).In(setIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDSet, &cond)
	if err != nil {
		blog.V(3).Infof("get sets by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get sets by id failed, err: %+v", err)
	}
	sets := make([]SetSimplify, 0)
	for _, cls := range result.Data.Info {
		set := SetSimplify{}
		_, err = set.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get sets by object failed, err: %+v", err)
		}
		sets = append(sets, set)
	}
	return sets, nil
}

func (am *AuthManager) extractBusinessIDFromSets(sets ...SetSimplify) (int64, error) {
	var businessID int64
	for idx, set := range sets {
		bizID := set.BKAppIDField
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from sets")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesBySet(header http.Header, action meta.Action, businessID int64, sets ...SetSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, set := range sets {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelSet,
				Name:       set.BKSetNameField,
				InstanceID: set.BKSetIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) AuthorizeBySetID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
	if !am.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	if !am.RegisterSetEnabled {
		return nil
	}

	sets, err := am.collectSetBySetIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("collect set by id failed, err: %+v", err)
	}
	return am.AuthorizeBySet(ctx, header, action, sets...)
}

func (am *AuthManager) AuthorizeBySet(ctx context.Context, header http.Header, action meta.Action, sets ...SetSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, sets: %+v, rid: %s", sets, rid)
		return nil
	}
	if !am.RegisterSetEnabled {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromSets(sets...)
	if err != nil {
		return fmt.Errorf("authorize sets failed, extract business id from sets failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesBySet(header, action, bizID, sets...)

	return am.batchAuthorize(ctx, header, resources...)
}
