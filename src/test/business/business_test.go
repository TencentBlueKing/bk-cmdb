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

package business_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/test/run"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func CopyHeader(h http.Header) http.Header {
	hd := make(map[string][]string)
	for k, v := range h {
		hd[k] = v
	}

	return http.Header(hd)
}

var _ = Describe("GetBusinessList", func() {
	// initialize http header
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, "0")
	header.Add(common.BKHTTPHeaderUser, "admin")
	header.Add("Content-Type", "application/json")

	Describe("get business list load test", func() {
		var cond map[string]interface{}
		BeforeEach(func() {
			cond = mapstr.MapStr{
				"condition": mapstr.MapStr{
					"bk_data_status": mapstr.MapStr{
						"$ne": "disabled",
					},
				},
			}
		})

		Measure("emit the request", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				_, err := clientSet.ApiServer().GetUserPrivilegeApp(context.Background(), header, "0", "admin", cond)
				Expect(err).Should(BeNil())
			})
			Expect(runtime.Seconds()).Should(BeNumerically("<", 0.05))
		}, 10)

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				h := CopyHeader(header)
				rsp, err := clientSet.ApiServer().GetUserPrivilegeApp(context.Background(), h, "0", "admin", cond)
				if err != nil {
					return err
				}
				if !rsp.Result {
					return errors.New("get app list failed")
				}
				return nil
			})
			fmt.Printf("get app list perform: \n" + m.Format())
		})

	})
})
