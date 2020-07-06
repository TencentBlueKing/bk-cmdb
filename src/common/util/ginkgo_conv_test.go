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

package util_test

import (
	"configcenter/src/common/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test SliceInterfaceToString,Int64,Bool", func() {
	Context("Test SliceInterfaceToString ", func() {
		It("number1", func() {
			var input = []interface{}{"abcd","1111",""}
			var shouldout = []string{"abcd","1111",""}
			results ,err := util.SliceInterfaceToString(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})

		It("number2", func() {
			var input = []interface{}{}
			var shouldout = []string{}
			results ,err := util.SliceInterfaceToString(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})

		It("shoule error", func() {
			var input = []interface{}{"abcd",12}
			_ ,err := util.SliceInterfaceToString(input)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Test SliceInterfaceToBool", func() {
		It("number1", func() {
			var input = []interface{}{true,false,true,true}
			var shouldout = []bool{true,false,true,true}
			results ,err := util.SliceInterfaceToBool(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})

		It("number2", func() {
			var input = []interface{}{}
			var shouldout = []bool{}
			results ,err := util.SliceInterfaceToBool(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})

		It("shoule error", func() {
			var input = []interface{}{"abcd",12}
			_ ,err := util.SliceInterfaceToBool(input)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Test SliceInterfaceToInt64", func() {
		It("number1", func() {
			var input = []interface{}{int64(1),int32(32),int(32),uint8(100)}
			var shouldout = []int64{1,32,32,100}
			results ,err := util.SliceInterfaceToInt64(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})

		It("number2", func() {
			var input = []interface{}{int64(1),int64(32),int64(32),int64(100)}
			var shouldout = []int64{1,32,32,100}
			results ,err := util.SliceInterfaceToInt64(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})

		It("number3", func() {
			var input = []interface{}{}
			var shouldout = []int64{}
			results ,err := util.SliceInterfaceToInt64(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To(Equal(results))
		})
		It("shoule error", func() {
			var input = []interface{}{"abcd",12}
			_ ,err := util.SliceInterfaceToInt64(input)
			Expect(err).To(HaveOccurred())
		})
	})
})
