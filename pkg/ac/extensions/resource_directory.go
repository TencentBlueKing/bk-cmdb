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
	meta2 "configcenter/pkg/ac/meta"
	"context"
	"fmt"
	"net/http"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/metadata"
	"configcenter/pkg/util"
)

/*
 * module instance
 */

func (am *AuthManager) collectResourceDirectoryByDirectoryIDs(ctx context.Context, header http.Header,
	directoryIDs ...int64) ([]ModuleSimplify, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	directoryIDs = util.IntArrayUnique(directoryIDs)

	cond := metadata.QueryCondition{
		Condition: map[string]interface{}{common.BKModuleIDField: map[string]interface{}{common.BKDBIN: directoryIDs}},
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDModule, &cond)
	if err != nil {
		blog.V(3).Infof("get directory by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get directory by id failed, err: %+v", err)
	}

	directoryArr := make([]ModuleSimplify, 0)
	for _, cls := range result.Info {
		directory := ModuleSimplify{}
		_, err = directory.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("parse directory failed, err: %+v", err)
		}
		directoryArr = append(directoryArr, directory)
	}
	return directoryArr, nil
}

// MakeResourcesByResourceDirectory TODO
func (am *AuthManager) MakeResourcesByResourceDirectory(header http.Header, action meta2.Action, directoryArr ...ModuleSimplify) []meta2.ResourceAttribute {
	resources := make([]meta2.ResourceAttribute, 0)
	for _, directory := range directoryArr {
		resource := meta2.ResourceAttribute{
			Basic: meta2.Basic{
				Action:     action,
				Type:       meta2.ResourcePoolDirectory,
				Name:       directory.BKModuleNameField,
				InstanceID: directory.BKModuleIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
		}

		resources = append(resources, resource)
	}
	return resources
}

// AuthorizeByResourceDirectoryID TODO
func (am *AuthManager) AuthorizeByResourceDirectoryID(ctx context.Context, header http.Header, action meta2.Action, ids ...int64) error {
	if !am.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	directoryArr, err := am.collectResourceDirectoryByDirectoryIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered directoryArr failed, get directoryArr by id failed, err: %+v", err)
	}
	return am.AuthorizeByResourceDirectory(ctx, header, action, directoryArr...)
}

// AuthorizeByResourceDirectory TODO
func (am *AuthManager) AuthorizeByResourceDirectory(ctx context.Context, header http.Header, action meta2.Action, directoryArr ...ModuleSimplify) error {
	if !am.Enabled() {
		return nil
	}

	if len(directoryArr) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByResourceDirectory(header, action, directoryArr...)

	return am.batchAuthorize(ctx, header, resources...)
}
