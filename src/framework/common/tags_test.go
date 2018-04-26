package common_test

import (
	"configcenter/src/framework/common"
	"testing"
)

type testObj struct {
	filed1 string `field:"filed_one"`
	filed2 string `field:"filed_two"`
}

func TestGetTags(t *testing.T) {
	obj := &testObj{}
	tags := common.GetTags(obj)
	t.Logf("tags:%v", tags)
}
