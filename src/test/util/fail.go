/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/util"

	"github.com/onsi/ginkgo"
)

var (
	response interface{} = nil
	rid      string
)

// RegisterResponse register api response for debugging
func RegisterResponse(rsp interface{}) {
	response = rsp
}

// RegisterResponseWithRid register api response and rid from header to output them for debugging, set new rid to header
func RegisterResponseWithRid(rsp interface{}, header http.Header) {
	response = rsp
	rid = util.GetHTTPCCRequestID(header)
	header.Set(common.BKHTTPCCRequestID, util.GenerateRID())
}

// Fail ginkgo test fail hook
func Fail(message string, callerSkip ...int) {
	msg := message
	if response != nil {
		j, _ := json.MarshalIndent(response, "", "\t")
		msg = fmt.Sprintf("Test failed, message:\n%v\n\nresponse:\n%s", message, j)
	}
	if rid != "" {
		msg += fmt.Sprintf("\nrid: %s\n", rid)
	}
	skip := 0
	if len(callerSkip) > 0 {
		skip = callerSkip[0] + 1
	}
	ginkgo.Fail(msg, skip)
}
