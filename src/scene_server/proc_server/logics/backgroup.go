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
	"time"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

func (lgc *Logics) reshReshInitChan(ctx context.Context, config *ProcHostInstConfig) {
	if nil == config {
		config = &ProcHostInstConfig{
			MaxEventCount:         0,
			MaxRefreshModuleCount: 0,
			GetModuleIDInterval:   0,
		}
	}
	if 0 != config.MaxEventCount {
		maxEventDataChan = config.MaxEventCount
	}
	if 0 != config.MaxRefreshModuleCount {
		maxRefreshModuleData = config.MaxRefreshModuleCount
	}
	handEventDataChan = make(chan chanItem, maxEventDataChan)
	refreshHostInstModuleIDChan = make(chan *refreshHostInstModuleID, maxRefreshModuleData)
	// get appID,moduleID from redis
	go lgc.getEventRefreshModuleItemFromRedis(config.GetModuleIDInterval)
	go lgc.backgroudHandleOpGseProcTaskResult(ctx, config.FetchGseOPProcResultInterval)
}

// getEventRefreshModuleItemFromRedis Run once at startup
func (lgc *Logics) getEventRefreshModuleItemFromRedis(interval time.Duration) {
	for {
		val, err := lgc.cache.SPop(common.RedisProcSrvHostInstanceRefreshModuleKey).Result()
		if redis.Nil == err {
			if 0 >= interval {
				interval = SPOPINTERVAL
			}
			time.Sleep(interval)
			continue
		}
		if nil != err {
			blog.Warnf("getEventRefreshModuleItemFromRedis error:%s,rid:%s", err.Error(), lgc.rid)
			continue
		}
		item := &refreshHostInstModuleID{}
		err = json.Unmarshal([]byte(val), item)
		if nil != err {
			blog.Warnf("getEventRefreshModuleItemFromRedis  error:%s,rid:%s", err.Error(), lgc.rid)
			continue
		}
		refreshHostInstModuleIDChan <- item

	}

}

func (lgc *Logics) backgroudHandleOpGseProcTaskResult(ctx context.Context, interval time.Duration) {
	go lgc.getGseOPProcTaskIDFromRedis(interval)
	go lgc.timedTriggerTaskInfoToRedis(ctx)
	for {
		select {
		case taskInfo := <-gseOPProcTaskChan:
			newLgc := lgc.NewFromHeader(taskInfo.Header)
			waitExecArr, _, requestErr := newLgc.handleOPProcTask(ctx, taskInfo.TaskID)
			if nil != requestErr || 0 < len(waitExecArr) {
				taskInfoByte, err := json.Marshal(taskInfo)
				if nil != err {
					blog.Warnf("backgroudHandleOpGseProcTaskResult json marshal error, raw data:%+v error:%s,rid:%s", taskInfo, err.Error(), newLgc.rid)
					continue
				}

				_, err = newLgc.cache.SAdd(common.RedisProcSrvQueryProcOPResultKey, string(taskInfoByte)).Result()
				if nil != err {
					blog.Warnf("backgroudHandleOpGseProcTaskResult cache task info  error, task info:%s error:%s,rid:%s", string(taskInfoByte), err.Error(), newLgc.rid)
				}
			}
		}
	}
}
