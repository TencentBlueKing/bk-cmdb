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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) LockHost(kit *rest.Kit, input *metadata.HostLockRequest) errors.CCError {

	hostLockResult, err := lgc.CoreAPI.CoreService().Host().LockHost(kit.Ctx, kit.Header, input)
	if nil != err {
		blog.Errorf("lock host, http request error, error:%s,input:%+v,logID:%s", err.Error(), input, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostLockResult.Result {
		blog.Errorf("lock host, add host lock  error, error code:%d error message:%s,input:%+v,logID:%s", hostLockResult.Code, hostLockResult.ErrMsg, input, kit.Rid)
		return kit.CCError.New(hostLockResult.Code, hostLockResult.ErrMsg)
	}
	return nil
}

func (lgc *Logics) UnlockHost(kit *rest.Kit, input *metadata.HostLockRequest) errors.CCError {

	hostUnlockResult, err := lgc.CoreAPI.CoreService().Host().UnlockHost(kit.Ctx, kit.Header, input)
	if nil != err {
		blog.Errorf("unlock host, http request error, error:%s,input:%+v,logID:%s", err.Error(), input, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostUnlockResult.Result {
		blog.Errorf("unlock host, release host lock  error, error code:%d error message:%s,input:%+v,logID:%s", hostUnlockResult.Code, hostUnlockResult.ErrMsg, input, kit.Rid)
		return kit.CCError.New(hostUnlockResult.Code, hostUnlockResult.ErrMsg)
	}
	return nil
}

func (lgc *Logics) QueryHostLock(kit *rest.Kit, input *metadata.QueryHostLockRequest) (map[int64]bool, errors.CCError) {

	hostLockResult, err := lgc.CoreAPI.CoreService().Host().QueryHostLock(kit.Ctx, kit.Header, input)
	if nil != err {
		blog.Errorf("query lock host, http request error, error:%s,input:%+v,logID:%s", err.Error(), input, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostLockResult.Result {
		blog.Errorf("query host lock  error, error code:%d error message:%s,input:%+v,logID:%s", hostLockResult.Code, hostLockResult.ErrMsg, input, kit.Rid)
		return nil, kit.CCError.New(hostLockResult.Code, hostLockResult.ErrMsg)
	}
	hostLockMap := make(map[int64]bool, 0)
	for _, id := range input.IDS {
		hostLockMap[id] = false
	}
	for _, hostLock := range hostLockResult.Data.Info {
		hostLockMap[hostLock.ID] = true
	}

	return hostLockMap, nil
}
