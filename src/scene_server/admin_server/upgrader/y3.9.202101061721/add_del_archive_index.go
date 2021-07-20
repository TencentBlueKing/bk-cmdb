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

package y3_9_202101061721

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var oidCollIndex = types.Index{
	Keys:       map[string]int32{"oid": 1, "coll": 1},
	Unique:     true,
	Background: true,
	Name:       "idx_oid_coll",
}

var collIndex = types.Index{
	Keys:       map[string]int32{"coll": 1},
	Unique:     false,
	Background: true,
	Name:       "idx_coll",
}

type archiveData struct {
	MongoID primitive.ObjectID `bson:"_id"`
}

// delPreviousDelArchiveData delete previous del archive data
func delPreviousDelArchiveData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for {
		dataArr := make([]archiveData, 0)
		err := db.Table(common.BKTableNameDelArchive).Find(nil).Fields("_id").Start(0).
			Limit(common.BKMaxPageSize).All(ctx, &dataArr)
		if err != nil {
			blog.Errorf("find previous del archive data failed, err: %v", err)
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
			blog.Errorf("delete previous del archive data failed, err: %v", err)
			return err
		}

		time.Sleep(time.Millisecond * 5)
	}

	return nil
}

// addDelArchiveIndex add unique index for coll and oid
func addDelArchiveIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	existIndexes, err := db.Table(common.BKTableNameDelArchive).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("find indexes for del archive table failed. err: %v", err)
		return err
	}

	for _, index := range existIndexes {
		if index.Name == oidCollIndex.Name || index.Name == collIndex.Name {
			return nil
		}
	}

	err = db.Table(common.BKTableNameDelArchive).CreateIndex(ctx, oidCollIndex)
	if err != nil {
		blog.ErrorJSON("add index %s for del archive table failed, err: %s", oidCollIndex, err)
		return err
	}

	err = db.Table(common.BKTableNameDelArchive).CreateIndex(ctx, collIndex)
	if err != nil {
		blog.ErrorJSON("add index %s for del archive table failed, err: %s", collIndex, err)
		return err
	}

	return nil
}
