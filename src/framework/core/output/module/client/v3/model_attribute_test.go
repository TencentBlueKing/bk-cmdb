package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client"
	//"configcenter/src/framework/core/types"
	//"fmt"
	"testing"
)

func TestSearchObjectAttributes(t *testing.T) {
	cli := client.NewForConfig(config.Config{"supplierAccount": "0", "user": "build_user", "http://test.apiserver:8080": "http://test.apiserver:8080"}, nil)
	cond := common.CreateCondition().Field("bk_obj_id").Like("host")

	dataMap, err := cli.CCV3().Attribute().SearchObjectAttributes(cond)

	if nil != err {
		t.Errorf("failed to search, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		t.Logf("success, data:%+v", item.String("bk_property_name"))
	}

}
