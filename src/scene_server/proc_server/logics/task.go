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

package logics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

var (
	needOpGseProcResultStatus = []metadata.ProcOpTaskStatus{metadata.ProcOpTaskStatusWaitOP, metadata.ProcOpTaskStatusExecuteing}
)

func (lgc *Logics) ModifyTaskInfo(ctx context.Context, cond mapstr.MapStr, data mapstr.MapStr) errors.CCError {

	dat := new(metadata.UpdateParams)
	dat.Condition = cond
	dat.Data = data
	rsp, err := lgc.CoreAPI.ProcController().UpdateOperateTaskInfo(ctx, lgc.header, dat)
	if nil != err {
		blog.Errorf("ModifyTaskInfo http error:%s, input:%+v,rid:%s", err.Error(), dat, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("ModifyTaskInfo http reply error:%s, input:%+v,rid:%s", rsp.ErrMsg, dat, lgc.rid)
		return lgc.ccErr.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

// FilterGseTaskIDWaitResultByTaskID Filter out the need to go to gse to get the execution result of the gse task id
func (lgc *Logics) FilterGseTaskIDWaitResultByTaskID(ctx context.Context, taskID string) ([]string, errors.CCError) {

	dat := &metadata.QueryInput{}
	dat.Limit = common.BKNoLimit
	dat.Condition = mapstr.MapStr{
		common.BKTaskIDField: taskID,
		common.BKStatusField: mapstr.MapStr{common.BKDBIN: needOpGseProcResultStatus},
	}
	rsp, err := lgc.CoreAPI.ProcController().SearchOperateTaskInfo(ctx, lgc.header, dat)
	if nil != err {
		blog.Errorf("FilterGseTaskWaitResult http error:%s, input:%+v,rid:%s", err.Error(), dat, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("ModifyTaskStatus http reply error:%s,input:%+v, rid:%s", rsp.ErrMsg, dat, lgc.rid)
		return nil, lgc.ccErr.New(rsp.Code, rsp.ErrMsg)
	}

	gseTaskIDArr := make([]string, 0)
	for _, taskInfo := range rsp.Data.Info {
		gseTaskIDArr = append(gseTaskIDArr, taskInfo.GseTaskID)
	}

	return gseTaskIDArr, nil
}

// handleGseOPProcResult  backgroud handle gse operate process result
func (lgc *Logics) getGseOPProcTaskIDFromRedis(interval time.Duration) {
	for {
		val, err := lgc.cache.SPop(common.RedisProcSrvHostInstanceRefreshModuleKey).Result()
		if redis.Nil == err {
			if 0 >= interval {
				interval = GETTASKIDSPOPINTERVAL
			}
			time.Sleep(interval)
			continue
		}
		if nil != err {
			blog.Warnf("get timed trigger host instance event from redis,  error:%s", err.Error())
			continue
		}
		item := &opProcTask{}
		err = json.Unmarshal([]byte(val), item)
		if nil != err {
			blog.Warnf("get timed trigger host instance event from redis,  content not json,  error:%s", err.Error())
			continue
		}
		gseOPProcTaskChan <- item
	}
}

func (lgc *Logics) handleOPProcTask(ctx context.Context, taskID string) (waitExecArr []string, exceErrMap map[string]string, requestErr errors.CCError) {
	isErr := false
	waitExecArr = make([]string, 0)
	exceErrMap = make(map[string]string, 0)

	gseTaskIDArr, err := lgc.FilterGseTaskIDWaitResultByTaskID(ctx, taskID)
	if nil != err {
		blog.Errorf("handleOPProcTask query task info from gse  error, taskID:%s, error:%s logID:%s", taskID, err.Error(), lgc.rid)
		return nil, nil, err
	}

	for _, gseTaskID := range gseTaskIDArr {

		gseRet, err := lgc.esbServ.GseSrv().QueryProcOperateResult(ctx, lgc.header, gseTaskID)
		if err != nil {
			requestErr = err
			blog.Errorf("handleOPProcTask query task info from gse  error, taskID:%s, gseTaskID:%s, error:%s logID:%s", taskID, gseTaskID, err.Error(), lgc.rid)
			continue
		} else if !gseRet.Result {
			requestErr = lgc.ccErr.New(gseRet.Code, gseRet.Message)
			blog.Errorf("handleOPProcTask query task info from gse failed,  taskID:%s, gseTaskID:%s, gse return error:%s, error code:%d rid:%s", taskID, gseTaskID, gseRet.Message, gseRet.Code, lgc.rid)
			continue
		}

		for key, item := range gseRet.Data {
			if 0 != item.Errcode {
				if int(metadata.ProcOpTaskStatusExecuteing) == item.Errcode {
					waitExecArr = append(waitExecArr, key)
				} else {
					exceErrMap[key] = item.ErrMsg
					isErr = true
				}
			}
		}
		// same task group excute error,
		if isErr {
			conds := mapstr.MapStr{
				common.BKTaskIDField: taskID,
				common.BKStatusField: mapstr.MapStr{common.BKDBNE: metadata.ProcOpTaskStatusSucc},
			}
			data := mapstr.MapStr{common.BKStatusField: metadata.ProcOpTaskStatusErr, common.BKGseOpProcTaskDetailField: gseRet.Data}
			err := lgc.ModifyTaskInfo(ctx, conds, data)
			if nil != err {
				blog.Errorf("handleOPProcTask ModifyTaskStatus task detail  failed,  taskID:%s, gseTaskID:%s, gse return error:%s, error code:%d rid:%s", taskID, gseTaskID, gseRet.Message, gseRet.Code, lgc.rid)
				return nil, nil, err
			}

			conds = mapstr.MapStr{
				common.BKTaskIDField: taskID,
				common.BKStatusField: mapstr.MapStr{common.BKDBNE: metadata.ProcOpTaskStatusSucc},
			}
			data = mapstr.MapStr{common.BKStatusField: metadata.ProcOpTaskStatusErr}
			err = lgc.ModifyTaskInfo(ctx, conds, data)
			if nil != err {
				blog.Errorf("handleOPProcTask ModifyTaskStatus task status  failed,  taskID:%s, gseTaskID:%s, gse return error:%s, error code:%d rid:%s", taskID, gseTaskID, gseRet.Message, gseRet.Code, lgc.rid)
			}
			break
		}
		conds := mapstr.MapStr{
			common.BKTaskIDField:      taskID,
			common.BKGseOpTaskIDField: gseTaskID,
		}
		data := mapstr.MapStr{common.BKStatusField: metadata.ProcOpTaskStatusSucc, common.BKGseOpProcTaskDetailField: gseRet.Data}
		err = lgc.ModifyTaskInfo(ctx, conds, data)
		if nil != err {
			requestErr = err
			blog.Errorf("handleOPProcTask ModifyTaskStatus task status  failed,  taskID:%s, gseTaskID:%s, gse return error:%s, error code:%d rid:%s", taskID, gseTaskID, gseRet.Message, gseRet.Code, lgc.rid)
			continue
		}

	}
	if isErr {
		waitExecArr = make([]string, 0)
	}
	return waitExecArr, exceErrMap, requestErr

}

func (lgc *Logics) timedTriggerTaskInfoToRedis(ctx context.Context) {
	go func() {
		triggerChn := time.NewTicker(timedTriggerTaskTime)
		for range triggerChn.C {
			header := make(http.Header, 0)
			header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
			header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
			// use new header, so, new logics struct
			newLgc := lgc.NewFromHeader(header)

			dat := new(metadata.QueryInput)
			dat.Fields = fmt.Sprintf("%s,%s,http_header", common.BKTaskIDField, common.BKGseTaskIDField)
			dat.Limit = common.BKNoLimit
			dat.Condition = mapstr.MapStr{common.BKStatusField: mapstr.MapStr{common.BKDBIN: needOpGseProcResultStatus}}
			rsp, err := newLgc.CoreAPI.ProcController().SearchOperateTaskInfo(ctx, newLgc.header, dat)
			if nil != err {
				blog.V(5).Infof("timedTriggerTaskInfoToRedis http do error:%s,rid:%s", err.Error(), newLgc.rid)
				continue
			}
			if !rsp.Result {
				blog.V(5).Infof("timedTriggerTaskInfoToRedis http reply error:%s,rid:%s", rsp.ErrMsg, newLgc.rid)
				continue
			}
			if 0 == rsp.Data.Count {
				continue
			}
			opProcTaskMap := make(map[string]*opProcTask, 0)
			for _, taskInfo := range rsp.Data.Info {
				opProcTaskItem, ok := opProcTaskMap[taskInfo.TaskID]
				if !ok {
					opProcTaskItem = &opProcTask{
						TaskID: taskInfo.TaskID,
						Header: taskInfo.HTTPHeader,
					}
					opProcTaskMap[taskInfo.TaskID] = opProcTaskItem
				}
				opProcTaskItem.GseTaskIDArr = append(opProcTaskItem.GseTaskIDArr, taskInfo.GseTaskID)
			}
			cacheAllStr := make([]interface{}, 0)
			for _, item := range opProcTaskMap {
				byteInfo, err := json.Marshal(item)
				if nil != err {
					continue
				}
				cacheAllStr = append(cacheAllStr, string(byteInfo))
			}
			err = newLgc.cache.SAdd(common.RedisProcSrvQueryProcOPResultKey, cacheAllStr...).Err()
			if nil != err {
				blog.Warnf("timedTriggerTaskInfoToRedis set cache  error:%s,rid:%s", err.Error(), newLgc.rid)

			}

		}
	}()

}
