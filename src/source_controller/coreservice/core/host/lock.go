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

package host

import (
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (hm *hostManager) LockHost(params core.ContextParams, input *metadata.HostLockRequest) errors.CCError {
	fields := []string{common.BKHostIDField, common.BKHostInnerIPField}
	condition := mapstr.MapStr{
		common.BKCloudIDField:     input.CloudID,
		common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: input.IPS},
	}
	condition = util.SetQueryOwner(condition, params.SupplierAccount)
	hostInfos := make([]mapstr.MapStr, 0)
	limit := uint64(len(input.IPS))
	err := hm.DbProxy.Table(common.BKTableNameBaseHost).Find(condition).Fields(fields...).Limit(limit).All(params.Context, &hostInfos)
	if nil != err {
		blog.Errorf("lock host, query host from db error, condition: %+v, err: %+v, rid: %s", condition, err, params.ReqID)
		return params.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	diffIP := diffHostLockIP(input.IPS, hostInfos, params.ReqID)
	if 0 != len(diffIP) {
		blog.Errorf("lock host, not found, ip:%+v, rid:%s", diffIP, params.ReqID)
		return params.Error.Errorf(common.CCErrCommParamsIsInvalid, " ip_list["+strings.Join(diffIP, ",")+"]")
	}

	user := util.GetUser(params.Header)
	var insertDataArr []interface{}
	ts := time.Now().UTC()
	for _, ip := range input.IPS {
		conds := mapstr.MapStr{
			common.BKHostInnerIPField: ip,
			common.BKCloudIDField:     input.CloudID,
		}
		conds = util.SetQueryOwner(conds, params.SupplierAccount)
		cnt, err := hm.DbProxy.Table(common.BKTableNameHostLock).Find(conds).Count(params.Context)
		if nil != err {
			blog.Errorf("lock host, query host lock from db failed, err:%+v, rid:%s", err, params.ReqID)
			return params.Error.Errorf(common.CCErrCommDBSelectFailed)
		}
		if 0 == cnt {
			insertDataArr = append(insertDataArr, metadata.HostLockData{
				User:       user,
				IP:         ip,
				CloudID:    input.CloudID,
				CreateTime: ts,
				OwnerID:    util.GetOwnerID(params.Header),
			})
		}
	}

	if 0 < len(insertDataArr) {
		err := hm.DbProxy.Table(common.BKTableNameHostLock).Insert(params.Context, insertDataArr)
		if nil != err {
			blog.Errorf("lock host, save host lock to db failed, err: %+v, rid:%s", err, params.ReqID)
			return params.Error.Errorf(common.CCErrCommDBInsertFailed)
		}
	}
	return nil
}

func (hm *hostManager) UnlockHost(params core.ContextParams, input *metadata.HostLockRequest) errors.CCError {
	conds := mapstr.MapStr{
		common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: input.IPS},
		common.BKCloudIDField:     input.CloudID,
	}
	conds = util.SetModOwner(conds, params.SupplierAccount)
	err := hm.DbProxy.Table(common.BKTableNameHostLock).Delete(params.Context, conds)
	if nil != err {
		blog.Errorf("unlock host, delete host lock from db error, err: %+v, rid:%s", err, params.ReqID)
		return params.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

func (hm *hostManager) QueryHostLock(params core.ContextParams, input *metadata.QueryHostLockRequest) ([]metadata.HostLockData, errors.CCError) {
	hostLockInfoArr := make([]metadata.HostLockData, 0)
	conds := mapstr.MapStr{
		common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: input.IPS},
		common.BKCloudIDField:     input.CloudID,
	}
	conds = util.SetModOwner(conds, params.SupplierAccount)
	limit := uint64(len(input.IPS))
	err := hm.DbProxy.Table(common.BKTableNameHostLock).Find(conds).Limit(limit).All(params.Context, &hostLockInfoArr)
	if nil != err {
		blog.Errorf("query lock host, query host lock from db error, err: %+v, rid:%s", err, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	return hostLockInfoArr, nil
}

func diffHostLockIP(ips []string, hostInfos []mapstr.MapStr, rid string) []string {
	mapInnerIP := make(map[string]bool, 0)
	for _, hostInfo := range hostInfos {
		innerIP, err := hostInfo.String(common.BKHostInnerIPField)
		if nil != err {
			blog.ErrorJSON("different host lock IP not inner ip, %s, rid: %s", hostInfo, rid)
			continue
		}
		mapInnerIP[innerIP] = true
	}
	var diffIPS []string
	for _, ip := range ips {
		_, exist := mapInnerIP[ip]
		if !exist {
			diffIPS = append(diffIPS, ip)
		}
	}
	return diffIPS
}
