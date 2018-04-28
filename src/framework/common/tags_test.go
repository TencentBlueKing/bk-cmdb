package common_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
	"testing"
)

type testObj struct {
	Filed1 string `field:"field_one"`
	Filed2 bool   `field:"field_two"`
	Filed3 int    `field:"field_three"`
	Filed4 int64  `field:"field_four"`
}

func TestGetTags(t *testing.T) {
	obj := &testObj{}
	tags := common.GetTags(obj)
	t.Logf("tags:%v", tags)
}

func TestSetValueByTags(t *testing.T) {
	obj := &testObj{}
	data := types.MapStr{
		"field_one":   "test_one_value",
		"field_two":   true,
		"field_three": 3,
		"field_four":  4,
	}
	common.SetValueToStructByTags(obj, data)
	t.Logf("tags:%v", obj)
}
