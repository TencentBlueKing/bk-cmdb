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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_synchronizer/pkg/synchronizer/meta"
	"configcenter/src/scene_server/auth_synchronizer/pkg/utils"
)

// HandleBusinessSync do sync all business to iam
func (ih *IAMHandler) HandleBusinessSync(task *meta.WorkRequest) error {
	header := utils.NewListBusinessAPIHeader()
	condition := metadata.QueryCondition{}
	result, err := ih.CoreAPI.CoreService().Instance().ReadInstance(context.TODO(), *header, common.BKInnerObjIDApp, &condition)
	if err != nil {
		blog.Errorf("list business failed, err: %v, job: %+v", err, task)
		return err
	}
	businessIDArr := make([]int64, 0)
	for _, business := range result.Data.Info {
		businessID := int64(business[common.BKAppIDField].(float64))
		businessIDArr = append(businessIDArr, businessID)
	}
	blog.Info("list business businessIDArr: %+v", businessIDArr)

	return nil
}
