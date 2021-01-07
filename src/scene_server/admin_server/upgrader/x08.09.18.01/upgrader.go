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

package x08_09_18_01

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func fixedHostPlatAssocateRelation(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	type instAsstStruct struct {
		ID           int64     `bson:"id"`
		InstID       int64     `bson:"bk_inst_id"`
		ObjectID     string    `bson:"bk_obj_id"`
		AsstInstID   int64     `bson:"bk_asst_inst_id"`
		AsstObjectID string    `bson:"bk_asst_obj_id"`
		OwnerID      string    `bson:"bk_supplier_account"`
		CreateTime   time.Time `bson:"create_time"`
		LastTime     time.Time `bson:"last_time"`
	}

	instAsstArr := make([]instAsstStruct, 0)
	instAsstConditionMap := mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKAsstInstIDField: common.BKInnerObjIDPlat,
	}

	err = db.Table(common.BKTableNameInstAsst).Find(instAsstConditionMap).Fields(common.BKHostIDField).All(ctx, &instAsstArr)
	if nil != err && !db.IsNotFoundError(err) {
		return err
	}
	var exitsAsstHostIDArr []int64
	for _, instAsst := range instAsstArr {
		exitsAsstHostIDArr = append(exitsAsstHostIDArr, instAsst.InstID)
	}

	mongo, ok := db.(*local.Mongo)
	if !ok {
		return fmt.Errorf("db is not *local.Mongo type")
	}
	dbc := mongo.GetDBClient()

	hostInfoMap := make([]map[string]interface{}, 0)
	findOpts := &options.FindOptions{}
	findOpts.SetProjection(map[string]int{common.BKHostIDField: 1, common.BKCloudIDField: 1, common.BKOwnerIDField: 1})
	hostCondition := make(mapstr.MapStr)
	if 0 < len(exitsAsstHostIDArr) {
		hostCondition[common.BKHostIDField] = mapstr.MapStr{common.BKDBNIN: exitsAsstHostIDArr}
	}

	cursor, err := dbc.Database(mongo.GetDBName()).Collection(common.BKTableNameBaseHost).Find(ctx, hostCondition, findOpts)
	if err != nil && !db.IsNotFoundError(err) {
		return err
	}

	if err := cursor.All(ctx, &hostInfoMap); err != nil && !db.IsNotFoundError(err) {
		return err
	}

	nowTime := time.Now().UTC()
	for _, host := range hostInfoMap {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}
		platID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}
		ownerID := util.GetStrByInterface(host[common.BkSupplierAccount])

		instAsstConditionMap := mapstr.MapStr{
			common.BKObjIDField:     common.BKInnerObjIDHost,
			common.BKInstIDField:    hostID,
			common.BKAsstObjIDField: common.BKInnerObjIDPlat,
		}
		cnt, err := db.Table(common.BKTableNameInstAsst).Find(instAsstConditionMap).Count(ctx)
		if nil != err {
			return err
		}
		if 0 == cnt {
			id, err := db.NextSequence(ctx, common.BKTableNameInstAsst)
			if nil != err {
				return err
			}
			addAsstInst := instAsstStruct{
				ID:           int64(id),
				InstID:       hostID,
				ObjectID:     common.BKInnerObjIDHost,
				AsstInstID:   platID,
				AsstObjectID: common.BKInnerObjIDPlat,
				OwnerID:      ownerID,
				CreateTime:   nowTime,
				LastTime:     nowTime,
			}
			err = db.Table(common.BKTableNameInstAsst).Insert(ctx, addAsstInst)
			if nil != err {
				return err
			}
		}
	}
	return nil
}
