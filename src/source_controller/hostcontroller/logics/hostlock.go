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
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) LockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) errors.CCError {

	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	user := util.GetUser(header)

	fields := []string{common.BKHostIDField, common.BKHostInnerIPField}
	condition := mapstr.MapStr{common.BKCloudIDField: input.CloudID, common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: input.IPS}}
	hostInfos := make([]mapstr.MapStr, 0)
	err := lgc.Instance.Table(common.BKTableNameBaseHost).Find(condition).Fields(fields...).Limit(uint64(len(input.IPS))).All(ctx, &hostInfos)
	if nil != err {
		blog.Errorf("lcok host, query host from db error, error:%s ,logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return defErr.Errorf(common.CCErrCommDBSelectFailed)
	}

	diffIP := diffHostLockIP(input.IPS, hostInfos)
	if 0 != len(diffIP) {
		blog.Errorf("lock host, not found ip:%+v,logID:%s", diffIP, util.GetHTTPCCRequestID(header))
		return defErr.Errorf(common.CCErrCommParamsIsInvalid, " ip_list["+strings.Join(diffIP, ",")+"]")
	}

	var insertDataArr []interface{}
	ts := time.Now().UTC()
	for _, ip := range input.IPS {
		conds := mapstr.MapStr{common.BKHostInnerIPField: ip, common.BKCloudIDField: input.CloudID}
		cnt, err := lgc.Instance.Table(common.BKTableNameHostLock).Find(conds).Count(ctx)
		if nil != err {
			blog.Errorf("lcok host, query host lock from db error, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
			return defErr.Errorf(common.CCErrCommDBSelectFailed)
		}
		if 0 == cnt {
			insertDataArr = append(insertDataArr, metadata.HostLockData{
				User:       user,
				IP:         ip,
				CloudID:    input.CloudID,
				CreateTime: ts,
			})
		}
	}

	if 0 < len(insertDataArr) {
		err := lgc.Instance.Table(common.BKTableNameHostLock).Insert(ctx, insertDataArr)
		if nil != err {
			blog.Errorf("lcok host, save host lock to db error, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
			return defErr.Errorf(common.CCErrCommDBInsertFailed)
		}
	}
	return nil
}

func (lgc *Logics) UnlockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) errors.CCError {

	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	conds := mapstr.MapStr{common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: input.IPS}, common.BKCloudIDField: input.CloudID}
	err := lgc.Instance.Table(common.BKTableNameHostLock).Delete(ctx, conds)

	if nil != err {
		blog.Errorf("unlock host, delete host lock from db error, error:%s,logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return defErr.Errorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

func (lgc *Logics) QueryHostLock(ctx context.Context, header http.Header, input *metadata.QueryHostLockRequest) ([]metadata.HostLockData, errors.CCError) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	hostLockInfoArr := make([]metadata.HostLockData, 0)
	conds := mapstr.MapStr{common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: input.IPS}, common.BKCloudIDField: input.CloudID}
	err := lgc.Instance.Table(common.BKTableNameHostLock).Find(conds).Limit(uint64(len(input.IPS))).All(ctx, &hostLockInfoArr)
	if nil != err {
		blog.Errorf("query lcok host, query host lock from db error, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return nil, defErr.Errorf(common.CCErrCommDBSelectFailed)
	}
	return hostLockInfoArr, err
}

func diffHostLockIP(ips []string, hostInfos []mapstr.MapStr) []string {
	mapInnerIP := make(map[string]bool, 0)
	for _, hostInfo := range hostInfos {
		innerIP, err := hostInfo.String(common.BKHostInnerIPField)
		if nil != err {
			blog.Warnf("different host lock IP not inner ip")
			continue
		}
		mapInnerIP[innerIP] = true
	}
	var diffIPS []string
	for _, ip := range ips {
		_, ok := mapInnerIP[ip]
		if !ok {
			diffIPS = append(diffIPS, ip)
		}
	}
	return diffIPS
}
