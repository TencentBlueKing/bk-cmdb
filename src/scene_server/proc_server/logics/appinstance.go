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
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

func copyHeader(header http.Header) http.Header {
	newHeader := make(http.Header, 0)
	for key, values := range header {
		for _, v := range values {
			newHeader.Add(key, v)
		}
	}

	return newHeader
}

func (lgc *Logics) RefreshHostInstanceByApp(ctx context.Context, appID int64, appInfo mapstr.MapStr) error {
	ownerID, err := appInfo.String(common.BKOwnerIDField)
	if nil != err {
		blog.Errorf("RefreshHostInstanceByApp error  appID:%d, appInfo:%+v, error:%s", appID, appInfo, err.Error())
		return err
	}
	header := copyHeader(lgc.header)
	if nil == header {
		header = make(http.Header, 0)
	}
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, ownerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}
	// use new header, so, new logics struct
	newLgc := lgc.NewFromHeader(header)

	moduleIDs, err := newLgc.GetModueleIDByAppID(ctx, appID)
	if nil != err {
		blog.Errorf("RefreshHostInstanceByApp error   appID:%d, appInfo:%+v,  error:%s,rid:%s", appID, appInfo, err.Error(), newLgc.rid)
		return err
	}
	return newLgc.addEventRefreshModuleItems(ctx, appID, moduleIDs)
}

func (lgc *Logics) RefreshAllHostInstance(ctx context.Context) error {
	header := copyHeader(lgc.header)
	if nil == header {
		header = make(http.Header, 0)
	}
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}
	newLgc := lgc.NewFromHeader(header)
	fields := []string{common.BKAppIDField, common.BKOwnerIDField}
	appInfoArr, err := newLgc.GetAppList(ctx, fields)
	if nil != err {
		blog.Errorf("RefreshAllHostInstance error:%s,rid:%s", err.Error(), newLgc.rid)
		return err
	}
	for _, appInfo := range appInfoArr {
		appID, err := appInfo.Int64(common.BKAppIDField)
		if nil != err {
			blog.Warnf("RefreshAllHostInstance get appID by app Info:%+v error:%s,rid:%s", appInfo, err.Error(), newLgc.rid)
			continue
		}
		ownerID, err := appInfo.String(common.BKOwnerIDField)
		if nil != err {
			blog.Warnf("RefreshAllHostInstance get supplier accout by app Info:%+v error:%s,rid:%s", appInfo, err.Error(), newLgc.rid)
			continue
		}
		newHeader := copyHeader(lgc.header)
		header.Set(common.BKHTTPOwnerID, ownerID)
		newLgc := lgc.NewFromHeader(newHeader)
		err = newLgc.RefreshHostInstanceByApp(ctx, appID, appInfo)
		if nil != err {
			blog.Warnf("RefreshAllHostInstance RefreshHostInstanceByApp by app Info:%+v error:%s,rid:%s", appInfo, err.Error(), newLgc.rid)
			continue
		}
	}
	return nil
}

func (lgc *Logics) timedTriggerRefreshHostInstance(ctx context.Context) {
	go func() {
		triggerChn := time.NewTicker(timedTriggerTime)
		for range triggerChn.C {
			locked, err := lgc.cache.SetNX(common.RedisProcSrvHostInstanceAllRefreshLockKey, "", timedTriggerLockExpire).Result()
			if nil != err {
				blog.Errorf("locked refresh  error:%s", err.Error())
				continue
			}
			if locked {
				err := lgc.RefreshAllHostInstance(ctx)
				if nil != err {
					blog.Errorf("RefreshAllHostInstance error:%s", err.Error())
					continue
				}
			}
		}
	}()

}
