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

	"configcenter/src/test"
	"configcenter/src/test/run"

	. "github.com/onsi/ginkgo"
)

func CopyHeader(h http.Header) http.Header {
	hd := make(map[string][]string)
	for k, v := range h {
		hd[k] = v
	}

	return http.Header(hd)
}

var _ = Describe("Concurrent Test", func() {

	header := test.GetHeader()

	Describe("Concurrent add host test", func() {
		input := map[string]interface{}{
			"host_info": map[string]interface{}{
				"4": map[string]interface{}{
					"bk_host_innerip": "183.0.0.1",
					"bk_host_outerip": "10.10.10.12",
					"bk_cloud_id":     0,
				},
			},
		}

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				rsp, err := clientSet.ApiServer().AddHost(context.Background(), header, input)
				if err != nil {
					return err
				}
				if !rsp.Result {
					return errors.New("get app list failed")
				}
				return nil
			})
			fmt.Printf("Concurrent add host perform: \n" + m.Format())
		})
	})

})
