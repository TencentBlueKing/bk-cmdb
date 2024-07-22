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

package y3_13_202407191507

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func upgradeDelArchive(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := delExpiredData(ctx, db); err != nil {
		blog.Errorf("delete expired del archive data failed, err: %v", err)
		return err
	}

	if err := addTimeField(ctx, db); err != nil {
		blog.Errorf("add del archive time field failed, err: %v", err)
		return err
	}

	if err := addKubeDelArchiveTable(ctx, db); err != nil {
		blog.Errorf("add kube del archive table failed, err: %v", err)
		return err
	}

	if err := moveKubeData(ctx, db); err != nil {
		blog.Errorf("move kube data to kube del archive table failed, err: %v", err)
		return err
	}

	return nil
}

type mongoIDData struct {
	MongoID       primitive.ObjectID `bson:"_id"`
	mapstr.MapStr `bson:",inline"`
}

func delExpiredData(ctx context.Context, db dal.RDB) error {
	// generate an ObjectID with the time of a week ago
	weekAgo := time.Unix(time.Now().Unix()-7*24*60*60, 0)
	oid := primitive.NewObjectIDFromTimestamp(weekAgo)

	// delete data earlier than this oid
	filter := mapstr.MapStr{
		"_id": mapstr.MapStr{
			common.BKDBLT: oid,
		},
	}

	for {
		dataArr := make([]mongoIDData, 0)
		err := db.Table(common.BKTableNameDelArchive).Find(filter).Fields("_id").Limit(common.BKMaxPageSize).
			All(ctx, &dataArr)
		if err != nil {
			blog.Errorf("get expired del archive data failed, err: %v", err)
			return err
		}

		if len(dataArr) == 0 {
			return nil
		}

		delMongoIDs := make([]primitive.ObjectID, len(dataArr))
		for index, data := range dataArr {
			delMongoIDs[index] = data.MongoID
		}

		delCond := map[string]interface{}{
			"_id": map[string]interface{}{common.BKDBIN: delMongoIDs},
		}
		if err := db.Table(common.BKTableNameDelArchive).Delete(ctx, delCond); err != nil {
			blog.Errorf("delete kube del archive data failed, err: %v, cond: %+v", err, delCond)
			return err
		}

		time.Sleep(time.Millisecond * 5)
	}

	return nil
}

func addTimeField(ctx context.Context, db dal.RDB) error {
	// add time field for del archive data
	filter := mapstr.MapStr{
		"time": mapstr.MapStr{
			common.BKDBExists: false,
		},
	}

	for {
		dataArr := make([]mongoIDData, 0)
		err := db.Table(common.BKTableNameDelArchive).Find(filter).Fields("_id").Limit(common.BKMaxPageSize).
			All(ctx, &dataArr)
		if err != nil {
			blog.Errorf("get expired del archive data failed, err: %v", err)
			return err
		}

		if len(dataArr) == 0 {
			break
		}

		for _, data := range dataArr {
			updateCond := mapstr.MapStr{
				"_id": data.MongoID,
			}
			updateData := mapstr.MapStr{
				"time": data.MongoID.Timestamp(),
			}

			if err := db.Table(common.BKTableNameDelArchive).Update(ctx, updateCond, updateData); err != nil {
				blog.Errorf("update del archive failed, err: %v, cond: %+v, data: %+v", err, updateCond, updateData)
				return err
			}
		}

		time.Sleep(time.Millisecond * 5)
	}

	// add time index
	timeIndex := types.Index{
		Name:               common.CCLogicIndexNamePrefix + "time",
		Keys:               bson.D{{"time", -1}},
		Background:         true,
		ExpireAfterSeconds: 7 * 24 * 60 * 60,
	}

	existIndexArr, err := db.Table(common.BKTableNameDelArchive).Indexes(ctx)
	if err != nil {
		blog.Errorf("get %s exist indexes failed, err: %v, rid: %s", common.BKTableNameDelArchive, err)
		return err
	}

	for _, index := range existIndexArr {
		if index.Name == timeIndex.Name {
			return nil
		}
	}

	err = db.Table(common.BKTableNameDelArchive).CreateIndex(ctx, timeIndex)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.Errorf("create %s index(%+v) failed, err: %v, rid: %s", common.BKTableNameDelArchive, timeIndex, err)
		return err
	}

	return nil
}

func addKubeDelArchiveTable(ctx context.Context, db dal.RDB) error {
	// add kube del archive table
	table := common.BKTableNameKubeDelArchive
	exists, err := db.HasTable(ctx, table)
	if err != nil {
		blog.Errorf("check if %s table exists failed, err: %v", table, err)
		return err
	}

	if !exists {
		if err = db.CreateTable(ctx, table); err != nil {
			blog.Errorf("create %s table failed, err: %v", table, err)
			return err
		}
	}

	existIndexes, err := db.Table(table).Indexes(ctx)
	if err != nil {
		blog.Errorf("get %s indexes failed. err: %v", table, err)
		return err
	}

	existIndexMap := make(map[string]struct{})
	for _, index := range existIndexes {
		existIndexMap[index.Name] = struct{}{}
	}

	// add kube del archive index
	indexes := []types.Index{
		{
			Name:               common.CCLogicIndexNamePrefix + "time",
			Keys:               bson.D{{"time", -1}},
			Background:         true,
			ExpireAfterSeconds: 2 * 24 * 60 * 60,
		}, {
			Name: common.CCLogicIndexNamePrefix + "coll_oid",
			Keys: bson.D{
				{"coll", 1},
				{"oid", 1},
			},
			Unique:     true,
			Background: true,
		}, {
			Name: common.CCLogicIndexNamePrefix + "coll",
			Keys: bson.D{
				{"coll", 1},
			},
			Background: true,
		}, {
			Name: common.CCLogicIndexNamePrefix + "oid",
			Keys: bson.D{
				{"oid", 1},
			},
			Background: true,
		},
	}

	needCreateIndexes := make([]types.Index, 0)
	for _, index := range indexes {
		_, exists := existIndexMap[index.Name]
		if !exists {
			needCreateIndexes = append(needCreateIndexes, index)
		}
	}

	if len(needCreateIndexes) == 0 {
		return nil
	}

	err = db.Table(table).BatchCreateIndexes(ctx, needCreateIndexes)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.Errorf("create %s indexes %+v failed, err: %s", table, needCreateIndexes, err)
		return err
	}

	return nil
}

func moveKubeData(ctx context.Context, db dal.RDB) error {
	filter := mapstr.MapStr{
		"coll": mapstr.MapStr{
			common.BKDBIN: []string{kubetypes.BKTableNameBaseCluster, kubetypes.BKTableNameBaseNode,
				kubetypes.BKTableNameBaseNamespace, kubetypes.BKTableNameBaseWorkload,
				kubetypes.BKTableNameBaseDeployment, kubetypes.BKTableNameBaseStatefulSet,
				kubetypes.BKTableNameBaseDaemonSet, kubetypes.BKTableNameGameDeployment,
				kubetypes.BKTableNameGameStatefulSet, kubetypes.BKTableNameBaseCronJob, kubetypes.BKTableNameBaseJob,
				kubetypes.BKTableNameBasePodWorkload, kubetypes.BKTableNameBaseCustom, kubetypes.BKTableNameBasePod,
				kubetypes.BKTableNameBaseContainer, kubetypes.BKTableNameNsSharedClusterRel},
		},
	}
	opts := types.NewFindOpts().SetWithObjectID(true)

	twoDayAgo := time.Unix(time.Now().Unix()-2*24*60*60, 0)

	for {
		dataArr := make([]mongoIDData, 0)
		err := db.Table(common.BKTableNameDelArchive).Find(filter, opts).Sort("_id").Limit(common.BKMaxPageSize).
			All(ctx, &dataArr)
		if err != nil {
			blog.Errorf("get expired del archive data failed, err: %v", err)
			return err
		}

		if len(dataArr) == 0 {
			break
		}

		mongoIDs := make([]primitive.ObjectID, 0)
		needMovedData := make([]mongoIDData, 0)
		for _, data := range dataArr {
			mongoIDs = append(mongoIDs, data.MongoID)

			if data.MongoID.Timestamp().Before(twoDayAgo) {
				continue
			}

			needMovedData = append(needMovedData, data)
		}

		if len(needMovedData) > 0 {
			if err := db.Table(common.BKTableNameKubeDelArchive).Insert(ctx, needMovedData); err != nil {
				blog.Errorf("insert kube del archive data failed, err: %v, data: %+v", err, needMovedData)
				return err
			}
		}

		delCond := map[string]interface{}{
			"_id": map[string]interface{}{common.BKDBIN: mongoIDs},
		}
		if err := db.Table(common.BKTableNameDelArchive).Delete(ctx, delCond); err != nil {
			blog.Errorf("delete kube del archive data failed, err: %v, cond: %+v", err, delCond)
			return err
		}

		time.Sleep(time.Millisecond * 5)
	}

	return nil
}
