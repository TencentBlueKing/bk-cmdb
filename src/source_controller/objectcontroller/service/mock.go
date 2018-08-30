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
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/storage/mock"
	"github.com/emicklei/go-restful"
	"net/http"
	"net/http/httptest"
	"strings"
)

func NewRestfulTestCase(data string) (*Service, *restful.Request, *restful.Response) {
	req, resp := NewRestfulRequestResponse(data)
	return NewService(), req, resp
}

func NewRestfulRequestResponse(data string) (*restful.Request, *restful.Response) {
	bodyReader := strings.NewReader(data)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set(common.BKHTTPOwnerID, "")
	request := &restful.Request{Request: httpRequest}

	resp := &restful.Response{ResponseWriter: httptest.NewRecorder()}

	resp.SetRequestAccepts("application/json;application/xml")
	resp.WriteHeader(200)
	resp.Write([]byte("ok"))

	return request, resp
}

func NewService() *Service {
	core := &backbone.Engine{
		Language: language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:    errors.NewFromCtx(errors.EmptyErrorsSetting),
	}
	return &Service{Core: core, Instance: &mock.MockDB{}, Cache: nil}
}
