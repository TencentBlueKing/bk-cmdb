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
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/synchronize_server/app/options"
)

func getVersion() int64 {
	return time.Now().Unix()
}

func (lgc *Logics) TriggerSynchronize(ctx context.Context, config *options.Config) {
	if config == nil {
		blog.Errorf("TriggerSynchronize not config ")
		return
	}
	if len(config.Names) == 0 {
		blog.Errorf("TriggerSynchronize not config ")
		return
	}
	lgc = lgc.NewFromHeader(copyHeader(lgc.header))
	interval, err := util.GetInt64ByInterface(config.Trigger.Role)
	if err != nil {
		blog.Warnf("Trigger.Role %v not integer, err:%s", config.Trigger.Role, err.Error())
		if config.Trigger.IsTiming() {
			interval = 10
		} else {
			// default 60minute
			interval = 60
		}
	}
	if config.Trigger.IsTiming() {
		now := time.Now().Local()
		// caculate the time difference between th next trigger time
		// (now.Hour()*60 + now.Minute()) The time that has been consumed that day
		expendTime := int64(now.Hour()*60 + now.Minute())
		if expendTime == interval {
			interval = 1
		} else {
			interval = interval + (nextDayTrigger - expendTime)
		}
		//
		if interval < 1 {
			interval = 1
		}
	}
	if lgc.Engine.ServiceManageInterface.IsMaster() {
		lgc.Synchronize(ctx, config)
	}

	blog.V(4).Infof("synchronize ready")
	timeInterval := time.Duration(interval) * time.Minute

	for {
		ticker := time.NewTimer(timeInterval)
		<-ticker.C
		if config.Trigger.IsTiming() {
			timeInterval = time.Duration(nextDayTrigger) * time.Minute
		}
		if !lgc.Engine.ServiceManageInterface.IsMaster() {
			blog.Infof(" not master ")
			continue
		}
		lgc.Synchronize(ctx, config)
	}

}

// Synchronize synchronize manager
func (lgc *Logics) Synchronize(ctx context.Context, config *options.Config) {

	for idx := range config.ConifgItemArray {
		go lgc.SynchronizeItem(ctx, config.ConifgItemArray[idx])
	}

}

// SynchronizeItem  synchronize data
func (lgc *Logics) SynchronizeItem(ctx context.Context, syncConfig *options.ConfigItem) {
	version := getVersion()

	blog.InfoJSON("start synchonrize config:%s, verison:%s", syncConfig, version)
	// syncConfig can modify
	synchronizeItem := lgc.NewSynchronizeItem(version, syncConfig)

	exceptionMap := make(map[string][]metadata.ExceptionResult)
	var err error
	exceptionMap["model"], err = synchronizeItem.synchronizeModelTask(ctx) //lgc.synchronizeModelTask(ctx, syncConfig, version, nil)
	if err != nil {
		blog.Errorf("SynchronizeItem model error, config:%#v,err:%s,version:%d,rid:%s", syncConfig, err.Error(), version, lgc.rid)
	}

	exceptionMap["instance"], err = synchronizeItem.synchronizeInstanceTask(ctx) //(ctx, syncConfig, version, nil)
	if err != nil {
		blog.Errorf("SynchronizeItem instance error, config:%#v,err:%s,version:%d,rid:%s", syncConfig, err.Error(), version, lgc.rid)
	}

	exceptionMap["association"], err = synchronizeItem.synchronizeAssociationTask(ctx) //(ctx, syncConfig, version, nil)
	if err != nil {
		blog.Errorf("SynchronizeItem association error, config:%#v,err:%s,version:%d,rid:%s", syncConfig, err.Error(), version, lgc.rid)
	}
	exceptionMapClear, err := synchronizeItem.synchronizeItemClearData(ctx)
	if err != nil {
		blog.Errorf("SynchronizeItem synchronizeItemClearData error, config:%#v,err:%s,version:%d,rid:%s", syncConfig, err.Error(), version, lgc.rid)
	}
	for key, val := range exceptionMapClear {
		exceptionMap[key] = val
	}
	go synchronizeItem.synchronizeItemException(ctx, exceptionMap)

	blog.InfoJSON("end synchonrize config:%s, verison:%s", syncConfig, version)

}
