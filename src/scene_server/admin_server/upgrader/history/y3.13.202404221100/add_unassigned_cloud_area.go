/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_13_202404221100

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func addUnassignedCloudArea(ctx context.Context, db dal.RDB, conf *history.Config) error {
	cond := mapstr.MapStr{
		common.BKDBAND: []mapstr.MapStr{
			{common.BKCloudIDField: mapstr.MapStr{common.BKDBGTE: common.ReservedCloudAreaStartID}},
			{common.BKCloudIDField: mapstr.MapStr{common.BKDBLTE: common.ReservedCloudAreaEndID}},
			{common.BKCloudIDField: mapstr.MapStr{common.BKDBNIN: common.ReservedCloudAreaIDs}},
		},
	}

	count, err := db.Table(common.BKTableNameBasePlat).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("count cloud area failed, err: %v, cond: %+v", err, cond)
		return err
	}
	if count > 0 {
		msg := fmt.Sprintf("reserved cloud area data exists in db, range: [%d:%d]", common.ReservedCloudAreaStartID,
			common.ReservedCloudAreaEndID)
		blog.Errorf(msg)
		return errors.New(msg)
	}

	cond = mapstr.MapStr{common.BKCloudIDField: common.UnassignedCloudAreaID}
	result := make([]CloudArea, 0)
	if err = db.Table(common.BKTableNameBasePlat).Find(cond).All(ctx, &result); err != nil {
		blog.Errorf("find cloud area failed, cond: %+v, err: %v", cond, err)
		return err
	}

	if len(result) > 1 {
		msg := fmt.Sprintf("multiple cloud area have been found, cond: %+v, count: %d", cond, len(result))
		blog.Errorf(msg)
		return errors.New(msg)
	}

	if len(result) == 1 {
		data := result[0]
		if data.CloudName == common.UnassignedCloudAreaName && data.OwnerID == "0" &&
			data.Creator == conf.User && data.Default == int64(common.BuiltIn) {
			return nil
		}

		msg := fmt.Sprintf("cloud area[%d] already exists, data: %+v", common.UnassignedCloudAreaID, data)
		blog.Errorf(msg)
		return errors.New(msg)
	}

	cloudArea := &CloudArea{
		CloudID:    common.UnassignedCloudAreaID,
		CloudName:  common.UnassignedCloudAreaName,
		OwnerID:    "0",
		Creator:    conf.User,
		LastEditor: conf.User,
		CreateTime: time.Now(),
		LastTime:   time.Now(),
		Default:    int64(common.BuiltIn),
	}
	if err = db.Table(common.BKTableNameBasePlat).Insert(ctx, cloudArea); err != nil {
		blog.Errorf("create unassigned cloud area failed, data: %+v, err: %v", cloudArea, err)
		return err
	}

	return nil
}

// CloudArea 管控区域
type CloudArea struct {
	CloudID     int64     `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName   string    `json:"bk_cloud_name" bson:"bk_cloud_name"`
	Status      string    `json:"bk_status" bson:"bk_status"`
	CloudVendor string    `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	OwnerID     string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	VpcID       string    `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName     string    `json:"bk_vpc_name" bson:"bk_vpc_name"`
	Region      string    `json:"bk_region" bson:"bk_region"`
	AccountID   int64     `json:"bk_account_id" bson:"bk_account_id"`
	Creator     string    `json:"bk_creator" bson:"bk_creator"`
	LastEditor  string    `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	LastTime    time.Time `json:"last_time" bson:"last_time"`
	Default     int64     `json:"default" bson:"default"`
}
