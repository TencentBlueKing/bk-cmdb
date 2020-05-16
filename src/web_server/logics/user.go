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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	commonutil "configcenter/src/common/util"
	"configcenter/src/web_server/app/options"
)

func (lgc *Logics) GetUserList(ctx context.Context, header http.Header, params map[string]string, config *options.Config) ([]*metadata.LoginSystemUserInfo, errors.CCErrorCoder) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(commonutil.GetLanguage(header))

	switch config.LoginVersion {
	case common.BKBluekingLoginPluginVersion:
		return getBluekingLoginUserList(ctx, header, params)
	case common.BKOpenSourceLoginPluginVersion:
		return getOpenSourceLoginUserList(header, config, defErr)
	case common.BKSkipLoginPluginVersion:
		return getSkipLoginUserList()
	default:
		blog.Errorf("Unknown login version:%s, rid:%s", config.LoginVersion, commonutil.GetHTTPCCRequestID(header))
		return nil, defErr.CCErrorf(common.CCErrWebUnknownLoginVersion, config.LoginVersion)
	}

}

// getBluekingLoginUserList is used in the blueking login version
func getBluekingLoginUserList(ctx context.Context, header http.Header, params map[string]string) (
	[]*metadata.LoginSystemUserInfo, errors.CCErrorCoder) {
	rid := commonutil.GetHTTPCCRequestID(header)
	users := make([]*metadata.LoginSystemUserInfo, 0)
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

	for _, userInfo := range result.Data {
		user := &metadata.LoginSystemUserInfo{
			CnName: userInfo.DisplayName,
			EnName: userInfo.Username,
		}
		users = append(users, user)
	}

	return users, nil
}

// getOpenSourceLoginUserList is used in the open source login version
func getOpenSourceLoginUserList(header http.Header, config *options.Config, defErr errors.DefaultCCErrorIf) ([]*metadata.LoginSystemUserInfo, errors.CCErrorCoder) {
	rid := commonutil.GetHTTPCCRequestID(header)
	users := make([]*metadata.LoginSystemUserInfo, 0)
	if len(config.ConfigMap["session.user_info"]) == 0 {
		blog.Errorf("User name and password can't be found at session.user_info in config file common.conf, rid:%s", rid)
		return nil, defErr.CCError(common.CCErrWebNoUsernamePasswd)
	}
	userInfos := strings.Split(config.ConfigMap["session.user_info"], ",")
	for _, userInfo := range userInfos {
		userPasswd := strings.Split(userInfo, ":")
		if len(userPasswd) != 2 {
			blog.Errorf("The format of user name and password are wrong, please check session.user_info in config file common.conf, rid:%s", rid)
			return nil, defErr.CCError(common.CCErrWebUserinfoFormatWrong)
		}
		user := &metadata.LoginSystemUserInfo{
			CnName: userPasswd[0],
			EnName: userPasswd[0],
		}
		users = append(users, user)
	}

	return users, nil
}

// getSkipLoginUserList is used in the open source login version
func getSkipLoginUserList() ([]*metadata.LoginSystemUserInfo, errors.CCErrorCoder) {
	return []*metadata.LoginSystemUserInfo{
		{
			CnName: "admin",
			EnName: "admin",
		},
	}, nil
}
