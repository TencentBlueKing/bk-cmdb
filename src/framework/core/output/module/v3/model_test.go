package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/v3"
	"configcenter/src/framework/core/types"
	"fmt"
	"testing"
)

func TestCreateObject(t *testing.T) {

	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	id, err := cli.CreateObject(types.MapStr{
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
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("id").Eq(16)

	err := cli.DeleteObject(cond)

	if nil != err {
		t.Errorf("failed to delete, error info is %s", err.Error())
	}

	t.Log("success")
}

func TestUpdateObject(t *testing.T) {
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("id").Eq(16)

	err := cli.UpdateObject(types.MapStr{"bk_obj_name": "test_update"}, cond)

	if nil != err {
		t.Errorf("failed to update, error info is %s", err.Error())
	}

	t.Log("success")
}
func TestSearchObject(t *testing.T) {
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("bk_obj_name").Like("test")

	dataMap, err := cli.SearchObjects(cond)

	if nil != err {
		t.Errorf("failed to search, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		t.Logf("success, data:%+v", item.String("bk_obj_name"))
	}

}


