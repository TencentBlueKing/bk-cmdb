package inst_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/inst"
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

	clsItem, err := model.FindClassificationsByCondition(common.CreateCondition().Field("bk_classification_id").Eq("bk_host_manage"))
	if nil != err {
		t.Errorf("failed to find classifications, %s", err.Error())
		return
	}

	if nil == clsItem {
		t.Errorf("not found the host classification")
		return
	}

	clsItem.ForEach(func(item model.Classification) {

		modelIter, err := item.FindModelsByCondition(common.CreateCondition().Field("bk_obj_id").Eq("host"))
		if nil != err {
			t.Errorf("failed to search classification, %s", err.Error())
			return
		}

		if nil == modelIter {
			t.Log("not found the model")
			return
		}

		// deal host model
		modelIter.ForEach(func(modelItem model.Model) {

			// create host
			hostInst, err := inst.CreateInst(modelItem)
			if nil != err {
				t.Errorf("failed to create host ")
				return
			}

			// Only test
			t.Logf("model name:%s", hostInst.GetModel().GetName())

			// set host value
			err = hostInst.SetValue("test", "test")
			if nil != err {
				t.Errorf("failed to set value, %s", err.Error())
				return
			}

			// save host info
			err = hostInst.Save()

			if nil != err {
				t.Errorf("failed to save ,%s", err.Error())
			}

		})

	})

}
