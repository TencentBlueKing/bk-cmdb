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

package event

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// cleanDelArchiveData is to clean the table cc_DelArchive data which is a week ago.
// we do this everyday at a fixed time.
// we find the expired data with _id.
func (f *Flow) cleanDelArchiveData(ctx context.Context) {
	blog.Infof("start clean cc_DelArchive data job success.")
	go func() {
		for {
			if time.Now().Hour() != 1 {
				time.Sleep(5 * time.Minute)
				continue
			}
			rid := util.GenerateRID()
			blog.Infof("start do clean cc_DelArchive data job, rid: %s", rid)
			f.doClean(ctx, rid)
			blog.Infof("start do clean cc_DelArchive data job done, rid: %s", rid)
		}
	}()
}

func (f *Flow) doClean(ctx context.Context, rid string) {
	timeout := time.After(time.Hour)
	for {
		select {
		case <-timeout:
			blog.Errorf("do clean cc_DelArchive data job timeout, rid: %s", rid)
			return
		default:
		}
		time.Sleep(5 * time.Minute)

		if !f.isMaster.IsMaster() {
			blog.Infof("try to clean cc_DelArchive data job, but not master, skip, rid: %s", rid)
			continue
		}

		blog.Infof("do clean cc_DelArchive data job, rid: %s", rid)

		// it's time to do the clean job.
		// generate a ObjectID with a time.
		week := time.Now().Unix() - 7*24*60*60
		weekAgo := time.Unix(week, 0)
		oid := primitive.NewObjectIDFromTimestamp(weekAgo)

		// count the data older than this oid
		filter := mapstr.MapStr{
			"_id": mapstr.MapStr{
				common.BKDBLT: oid,
			},
		}

		count, err := f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).Count(ctx)
		if err != nil {
			blog.Errorf("clean cc_DelArchive data, but count expired data in %s failed. rid: %s", common.BKTableNameDelArchive, rid)
			continue
		}

		blog.Infof("do clean cc_DelArchive data job, found %d expired docs, rid: %s", count, rid)

		pageSize := 500
		success := true
		for start := 0; start < int(count); start += pageSize {
			docs := make([]archived, pageSize)
			err = f.ccDB.Table(common.BKTableNameDelArchive).Find(filter).Fields("oid").All(ctx, &docs)
			if err != nil {
				blog.Errorf("clean cc_DelArchive data, but find expired data failed, err: %v, rid: %s", err, rid)
				time.Sleep(10 * time.Second)
				success = false
				continue
			}

			oids := make([]string, len(docs))
			for idx, doc := range docs {
				oids[idx] = doc.Oid
			}

			delFilter := mapstr.MapStr{
				"oid": mapstr.MapStr{
					common.BKDBIN: oids,
				},
			}

			err = f.ccDB.Table(common.BKTableNameDelArchive).Delete(ctx, delFilter)
			if err != nil {
				blog.Errorf("clean cc_DelArchive data, but delete data failed, err: %v, rid: %s", err, rid)
				time.Sleep(10 * time.Second)
				success = false
				continue
			}
			// sleep a while
			time.Sleep(10 * time.Second)
		}

		if success {
			blog.Infof("clean cc_DelArchive data success, delete %d docs, rid: %s", count, rid)
		} else {
			blog.Infof("clean cc_DelArchive data job done, but part of it is failed, rid: %s", rid)
		}

		// finished the for loop.
		return
	}
}

type archived struct {
	Oid string `bson:"oid"`
}
