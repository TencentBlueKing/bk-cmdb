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
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

func (lgc *Logics) RefreshHostInstanceByApp(ctx context.Context, header http.Header, appID int64, appInfo mapstr.MapStr) error {
	ownerID, err := appInfo.String(common.BKOwnerIDField)
	if nil != err {
		blog.Errorf("RefreshHostInstanceByApp error  appID:%d, appInfo:%+v, error:%s", appID, appInfo, err.Error())
		return err
	}
	if nil == header {
		header = make(http.Header, 0)
	}
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, ownerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}

	moduleIDs, err := lgc.GetModueleIDByAppID(ctx, header, appID)
	if nil != err {
		blog.Errorf("RefreshHostInstanceByApp error   appID:%d, appInfo:%+v,  error:%s", appID, appInfo, err.Error())
		return err
	}
	return lgc.addEventRefreshModuleItems(appID, moduleIDs, header)
}

func (lgc *Logics) RefreshAllHostInstance(ctx context.Context, header http.Header) error {
	if nil == header {
		header = make(http.Header, 0)
	}
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}
	fields := fmt.Sprintf("%s,%s", common.BKAppIDField, common.BKOwnerIDField)
	appInfoArr, err := lgc.GetAppList(ctx, header, fields)
	if nil != err {
		blog.Errorf("RefreshAllHostInstance error:%s", err.Error())
		return err
	}
	for _, appInfo := range appInfoArr {
		appID, err := appInfo.Int64(common.BKAppIDField)
		if nil != err {
			blog.Warnf("RefreshAllHostInstance get appID by app Info:%+v error:%s", appInfo, err.Error())
			continue
		}
		ownerID, err := appInfo.String(common.BKOwnerIDField)
		if nil != err {
			blog.Warnf("RefreshAllHostInstance get supplier accout by app Info:%+v error:%s", appInfo, err.Error())
			continue
		}
		header.Set(common.BKHTTPOwnerID, ownerID)
		err = lgc.RefreshHostInstanceByApp(ctx, header, appID, appInfo)
		if nil != err {
			blog.Warnf("RefreshAllHostInstance RefreshHostInstanceByApp by app Info:%+v error:%s", appInfo, err.Error())
			continue
		}
	}
	return nil
}

func (lgc *Logics) timedTriggerRefreshHostInstance() {
	go func() {
		triggerChn := time.NewTicker(timedTriggerTime)
		for range triggerChn.C {
			lgc.cache.Del(common.RedisProcSrvHostInstanceAllRefreshLockKey)
			locked, err := lgc.cache.SetNX(common.RedisProcSrvHostInstanceAllRefreshLockKey, "", timedTriggerLockExpire).Result()
			if nil != err {
				blog.Errorf("locked refresh  error:%s", err.Error())
				continue
			}
			if locked {
				err := lgc.RefreshAllHostInstance(context.Background(), nil)
				if nil != err {
					blog.Errorf("RefreshAllHostInstance error:%s", err.Error())
					continue
				}
			}
		}
	}()

}
