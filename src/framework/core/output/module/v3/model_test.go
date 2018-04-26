package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/v3"
	"configcenter/src/framework/core/types"
	"fmt"
	"testing"
)

func TestCreateClassification(t *testing.T) {

	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	id, err := cli.CreateClassification(types.MapStr{
		"bk_classification_id":   common.UUID(),
		"bk_classification_name": fmt.Sprintf("test_%s", common.UUID()),
	})

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Logf("id:%d", id)
}

func TestDeleteClassification(t *testing.T) {
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("id").Eq(9)

	err := cli.DeleteClassification(cond)

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Log("success")
}

func TestSearchClassification(t *testing.T) {
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("bk_classification_name").Like("test_")

	dataMap, err := cli.SearchClassifications(cond)

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		t.Logf("success, data:%+v", item.String("bk_classification_name"))
	}

}

func TestSearchClassificationWithObjects(t *testing.T) {
	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("bk_classification_name").Like("业务")

	dataMap, err := cli.SearchClassificationWithObjects(cond)

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Logf("data:%+v", dataMap)

	for _, item := range dataMap {
		data, dataErr := item.MapStrArray("bk_objects")
		if nil != dataErr {
			t.Errorf("the error is %s", dataErr.Error())
			continue
		}

		//t.Logf("success,time:%v data:%+v ", tm, data)
		for _, subItem := range data {
			tm, tmErr := subItem.Time("create_time")
			if nil != tmErr {
				t.Errorf("the error info is %s", tmErr.Error())
				continue
			}
			t.Logf("success,time:%v ", tm)
		}
	}

}
