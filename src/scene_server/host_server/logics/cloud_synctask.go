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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	com "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

var (
	taskChan            = make(map[int64]chan bool)
	checkDuration int64 = 5
)

func (lgc *Logics) AddCloudTask(ctx context.Context, taskList *meta.CloudTaskList) error {
	// TaskName Uniqueness check
	resp, err := lgc.CoreAPI.CoreService().Cloud().CheckTaskNameUnique(ctx, lgc.header, taskList)
	if err != nil {
		return err
	}

	if resp.Count != 0 {
		blog.Errorf("add task failed, task name %s already exits, rid: %s", taskList.TaskName, lgc.rid)
		return lgc.ccErr.Error(1110038)
	}

	// Encode secretKey
	taskList.SecretKey = base64.StdEncoding.EncodeToString([]byte(taskList.SecretKey))

	if _, err := lgc.CoreAPI.CoreService().Cloud().CreateCloudSyncTask(ctx, lgc.header, taskList); err != nil {
		blog.Errorf("add cloud task failed, err: %v, rid: %s", err, lgc.rid)
		return err
	}

	return nil
}

func (lgc *Logics) TimerTriggerCheckStatus(ctx context.Context) {
	go func() {
		lgc.SyncTaskDBManager(ctx)
		timer := time.NewTicker(time.Duration(checkDuration) * time.Minute)
		for range timer.C {
			lgc.CompareRedisWithDB(ctx)
			lgc.CheckSyncAlive(ctx)
			lgc.SyncTaskRedisStopManager(ctx)
		}
	}()
	go lgc.SyncTaskRedisStartManager(ctx)
	go lgc.ListenRedisSubscribe(ctx)
}

func (lgc *Logics) SyncTaskDBManager(ctx context.Context) {
	mutex := &sync.Mutex{}
	if isMaster := lgc.Engine.ServiceManageInterface.IsMaster(); !isMaster {
		blog.Errorf("not master, stop syncTaskDBManager, rid: %v", lgc.rid)
		return
	}
	opt := make(map[string]interface{}, 0)
	resp, err := lgc.CoreAPI.CoreService().Cloud().SearchCloudSyncTask(ctx, lgc.header, opt)
	if err != nil {
		blog.Errorf("get cloud sync task instance failed, err: %v, rid: %s", err, lgc.rid)
		return
	}

	for _, taskInfo := range resp.Info {
		if taskInfo.Status {
			newHeader := make(http.Header, 0)
			ownerID := taskInfo.OwnerID
			newHeader.Set(common.BKHTTPOwnerID, ownerID)
			newHeader.Set(common.BKHTTPHeaderUser, taskInfo.User)

			taskID := taskInfo.TaskID
			mutex.Lock()
			if _, ok := taskChan[taskID]; ok {
				continue
			}
			mutex.Unlock()
			nextTrigger := lgc.NextTrigger(ctx, taskInfo.PeriodType, taskInfo.Period)
			taskInfoItem := &meta.TaskInfo{
				Method:      taskInfo.PeriodType,
				NextTrigger: nextTrigger,
				Args:        taskInfo,
			}

			info := meta.CloudSyncRedisPendingStart{TaskID: taskID, TaskItemInfo: *taskInfoItem, OwnerID: ownerID, NewHeader: newHeader}
			pendingStartTaskInfo, err := json.Marshal(info)
			if err != nil {
				blog.Errorf("add redis failed taskID: %v, accountAdmin: %v, rid: %s", taskInfo.TaskID, taskInfo.AccountAdmin, lgc.rid)
				continue
			}

			if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStart, pendingStartTaskInfo).Err(); err != nil {
				blog.Errorf("add cloud task item to redis fail, info: %v err: %v, rid: %s", info, err, lgc.rid)
				continue
			}
		}
	}

	return
}

func (lgc *Logics) FrontEndSyncSwitch(ctx context.Context, opt map[string]interface{}, update bool) error {
	response, err := lgc.CoreAPI.CoreService().Cloud().SearchCloudSyncTask(ctx, lgc.header, opt)
	if err != nil {
		blog.Errorf("search cloud task instance failed, err: %v, rid: %s", err, lgc.rid)
		return lgc.ccErr.Error(1110036)
	}

	if response.Count > 0 {
		taskInfo := response.Info[0]
		status := taskInfo.Status
		taskID := taskInfo.TaskID

		if status {
			nextTrigger := lgc.NextTrigger(ctx, taskInfo.PeriodType, taskInfo.Period)
			taskInfoItem := meta.TaskInfo{
				Method:      taskInfo.PeriodType,
				NextTrigger: nextTrigger,
				Args:        taskInfo,
			}

			ownerID := util.GetOwnerID(lgc.header)
			info := meta.CloudSyncRedisPendingStart{TaskID: taskID, TaskItemInfo: taskInfoItem, OwnerID: ownerID, NewHeader: lgc.header, Update: update}

			pendingStartTaskInfo, err := json.Marshal(info)
			if err != nil {
				blog.Errorf("add redis failed taskID: %v, accountAdmin: %v， rid: %s, error: %v", taskInfo.TaskID, taskInfo.AccountAdmin, lgc.rid, err)
				return err
			}
			if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStart, pendingStartTaskInfo).Err(); err != nil {
				blog.Errorf("add cloud task redis item fail, info: %v, err: %v, rid: %s", info, err, lgc.rid)
				return err
			}
		} else {
			if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStop, taskInfo.TaskID).Err(); err != nil {
				blog.Errorf("add cloud task redis item fail, info: %v, err: %v, rid: %s", taskInfo.TaskID, err, lgc.rid)
				return err
			}
		}
	}

	return nil
}

// ListenRedisSubscribe subscribe redis channel to stop the started sync task
func (lgc *Logics) ListenRedisSubscribe(ctx context.Context) {
	var mutex = &sync.Mutex{}
	newClient := *lgc.cache

	for {
		pub, err := newClient.Subscribe("stop")
		if err != nil {
			time.Sleep(5 * time.Second)
			blog.Errorf("redis subscribe fail, err: %v, rid: %v", err, lgc.rid)
			continue
		}
		for {
			receive, err := pub.ReceiveMessage()
			if err != nil {
				blog.Errorf("redis subscribe get value fail, err: %v, rid: %v", err, lgc.rid)
				continue
			}

			taskID, err := strconv.ParseInt(receive.Payload, 10, 64)
			if err != nil {
				blog.Errorf("interface convert to int64 fail, err: %v, rid: %v", err, lgc.rid)
				continue
			}
			mutex.Lock()
			if _, ok := taskChan[taskID]; ok {
				taskChan[taskID] <- true
			}
			mutex.Unlock()
		}
	}
}

func (lgc *Logics) SyncTaskRedisStartManager(ctx context.Context) {
	var mutex = &sync.Mutex{}

	for {
		val, err := lgc.cache.BLPop(0, common.RedisCloudSyncInstancePendingStart).Result()
		if err != nil {
			blog.Warnf("get task pending start item from redis fail, taskInfo: %s, err:%v, rid: %s", val, err.Error(), lgc.rid)
			continue
		}
		if len(val) == 0 {
			continue
		}
		item := val[1]
		pendingStartItem := meta.CloudSyncRedisPendingStart{}
		if err := json.Unmarshal([]byte(item), &pendingStartItem); err != nil {
			blog.Warnf("get task pending start item from redis fail, taskInfo: %s, err:%v, rid: %s", item, err.Error(), lgc.rid)
			continue
		}
		taskID := pendingStartItem.TaskID

		if pendingStartItem.Update {
			mutex.Lock()
			if _, ok := taskChan[taskID]; ok {
				taskChan[taskID] <- true
			} else {
				lgc.cache.Publish("stop", strconv.FormatInt(taskID, 10))
			}
			mutex.Unlock()
		}
		taskInfoItem := pendingStartItem.TaskItemInfo

		ownerID := util.GetOwnerID(pendingStartItem.NewHeader)
		newLgc := lgc.NewFromHeader(pendingStartItem.NewHeader)
		taskChannel := make(chan bool, 0)

		mutex.Lock()
		taskChan[taskID] = taskChannel
		mutex.Unlock()

		waitStopItems, err := lgc.cache.LRange(common.RedisCloudSyncInstancePendingStop, 0, -1).Result()
		if err != nil {
			blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
		}

		if len(waitStopItems) > 0 {
			for _, stopTaskID := range waitStopItems {
				intStopItem, err := strconv.ParseInt(stopTaskID, 10, 64)
				if err != nil {
					blog.Errorf("string convert to int64 fail, taskID: %v, rid: %v", intStopItem, lgc.rid)
					continue
				}
				if taskID == intStopItem {
					if err := lgc.cache.LRem(common.RedisCloudSyncInstancePendingStop, 1, stopTaskID).Err(); err != nil {
						blog.Errorf("remove stop task item fail, taskInfo: %s, error: %v, rid: %s", item, err, lgc.rid)
						continue
					}

					lgc.deleteStartedTaskRedis(ctx, taskID)
					info := meta.CloudSyncRedisAlreadyStarted{TaskID: taskID, TaskItemInfo: taskInfoItem, OwnerID: ownerID, LastSyncTime: time.Now(), NewHeader: pendingStartItem.NewHeader}
					startedTaskInfo, err := json.Marshal(info)
					if err != nil {
						blog.Errorf("add redis failed, info: %v, err: %v, rid: %s", info, err, lgc.rid)
						continue
					}
					if err := lgc.cache.RPush(common.RedisCloudSyncInstanceStarted, startedTaskInfo).Err(); err != nil {
						blog.Errorf("add cloud task item to redis fail, info: %v, err: %v, rid: %s", info, err, lgc.rid)
						continue
					}
					newLgc.CloudSyncSwitch(ctx, &taskInfoItem)
					continue
				}
			}
			continue
		}

		info := meta.CloudSyncRedisAlreadyStarted{TaskID: taskID, TaskItemInfo: taskInfoItem, OwnerID: ownerID, LastSyncTime: time.Now(), NewHeader: pendingStartItem.NewHeader}
		startedTaskInfo, err := json.Marshal(info)
		if err != nil {
			blog.Errorf("add redis failed, info: %v, err: %v, rid: %s", info, err, lgc.rid)
			continue
		}
		if err := lgc.cache.RPush(common.RedisCloudSyncInstanceStarted, startedTaskInfo).Err(); err != nil {
			blog.Errorf("add cloud task item to redis fail, info: %v, err: %v, rid: %s", info, err, lgc.rid)
			continue
		}

		newLgc.CloudSyncSwitch(ctx, &taskInfoItem)
	}
}

func (lgc *Logics) SyncTaskRedisStopManager(ctx context.Context) {
	var mutex = &sync.Mutex{}

	redisTaskItems, err := lgc.cache.LRange(common.RedisCloudSyncInstancePendingStop, 0, -1).Result()
	if err != nil {
		blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
		return
	}

	if len(redisTaskItems) == 0 {
		return
	}

	for _, item := range redisTaskItems {
		stopTaskId, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			blog.Errorf("string convert to int64 fail, taskID: %v, rid: %v", stopTaskId, lgc.rid)
			continue
		}

		for key := range taskChan {
			if stopTaskId == key {
				if err := lgc.cache.LRem(common.RedisCloudSyncInstancePendingStop, 1, item).Err(); err != nil {
					blog.Errorf("remove stop task item fail, taskInfo: %s, error: %v, rid: %s", item, err, lgc.rid)
					continue
				}
				mutex.Lock()
				taskChan[key] <- true
				mutex.Lock()
			}
		}
	}
}

func (lgc *Logics) CheckSyncAlive(ctx context.Context) {
	startedTaskItems, err := lgc.cache.LRange(common.RedisCloudSyncInstanceStarted, 0, -1).Result()
	if err != nil {
		blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
		return
	}

	if len(startedTaskItems) == 0 {
		return
	}

	needStartAgain := make([]meta.CloudSyncRedisAlreadyStarted, 0)
	for _, item := range startedTaskItems {
		startedItem := meta.CloudSyncRedisAlreadyStarted{}
		if err := json.Unmarshal([]byte(item), &startedItem); err != nil {
			blog.Warnf("get task started item from redis fail, taskInfo: %s, error:%s, rid: %s", item, err.Error(), lgc.rid)
			continue
		}

		timeInterval := time.Now().Unix() - startedItem.LastSyncTime.Unix()

		switch startedItem.TaskItemInfo.Method {
		case "day":
			if timeInterval > 90000 {
				needStartAgain = append(needStartAgain, startedItem)
			}
		case "hour":
			if timeInterval > 5400 {
				needStartAgain = append(needStartAgain, startedItem)
			}
		case "minute":
			if timeInterval > 600 {
				needStartAgain = append(needStartAgain, startedItem)
			}
		}
	}

	if len(needStartAgain) == 0 {
		return
	}

	blog.V(5).Info("needStartAgain: %v, rid: %s", needStartAgain, lgc.rid)

	for _, item := range needStartAgain {
		if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStop, item.TaskID).Err(); err != nil {
			blog.Errorf("add cloud task redis item fail, taskID: %v, err: %v, rid: %s", item.TaskID, err, lgc.rid)
			continue
		}

		startInfo := meta.CloudSyncRedisPendingStart{TaskID: item.TaskID, TaskItemInfo: item.TaskItemInfo, OwnerID: item.OwnerID, NewHeader: item.NewHeader}
		pendingStartTaskInfo, err := json.Marshal(startInfo)
		if err != nil {
			blog.Errorf("add redis failed taskID: %v, rid: %s", item.TaskID, lgc.rid)
			continue
		}
		if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStart, pendingStartTaskInfo).Err(); err != nil {
			blog.Errorf("add cloud task redis item fail, taskID: %v, err: %v, rid: %s", item.TaskID, err, lgc.rid)
			continue
		}

		lgc.deleteStartedTaskRedis(ctx, item.TaskID)
	}
}

func (lgc *Logics) CompareRedisWithDB(ctx context.Context) {
	// master 才可以往redis写数据，避免写入重复数据
	if ok := lgc.Engine.ServiceManageInterface.IsMaster(); !ok {
		return
	}

	header := copyHeader(ctx, lgc.header)
	if nil == header {
		header = make(http.Header, 0)
	}
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}
	newLgc := lgc.NewFromHeader(header)

	opt := make(map[string]interface{})
	response, err := newLgc.CoreAPI.CoreService().Cloud().SearchCloudSyncTask(ctx, lgc.header, opt)
	if err != nil {
		blog.Errorf("search cloud task info fail, err: %v, rid: %s", err, lgc.rid)
		return
	}

	if response.Count == 0 {
		return
	}

	pendingStopItems, err := lgc.cache.LRange(common.RedisCloudSyncInstancePendingStop, 0, -1).Result()
	if err != nil {
		blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
		return
	}

	pendingStartItems, err := lgc.cache.LRange(common.RedisCloudSyncInstancePendingStart, 0, -1).Result()
	if err != nil {
		blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
		return
	}

	startedItems, err := lgc.cache.LRange(common.RedisCloudSyncInstanceStarted, 0, -1).Result()
	if err != nil {
		blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
		return
	}

	allStartItems := make([]string, 0)
	allStartItems = append(allStartItems, pendingStartItems...)
	allStartItems = append(allStartItems, startedItems...)

	startTaskArr := make([]meta.CloudSyncRedisPendingStart, 0)
	for _, item := range allStartItems {
		startItem := meta.CloudSyncRedisPendingStart{}
		if err := json.Unmarshal([]byte(item), &startItem); err != nil {
			blog.Warnf("get task pending start item from redis fail, error:%s, rid: %s", err.Error(), lgc.rid)
			continue
		}
		startTaskArr = append(startTaskArr, startItem)
	}

	shouldStartItems := make([]meta.CloudSyncRedisPendingStart, 0)
	shouldStopItems := make([]int64, 0)
	var mutex = &sync.Mutex{}

	for _, dbItem := range response.Info {
		if dbItem.Status {
			for _, item := range startTaskArr {
				if dbItem.TaskID == item.TaskID {
					continue
				}
				shouldStartItems = append(shouldStartItems, item)
				itemChan := make(chan bool, 0)

				mutex.Lock()
				taskChan[dbItem.TaskID] = itemChan
				mutex.Unlock()
			}
		} else {
			for _, item := range pendingStopItems {
				int64Item, err := strconv.ParseInt(item, 10, 64)
				if err != nil {
					blog.Errorf("string convert to int64 fail,taskID: %v ,err: %v, rid: %v", item, err, lgc.rid)
					continue
				}
				if dbItem.TaskID == int64Item {
					continue
				}
				shouldStopItems = append(shouldStopItems, int64Item)
			}
		}
	}

	if len(shouldStartItems) > 0 {
		for _, item := range shouldStartItems {
			waitStart, err := json.Marshal(item)
			if err != nil {
				blog.Errorf("add redis failed taskID: %v, rid: %s", item.TaskID, lgc.rid)
				return
			}
			if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStart, waitStart).Err(); err != nil {
				blog.Errorf("add cloud task item to redis fail, err: %v, rid: %s", err, lgc.rid)
				return
			}
		}
	}
	if len(shouldStopItems) > 0 {
		for _, item := range shouldStopItems {
			if err := lgc.cache.RPush(common.RedisCloudSyncInstancePendingStop, item).Err(); err != nil {
				blog.Errorf("add cloud task item to redis fail, err: %v, rid: %s", err, lgc.rid)
				return
			}
		}
	}

	return
}

func (lgc *Logics) CloudSyncSwitch(ctx context.Context, taskInfoItem *meta.TaskInfo) {
	mutex := &sync.Mutex{}

	timer := time.NewTimer(time.Duration(taskInfoItem.NextTrigger) * time.Minute)
	go func() {
		for {
			select {
			case <-timer.C:
				lgc.ExecSync(ctx, taskInfoItem.Args)
				switch taskInfoItem.Method {
				case "day":
					taskInfoItem.NextTrigger = 1440
				case "hour":
					taskInfoItem.NextTrigger = 60
				case "minute":
					taskInfoItem.NextTrigger = 5
				}
				timer.Reset(time.Duration(taskInfoItem.NextTrigger) * time.Minute)
			case <-taskChan[taskInfoItem.Args.TaskID]:
				mutex.Lock()
				close(taskChan[taskInfoItem.Args.TaskID])
				delete(taskChan, taskInfoItem.Args.TaskID)
				mutex.Unlock()
				lgc.deleteStartedTaskRedis(ctx, taskInfoItem.Args.TaskID)
				return
			}
		}
	}()
}

func (lgc *Logics) ExecSync(ctx context.Context, taskInfo meta.CloudTaskInfo) {
	cloudHistory := new(meta.CloudHistory)
	cloudHistory.ObjID = taskInfo.ObjID
	cloudHistory.TaskID = taskInfo.TaskID
	startTime := time.Now().Unix()

	var errOrigin error
	defer func() {
		if errOrigin != nil {
			cloudHistory.Status = "fail"
			errString := fmt.Sprintf("%s", errOrigin)
			if strings.Contains(errString, "AuthFailure") {
				cloudHistory.FailReason = "AuthFailure"
			} else {
				cloudHistory.FailReason = "else"
			}
		}
		lgc.CloudSyncHistory(ctx, taskInfo.TaskID, startTime, cloudHistory)
	}()

	// obtain the hosts from cc_HostBase
	body := new(meta.HostCommonSearch)
	host, err := lgc.SearchHost(ctx, body, false)
	if err != nil {
		blog.Errorf("search host failed, err: %v, rid: %s", err, lgc.rid)
		errOrigin = err
		return
	}

	existHostList := make([]string, 0)
	for i := 0; i < host.Count; i++ {
		hostInfo, err := mapstr.NewFromInterface(host.Info[i]["host"])
		if err != nil {
			blog.Errorf("get hostInfo failed with err: %v, rid: %s", err, lgc.rid)
			errOrigin = err
			return
		}

		ip, err := hostInfo.String(common.BKHostInnerIPField)
		if err != nil {
			blog.Errorf("get hostIp failed with err: %v, rid: %s", err, lgc.rid)
			errOrigin = err
			return
		}

		existHostList = append(existHostList, ip)
	}

	// obtain hosts from TencentCloud needs secretID and secretKey
	decodeBytes, err := base64.StdEncoding.DecodeString(taskInfo.SecretKey)
	if err != nil {
		blog.Errorf("Base64 decode secretKey failed, rid: %s", lgc.rid)
		errOrigin = err
		return
	}
	secretKey := string(decodeBytes)
	secretID := taskInfo.SecretID

	// ObtainCloudHosts obtain cloud hosts
	cloudHostInfo, err := lgc.ObtainCloudHosts(ctx, secretID, secretKey)
	if err != nil {
		blog.Errorf("obtain cloud hosts failed with err: %v, rid: %s", err, lgc.rid)
		errOrigin = err
		return
	}

	// pick out the new add cloud hosts
	newAddHost := make([]string, 0)
	newCloudHost := make([]mapstr.MapStr, 0)
	for _, hostInfo := range cloudHostInfo {
		newHostInnerip, ok := hostInfo[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Errorf("interface convert to string failed, rid: %s", lgc.rid)
		}
		if !util.InStrArr(existHostList, newHostInnerip) {
			newAddHost = append(newAddHost, newHostInnerip)
			newCloudHost = append(newCloudHost, hostInfo)
		}
	}

	// pick out the hosts that has changed attributes
	cloudHostAttr := make([]mapstr.MapStr, 0)
	for _, hostInfo := range cloudHostInfo {
		newHostInnerip, ok := hostInfo[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Errorf("interface convert to string failed, err: %v, rid: %s", err, lgc.rid)
			continue
		}
		newHostOuterip, ok := hostInfo[common.BKHostOuterIPField].(string)
		if !ok {
			blog.Errorf("interface convert to string failed, err: %v, rid: %s", err, lgc.rid)
			continue
		}
		newHostOsname, ok := hostInfo[common.BKOSNameField].(string)
		if !ok {
			blog.Errorf("interface convert to string failed, err: %v, rid: %s", err, lgc.rid)
			continue
		}

		for i := 0; i < host.Count; i++ {
			existHostInfo, err := mapstr.NewFromInterface(host.Info[i]["host"])
			if err != nil {
				blog.Errorf("get hostInfo failed with err: %v, rid: %s", err, lgc.rid)
				errOrigin = err
				return
			}

			existHostIp, ok := existHostInfo.String(common.BKHostInnerIPField)
			if ok != nil {
				blog.Errorf("get hostIp failed with err: %v, rid: %s", ok, lgc.rid)
				errOrigin = ok
				break
			}
			existHostOsname, ok := existHostInfo.String(common.BKOSNameField)
			if ok != nil {
				blog.Errorf("get os name failed with err: %v, rid: %s", ok, lgc.rid)
				errOrigin = ok
				break
			}

			existHostOuterip, ok := existHostInfo.String(common.BKHostOuterIPField)
			if ok != nil {
				blog.Errorf("get outerip failed with, rid: %s", lgc.rid)
				errOrigin = ok
				break
			}

			existHostID, ok := existHostInfo.String(common.BKHostIDField)
			if ok != nil {
				blog.Errorf("get hostID failed, rid: %s", lgc.rid)
				errOrigin = ok
				break
			}

			if existHostIp == newHostInnerip {
				if existHostOsname != newHostOsname || existHostOuterip != newHostOuterip {
					hostInfo[common.BKHostIDField] = existHostID
					cloudHostAttr = append(cloudHostAttr, hostInfo)
				}
			}
		}
	}

	cloudHistory.NewAdd = len(newAddHost)
	cloudHistory.AttrChanged = len(cloudHostAttr)

	attrConfirm := taskInfo.AttrConfirm
	resourceConfirm := taskInfo.ResourceConfirm

	if !resourceConfirm && !attrConfirm {
		if len(newCloudHost) > 0 {
			err := lgc.AddCloudHosts(ctx, newCloudHost)
			if err != nil {
				blog.Errorf("add cloud hosts failed, err: %v, rid: %s", err, lgc.rid)
				errOrigin = err
				return
			}
		}
		if len(cloudHostAttr) > 0 {
			err := lgc.UpdateCloudHosts(ctx, cloudHostAttr)
			if err != nil {
				blog.Errorf("update cloud hosts failed, err: %v, rid: %s", err, lgc.rid)
				errOrigin = err
				return
			}
		}
	}

	if resourceConfirm {
		newAddNum, err := lgc.NewAddConfirm(ctx, taskInfo, newCloudHost)
		cloudHistory.NewAdd = newAddNum
		if err != nil {
			blog.Errorf("newly add cloud resource confirm failed, err: %v, rid: %s", err, lgc.rid)
			errOrigin = err
			return
		}
	}

	if attrConfirm && len(cloudHostAttr) > 0 {
		blog.V(5).Info("attr chang, rid: %s", lgc.rid)

		for _, host := range cloudHostAttr {
			resourceConfirm := mapstr.MapStr{}
			resourceConfirm["bk_obj_id"] = taskInfo.ObjID
			innerIp, err := host.String(common.BKHostInnerIPField)
			if err != nil {
				blog.Errorf("mapstr.Map convert to string failed, rid: %s", lgc.rid)
				errOrigin = err
				return
			}
			outerIp, err := host.String(common.BKHostOuterIPField)
			if err != nil {
				blog.Error("mapstr.Map convert to string failed, rid: %s", lgc.rid)
				errOrigin = err
				return
			}
			osName, err := host.String(common.BKOSNameField)
			if err != nil {
				blog.Error("mapstr.Map convert to string failed, rid: %s", lgc.rid)
				errOrigin = err
				return
			}

			resourceConfirm[common.BKHostInnerIPField] = innerIp
			resourceConfirm[common.BKHostOuterIPField] = outerIp
			resourceConfirm[common.BKOSNameField] = osName
			resourceConfirm[common.BKCloudTaskID] = taskInfo.TaskID
			resourceConfirm[common.BKAttrConfirm] = attrConfirm
			resourceConfirm[common.BKCloudConfirm] = false
			resourceConfirm[common.BKCloudSyncTaskName] = taskInfo.TaskName
			resourceConfirm[common.BKCloudAccountType] = taskInfo.AccountType
			resourceConfirm[common.BKCloudSyncAccountAdmin] = taskInfo.AccountAdmin
			resourceConfirm[common.BKResourceType] = "change"

			if _, err := lgc.CoreAPI.CoreService().Cloud().CreateConfirm(ctx, lgc.header, resourceConfirm); err != nil {
				blog.Errorf("add resource confirm failed with confirmInfo: %#v, err: %v, rid: %s", resourceConfirm, err, lgc.rid)
				errOrigin = err
				return
			}
		}
		return
	}

	cloudHistory.Status = "success"
	blog.V(3).Info("finish sync, rid: %s", lgc.rid)
	return
}

func (lgc *Logics) AddCloudHosts(ctx context.Context, newCloudHost []mapstr.MapStr) error {
	hostList := new(meta.HostList)
	hostInfoMap := make(map[int64]map[string]interface{}, 0)
	appID := hostList.ApplicationID

	if appID == 0 {
		// get default app id
		var err error
		appID, err = lgc.GetDefaultAppIDWithSupplier(ctx)
		if err != nil {
			blog.Errorf("add host, but get default appid failed, err: %v, rid: %s", err, lgc.rid)
			return err
		}
	}

	cond := hutil.NewOperation().WithModuleName(common.DefaultResModuleName).WithAppID(appID).Data()
	cond[common.BKDefaultField] = common.DefaultResModuleFlag
	moduleID, err := lgc.GetResourcePoolModuleID(ctx, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s, rid: %s", err.Error(), lgc.rid)
		return err
	}

	blog.V(5).Infof("resource confirm add new hosts, rid: %s", lgc.rid)
	for index, hostInfo := range newCloudHost {
		if _, ok := hostInfoMap[int64(index)]; !ok {
			hostInfoMap[int64(index)] = make(map[string]interface{}, 0)
		}

		hostInfoMap[int64(index)][common.BKHostInnerIPField] = hostInfo[common.BKHostInnerIPField]
		hostInfoMap[int64(index)][common.BKHostOuterIPField] = hostInfo[common.BKHostOuterIPField]
		hostInfoMap[int64(index)][common.BKOSNameField] = hostInfo[common.BKOSNameField]
		hostInfoMap[int64(index)][common.BKImportFrom] = "3"
		hostInfoMap[int64(index)][common.BKCloudIDField] = 1
	}

	hostIDs, succ, updateErrRow, errRow, ok := lgc.AddHost(ctx, appID, []int64{moduleID}, util.GetOwnerID(lgc.header), hostInfoMap, hostList.InputType)
	if ok != nil {
		blog.Errorf("add host failed, hostIDs: %+v, succ: %v, update: %v, err: %v, %v, rid: %s", hostIDs, succ, updateErrRow, ok, errRow, lgc.rid)
		return ok
	}

	return nil
}

func (lgc *Logics) UpdateCloudHosts(ctx context.Context, cloudHostAttr []mapstr.MapStr) error {
	for _, hostInfo := range cloudHostAttr {
		hostID, err := hostInfo.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("hostID convert to string failed, hostInfo: %#v, err: %v, rid: %s", hostInfo, err, lgc.rid)
			return err
		}

		delete(hostInfo, common.BKHostIDField)
		delete(hostInfo, common.BKCloudConfirm)
		delete(hostInfo, common.BKAttrConfirm)

		updateParam := &meta.UpdateOption{
			Data:      hostInfo,
			Condition: mapstr.MapStr{common.BKHostIDField: hostID},
		}
		result, err := lgc.CoreAPI.CoreService().Instance().UpdateInstance(ctx, lgc.header, common.BKInnerObjIDHost, updateParam)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("update host batch failed, ids[%v], err: %v, %v, rid: %s", hostID, err, result.ErrMsg, lgc.rid)
			return err
		}
	}
	return nil
}

func (lgc *Logics) NewAddConfirm(ctx context.Context, taskInfo meta.CloudTaskInfo, newCloudHost []mapstr.MapStr) (int, error) {
	// Check whether the host is already exist in resource confirm.
	opt := make(map[string]interface{})
	confirmHosts, err := lgc.CoreAPI.CoreService().Cloud().SearchConfirm(ctx, lgc.header, opt)
	if err != nil {
		blog.Errorf("get confirm info failed with err: %v, rid: %s", err, lgc.rid)
		return 0, err
	}

	confirmIpList := make([]string, 0)
	if confirmHosts.Count > 0 {
		for _, confirmInfo := range confirmHosts.Info {
			ip, ok := confirmInfo[common.BKHostInnerIPField].(string)
			if !ok {
				continue
			}
			confirmIpList = append(confirmIpList, ip)
		}
	}

	newHostIp := make([]string, 0)
	for _, host := range newCloudHost {
		innerIp, err := host.String(common.BKHostInnerIPField)
		if err != nil {
			blog.Errorf("mapstr.Map convert to string failed, err: %v, rid: %s", err, lgc.rid)
			return 0, err
		}
		if !util.InStrArr(confirmIpList, innerIp) {
			newHostIp = append(newHostIp, innerIp)
		}
	}

	// newly added cloud hosts confirm
	if len(newHostIp) > 0 {
		for _, host := range newCloudHost {
			innerIp, err := host.String(common.BKHostInnerIPField)
			if err != nil {
				blog.Errorf("mapstr.Map convert to string failed, err: %v, rid: %s", err, lgc.rid)
				return 0, err
			}
			outerIp, err := host.String(common.BKHostOuterIPField)
			if err != nil {
				blog.Error("mapstr.Map convert to string failed, err: %v, rid: %s", err, lgc.rid)
				return 0, err
			}
			osName, err := host.String(common.BKOSNameField)
			if err != nil {
				blog.Error("mapstr.Map convert to string failed, err: %v, rid: %s", err, lgc.rid)
				return 0, err
			}
			resourceConfirm := mapstr.MapStr{}
			resourceConfirm[common.BKObjIDField] = taskInfo.ObjID
			resourceConfirm[common.BKHostInnerIPField] = innerIp
			resourceConfirm[common.BKCloudTaskID] = taskInfo.TaskID
			resourceConfirm[common.BKOSNameField] = osName
			resourceConfirm[common.BKHostOuterIPField] = outerIp
			resourceConfirm[common.BKCloudConfirm] = true
			resourceConfirm[common.BKAttrConfirm] = false
			resourceConfirm[common.BKCloudSyncTaskName] = taskInfo.TaskName
			resourceConfirm[common.BKCloudAccountType] = taskInfo.AccountType
			resourceConfirm[common.BKCloudSyncAccountAdmin] = taskInfo.AccountAdmin
			resourceConfirm[common.BKResourceType] = common.BKNewAddHost

			if _, err := lgc.CoreAPI.CoreService().Cloud().CreateConfirm(ctx, lgc.header, resourceConfirm); err != nil {
				blog.Errorf("add resource confirm failed with err: confirmInfo: %#v, %v, rid: %s", resourceConfirm, err, lgc.rid)
				return 0, err
			}
		}
	}
	num := len(newHostIp)
	return num, nil
}

func (lgc *Logics) NextTrigger(ctx context.Context, periodType string, period string) int64 {
	toBeCharge := period
	var unixSubtract int64
	nowStr := time.Unix(time.Now().Unix(), 0).Format(common.TimeTransferModel)

	if periodType == "day" {
		intHour, _ := strconv.Atoi(toBeCharge[:2])
		intMinute, _ := strconv.Atoi(toBeCharge[3:])
		if intHour > time.Now().Hour() {
			toBeCharge = fmt.Sprintf("%s%s%s", nowStr[:11], toBeCharge, ":00")
		}
		if intHour < time.Now().Hour() {
			toBeCharge = fmt.Sprintf("%s%d %s%s", nowStr[:8], time.Now().Day()+1, toBeCharge, ":00")
		}
		if intHour == time.Now().Hour() && intMinute > time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%s%s", nowStr[:11], toBeCharge, ":00")
		}
		if intHour == time.Now().Hour() && intMinute <= time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d %s%s", nowStr[:8], time.Now().Day()+1, toBeCharge, ":00")
		}

		loc, _ := time.LoadLocation("Local")
		theTime, _ := time.ParseInLocation(common.TimeTransferModel, toBeCharge, loc)
		sr := theTime.Unix()
		unixSubtract = sr - time.Now().Unix()
	}

	if periodType == "hour" {
		intToBeCharge, err := strconv.Atoi(toBeCharge)
		if err != nil {
			blog.Errorf("period transfer to int failed with err: %v, rid: %s", err, lgc.rid)
			return 0
		}

		if intToBeCharge >= 10 && intToBeCharge > time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:%s:%s", nowStr[:11], time.Now().Hour(), toBeCharge, "00")
		}
		if intToBeCharge >= 10 && intToBeCharge < time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:%s:%s", nowStr[:11], time.Now().Hour()+1, toBeCharge, "00")
		}
		if intToBeCharge < 10 && intToBeCharge > time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:0%s:%s", nowStr[:11], time.Now().Hour(), toBeCharge, "00")
		}
		if intToBeCharge < 10 && intToBeCharge < time.Now().Minute() {
			toBeCharge = fmt.Sprintf("%s%d:0%s:%s", nowStr[:11], time.Now().Hour()+1, toBeCharge, "00")
		}

		loc, _ := time.LoadLocation("Local")
		theTime, _ := time.ParseInLocation(common.TimeTransferModel, toBeCharge, loc)
		sr := theTime.Unix()
		unixSubtract = sr - time.Now().Unix()
	}

	if periodType == "minute" {
		unixSubtract = 300
	}

	minuteNextTrigger := unixSubtract / 60
	return minuteNextTrigger
}

func (lgc *Logics) CloudSyncHistory(ctx context.Context, taskID int64, startTime int64, cloudHistory *meta.CloudHistory) {
	finishTime := time.Now().Unix()
	timeConsumed := finishTime - startTime
	if timeConsumed > 60 {
		minute := timeConsumed / 60
		seconds := timeConsumed % 60
		cloudHistory.TimeConsume = fmt.Sprintf("%dmin%ds", minute, seconds)
	} else {
		cloudHistory.TimeConsume = fmt.Sprintf("%ds", timeConsumed)
	}

	startTimeStr := time.Unix(startTime, 0).Format(common.TimeTransferModel)
	cloudHistory.StartTime = startTimeStr

	blog.V(3).Info("cloudHistory.TimeConsume: %+v, rid: %s", cloudHistory.TimeConsume, lgc.rid)

	updateData := mapstr.MapStr{}
	updateTime := time.Now()
	updateData[common.BKLastTimeCloudSync] = updateTime
	updateData[common.BKCloudTaskID] = taskID
	updateData[common.BKSyncStatus] = cloudHistory.Status
	updateData[common.BKNewAddHost] = cloudHistory.NewAdd
	updateData[common.BKAttrChangedHost] = cloudHistory.AttrChanged

	if _, err := lgc.CoreAPI.CoreService().Cloud().UpdateCloudSyncTask(ctx, lgc.header, updateData); err != nil {
		blog.Errorf("update task failed, taskInfo: %#v, err: %v, rid: %s", updateData, err, lgc.rid)
		return
	}

	if _, err := lgc.CoreAPI.CoreService().Cloud().CreateSyncHistory(ctx, lgc.header, cloudHistory); err != nil {
		blog.Errorf("add cloud history table failed, history: %v, err: %v, rid: %s", cloudHistory, err, lgc.rid)
		return
	}

	return
}

func (lgc *Logics) ObtainCloudHosts(ctx context.Context, secretID string, secretKey string) ([]map[string]interface{}, error) {
	credential := com.NewCredential(
		secretID,
		secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = common.BKHttpGet
	cpf.HttpProfile.ReqTimeout = common.BKTencentCloudTimeOut
	cpf.HttpProfile.Endpoint = common.TencentCloudUrl
	cpf.SignMethod = common.TencentCloudSignMethod

	ClientRegion, _ := cvm.NewClient(credential, regions.Guangzhou, cpf)
	regionRequest := cvm.NewDescribeRegionsRequest()
	Response, err := ClientRegion.DescribeRegions(regionRequest)

	if err != nil {
		return nil, err
	}

	data := Response.ToJsonString()
	regionResponse := new(meta.RegionResponse)
	if err := json.Unmarshal([]byte(data), regionResponse); err != nil {
		blog.Errorf("json unmarsha1 error :%v, rid: %v", err, lgc.rid)
		return nil, err
	}

	cloudHostInfo := make([]map[string]interface{}, 0)
	for _, region := range regionResponse.Response.Data {
		var inneripList string
		var outeripList string
		var osName string
		regionHosts := make(map[string]interface{})

		client, _ := cvm.NewClient(credential, region.Region, cpf)
		instRequest := cvm.NewDescribeInstancesRequest()
		response, err := client.DescribeInstances(instRequest)

		if _, ok := err.(*cErrors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s, rid: %v", err, lgc.rid)
			return nil, err
		}
		if err != nil {
			blog.Error("obtain cloud hosts failed, err: %v, rid: %v", err, lgc.rid)
			return nil, err
		}

		data := response.ToJsonString()
		Hosts := meta.HostResponse{}
		if err := json.Unmarshal([]byte(data), &Hosts); err != nil {
			fmt.Printf("json unmarsha1 error :%v\n, rid: %v", err, lgc.rid)
		}

		instSet := Hosts.HostResponse.InstanceSet
		for _, obj := range instSet {
			osName = obj.OsName
			if len(obj.PrivateIpAddresses) > 0 {
				inneripList = obj.PrivateIpAddresses[0]
			}
		}

		for _, obj := range instSet {
			if len(obj.PublicIpAddresses) > 0 {
				outeripList = obj.PublicIpAddresses[0]
			}
		}

		if len(instSet) > 0 {
			regionHosts[common.BKHostCloudRegionField] = region.Region
			regionHosts[common.BKHostInnerIPField] = inneripList
			regionHosts[common.BKHostOuterIPField] = outeripList
			regionHosts[common.BKOSNameField] = osName
			cloudHostInfo = append(cloudHostInfo, regionHosts)
		}
	}
	return cloudHostInfo, nil
}

func copyHeader(ctx context.Context, header http.Header) http.Header {
	newHeader := make(http.Header, 0)
	for key, values := range header {
		for _, v := range values {
			newHeader.Add(key, v)
		}
	}

	return newHeader
}

func (lgc *Logics) deleteStartedTaskRedis(ctx context.Context, taskID int64) {
	startedTaskItems, err := lgc.cache.LRange(common.RedisCloudSyncInstanceStarted, 0, -1).Result()
	if err != nil {
		blog.Errorf("get task item from redis fail, error: %v, rid: %s", err, lgc.rid)
	}

	if len(startedTaskItems) == 0 {
		return
	}

	for _, item := range startedTaskItems {
		startedItem := meta.CloudSyncRedisAlreadyStarted{}
		if err := json.Unmarshal([]byte(item), &startedItem); err != nil {
			blog.Warnf("get task started item from redis fail, taskInfo: %s, error:%s, rid: %s", item, err.Error(), lgc.rid)
			continue
		}
		if taskID == startedItem.TaskID {
			if err := lgc.cache.LRem(common.RedisCloudSyncInstanceStarted, 1, item).Err(); err != nil {
				blog.Errorf("remove stop task item fail, taskInfo: %s, error: %v, rid: %s", item, err, lgc.rid)
			}
		}
	}
	return
}
