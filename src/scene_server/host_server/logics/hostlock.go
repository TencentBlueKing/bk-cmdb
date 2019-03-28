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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) LockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) errors.CCError {

	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	rid := util.GetHTTPCCRequestID(header)
	hostLockResult, err := lgc.CoreAPI.HostController().Host().LockHost(ctx, header, input)
	if nil != err {
		blog.Errorf("lock host, http request error, error:%s,logID:%s", err.Error(), rid)
		return defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostLockResult.Result {
		blog.Errorf("lock host, add host lock  error, error code:%d error message:%s,logID:%s", hostLockResult.Code, hostLockResult.ErrMsg, rid)
		return defErr.New(hostLockResult.Code, hostLockResult.ErrMsg)
	}
	return nil
}

func (lgc *Logics) UnlockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) errors.CCError {

	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	hostUnlockResult, err := lgc.CoreAPI.HostController().Host().UnlockHost(ctx, header, input)
	if nil != err {
		blog.Errorf("unlock host, http request error, error:%s,logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostUnlockResult.Result {
		blog.Errorf("unlock host, release host lock  error, error code:%d error message:%s,logID:%s", hostUnlockResult.Code, hostUnlockResult.ErrMsg, util.GetHTTPCCRequestID(header))
		return defErr.New(hostUnlockResult.Code, hostUnlockResult.ErrMsg)
	}
	return nil
}

func (lgc *Logics) QueryHostLock(ctx context.Context, header http.Header, input *metadata.QueryHostLockRequest) (map[string]bool, errors.CCError) {

	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	hostLockResult, err := lgc.CoreAPI.HostController().Host().QueryHostLock(ctx, header, input)
	if nil != err {
		blog.Errorf("query lock host, http request error, error:%s,logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return nil, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostLockResult.Result {
		blog.Errorf("unlock host, query host lock  error, error code:%d error message:%s,logID:%s", hostLockResult.Code, hostLockResult.ErrMsg, util.GetHTTPCCRequestID(header))
		return nil, defErr.New(hostLockResult.Code, hostLockResult.ErrMsg)
	}
	hostLockMap := make(map[string]bool, 0)
	for _, ip := range input.IPS {
		hostLockMap[ip] = false
	}
	for _, hostLock := range hostLockResult.Data.Info {
		hostLockMap[hostLock.IP] = true
	}

	return hostLockMap, nil
}

func diffHostLockIP(ips []string, hostInfos []map[string]interface{}) []string {
	mapInnerIP := make(map[string]bool, 0)
	for _, hostInfo := range hostInfos {
		innerIP, ok := hostInfo[common.BKHostInnerIPField].(string)
		if ok {
			mapInnerIP[innerIP] = true
		}
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
