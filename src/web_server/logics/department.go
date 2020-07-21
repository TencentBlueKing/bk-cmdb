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
	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	commonutil "configcenter/src/common/util"
	"configcenter/src/web_server/app/options"

	"github.com/gin-gonic/gin"
)

// GetDepartment get department info from paas
func (lgc *Logics) GetDepartment(c *gin.Context, config *options.Config) (*metadata.DepartmentData, errors.CCErrorCoder) {
	// if no esb config, return
	if !configcenter.IsExist("webServer.esb.addr") {
		return &metadata.DepartmentData{}, nil
	}

	header := c.Request.Header
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(commonutil.GetLanguage(header))
	rid := commonutil.GetHTTPCCRequestID(header)

	result, esbErr := esb.EsbClient().User().GetDepartment(c.Request.Context(), c.Request)
	if esbErr != nil {
		blog.Errorf("get department by esb client failed, http failed, err: %+v, rid: %s", esbErr, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get department by esb client failed, result is false, err: %+v, rid: %s", result, rid)
		return nil, errors.NewCCError(result.Code, result.Message)
	}
	return &result.Data, nil

}

// GetDepartmentProfile get department profile from paas
func (lgc *Logics) GetDepartmentProfile(c *gin.Context, config *options.Config) (*metadata.DepartmentProfileData, errors.CCErrorCoder) {
	// if no esb config, return
	if !configcenter.IsExist("webServer.esb.addr") {
		return &metadata.DepartmentProfileData{}, nil
	}

	header := c.Request.Header
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(commonutil.GetLanguage(header))
	rid := commonutil.GetHTTPCCRequestID(header)

	result, esbErr := esb.EsbClient().User().GetDepartmentProfile(c.Request.Context(), c.Request)
	if esbErr != nil {
		blog.Errorf("get department by esb client failed, http failed, err: %+v, rid: %s", esbErr, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get department by esb client failed, result is false, err: %+v, rid: %s", result, rid)
		return nil, errors.NewCCError(result.Code, result.Message)
	}
	return &result.Data, nil
}
