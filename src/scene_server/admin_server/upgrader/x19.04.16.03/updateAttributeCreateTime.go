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

package x19_04_16_03

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updateAttributeCreateTime(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	var start uint64

	type Attribute struct {
		CreaetTime *time.Time `bson:"creaet_time"`
		CreateTime *time.Time `bson:"create_time"`
		ID         uint64     `bson:"id"`
	}

	now := time.Now()
	for {
		attrs := []Attribute{}
		err := db.Table(common.BKTableNameObjAttDes).Find(nil).Start(start).Limit(50).All(ctx, &attrs)
		if err != nil {
			return err
		}
		if len(attrs) <= 0 {
			break
		}
		start += 50

		for _, attr := range attrs {
			if attr.CreateTime == nil {
				createTime := attr.CreaetTime
				if createTime == nil {
					createTime = &now
				}
				if attr.CreateTime != nil {
					createTime = attr.CreateTime
				}

				cond := condition.CreateCondition()
				cond.Field(common.BKFieldID).Eq(attr.ID)
				err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond.ToMapStr(), mapstr.MapStr{common.CreateTimeField: createTime})
				if err != nil {
					return err
				}
			}
		}
	}

	return db.Table(common.BKTableNameObjAttDes).DropColumn(ctx, "creaet_time")
}
