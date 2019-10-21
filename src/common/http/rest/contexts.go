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

package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"github.com/emicklei/go-restful"
)

type Kit struct {
	Rid             string
	Header          http.Header
	Ctx             context.Context
	CCError         errors.DefaultCCErrorIf
	User            string
	SupplierAccount string
}

type Contexts struct {
	Kit            *Kit
	Request        *restful.Request
	resp           *restful.Response
	respStatusCode int
}

func (c *Contexts) DecodeInto(to interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Request.Body)
	if err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, read request body failed, err: %v", c.Kit.Rid, err))
		return c.Kit.CCError.Error(common.CCErrCommHTTPReadBodyFailed)
	}

	if err := json.Unmarshal(body, to); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, unmarshal request body failed, err: %v, body: %s", c.Kit.Rid, err, string(body)))
		return c.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	return nil
}

func (c *Contexts) RespEntity(data interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	c.resp.Header().Add(common.BKHTTPCCRequestID, c.Kit.Rid)
	if err := c.resp.WriteAsJson(metadata.NewSuccessResp(data)); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, response http request failed, err: %v", c.Kit.Rid, err))
	}
}

type CountInfo struct {
	Count int64       `json:"count"`
	Info  interface{} `json:"info"`
}

func (c *Contexts) RespEntityWithCount(count int64, info interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	c.resp.Header().Add(common.BKHTTPCCRequestID, c.Kit.Rid)
	resp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data: CountInfo{
			Count: count,
			Info:  info,
		},
	}
	if err := c.resp.WriteAsJson(resp); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, response http request failed, err: %v", c.Kit.Rid, err))
	}
}

func (c *Contexts) WithStatusCode(statusCode int) *Contexts {
	c.respStatusCode = statusCode
	return c
}

// WriteError is used to write a error response to the http client, which means the request occur an error.
// It receive an err and an optional error code parameter.
// It will testify the err, if the err is a CCErrorCoder, then the error code inside it will be used.
// Otherwise, if errCode is set and > 0, then errCode value is used.
// Finally, if error code is not set and err is not CCErrorCoder, then it will be set with a default
// CCSystemBusy code.
// This function will also write a log when it's called which contains the request id field.
func (c *Contexts) RespWithError(err error, errCode int, format string, args ...interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, %s, err: %v", c.Kit.Rid, fmt.Sprintf(format, args), err))

	var code int
	var errMsg string
	if err != nil {
		t, yes := err.(errors.CCErrorCoder)
		if yes {
			code = t.GetCode()
			errMsg = t.Error()
		} else {
			if errCode > 0 {
				code = errCode
				errMsg = c.Kit.CCError.Error(code).Error()
			} else {
				code = common.CCErrorUnknownOrUnrecognizedError
				errMsg = c.Kit.CCError.Error(code).Error()
			}
		}
		// log the error

	} else {
		code = common.CCErrorUnknownOrUnrecognizedError
		errMsg = c.Kit.CCError.Error(code).Error()
	}

	c.resp.Header().Set("Content-Type", "application/json")
	c.resp.Header().Add(common.BKHTTPCCRequestID, c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: errMsg,
			Code:   code,
		},
		Data: nil,
	}

	if err := c.resp.WriteAsJson(body); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, response http request with error failed, err: %v", c.Kit.Rid, err))
		return
	}
}

func (c *Contexts) RespAutoError(err error) {
	blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, err: %v", c.Kit.Rid, err))
	var code int
	var errMsg string
	if err != nil {
		t, yes := err.(errors.CCErrorCoder)
		if yes {
			code = t.GetCode()
			errMsg = t.Error()
		} else {
			code = common.CCErrorUnknownOrUnrecognizedError
			errMsg = c.Kit.CCError.Error(code).Error()
		}
	} else {
		code = common.CCErrorUnknownOrUnrecognizedError
		errMsg = c.Kit.CCError.Error(code).Error()
	}

	c.resp.Header().Set("Content-Type", "application/json")
	c.resp.Header().Add(common.BKHTTPCCRequestID, c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: errMsg,
			Code:   code,
		},
		Data: "",
	}

	if err := c.resp.WriteAsJson(body); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, response http request with error failed, err: %v", c.Kit.Rid, err))
		return
	}
}

// WriteErrorf used to write a error response to the request client.
// it will wrapper the error with error code and other errorf args.
// errorf is used to format multiple-language error message.
// it will also will log the error at the same time with logMsg.
func (c *Contexts) RespErrorCodeF(errCode int, logMsg string, errorf ...interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	blog.ErrorfDepth(1, fmt.Errorf("rid: %s, %s", c.Kit.Rid, logMsg))

	c.resp.Header().Set("Content-Type", "application/json")
	c.resp.Header().Add(common.BKHTTPCCRequestID, c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: c.Kit.CCError.CCErrorf(errCode, errorf).Error(),
			Code:   errCode,
		},
		Data: "",
	}

	if err := c.resp.WriteAsJson(body); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, response http request with error failed, err: %v", c.Kit.Rid, err))
		return
	}
}

func (c *Contexts) RespErrorCodeOnly(errCode int, format string, args ...interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, %s", c.Kit.Rid, fmt.Sprintf(format, args)))

	c.resp.Header().Set("Content-Type", "application/json")
	c.resp.Header().Add(common.BKHTTPCCRequestID, c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: c.Kit.CCError.Error(errCode).Error(),
			Code:   errCode,
		},
		Data: "",
	}

	if err := c.resp.WriteAsJson(body); err != nil {
		blog.ErrorfDepth(1, fmt.Sprintf("rid: %s, response http request with error failed, err: %v", c.Kit.Rid, err))
		return
	}
}
