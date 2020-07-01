package local

import (
	"time"

	"configcenter/src/common/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)



var _ = Describe("Decode", func() {
	Context("[]int test", func() {
		It("", func() {
			var input = []interface{}{1,2,3,4,100,-20,10}
			var shouldout = []int64{1,2,3,4,100,-20,10}
			var results []int64

			// 获得结果,且 err == nil
			ret,err := decodeDistinctIntoSlice(input)
			Expect(err).NotTo(HaveOccurred())

			// 转化到int64,err == nil,且deep-equal
			results,err = util.SliceInterfaceToInt64(ret)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To( Equal(results))
		})
	})
	Context("[]uint test", func() {
		It("", func() {
			var input = []interface{}{uint(1),uint(2),uint(3),uint(4),uint(100),uint(111111111),uint(2222222),uint(1844674407370955161)}
			var shouldout = []int64{1,2,3,4,100,111111111,2222222,1844674407370955161}
			var results []int64

			// 获得结果,且 err == nil
			ret,err := decodeDistinctIntoSlice(input)
			Expect(err).NotTo(HaveOccurred())

			// 转化到int64,err == nil,且deep-equal
			results,err = util.SliceInterfaceToInt64(ret)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To( Equal(results))
		})

	})
	Context("[]string test", func() {
		It("", func() {
			var input = []interface{}{"a","b","c","hello world",""}
			var shouldout = []string{"a","b","c","hello world",""}
			var results []string

			// 获得结果,且 err == nil
			ret,err := decodeDistinctIntoSlice(input)
			Expect(err).NotTo(HaveOccurred())

			// 转化到string,err == nil，且deep-equal
			results,err = util.SliceInterfaceToString(ret)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To( Equal(results))
		})
	})
	Context("[]bool test", func() {
		It("", func() {
			var input = []interface{}{true,false,true}
			var shouldout = []bool{true,false,true}
			var results []bool

			// 获得结果,且 err == nil
			ret,err := decodeDistinctIntoSlice(input)
			Expect(err).NotTo(HaveOccurred())

			// 转化到Bool,err == nil，且deep-equal
			results,err = util.SliceInterfaceToBool(ret)
			Expect(err).NotTo(HaveOccurred())
			Expect(shouldout).To( Equal(results))
		})
	})
	Context("nil text", func() {
		It("", func() {
			results,err := decodeDistinctIntoSlice(nil)
			Expect(err).NotTo(HaveOccurred())

			// 结果不等于nil,且长度为0
			Expect(results).NotTo(BeNil())
			Expect(len(results)).To(Equal(0))
		})
	})

	Context("not convert", func() {
		It("", func() {
			// 因为decode检测的是别名是否符合,所以虽然time.Second的实际类型是int64,但别名是Duration，所以不能转换
			var badinput = []interface{}{time.Microsecond,time.Second,time.Millisecond}

			//错误发生
			_,err := decodeDistinctIntoSlice(badinput)
			Expect(err).To(HaveOccurred())
		})
	})

})
