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
	"bytes"
	"fmt"
	"strconv"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/web_server/app/options"

	"github.com/gin-gonic/gin"
)

// GetDepartment get department info from paas
func (lgc *Logics) GetDepartment(c *gin.Context, config *options.Config) (*metadata.DepartmentData,
	errors.CCErrorCoder) {
	if config.LoginVersion == common.BKOpenSourceLoginPluginVersion ||
		config.LoginVersion == common.BKSkipLoginPluginVersion {
		return &metadata.DepartmentData{}, nil
	}

	header := c.Request.Header
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(header))
	rid := httpheader.GetRid(header)

	result, esbErr := esb.EsbClient().User().GetDepartment(c.Request.Context(), c.Request.Header, c.Request.URL)
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
func (lgc *Logics) GetDepartmentProfile(c *gin.Context, config *options.Config) (*metadata.DepartmentProfileData,
	errors.CCErrorCoder) {
	if config.LoginVersion == common.BKOpenSourceLoginPluginVersion ||
		config.LoginVersion == common.BKSkipLoginPluginVersion {
		return &metadata.DepartmentProfileData{}, nil
	}

	header := c.Request.Header
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(header))
	rid := httpheader.GetRid(header)

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

// GetAllDepartment get department info from paas
func (lgc *Logics) GetAllDepartment(c *gin.Context, config *options.Config, orgIDs []int64) (*metadata.DepartmentData,
	errors.CCErrorCoder) {
	if config.LoginVersion == common.BKOpenSourceLoginPluginVersion ||
		config.LoginVersion == common.BKSkipLoginPluginVersion {
		return &metadata.DepartmentData{}, nil
	}

	defErr := lgc.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(c.Request.Header))
	rid := httpheader.GetRid(c.Request.Header)

	orgIDList := lgc.getOrgListStr(orgIDs)
	departments := &metadata.DepartmentData{}
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr error
	pipeline := make(chan bool, 10)

	for _, subStr := range orgIDList {
		pipeline <- true
		wg.Add(1)
		go func(subStr string) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			params := make(map[string]string)
			params["exact_lookups"] = subStr
			result, esbErr := esb.EsbClient().User().GetAllDepartment(c.Request.Context(), c.Request.Header, params)
			if esbErr != nil {
				firstErr = esbErr
				blog.Errorf("get department by esb client failed, http failed, params: %+v, err: %v, rid: %s",
					params, esbErr, rid)
				return
			}
			if !result.Result {
				blog.Errorf("get department by esb client failed, result is false, params: %+v, rid: %s", params,
					rid)
				firstErr = fmt.Errorf("get department by esb failed, result is false, params: %v", params)
				return
			}

			lock.Lock()
			departments.Count += result.Data.Count
			departments.Results = append(departments.Results, result.Data.Results...)
			lock.Unlock()
		}(subStr)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return departments, nil
}

const getOrganizationMaxLength = 500

// getOrgListStr get org list str
func (lgc *Logics) getOrgListStr(orgIDList []int64) []string {
	orgListStr := make([]string, 0)

	orgBuffer := bytes.Buffer{}
	for _, orgID := range orgIDList {
		if orgBuffer.Len()+len(strconv.FormatInt(orgID, 10)) > getOrganizationMaxLength {
			orgBuffer.WriteString(strconv.FormatInt(orgID, 10))
			orgStr := orgBuffer.String()
			orgListStr = append(orgListStr, orgStr)
			orgBuffer.Reset()
			continue
		}

		orgBuffer.WriteString(strconv.FormatInt(orgID, 10))
		orgBuffer.WriteByte(',')
	}

	if orgBuffer.Len() == 0 {
		return []string{}
	}

	orgStr := orgBuffer.String()
	orgListStr = append(orgListStr, orgStr[:len(orgStr)-1])

	return orgListStr
}
