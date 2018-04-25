package common_test

import (
	"configcenter/src/framework/common"
	"testing"
)

func TestCondition(t *testing.T) {

	cond := common.CreateCondition()
	cond.Field("test_field").Eq(1024).Field("test_field2").In([]int{0, 1, 2, 3})
	result := cond.ToMapStr()

	t.Logf("the result:%+v", string(result.ToJSON()))
}
