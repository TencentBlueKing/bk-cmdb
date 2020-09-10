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

package y3_9_202009101936

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"

	"go.mongodb.org/mongo-driver/bson"
)

func migrateHistory(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// query history dynamic group data that storaged in "UserAPI" mode.
	opt := mapstr.MapStr{}
	userConfigs := make([]metadata.UserConfig, 0)

	if err := db.Table(common.BKTableNameUserAPI).Find(opt).All(ctx, userConfigs); err != nil {
		return err
	}

	// range user configs and convert to "DynamicGroup" mode.
	for _, userConfig := range userConfigs {
		dynamicGroup := &metadata.DynamicGroup{
			AppID:      userConfig.AppID,
			Name:       userConfig.Name,
			ID:         userConfig.ID,
			CreateTime: userConfig.CreateTime,
			UpdateTime: userConfig.UpdateTime,
			CreateUser: userConfig.CreateUser,
			ModifyUser: userConfig.ModifyUser,
			Info:       metadata.DynamicGroupInfo{},
		}

		if err := json.Unmarshal([]byte(userConfig.Info), &dynamicGroup.Info); err != nil {
			blog.Errorf("unmarshal new dynamic group info error: %+v", err)
			return err
		}

		if len(dynamicGroup.Name) == 0 || len(dynamicGroup.ID) == 0 || len(dynamicGroup.Info.Condition) == 0 {
			blog.Errorf("invalid new dynamic group object, %+v", dynamicGroup)
			return fmt.Errorf("invalid new dynamic group object, %+v", dynamicGroup)
		}

		for idxInfo, infoCond := range dynamicGroup.Info.Condition {
			if len(infoCond.ObjID) == 0 || len(infoCond.Condition) == 0 {
				blog.Errorf("invalid new dynamic group info condition, %+v", infoCond)
				return fmt.Errorf("invalid new dynamic group info condition, %+v", infoCond)
			}

			for idxCond, cond := range infoCond.Condition {
				if cond.Operator == common.BKDBMULTIPLELike {
					dynamicGroup.Info.Condition[idxInfo].Condition[idxCond].Operator = common.BKDBIN
				}
			}
		}

		data := map[string]interface{}{}
		row, err := bson.Marshal(dynamicGroup)
		if err != nil {
			blog.Errorf("marshal new dynamic group error:%v", err)
			return err
		}
		if err = bson.Unmarshal(row, data); err != nil {
			blog.Errorf("unmarshal new dynamic group error:%v", err)
			return err
		}

		condition := map[string]interface{}{
			common.BKAppIDField: data[common.BKAppIDField],
			common.BKFieldName:  data[common.BKFieldName],
		}

		count, err := db.Table(common.BKTableNameDynamicGroup).Find(condition).Count(ctx)
		if err != nil {
			blog.Errorf("find count error:%v", err)
			return err
		}
		if count > 0 {
			continue
		}

		err = db.Table(common.BKTableNameDynamicGroup).Insert(ctx, data)
		if err != nil {
			blog.Errorf("insert error %v", err)
			return err
		}
	}

	return nil
}
