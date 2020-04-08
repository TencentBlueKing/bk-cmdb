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
	"configcenter/src/common/resource/esb"
	commonutil "configcenter/src/common/util"
)

func (lgc *Logics) GetUserList(ctx context.Context, header http.Header, params map[string]string) ([]*metadata.LoginSystemUserInfo, errors.CCErrorCoder) {
	rid := commonutil.GetHTTPCCRequestID(header)

	// try to use esb user list api
	result, err := esb.EsbClient().User().ListUsers(ctx, header, params)
	if err != nil {
		blog.Errorf("get users by esb client failed, http failed, err: %+v, rid: %s", err, rid)
		return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !result.Result {
		blog.Errorf("request esb, get user list failed, err: %v, rid: %s", result.Message, result.EsbRequestID)
		return nil, errors.New(result.Code, result.Message)
	}

	users := make([]*metadata.LoginSystemUserInfo, 0)
	for _, userInfo := range result.Data {
		user := &metadata.LoginSystemUserInfo{
			CnName: userInfo.DisplayName,
			EnName: userInfo.Username,
		}
		users = append(users, user)
	}
	return users, nil

}
