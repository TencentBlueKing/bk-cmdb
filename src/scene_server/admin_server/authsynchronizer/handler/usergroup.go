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

package handler

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
)

var (
	userGroups = []string{
		common.BKMaintainersField,
		common.BKProductPMField,
		common.BKTesterField,
		common.BKDeveloperField,
		common.BKOperatorField,
	}
)

func (ih *IAMHandler) HandleUserGroupSync(task *meta.WorkRequest) error {
	biz := task.Data.(extensions.BusinessSimplify)
	h := make(http.Header)
	// list biz members from auth
	members, err := ih.authManager.Authorize.GetUserGroupMembers(context.TODO(), h, biz.BKAppIDField, userGroups)
	if err != nil {
		return fmt.Errorf("sync biz: %d, name: %s user group from iam failed, err: %v", biz.BKAppIDField, biz.BKAppNameField, err)
	}

	changedFields := mapstr.MapStr{}
	for _, m := range members {
		switch m.Name {
		case common.BKMaintainersField:
			if !isUserDifferent(m.Users, strings.Split(biz.Maintainer, ",")) {
				changedFields[common.BKMaintainersField] = strings.Join(m.Users, ",")
				blog.Warnf("sync user group with biz: %s,  %s has changed from: %s to %+v.", biz.BKAppNameField,
					common.BKMaintainersField, biz.Maintainer, m.Users)
			}
		case common.BKProductPMField:
			if !isUserDifferent(m.Users, strings.Split(biz.Producer, ",")) {
				changedFields[common.BKProductPMField] = strings.Join(m.Users, ",")
				blog.Warnf("sync user group with biz: %s,  %s has changed from: %s to %+v.", biz.BKAppNameField,
					common.BKProductPMField, biz.Producer, m.Users)
			}
		case common.BKTesterField:
			if !isUserDifferent(m.Users, strings.Split(biz.Tester, ",")) {
				changedFields[common.BKTesterField] = strings.Join(m.Users, ",")
				blog.Warnf("sync user group with biz: %s,  %s has changed from: %s to %+v.", biz.BKAppNameField,
					common.BKTesterField, biz.Tester, m.Users)
			}
		case common.BKDeveloperField:
			if !isUserDifferent(m.Users, strings.Split(biz.Developer, ",")) {
				changedFields[common.BKDeveloperField] = strings.Join(m.Users, ",")
				blog.Warnf("sync user group with biz: %s,  %s has changed from: %s to %+v.", biz.BKAppNameField,
					common.BKDeveloperField, biz.Developer, m.Users)
			}
		case common.BKOperatorField:
			if !isUserDifferent(m.Users, strings.Split(biz.Operator, ",")) {
				changedFields[common.BKOperatorField] = strings.Join(m.Users, ",")
				blog.Warnf("sync user group with biz: %s,  %s has changed from: %s to %+v.", biz.BKAppNameField,
					common.BKOperatorField, biz.Operator, m.Users)
			}
		default:
			return fmt.Errorf("sync user group from auth center, but got unsupported user group: %s", m.Name)
		}
	}

	if len(changedFields) == 0 {
		// nothing is changed, return now.
		return nil
	}

	// user group has changed, need to sync to cmdb now.
	op := metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKAppIDField: biz.BKAppIDField,
		},
		Data: changedFields,
	}
	h.Set(common.BKHTTPOwner, "0")
	h.Set(common.BKHTTPHeaderUser, "cc_system")
	result, err := ih.clientSet.CoreService().Instance().UpdateInstance(context.TODO(), h, "biz", &op)
	if err != nil {
		return fmt.Errorf("sync user group, usr has changed, but update: %+v to biz: %d, name: %s failed, err: %v",
			changedFields, biz.BKAppIDField, biz.BKAppNameField, err)
	}

	if !result.Result {
		return fmt.Errorf("sync user group, usr has changed, but update: %+v to biz: %d, name: %s failed, err: %v",
			changedFields, biz.BKAppIDField, biz.BKAppNameField, result.ErrMsg)
	}
	blog.Warnf("sync user group with biz: %s success", biz.BKAppNameField)
	return nil
}

func isUserDifferent(from []string, to []string) bool {
	fromMap, toMap := make(map[string]bool), make(map[string]bool)
	for _, f := range from {
		fromMap[f] = true
	}

	for _, t := range to {
		toMap[t] = true
	}

	return reflect.DeepEqual(fromMap, toMap)
}
