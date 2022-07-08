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

package service

import (
	"encoding/json"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type metaId struct {
	MongoID primitive.ObjectID `bson:"_id"`
}

const (
	// BKMaxDelPageSize max limit of delete operation.
	maxDelBatchLimit = 200
	// BKMaxDelDocPageLimit maximum number of documents deleted consecutively quantity.
	maxDelDocPageLimit = 10000
)

type deleteAuditLogReq struct {
	// delete logs before this day,the date format like '2021-08-19'.
	BeforeDay string `json:"beforeDay"`
}
type deleteAuditLogRsp struct {
	Num int `json:"num"`
}

// getMinObjIDAndMinDay 获取需要删除的最小日期和对应生成的objId
func (s *Service) getMinObjIDAndMinDay(baseDay int64, rid string) (primitive.ObjectID, int64, error) {

	// 根据指定删除的时间点(注意时间点是当天的0点，例如如果指定的是2021-08-19，指的是8月19日的0点)生成 objectId. 后续流程会将小于此
	// 时间戳的数据全部删掉.
	objId := primitive.NewObjectIDFromTimestamp(time.Unix(baseDay, 0))

	for {
		cond := map[string]interface{}{
			"_id": map[string]interface{}{
				common.BKDBLT: objId,
			},
		}
		count, err := s.db.Table(common.BKTableNameAuditLog).Find(cond).Count(s.ctx)
		if err != nil {
			return primitive.ObjectID{}, 0, err
		}
		if count > 0 {
			baseDay -= 24 * 60 * 60
			dayAgo := time.Unix(baseDay, 0)
			objId = primitive.NewObjectIDFromTimestamp(dayAgo)
		} else {
			break
		}
	}
	blog.Infof("getMinObjIDAndMinDay,the min day is: %s,rid: %s", baseDay, rid)

	return objId, baseDay, nil
}

// DeleteAuditLog delete user specified audit logs.
// 删除策略: 1、首先找到最早一天的审计日志，从前向后一天一天的删除审计日志。
//          2、每次批量删除200条日志。为了防止删除导致的cpu和磁盘io过高，每删除10000条日志需要sleep 5秒钟。
func (s *Service) DeleteAuditLog(req *restful.Request, resp *restful.Response) {

	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	param := new(deleteAuditLogReq)
	response := new(deleteAuditLogRsp)

	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("deleteAuditLog, decode body failed, err: %+v,rid: %s", err, rid)
		errInfo := metadata.RespError{Msg: defErr.CCError(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, &errInfo)
		return
	}
	blog.Infof("deleteAuditLog,the user specified date is: %s,rid: %s", param.BeforeDay, rid)

	// convert string format to timestamp.
	baseDay := util.TimeStrToUnixSecondDefault(param.BeforeDay)
	objId, minDay, err := s.getMinObjIDAndMinDay(baseDay, rid)
	if err != nil {
		blog.Errorf("get the earliest date to delete failed, err: %+v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	var cnt, total int

	for {
		// delete the data before the time point specified by the user.
		if minDay > baseDay {
			break
		}
		cond := map[string]interface{}{
			"_id": map[string]interface{}{
				common.BKDBLT: objId,
			},
		}
		metaIdList := make([]metaId, 0)

		// find docs for the specified date.
		err := s.db.Table(common.BKTableNameAuditLog).Find(cond).Fields("_id").
			Limit(uint64(maxDelBatchLimit)).All(s.ctx, &metaIdList)
		if err != nil {
			blog.Errorf("search auditLog failed, err: %+v, rid: %s", err, rid)
			_ = resp.WriteError(http.StatusInternalServerError, err)
			return
		}

		// the document of the specified date has been deleted，Find the content of the next day.
		if len(metaIdList) <= 0 {
			// convert time from day to seconds.
			minDay += 24 * 60 * 60
			dayAgo := time.Unix(minDay, 0)
			objId = primitive.NewObjectIDFromTimestamp(dayAgo)
			continue
		}

		mongoIDs := make([]primitive.ObjectID, len(metaIdList))
		for index, data := range metaIdList {
			mongoIDs[index] = data.MongoID
		}
		delCond := map[string]interface{}{"_id": map[string]interface{}{common.BKDBIN: mongoIDs}}

		if err := s.db.Table(common.BKTableNameAuditLog).Delete(s.ctx, delCond); err != nil {
			blog.Errorf("search auditLog failed, objIds: %v,err: %+v, rid: %s", mongoIDs, err, rid)
			errInfo := metadata.RespError{Msg: err}
			_ = resp.WriteError(http.StatusBadRequest, &errInfo)
			return
		}

		cnt += len(metaIdList)
		total += len(metaIdList)
		if cnt >= maxDelDocPageLimit {
			time.Sleep(5 * time.Second)
			cnt = 0
			t := time.Unix(minDay, 0)
			dateStr := t.Format("2006-01-02")

			blog.Infof("the delete date is: %s,the number of deleted items is: %d,rid: %s",
				dateStr, total, rid)
		}
	}
	blog.Infof(" delete all completed audit logs,the num: %d,rid: %s", total, rid)

	response.Num = total
	_ = resp.WriteEntity(metadata.NewSuccessResp(response))
	return
}
