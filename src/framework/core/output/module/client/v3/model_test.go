package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/types"
	"fmt"
	"testing"
)

func TestCreateObject(t *testing.T) {

	cli := client.NewForConfig(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	id, err := cli.CCV3().Model().CreateObject(types.MapStr{
		"bk_supplier_account":  "0",
		"bk_obj_id":            common.UUID(),
		"bk_classification_id": "bk_biz_topo",
		"bk_obj_name":          fmt.Sprintf("test_%s", common.UUID()),
	})

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Logf("id:%d", id)
}

func TestDeleteObject(t *testing.T) {
	cli := client.NewForConfig(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)
	cond := common.CreateCondition().Field("id").Eq(16)

	err := cli.CCV3().Model().DeleteObject(cond)

	if nil != err {
		t.Errorf("failed to delete, error info is %s", err.Error())
	}

	t.Log("success")
}

func TestUpdateObject(t *testing.T) {
	cli := client.NewForConfig(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("id").Eq(16)

	err := cli.CCV3().Model().UpdateObject(types.MapStr{"bk_obj_name": "test_update"}, cond)

	if nil != err {
		t.Errorf("failed to update, error info is %s", err.Error())
	}

	t.Log("success")
}
func TestSearchObject(t *testing.T) {
	cli := client.NewForConfig(config.Config{"supplierAccount": "0", "user": "build_user", "ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_obj_id").Like("host")

	dataMap, err := cli.CCV3().Model().SearchObjects(cond)

	if nil != err {
		t.Errorf("failed to search, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		t.Logf("success, data:%+v", item.String("bk_obj_name"))
	}

}
