package local

import (
	cc "context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

/*
decodeDistinctIntoSlice(ctx context.Context, dbResults []interface{}) (results interface{},error)
 if dbResults.Elem().Type is：
	string or bool or int64  		than：direct return
	int... or uint...    than： convert to int64
	others	than：return nil,error
*/

var _ = PDescribe("Decode", func() {
	var ctx  = cc.TODO()
	Context("[]int test", func() {
		It("", func() {
			var src = []int{1,2,3,4,100,-20,10}
			results,err := decodeDistinctIntoSlice(ctx,convert2Sliceface(src))

			//err == nil
			Expect(err).NotTo(HaveOccurred())
			//length should be equal
			Expect(len(results)).To(Equal(len(src)))

			for i,item := range results{

				val,ok := item.(int64)

				Expect(ok).To(BeTrue())

				Expect(val).To(Equal(int64(src[i])))
			}
		})
	})
	Context("[]uint test", func() {
		It("", func() {
			var src = []uint{1,2,3,4,100,111111111,2222222,18446744073709551615}

			results,err := decodeDistinctIntoSlice(ctx,convert2Sliceface(src))

			//err == nil
			Expect(err).NotTo(HaveOccurred())
			//length should be equal
			Expect(len(results)).To(Equal(len(src)))

			for i,item := range results{

				val,ok := item.(int64)

				Expect(ok).To(BeTrue())

				Expect(val).To(Equal(int64(src[i])))
			}
		})

	})
	Context("[]string test", func() {
		It("", func() {
			var src = []string{"a","b","c","hello world",""}

			results,err := decodeDistinctIntoSlice(ctx,convert2Sliceface(src))

			//err == nil
			Expect(err).NotTo(HaveOccurred())
			//length should be equal
			Expect(len(results)).To(Equal(len(src)))

			for i,item := range results{

				val,ok := item.(string)

				Expect(ok).To(BeTrue())

				Expect(val).To(Equal(src[i]))
			}
		})
	})
	Context("[]bool test", func() {
		It("", func() {
			var src = []bool{true,false,true}

			results,err := decodeDistinctIntoSlice(ctx,convert2Sliceface(src))

			//err == nil
			Expect(err).NotTo(HaveOccurred())
			//length should be equal
			Expect(len(results)).To(Equal(len(src)))

			for i,item := range results{

				val,ok := item.(bool)

				Expect(ok).To(BeTrue())

				Expect(val).To(Equal(src[i]))
			}
		})
	})
	Context("nil text", func() {
		It("", func() {

			results,err := decodeDistinctIntoSlice(ctx,nil)

			//err == nil
			Expect(err).NotTo(HaveOccurred())
			//length should be equal

			Expect(results).To(BeNil())
		})
	})
	Context("not convert", func() {
		It("", func() {
			var bad = []time.Duration{
				time.Second,
				time.Millisecond,
				time.Microsecond,
			}
			var src = make([]interface{},len(bad))
			for i := range bad{
				src[i] = bad[i]
			}
			_,err := decodeDistinctIntoSlice(ctx,src)

			Expect(err).To(HaveOccurred())
		})
	})

})
func convert2Sliceface(who interface{}) []interface{}{
	if who == nil {
		return nil
	}
	var ans []interface{}
	switch ve := who.(type) {
	case []string:
		ans = make([]interface{},len(ve))
		for i:=0;i < len(ve);i++{
			ans[i] = ve[i]
		}

	case []int:
		ans = make([]interface{},len(ve))
		for i:=0;i < len(ve);i++{
			ans[i] = ve[i]
		}
	case []uint:
		ans = make([]interface{},len(ve))
		for i:=0;i < len(ve);i++{
			ans[i] = ve[i]
		}
	case []bool:
		ans = make([]interface{},len(ve))
		for i:=0;i < len(ve);i++{
			ans[i] = ve[i]
		}
	default:
		Fail("only support []int,[]string,[]bool")

	}
	return ans
}
