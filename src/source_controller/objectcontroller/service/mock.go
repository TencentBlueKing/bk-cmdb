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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/storage/dal/mongo"
)

type TestCase struct {
	RequestBody        string
	ExpectResponseBody string
	ExpectStatus       int
	Callback           func(responseBody string, status int) error
}

func AssertEqual(t *testing.T, svc func(request *restful.Request, response *restful.Response), requestBody string, expectResponseBody string, expectStatus int) {
	resp, body := CallService(svc, requestBody)
	if strings.TrimSpace(expectResponseBody) != strings.TrimSpace(body) {
		t.Fail()
	}
	if expectStatus != resp.StatusCode() {
		t.Fail()
	}
}

func AssertCallback(t *testing.T, svc func(request *restful.Request, response *restful.Response), requestBody string, callback func(responseBody string, status int) error) {
	resp, body := CallService(svc, requestBody)
	if err := callback(body, resp.StatusCode()); err != nil {
		t.Fail()
	}
}

func AssertCases(t *testing.T, svc func(request *restful.Request, response *restful.Response), cases []*TestCase) {
	for _, c := range cases {
		if c.Callback == nil {
			AssertEqual(t, svc, c.RequestBody, c.ExpectResponseBody, c.ExpectStatus)
		} else {
			AssertCallback(t, svc, c.RequestBody, c.Callback)
		}
	}
}

func CallService(svc func(request *restful.Request, response *restful.Response), requestBody string) (response *restful.Response, responseBody string) {

	// build request
	bodyReader := strings.NewReader(requestBody)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/json")
	//httpRequest.Header.Set("Accept", "application/json;application/xml")
	httpRequest.Header.Set(common.BKHTTPOwnerID, "")
	req := &restful.Request{Request: httpRequest}

	// build response
	recorder := httptest.NewRecorder()
	response = &restful.Response{ResponseWriter: recorder}
	response.SetRequestAccepts("application/json;application/xml")
	response.WriteHeader(200)

	// call real service
	svc(req, response)

	// get response body
	body, _ := ioutil.ReadAll(recorder.Result().Body)
	responseBody = string(body)

	return
}

func NewMockService() *Service {
	core := &backbone.Engine{
		Language: language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:    errors.NewFromCtx(errors.EmptyErrorsSetting),
	}
	return &Service{Core: core, Instance: mongo.NewMock(), Cache: nil}
}
