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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) LockHost(ctx context.Context, input *metadata.HostLockRequest) errors.CCError {

	hostLockResult, err := lgc.CoreAPI.CoreService().Host().LockHost(ctx, lgc.header, input)
	if nil != err {
		blog.Errorf("lock host, http request error, error:%s,input:%+v,logID:%s", err.Error(), input, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostLockResult.Result {
		blog.Errorf("lock host, add host lock  error, error code:%d error message:%s,input:%+v,logID:%s", hostLockResult.Code, hostLockResult.ErrMsg, input, lgc.rid)
		return lgc.ccErr.New(hostLockResult.Code, hostLockResult.ErrMsg)
	}
	return nil
}

func (lgc *Logics) UnlockHost(ctx context.Context, input *metadata.HostLockRequest) errors.CCError {

	hostUnlockResult, err := lgc.CoreAPI.CoreService().Host().UnlockHost(ctx, lgc.header, input)
	if nil != err {
		blog.Errorf("unlock host, http request error, error:%s,input:%+v,logID:%s", err.Error(), input, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostUnlockResult.Result {
		blog.Errorf("unlock host, release host lock  error, error code:%d error message:%s,input:%+v,logID:%s", hostUnlockResult.Code, hostUnlockResult.ErrMsg, input, lgc.rid)
		return lgc.ccErr.New(hostUnlockResult.Code, hostUnlockResult.ErrMsg)
	}
	return nil
}

func (lgc *Logics) QueryHostLock(ctx context.Context, input *metadata.QueryHostLockRequest) (map[string]bool, errors.CCError) {

	hostLockResult, err := lgc.CoreAPI.CoreService().Host().QueryHostLock(ctx, lgc.header, input)
	if nil != err {
		blog.Errorf("query lock host, http request error, error:%s,input:%+v,logID:%s", err.Error(), input, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostLockResult.Result {
		blog.Errorf("query host lock  error, error code:%d error message:%s,input:%+v,logID:%s", hostLockResult.Code, hostLockResult.ErrMsg, input, lgc.rid)
		return nil, lgc.ccErr.New(hostLockResult.Code, hostLockResult.ErrMsg)
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
