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

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/hooks/process"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

// getReturnStr get return string
func getReturnStr(code int, message string, data interface{}) string {
	ret := make(map[string]interface{})
	ret["bk_error_code"] = code
	if 0 == code {
		ret["result"] = true
	} else {
		ret["result"] = false
	}
	ret["bk_error_msg"] = message
	ret["data"] = data
	msg, _ := json.Marshal(ret)

	return string(msg)

}

// ListenIPOptions TODO
func (s *Service) ListenIPOptions(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(c.Request.Header))

	hostIDStr := c.Param("bk_host_id")
	hostID, err := strconv.ParseInt(hostIDStr, 10, 64)
	if err != nil {
		blog.Infof("host id invalid, convert to int failed, hostID: %s, err: %+v, rid: %s", hostID, err, rid)
		result := metadata.BaseResp{Result: false, Code: common.CCErrCommParamsInvalid,
			ErrMsg: defErr.Errorf(common.CCErrCommParamsInvalid, common.BKHostIDField).Error()}
		c.JSON(http.StatusOK, result)
		return
	}

	option := metadata.ListHostsWithNoBizParameter{
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{
				Condition: querybuilder.ConditionAnd,
				Rules: []querybuilder.Rule{
					querybuilder.AtomRule{
						Field:    common.BKHostIDField,
						Operator: querybuilder.OperatorEqual,
						Value:    hostID,
					},
				},
			},
		},
		Fields: []string{
			common.BKHostIDField,
			common.BKHostNameField,
			common.BKHostInnerIPField,
			common.BKHostOuterIPField,
			common.BKHostInnerIPv6Field,
			common.BKHostOuterIPv6Field,
		},
		Page: metadata.BasePage{
			Start: 0,
			Limit: 1,
		},
	}
	resp, err := s.CoreAPI.ApiServer().ListHostWithoutApp(ctx, c.Request.Header, option)
	if err != nil {
		blog.Errorf("get host by id failed, hostID: %d, err: %+v, rid: %s", hostID, err, rid)
		result := metadata.BaseResp{Result: false, Code: common.CCErrHostGetFail,
			ErrMsg: defErr.Error(common.CCErrHostGetFail).Error()}
		c.JSON(http.StatusOK, result)
		return
	}
	if resp.Code != 0 || resp.Result == false {
		blog.Errorf("got host by id failed, hostID: %d, response: %+v, rid: %s", hostID, resp, rid)
		c.JSON(http.StatusOK, resp)
		return
	}
	if len(resp.Data.Info) == 0 {
		blog.Errorf("host not found, hostID: %d, rid: %s", hostID, rid)
		result := metadata.BaseResp{Result: false, Code: common.CCErrCommNotFound,
			ErrMsg: defErr.Error(common.CCErrCommNotFound).Error()}
		c.JSON(http.StatusOK, result)
		return
	}
	type hostBase struct {
		HostID    int64  `json:"bk_host_id"`
		HostName  string `json:"bk_host_name"`
		InnerIP   string `json:"bk_host_innerip"`
		InnerIPv6 string `json:"bk_host_innerip_v6"`
		OuterIP   string `json:"bk_host_outerip"`
		OuterIPv6 string `json:"bk_host_outerip_v6"`
	}
	host := hostBase{}
	raw := resp.Data.Info[0]
	if err := mapstr.DecodeFromMapStr(&host, raw); err != nil {
		msg := fmt.Sprintf("decode response data into host failed, raw: %+v, err: %+v, rid: %s", raw, err, rid)
		blog.Error(msg)
		result := metadata.BaseResp{Result: false, Code: common.CCErrCommJSONUnmarshalFailed,
			ErrMsg: defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error()}
		c.JSON(http.StatusOK, result)
		return
	}

	ipOptions := make([]string, 0)
	ipOptions = append(ipOptions, "127.0.0.1")
	ipOptions = append(ipOptions, "0.0.0.0")
	if len(host.InnerIP) > 0 {
		ipOptions = append(ipOptions, host.InnerIP)
	}
	if len(host.OuterIP) > 0 {
		ipOptions = append(ipOptions, host.OuterIP)
	}

	// add process ipv6 options if needed
	if process.NeedIPv6OptionsHook() {
		ipOptions = append(ipOptions, "::1")
		ipOptions = append(ipOptions, "::")
		if len(host.InnerIPv6) > 0 {
			ipOptions = append(ipOptions, host.InnerIPv6)
		}
		if len(host.OuterIPv6) > 0 {
			ipOptions = append(ipOptions, host.OuterIPv6)
		}
	}

	result := metadata.ResponseDataMapStr{
		BaseResp: metadata.BaseResp{Result: true, Code: 0},
		Data: map[string]interface{}{
			"options": ipOptions,
		},
	}
	c.JSON(http.StatusOK, result)
	return
}
