package inst_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/output/module/v3"
	//"configcenter/src/framework/core/types"
	"testing"
)

func TestHostManager(t *testing.T) {

	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	items, err := model.FindClassificationsByCondition(common.CreateCondition().Field("bk_classification_id").Eq("bk_host_manage"))
	if nil != err {
		t.Errorf("failed to find classifications, %s", err.Error())
		return
	}

	for {
		item, err := items.Next()
		if nil != err {
			t.Errorf("failed to get next classification, %s ", err.Error())
			break
		}
		if nil == item {
			t.Log("exit")
			break
		}
		t.Logf("the classifications:%+v", item)

		modelIterator, err := item.FindModelsByCondition(common.CreateCondition().Field("bk_obj_id").Eq("host"))
		for {

			item, err := modelIterator.Next()
			if nil != err {
				t.Errorf("failed to get model, %s", err.Error())
				break
			}

			if nil == item {
				t.Log("no model")
				break
			}

			t.Logf("the model:%+v", item.GetName())

			attrs, _ := item.Attributes()
			for _, attr := range attrs {
				t.Logf("the attribute:%+v", attr.GetName())
			}

		}

	}

}
