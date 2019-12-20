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

	"github.com/onsi/ginkgo"
)

var (
	response interface{} = nil
)

func RegisterResponse(rsp interface{}) {
	response = rsp
}

func Fail(message string, callerSkip ...int) {
	msg := message
	if response != nil {
		j, _ := json.MarshalIndent(response, "", "\t")
		msg = fmt.Sprintf("Test failed, message:\n%v\n\nresponse:\n%s", message, j)
	}
	skip := 0
	if len(callerSkip) > 0 {
		skip = callerSkip[0] + 1
	}
	ginkgo.Fail(msg, skip)
}
