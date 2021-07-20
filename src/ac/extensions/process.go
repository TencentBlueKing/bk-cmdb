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
 * process
 */

func (am *AuthManager) collectProcessesByIDs(ctx context.Context, header http.Header, ids ...int64) ([]ProcessSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKProcIDField).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDProc, &cond)
	if err != nil {
		blog.Errorf("get processes by id %+v failed, err: %+v, rid: %s", ids, err, rid)
		return nil, fmt.Errorf("get processes by id failed, err: %+v", err)
	}
	processes := make([]ProcessSimplify, 0)
	for _, item := range result.Data.Info {
		process := ProcessSimplify{}
		_, err = process.Parse(item)
		if err != nil {
			blog.Errorf("collectProcessesByIDs by id %+v failed, parse process %+v failed, err: %+v, rid: %s", ids, item, err, rid)
			return nil, fmt.Errorf("parse process from db data failed, err: %+v", err)
		}
		processes = append(processes, process)
	}
	return processes, nil
}

func (am *AuthManager) MakeResourcesByProcesses(header http.Header, action meta.Action, businessID int64, processes ...ProcessSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, process := range processes {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Process,
				Name:       process.ProcessName,
				InstanceID: process.ProcessID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) GenProcessNoPermissionResp(ctx context.Context, header http.Header, businessID int64) (*metadata.BaseResp, error) {
	// process read authorization is skipped
	resp := metadata.NewNoPermissionResp(nil)
	return &resp, nil
}

func (am *AuthManager) extractBusinessIDFromProcesses(processes ...ProcessSimplify) (int64, error) {
	var businessID int64
	for idx, process := range processes {
		bizID := process.BKAppIDField
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from processes")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) AuthorizeByProcesses(ctx context.Context, header http.Header, action meta.Action, processes ...ProcessSimplify) error {
	if !am.Enabled() {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromProcesses(processes...)
	if err != nil {
		return fmt.Errorf("authorize processes failed, extract business id from processes failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByProcesses(header, action, bizID, processes...)

	return am.batchAuthorize(ctx, header, resources...)
}

func (am *AuthManager) AuthorizeByProcessID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
	if !am.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	processes, err := am.collectProcessesByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("authorize processes failed, collect process by id failed, id: %+v, err: %+v", ids, err)
	}

	return am.AuthorizeByProcesses(ctx, header, action, processes...)
}
