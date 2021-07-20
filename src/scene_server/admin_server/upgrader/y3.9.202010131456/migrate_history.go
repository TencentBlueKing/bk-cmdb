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

package y3_9_202010131456

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func migrateHistory(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// query history dynamic group data that storaged in "UserAPI" mode.
	opt := mapstr.MapStr{}
	userConfigs := make([]UserConfig, 0)

	if err := db.Table(common.BKTableNameUserAPI).Find(opt).All(ctx, &userConfigs); err != nil {
		return err
	}

	// range user configs and convert to "DynamicGroup" mode.
	for _, userConfig := range userConfigs {
		dynamicGroup := make(map[string]interface{})

		dynamicGroupInfo := metadata.DynamicGroupInfo{}
		if err := json.Unmarshal([]byte(userConfig.Info), &dynamicGroupInfo); err != nil {
			blog.Errorf("unmarshal new dynamic group info error: %+v", err)
			return err
		}

		dynamicGroup[common.BKAppIDField] = userConfig.AppID
		dynamicGroup[common.BKFieldID] = userConfig.ID
		dynamicGroup[common.BKFieldName] = userConfig.Name
		dynamicGroup[common.BKObjIDField] = common.BKInnerObjIDHost
		dynamicGroup["create_user"] = userConfig.CreateUser
		dynamicGroup["create_time"] = userConfig.CreateTime
		dynamicGroup["modify_user"] = userConfig.ModifyUser
		dynamicGroup["last_time"] = userConfig.UpdateTime

		dynamicGroupInfoMeta := make(map[string]interface{})
		dynamicGroupInfoCondition := []metadata.DynamicGroupInfoCondition{}

		for idxInfo, infoCond := range dynamicGroupInfo.Condition {
			if len(infoCond.ObjID) == 0 || len(infoCond.Condition) == 0 {
				blog.Warnf("invalid new dynamic group info condition, %+v", infoCond)
				continue
			}

			if infoCond.ObjID != common.BKInnerObjIDHost &&
				infoCond.ObjID != common.BKInnerObjIDModule &&
				infoCond.ObjID != common.BKInnerObjIDSet {
				blog.Warnf("ignore old dynamic group condition, %+v", infoCond)
				continue
			}

			for idxCond, cond := range infoCond.Condition {
				if cond.Operator == common.BKDBMULTIPLELike {
					dynamicGroupInfo.Condition[idxInfo].Condition[idxCond].Operator = common.BKDBIN
				}
			}
			dynamicGroupInfoCondition = append(dynamicGroupInfoCondition, infoCond)
		}

		dynamicGroupInfoMeta["condition"] = dynamicGroupInfoCondition
		dynamicGroup["info"] = dynamicGroupInfoMeta

		condition := map[string]interface{}{common.BKFieldID: userConfig.ID}
		count, err := db.Table(common.BKTableNameDynamicGroup).Find(condition).Count(ctx)
		if err != nil {
			blog.Errorf("find count error:%v", err)
			return err
		}
		if count > 0 {
			continue
		}

		err = db.Table(common.BKTableNameDynamicGroup).Insert(ctx, dynamicGroup)
		if err != nil {
			blog.Errorf("insert error %v", err)
			return err
		}
	}

	return nil
}
