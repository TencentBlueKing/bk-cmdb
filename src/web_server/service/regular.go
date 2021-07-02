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
	"fmt"
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
	regular := c.PostForm("regular")

	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	if len(regular) == 0 {
		blog.Infof("params invalid regular")
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

	var isValid bool
	_, err := regexp.Compile(regular)
	if err == nil {
		isValid = true
	}

	c.JSON(http.StatusBadRequest, metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: true,
			Code:   0,
		},
		Data: map[string]bool{"is_valid": isValid},
	})

}

// VerifyRegularContentBatch verify regular content batch
func (s *Service) VerifyRegularContentBatch(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	requestBody := new(VerifyRegularContentBatchBody)
	err := c.BindJSON(&requestBody)
	if err != nil {
		blog.Errorf("verify regular content batch failed, but unmarshal body to json failed, err: %v, rid: %s", err, rid)
		msg := fmt.Sprintf("invalid body, parse json failed, err: %+v", err)
		c.String(http.StatusBadRequest, msg)
		return
	}

	if len(requestBody.Items) == 0 {
		blog.Errorf("need data items is null")
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        common.CCErrCommInstDataNil,
			ErrMsg:      defErr.Error(common.CCErrCommInstDataNil).Error(),
			Permissions: nil,
		})
		return
	}
	if len(requestBody.Items) > 100 {
		blog.Errorf("item field exceeds maximum value %v", len(requestBody.Items))
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        common.CCErrCommValExceedMaxFailed,
			ErrMsg:      defErr.Error(common.CCErrCommValExceedMaxFailed).Error(),
			Permissions: nil,
		})
		return
	}

	var regexpResult []int8
	result := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: true,
			Code:   0,
		},
	}

	var isValid int8
	for _, item := range requestBody.Items {
		isValid = 1
		_, err = regexp.Compile(item.Regular)
		if err != nil {
			isValid = 0
		} else {
			str := item.Content
			re := regexp.MustCompile(item.Regular)
			if len(re.FindAllString(str, -1)) < 1 {
				isValid = 0
			}
		}
		regexpResult = append(regexpResult, isValid)
	}
	result.Data = regexpResult

	c.JSON(http.StatusOK, result)
}

type VerifyRegularContentBatchBody struct {
	Items []struct {
		Regular string `json:"regular"`
		Content string `json:"content"`
	} `json:"items"`
}
