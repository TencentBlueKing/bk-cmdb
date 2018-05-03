package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client/v3"
	//"configcenter/src/framework/core/types"
	//"fmt"
	"testing"
)

func TestSearchObjectAttributes(t *testing.T) {
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("bk_obj_id").Like("host")

	dataMap, err := cli.SearchObjectAttributes(cond)

	if nil != err {
		t.Errorf("failed to search, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		t.Logf("success, data:%+v", item.String("bk_property_name"))
	}

}
