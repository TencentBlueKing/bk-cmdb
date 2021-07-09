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
	"net/http"
	"regexp"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

// VerifyRegularExpress verify regular express
func (s *Service) VerifyRegularExpress(c *gin.Context) {
	requestBody := new(VerifyRegularExpressRequest)
	err := c.BindJSON(requestBody)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	if len(requestBody.Regular) == 0 {
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommParamsInvalid,
				ErrMsg: defErr.Errorf(common.CCErrCommParamsInvalid, "Regular").Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}

	isValid := true
	invalidReason := ""
	_, err = regexp.Compile(requestBody.Regular)
	if err != nil {
		isValid = false
		invalidReason = err.Error()
	}

	c.JSON(http.StatusOK, metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: true,
			Code:   0,
		},
		Data: verifyRegularExpressDataResult{
			isValid,
			invalidReason,
		},
	})

}

// VerifyRegularContentBatch verify regular content batch
func (s *Service) VerifyRegularContentBatch(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	requestBody := new(VerifyRegularContentBatchRequest)
	err := c.BindJSON(requestBody)
	if err != nil {
		blog.Errorf("verify regular content batch failed, but unmarshal body to json failed, err: %v, rid: %s", err, rid)
		c.String(http.StatusBadRequest, "invalid body, parse json failed, err: %+v", err)
		return
	}

	if len(requestBody.Items) == 0 {
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        common.CCErrCommInstDataNil,
			ErrMsg:      defErr.Error(common.CCErrCommInstDataNil).Error(),
			Permissions: nil,
		})
		return
	}
	if len(requestBody.Items) > 100 {
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        common.CCErrCommValExceedMaxFailed,
			ErrMsg:      defErr.Error(common.CCErrCommValExceedMaxFailed).Error(),
			Permissions: nil,
		})
		return
	}

	var regexpResult []bool
	result := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: true,
			Code:   0,
		},
	}

	var isValid bool
	for _, item := range requestBody.Items {
		isValid = true
		_, err = regexp.Compile(item.Regular)
		if err != nil {
			isValid = false
		} else {
			if len(regexp.MustCompile(item.Regular).FindAllString(item.Content, -1)) < 1 {
				isValid = false
			}
		}
		regexpResult = append(regexpResult, isValid)
	}
	result.Data = regexpResult

	c.JSON(http.StatusOK, result)
}

type verifyRegularExpressDataResult struct {
	IsValid       bool   `json:"is_valid"`
	InvalidReason string `json:"invalid_reason"`
}

type VerifyRegularContentBatchRequest struct {
	Items []struct {
		Regular string `json:"regular"`
		Content string `json:"content"`
	} `json:"items"`
}

type VerifyRegularExpressRequest struct {
	Regular string `json:"regular"`
}
