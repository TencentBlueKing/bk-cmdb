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
	"errors"
	"testing"
)

func TestService_CreateUserGroup(t *testing.T) {

	// assert with response for custom
	resp, respBody := CallService(NewMockService().CreateUserGroup, `abcd`)
	if resp.StatusCode() != 400 {
		t.Fail()
	}
	if respBody == `{"result":false,"bk_error_code":1199000,"bk_error_msg":"invalid character 'a' looking for beginning of value","data":null}
` {
		t.Fail()
	}

	// assert with case array
	AssertCases(t, NewMockService().CreateUserGroup, []*TestCase{
		// case with expect response
		&TestCase{`{"k":"v"}`, `{"result":true,"bk_error_code":0,"bk_error_msg":"success","data":null}`, 200, nil},

		// case with callback
		&TestCase{`[1,2]`, ``, 0, func(responseBody string, status int) error {
			if status != 400 {
				return errors.New("bad status")
			}
			return nil
		}},
	})

	// assert with callback
	AssertCallback(t, NewMockService().CreateUserGroup, `{"k":"v"}`, func(responseBody string, status int) error {
		if status != 200 {
			return errors.New("bad status")
		}
		return nil
	})

	// assert with expect
	AssertEqual(t, NewMockService().CreateUserGroup, `{"k":"v"}`, `{"result":true,"bk_error_code":0,"bk_error_msg":"success","data":null}`, 200)

}
