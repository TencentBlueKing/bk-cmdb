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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// IsPlatAllExist is plat all exist
func (lgc *Logics) IsPlatAllExist(kit *rest.Kit, cloudIDs []int64) (bool, errors.CCError) {
	cloudIDs = util.IntArrayUnique(cloudIDs)
	cond := mapstr.MapStr{
		common.BKCloudIDField: map[string]interface{}{
			common.BKDBIN: cloudIDs,
		},
	}
	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Fields:    []string{common.BKCloudIDField},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDPlat,
		query)
	if err != nil {
		blog.Errorf("find plat failed, cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if result.Count == len(cloudIDs) {
		return true, nil
	}

	return false, nil
}
