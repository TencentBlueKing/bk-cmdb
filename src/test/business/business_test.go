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
	"time"

	//"configcenter/src/common/mapstr"
	"configcenter/src/test"
	//"configcenter/src/common"
	params "configcenter/src/common/paraparse"
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

var _ = Describe("Business Test", func() {

	header := test.GetHeader()

	Describe("get business list load test", func() {
		var cond *params.SearchParams
		BeforeEach(func() {
			cond = &params.SearchParams{
				Condition: map[string]interface{}{
					"bk_data_status": map[string]interface{}{
						"$ne": "disabled",
					},
				},
			}
		})

		Measure("emit the request", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				_, err := clientSet.ApiServer().SearchBiz(context.Background(), "0", header, cond)
				Expect(err).Should(BeNil())
			})
			Expect(runtime.Seconds()).Should(BeNumerically("<", 0.07))
		}, 10)

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				//h := CopyHeader(header)

				rsp, err := clientSet.ApiServer().SearchBiz(context.Background(), "0", header, cond)
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

	Describe("create business load test", func() {
		//var header = test.GetHeader()
		Measure("emit the request", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {

				ts := fmt.Sprintf("%d", time.Now().UnixNano())
				morediff := fmt.Sprintf("%d", time.Now().UnixNano())
				bizSuffix := ts + morediff
				input := map[string]interface{}{
					"life_cycle":        "2",
					"language":          "1",
					"bk_biz_maintainer": "admin",
					"bk_biz_productor":  "",
					"bk_biz_tester":     "",
					"bk_biz_developer":  "",
					"operator":          "",
					"bk_biz_name":       "testBiz" + bizSuffix,
					"time_zone":         "Africa/Accra",
				}
				_, err := clientSet.ApiServer().CreateBiz(context.Background(), "0", header, input)
				Expect(err).Should(BeNil())
			})
			Expect(runtime.Seconds()).Should(BeNumerically("<", 0.9))
		}, 10)

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				ts := fmt.Sprintf("%d", time.Now().UnixNano())
				morediff := fmt.Sprintf("%d", time.Now().UnixNano())
				bizSuffix := ts + morediff
				input := map[string]interface{}{
					"life_cycle":        "2",
					"language":          "1",
					"bk_biz_maintainer": "admin",
					"bk_biz_productor":  "",
					"bk_biz_tester":     "",
					"bk_biz_developer":  "",
					"operator":          "",
					"bk_biz_name":       "testBiz" + bizSuffix,
					"time_zone":         "Africa/Accra",
				}
				rsp, err := clientSet.ApiServer().CreateBiz(context.Background(), "0", header, input)
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
