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

package x18_12_12_05

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"github.com/rs/xid"
)

var (
	limit = uint64(2000)
)

func removeDeletedInstAsst(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	instAsstQueue := make(chan []metadata.InstAsst, 1000)

	assts := make([]*metadata.Association, 0)
	if err := db.Table(common.BKTableNameObjAsst).Find(nil).All(ctx, &assts); err != nil {
		blog.ErrorJSON("find %s error. error:%s", common.BKTableNameObjAsst, err.Error())
		return err
	}
	asstNameChan := make(chan string, len(assts))
	for _, asst := range assts {
		asstNameChan <- asst.AssociationName
	}
	if len(asstNameChan) == 0 {
		return nil
	}

	errChan := make(chan error, 1000)

	// 将所有需要处理的关联关系放到一个队列中
	queryInstAsstFunc := func() {
		for {
			asstName, ok := <-asstNameChan
			blog.InfoJSON("handle execute %s, has %s", asstName, len(asstNameChan))
			if !ok {
				return
			}
			maxID := int64(0)
			for {

				var assts []metadata.InstAsst
				var err error
				preMaxID := maxID
				assts, maxID, err = getInstAsst(ctx, db, asstName, maxID)

				if err != nil {
					blog.Errorf("get inst asst err:%s,flag:errChan", err.Error())
					errChan <- err
					return
				}
				if len(assts) == 0 {
					break
				}
				blog.Infof("get instacne association info. association name:%s, maxID:%d, next maxID:%d, info count:%d, queue:%d", asstName, preMaxID, maxID, len(assts), len(instAsstQueue))
				instAsstQueue <- assts
				if uint64(len(assts)) < limit {
					break
				}
			}

		}
	}

	handleAsstFunc := func() {
		rid := xid.New().String()
		for {
			blog.Infof("handle handleAsstFunc instAsstQueue len(%d), err:(%d). rid:%s", len(instAsstQueue), len(errChan), rid)
			if len(errChan) > 0 {
				return
			}
			blog.Infof("handle handleAsstFunc instAsstQueue start. rid:%s", rid)
			asst, ok := <-instAsstQueue
			if !ok {
				return
			}
			// 改成批量操作
			blog.Infof("handle handleAsstFunc goroutine run start. rid:%s", rid)
			if err := handleAsst(ctx, db, asst); err != nil {
				blog.Infof("handleAsst error. error:%s,flag:errChan", err.Error())
				errChan <- err
			}
		}

	}

	queryWait := handleCPUNoGoroutine(queryInstAsstFunc)
	handlehWait := handleCPUNoGoroutine(handleAsstFunc)

	go func() {
		timer := time.NewTicker(time.Second)
		for range timer.C {
			if len(asstNameChan) == 0 {
				close(asstNameChan)
				break
			}
		}
	}()

	blog.InfoJSON("query wait")
	select {
	case <-queryWait:
	case err := <-errChan:
		return err
	}
	blog.InfoJSON("query wait end")

	go func() {
		timer := time.NewTicker(time.Second)
		for range timer.C {
			blog.InfoJSON("query wait instAsstQueue %s", len(instAsstQueue))
			if len(instAsstQueue) == 0 {
				close(instAsstQueue)
				break
			}
		}
	}()

	blog.InfoJSON("handle wait")
	select {
	case <-handlehWait:
	case err := <-errChan:
		return err
	}

	return nil
}

// handleCPUNoGoroutine 开启2*cpu个数的协助
func handleCPUNoGoroutine(f func()) chan bool {
	var handlehWaitGroup sync.WaitGroup
	goroutineNu := runtime.NumCPU() * 2
	handlehWaitGroup.Add(goroutineNu)
	for idx := 0; idx < goroutineNu; idx++ {
		go func() {
			defer handlehWaitGroup.Done()
			f()
		}()
	}
	resultChan := make(chan bool, 10)
	go func() {
		handlehWaitGroup.Wait()
		resultChan <- true
	}()

	return resultChan
}

func handleAsst(ctx context.Context, db dal.RDB, asstArr []metadata.InstAsst) error {
	var objID string
	var asstObjID string
	var instIDArr []int64
	var asstInstIDArr []int64
	for _, asst := range asstArr {
		objID = asst.ObjectID
		asstObjID = asst.AsstObjectID
		instIDArr = append(instIDArr, asst.InstID)
		asstInstIDArr = append(asstInstIDArr, asst.AsstInstID)
	}

	instMap, err := getInst(ctx, db, objID, instIDArr)
	if err != nil {
		return err
	}

	asstInstMap, err := getInst(ctx, db, asstObjID, asstInstIDArr)
	if err != nil {
		return err
	}
	var delInstAsstID []int64
	for _, asst := range asstArr {
		if _, ok := instMap[asst.InstID]; !ok {
			delInstAsstID = append(delInstAsstID, asst.ID)
			continue
		}
		if _, ok := asstInstMap[asst.AsstInstID]; !ok {
			delInstAsstID = append(delInstAsstID, asst.ID)
			continue
		}
	}
	if len(delInstAsstID) == 0 {
		return nil
	}

	deleteCond := condition.CreateCondition()
	deleteCond.Field(common.BKFieldID).In(delInstAsstID)

	blog.InfoJSON("start delete cond:%s", deleteCond.ToMapStr())
	if err := db.Table(common.BKTableNameInstAsst).Delete(ctx, deleteCond.ToMapStr()); err != nil {
		blog.ErrorJSON("delete inst asst info error. err:%s", err.Error())
		return err
	}
	return nil
}

func getInst(ctx context.Context, db dal.RDB, objID string, instID []int64) (map[int64]bool, error) {
	idField := common.GetInstIDField(objID)
	insttable := common.GetInstTableName(objID)
	cond := condition.CreateCondition()
	cond.Field(idField).In(instID)
	if insttable == common.BKTableNameBaseInst {
		cond.Field(common.BKObjIDField).Eq(objID)
	}
	instArr := make([]map[string]int64, 0)
	rid := xid.New().String()

	blog.InfoJSON("get inst asst handle. obj:%s, idField:%s,cond:%s id:%s", objID, idField, cond.ToMapStr(), rid)
	err := db.Table(insttable).Find(cond.ToMapStr()).Fields(idField).All(ctx, &instArr)
	if err != nil {
		blog.ErrorJSON("get inst asst handle. cond:%s, err:%s", cond.ToMapStr(), err.Error())
		return nil, err
	}
	blog.InfoJSON("get inst asst handle. obj:%s, idField:%s, id:%s end", objID, idField, rid)
	instMap := make(map[int64]bool, 0)
	for _, inst := range instArr {
		id, ok := inst[idField]
		if !ok {
			blog.ErrorJSON("not found %s field. inst:%s", idField, inst)
			return nil, fmt.Errorf("object %s not found %s field", objID, idField)
		}
		instMap[id] = true
	}

	return instMap, nil

}

func getInstAsst(ctx context.Context, db dal.RDB, asstName string, maxID int64) ([]metadata.InstAsst, int64, error) {

	assts := []metadata.InstAsst{}
	cond := condition.CreateCondition()

	// 跟去ID来过滤数据
	cond.Field(common.BKFieldID).Gt(maxID)
	cond.Field(common.AssociationObjAsstIDField).Eq(asstName)
	blog.InfoJSON("get inst asst data, cond:%s", cond.ToMapStr())
	err := db.Table(common.BKTableNameInstAsst).Find(cond.ToMapStr()).Sort("+id").Limit(limit).All(ctx, &assts)
	if err != nil {
		blog.ErrorJSON("getInstAsst error:%s", err.Error())
		return assts, maxID, err
	}
	if len(assts) == 0 {
		return assts, maxID, nil
	}
	for _, asst := range assts {
		if maxID < asst.ID {
			maxID = asst.ID
		}
	}

	return assts, maxID, nil

}

func createInstanceAssociationIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	idxArr, err := db.Table(common.BKTableNameInstAsst).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error. err:%s", common.BKTableNameInstAsst, err.Error())
		return err
	}

	createIdxArr := []types.Index{
		types.Index{Name: "idx_objID_asstObjID_asstID", Keys: map[string]int32{"bk_obj_id": -1, "bk_asst_obj_id": -1, "bk_asst_id": -1}},
		types.Index{Name: "idx_asstID_id", Keys: map[string]int32{common.AssociationObjAsstIDField: -1, common.BKFieldID: -1}, Background: true, Unique: false},
	}
	for _, idx := range createIdxArr {
		exist := false
		for _, existIdx := range idxArr {
			if existIdx.Name == idx.Name {
				exist = true
				break
			}
		}
		// index already exist, skip create
		if exist {
			continue
		}
		if err := db.Table(common.BKTableNameInstAsst).CreateIndex(ctx, idx); err != nil {
			blog.ErrorJSON("create index to cc_InstAsst error, err:%s, current index:%s, all create index:%s", err.Error(), idx, createIdxArr)
			return err
		}

	}

	return nil

}
