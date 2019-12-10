package topo_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("object test", func() {
	var bizId string
	var childInstId string
	var setId string
	var childInstIdInt, bizIdInt int64
	objectClient := topoServerClient.Object()
	instClient := topoServerClient.Instance()

	Describe("mainline object test", func() {
		It("create business bk_biz_name = 'abc'", func() {
			test.ClearDatabase()
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_productor":  "",
				"bk_biz_tester":     "",
				"bk_biz_developer":  "",
				"operator":          "",
				"bk_biz_name":       "abc",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).To(ContainElement("abc"))
			bizIdInt = int64(rsp.Data["bk_biz_id"].(float64))
			bizId = strconv.FormatInt(bizIdInt, 10)
		})

		It("create mainline object bk_obj_id = 'test_object' and bk_obj_name='test_object'", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := map[string]interface{}{
				"bk_classification_id": "bk_biz_topo",
				"bk_obj_icon":          "icon-cc-business",
				"bk_obj_id":            "test_object",
				"bk_obj_name":          "test_object",
				"bk_supplier_account":  "0",
			}
			for k, v := range data {
				Expect(rsp.Data.(map[string]interface{})).To(HaveKeyWithValue(k, v))
			}
		})

		It("create same mainline object", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object nonexist bk_asst_obj_id", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "xxx",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object empty bk_asst_obj_id", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object less bk_asst_obj_id", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object less bk_obj_id", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object less bk_obj_name", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:   "bk_biz_topo",
					ObjectID: "test_object",
					OwnerID:  "0",
					ObjIcon:  "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object less bk_obj_icon", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object less bk_classification_id", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete mainline object bk_obj_id = 'test_object'", func() {
			rsp, err := objectClient.DeleteModel(context.Background(), "0", "test_object", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete mainline object twice", func() {
			rsp, err := objectClient.DeleteModel(context.Background(), "0", "test_object", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create mainline object bk_obj_id = 'cc_test_object' and bk_obj_name='cc_test_object'", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "cc_test_object",
					ObjectName: "cc_test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := map[string]interface{}{
				"bk_classification_id": "bk_biz_topo",
				"bk_obj_icon":          "icon-cc-business",
				"bk_obj_id":            "cc_test_object",
				"bk_obj_name":          "cc_test_object",
				"bk_supplier_account":  "0",
			}
			for k, v := range data {
				Expect(rsp.Data.(map[string]interface{})).To(HaveKeyWithValue(k, v))
			}
		})

		It("search mainline object", func() {
			rsp, err := objectClient.SelectModel(context.Background(), "0", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data1 := metadata.MainlineObjectTopo{
				ObjID:      "biz",
				ObjName:    "业务",
				OwnerID:    "0",
				NextObj:    "cc_test_object",
				NextName:   "cc_test_object",
				PreObjID:   "",
				PreObjName: "",
			}
			Expect(rsp.Data).To(ContainElement(data1))
			data2 := metadata.MainlineObjectTopo{
				ObjID:      "cc_test_object",
				ObjName:    "cc_test_object",
				OwnerID:    "0",
				NextObj:    "set",
				NextName:   "集群",
				PreObjID:   "biz",
				PreObjName: "业务",
			}
			Expect(rsp.Data).To(ContainElement(data2))
			data3 := metadata.MainlineObjectTopo{
				ObjID:      "set",
				ObjName:    "集群",
				OwnerID:    "0",
				NextObj:    "module",
				NextName:   "模块",
				PreObjID:   "cc_test_object",
				PreObjName: "cc_test_object",
			}
			Expect(rsp.Data).To(ContainElement(data3))
		})
	})

	Describe("instance topo test", func() {
		It("search instance topo", func() {
			rsp, err := objectClient.SelectInst(context.Background(), "0", bizId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			res := []map[string]interface{}{}
			json.Unmarshal(j, &res)
			data := res[0]
			Expect(data["bk_inst_name"].(string)).To(Equal("abc"))
			Expect(data["bk_obj_id"].(string)).To(Equal("biz"))
			Expect(data["bk_obj_name"].(string)).To(Equal("业务"))
			js, err := json.Marshal(data["child"].([]interface{})[0])
			child := map[string]interface{}{}
			json.Unmarshal(js, &child)
			Expect(child["bk_inst_name"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_id"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_name"].(string)).To(Equal("cc_test_object"))
			Expect(len(child["child"].([]interface{}))).To(Equal(0))
			childInstIdInt = int64(child["bk_inst_id"].(float64))
			childInstId = strconv.FormatInt(childInstIdInt, 10)
		})

		It("search instance topo child", func() {
			rsp, err := objectClient.SelectInstChild(context.Background(), "0", "cc_test_object", bizId, childInstId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []map[string]interface{}{}
			json.Unmarshal(j, &data)
			child := data[0]
			Expect(child["bk_inst_name"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_id"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_name"].(string)).To(Equal("cc_test_object"))
			Expect(len(child["child"].([]interface{}))).To(Equal(0))
		})

		It("search instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", bizId, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.SetName).To(Equal("空闲机池"))
			Expect(len(rsp.Data.Module)).To(Equal(3))
		})
	})

	Describe("classification test", func() {
		var clsId, clsId2 string

		It("create classification", func() {
			input := &metadata.Classification{
				ClassificationID:   "cc_class",
				ClassificationName: "新测试分类",
				ClassificationIcon: "icon-cc-middleware",
				OwnerID:            "0",
			}
			rsp, err := objectClient.CreateClassification(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data["bk_classification_icon"].(string)).To(Equal("icon-cc-middleware"))
			Expect(data["bk_classification_id"].(string)).To(Equal("cc_class"))
			Expect(data["bk_classification_name"].(string)).To(Equal("新测试分类"))
			Expect(data["bk_classification_type"].(string)).To(Equal(""))
			clsId = strconv.FormatInt(int64(data["id"].(float64)), 10)
		})

		It("create classification", func() {
			input := &metadata.Classification{
				ClassificationID:   "cc_est_object",
				ClassificationName: "cc_est_object",
				ClassificationIcon: "icon-cc-default-class",
				OwnerID:            "0",
			}
			rsp, err := objectClient.CreateClassification(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data["bk_classification_icon"].(string)).To(Equal("icon-cc-default-class"))
			Expect(data["bk_classification_id"].(string)).To(Equal("cc_est_object"))
			Expect(data["bk_classification_name"].(string)).To(Equal("cc_est_object"))
			Expect(data["bk_classification_type"].(string)).To(Equal(""))
			clsId2 = strconv.FormatInt(int64(data["id"].(float64)), 10)
		})

		It("create classification same ClassificationID", func() {
			input := &metadata.Classification{
				ClassificationID:   "cc_class",
				ClassificationName: "测试分类",
				ClassificationIcon: "icon-cc-middleware",
				OwnerID:            "0",
			}
			rsp, err := objectClient.CreateClassification(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create classification same ClassificationName", func() {
			input := &metadata.Classification{
				ClassificationID:   "cc_class1",
				ClassificationName: "cc_est_object",
				ClassificationIcon: "icon-cc-middleware",
				OwnerID:            "0",
			}
			rsp, err := objectClient.CreateClassification(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update classification", func() {
			input := map[string]interface{}{
				"bk_classification_name": "cc模型分类",
				"bk_classification_icon": "icon-cc-default-class",
			}
			rsp, err := objectClient.UpdateClassification(context.Background(), clsId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update classification same ClassificationName", func() {
			input := map[string]interface{}{
				"bk_classification_name": "cc_est_object",
			}
			rsp, err := objectClient.UpdateClassification(context.Background(), clsId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete classification", func() {
			input := map[string]interface{}{}
			rsp, err := objectClient.DeleteClassification(context.Background(), clsId2, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search classification", func() {
			input := map[string]interface{}{}
			rsp, err := objectClient.SelectClassificationWithParams(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).To(ContainSubstring("\"bk_classification_icon\":\"icon-cc-default-class\""))
			Expect(j).To(ContainSubstring("\"bk_classification_id\":\"cc_class\""))
			Expect(j).To(ContainSubstring("\"bk_classification_name\":\"cc模型分类\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_id\":\"cc_est_object\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_name\":\"cc_est_object\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_name\":\"新测试分类\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_name\":\"测试分类\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_id\":\"cc_class1\""))
		})
	})

	Describe("object test", func() {
		var objId string

		It("create object bk_classification_id = 'cc_class' and bk_obj_id='cc_obj'", func() {
			input := metadata.Object{
				ObjCls:     "cc_class",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "cc_obj",
				ObjectName: "cc模型",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.Object{}
			json.Unmarshal(j, &data)
			Expect(data.ObjCls).To(Equal(input.ObjCls))
			Expect(data.ObjIcon).To(Equal(input.ObjIcon))
			Expect(data.ObjectID).To(Equal(input.ObjectID))
			Expect(data.ObjectName).To(Equal(input.ObjectName))
			Expect(data.OwnerID).To(Equal(input.OwnerID))
			Expect(data.Creator).To(Equal(input.Creator))
		})

		It("delete classification with object", func() {
			input := map[string]interface{}{
				"bk_classification_id": "cc_class",
			}
			rsp, err := objectClient.DeleteClassification(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create object same bk_obj_id", func() {
			input := metadata.Object{
				ObjCls:     "bk_network",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "cc_obj",
				ObjectName: "cc",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create object same bk_obj_name", func() {
			input := metadata.Object{
				ObjCls:     "bk_network",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "cc",
				ObjectName: "cc模型",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create object invalid bk_classification_id", func() {
			input := metadata.Object{
				ObjCls:     "124334452",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "cc123",
				ObjectName: "cc123",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create object invalid bk_obj_name", func() {
			input := metadata.Object{
				ObjCls:     "bk_network",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "cc1234",
				ObjectName: "~!@#$%^&*()",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create object bk_classification_id = 'cc_class' and bk_obj_id='test_obj'", func() {
			input := metadata.Object{
				ObjCls:     "cc_class",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "test_obj",
				ObjectName: "test_obj",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.Object{}
			json.Unmarshal(j, &data)
			Expect(data.ObjCls).To(Equal(input.ObjCls))
			Expect(data.ObjIcon).To(Equal(input.ObjIcon))
			Expect(data.ObjectID).To(Equal(input.ObjectID))
			Expect(data.ObjectName).To(Equal(input.ObjectName))
			Expect(data.OwnerID).To(Equal(input.OwnerID))
			Expect(data.Creator).To(Equal(input.Creator))
			objId = strconv.FormatInt(data.ID, 10)
		})

		It("search classifications objects", func() {
			input := map[string]interface{}{
				"bk_classification_id": "cc_class",
			}
			rsp, err := objectClient.SelectClassificationWithObjects(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data[0]["bk_classification_id"].(string)).To(Equal("cc_class"))
			js, err := json.Marshal(data[0]["bk_objects"])
			objects := []metadata.Object{}
			json.Unmarshal(js, &objects)
			Expect(len(objects)).To(Equal(2))
			Expect(objects[0].ObjectID).To(Equal("cc_obj"))
			Expect(objects[1].ObjectID).To(Equal("test_obj"))
		})

		It("update object", func() {
			input := map[string]interface{}{
				"bk_obj_name": "test_obj1",
			}
			rsp, err := objectClient.UpdateObject(context.Background(), objId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update object same bk_obj_name", func() {
			input := map[string]interface{}{
				"bk_obj_name": "cc模型",
			}
			rsp, err := objectClient.UpdateObject(context.Background(), objId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update object invalid bk_obj_name", func() {
			input := map[string]interface{}{
				"bk_obj_name": "~!@#$%^&*()",
			}
			rsp, err := objectClient.UpdateObject(context.Background(), objId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search objects", func() {
			input := map[string]interface{}{
				"bk_obj_id":           "test_obj",
				"bk_supplier_account": "0",
			}
			rsp, err := objectClient.SelectObjectWithParams(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data[0]["bk_obj_id"].(string)).To(Equal("test_obj"))
			Expect(data[0]["bk_obj_name"].(string)).To(Equal("test_obj1"))
		})

		It("delete object", func() {
			input := map[string]interface{}{}
			rsp, err := objectClient.DeleteObject(context.Background(), objId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search object topo", func() {
			input := map[string]interface{}{
				"bk_classification_id": "cc_class",
			}
			rsp, err := objectClient.SelectObjectTopo(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search object topo graphics", func() {
			rsp, err := objectClient.SelectObjectTopoGraphics(context.Background(), "global", "0", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).To(ContainSubstring("\"node_name\":\"cc模型\""))
			Expect(j).NotTo(ContainSubstring("\"node_name\":\"test_obj1\""))
		})

		It("update object topo graphics", func() {
			input := map[string]interface{}{
				"bk_obj_id":  "cc_obj",
				"bk_inst_id": 0,
				"node_type":  "obj",
				"position": map[string]interface{}{
					"x": -75,
					"y": 108,
				},
			}
			rsp, err := objectClient.UpdateObjectTopoGraphics(context.Background(), "global", "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search object topo graphics", func() {
			rsp, err := objectClient.SelectObjectTopoGraphics(context.Background(), "global", "0", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).NotTo(ContainSubstring("\"position\":{\"x\":-75,\"y\":108}"))
		})
	})

	Describe("object attribute group test", func() {
		var groupId int64

		Describe("group test", func() {
			var group metadata.Group

			It("create group bk_obj_id='cc_obj'", func() {
				input := metadata.Group{
					GroupID:    "12",
					GroupName:  "1234",
					GroupIndex: 10,
					ObjectID:   "cc_obj",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Group{}
				json.Unmarshal(j, &data)
				Expect(data.GroupID).To(Equal(input.GroupID))
				Expect(data.GroupName).To(Equal(input.GroupName))
				Expect(data.GroupIndex).To(Equal(input.GroupIndex))
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				groupId = data.ID
			})

			It("create group bk_obj_id='cc_obj'", func() {
				input := metadata.Group{
					GroupID:    "1",
					GroupName:  "123",
					GroupIndex: 1,
					ObjectID:   "cc_obj",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Group{}
				json.Unmarshal(j, &data)
				Expect(data.GroupID).To(Equal(input.GroupID))
				Expect(data.GroupName).To(Equal(input.GroupName))
				Expect(data.GroupIndex).To(Equal(input.GroupIndex))
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				group = data
			})

			It("create group same GroupID", func() {
				input := metadata.Group{
					GroupID:    "1",
					GroupName:  "12345",
					GroupIndex: 2,
					ObjectID:   "cc_obj",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create group same GroupName", func() {
				input := metadata.Group{
					GroupID:    "2",
					GroupName:  "123",
					GroupIndex: 3,
					ObjectID:   "cc_obj",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create group invalid ObjectID", func() {
				input := metadata.Group{
					GroupID:    "123456",
					GroupName:  "123456",
					GroupIndex: 4,
					ObjectID:   "123456",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("update group", func() {
				input := &metadata.PropertyGroupCondition{
					Condition: map[string]interface{}{
						"id": group.ID,
					},
					Data: map[string]interface{}{
						"bk_group_name": "456",
					},
				}
				rsp, err := objectClient.UpdatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				group.GroupName = "456"
			})

			It("update group same bk_group_name", func() {
				input := &metadata.PropertyGroupCondition{
					Condition: map[string]interface{}{
						"id": group.ID,
					},
					Data: map[string]interface{}{
						"bk_group_name": "1234",
					},
				}
				rsp, err := objectClient.UpdatePropertyGroup(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("delete group", func() {
				rsp, err := objectClient.DeletePropertyGroup(context.Background(), strconv.FormatInt(groupId, 10), header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				groupId = group.ID
			})

			It("search group bk_obj_id='cc_obj'", func() {
				input := map[string]interface{}{}
				rsp, err := objectClient.SelectPropertyGroupByObjectID(context.Background(), "0", "cc_obj", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []metadata.Group{}
				json.Unmarshal(j, &data)
				exist := false
				for _, grp := range data {
					if grp.ID == group.ID {
						exist = true
						Expect(grp.GroupID).To(Equal(group.GroupID))
						Expect(grp.GroupName).To(Equal(group.GroupName))
						Expect(grp.GroupIndex).To(Equal(group.GroupIndex))
						Expect(grp.ObjectID).To(Equal(group.ObjectID))
						Expect(grp.OwnerID).To(Equal(group.OwnerID))
					}
				}
				Expect(exist).To(Equal(true))
				Expect(j).NotTo(ContainSubstring("123"))
			})
		})

		Describe("object attribute test", func() {
			var attrId, attrId1 string

			It("create object attribute bk_obj_id='cc_obj' and bk_property_id='test_sglchar' and bk_property_name='test_sglchar'", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "test_sglchar",
						PropertyName:  "test_sglchar",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Attribute{}
				json.Unmarshal(j, &data)
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.PropertyID).To(Equal(input.PropertyID))
				Expect(data.PropertyName).To(Equal(input.PropertyName))
				Expect(data.PropertyGroup).To(Equal(input.PropertyGroup))
				Expect(data.IsEditable).To(Equal(input.IsEditable))
				Expect(data.PropertyType).To(Equal(input.PropertyType))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				attrId = strconv.FormatInt(data.ID, 10)
			})

			It("create object attribute with same bk_property_id", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "test_sglchar",
						PropertyName:  "sglchar",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute with same bk_property_name", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "sglchar",
						PropertyName:  "test_sglchar",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute invalid ObjectID", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "123456",
						PropertyID:    "sglchar",
						PropertyName:  "123456",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute long PropertyID", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
						PropertyName:  "1234567",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute invalid PropertyID", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "~!@#$%^%^",
						PropertyName:  "12345678",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute invalid PropertyName", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "cc1",
						PropertyName:  "~!@#$%^%^",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute invalid PropertyType", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "cc2",
						PropertyName:  "123456789",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "xxxxxxxxxxxxxxx",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("create object attribute bk_obj_id='cc_obj' and bk_property_id='test_singlechar' and bk_property_name='test_singlechar' and invalid PropertyGroup", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "test_singlechar",
						PropertyName:  "test_singlechar",
						PropertyGroup: "abcdefg",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Attribute{}
				json.Unmarshal(j, &data)
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.PropertyID).To(Equal(input.PropertyID))
				Expect(data.PropertyName).To(Equal(input.PropertyName))
				Expect(data.PropertyGroup).To(Equal("bizdefault"))
				Expect(data.IsEditable).To(Equal(input.IsEditable))
				Expect(data.PropertyType).To(Equal(input.PropertyType))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				attrId1 = strconv.FormatInt(data.ID, 10)
			})

			It("update object attribute id="+attrId1, func() {
				input := map[string]interface{}{
					"bk_property_name": "ayayyaya",
				}
				rsp, err := apiServerClient.UpdateObjectAtt(context.Background(), attrId1, header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("update object attribute same bk_property_name", func() {
				input := map[string]interface{}{
					"bk_property_name": "test_sglchar",
				}
				rsp, err := apiServerClient.UpdateObjectAtt(context.Background(), attrId1, header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("update object attribute invalid bk_property_name", func() {
				input := map[string]interface{}{
					"bk_property_name": "~!@#$%^%^",
				}
				rsp, err := apiServerClient.UpdateObjectAtt(context.Background(), attrId1, header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("delete object attribute id="+attrId, func() {
				rsp, err := apiServerClient.DeleteObjectAtt(context.Background(), attrId, header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search object attribute", func() {
				input := mapstr.MapStr{
					"bk_obj_id": "cc_obj",
				}
				rsp, err := apiServerClient.GetObjectAttr(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []map[string]interface{}{}
				json.Unmarshal(j, &data)
				Expect(len(data)).To(Equal(2))
				Expect(data).To(ContainElement(HaveKeyWithValue("bk_property_name", "ayayyaya")))
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_name", "test_singlechar")))
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_name", "test_sglchar")))
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_name", "sglchar")))
			})
		})

		Describe("object attribute group test", func() {
			It("update object attribute property group", func() {
				arr := []metadata.PropertyGroupObjectAtt{
					metadata.PropertyGroupObjectAtt{},
				}
				arr[0].Condition.ObjectID = "cc_obj"
				arr[0].Condition.PropertyID = "test_singlechar"
				arr[0].Condition.OwnerID = "0"
				arr[0].Data.PropertyGroupID = "1"
				input := map[string]interface{}{
					"data": arr,
				}
				rsp, err := objectClient.UpdatePropertyGroupObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("update nonexist object attribute property group", func() {
				arr := []metadata.PropertyGroupObjectAtt{
					metadata.PropertyGroupObjectAtt{},
				}
				arr[0].Condition.ObjectID = "cc_obj"
				arr[0].Condition.PropertyID = "test_singlechar"
				arr[0].Condition.OwnerID = "0"
				arr[0].Data.PropertyGroupID = "10000"
				input := map[string]interface{}{
					"data": arr,
				}
				rsp, err := objectClient.UpdatePropertyGroupObjectAtt(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("search object attribute", func() {
				input := mapstr.MapStr{
					"bk_obj_id":      "cc_obj",
					"bk_property_id": "test_singlechar",
				}
				rsp, err := apiServerClient.GetObjectAttr(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				Expect(j).To(ContainSubstring("\"bk_property_group_name\":\"456\""))
				Expect(j).To(ContainSubstring("\"bk_property_group\":\"1\""))
			})

			It("delete object attribute property group with object", func() {
				rsp, err := objectClient.DeletePropertyGroup(context.Background(), strconv.FormatInt(groupId, 10), header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("search group bk_obj_id='cc_obj'", func() {
				input := map[string]interface{}{}
				rsp, err := objectClient.SelectPropertyGroupByObjectID(context.Background(), "0", "cc_obj", header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				Expect(j).To(ContainSubstring("\"bk_group_name\":\"456\""))
				Expect(j).To(ContainSubstring("\"bk_group_id\":\"1\""))
			})

			It("delete object attribute property group", func() {
				rsp, err := objectClient.DeletePropertyGroupObjectAtt(context.Background(), "0", "cc_obj", "test_singlechar", "1", header)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search object attribute", func() {
				input := mapstr.MapStr{
					"bk_obj_id":      "cc_obj",
					"bk_property_id": "test_singlechar",
				}
				rsp, err := apiServerClient.GetObjectAttr(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []map[string]interface{}{}
				json.Unmarshal(j, &data)
				Expect(j).NotTo(ContainSubstring("\"bk_property_group_name\":\"456\""))
				Expect(j).NotTo(ContainSubstring("\"bk_property_group\":\"1\""))
			})
		})
	})

	Describe("set test", func() {
		var setId1 string

		It("create set bk_biz_id="+bizId+" and bk_parent_id="+childInstId, func() {
			input := mapstr.MapStr{
				"bk_set_name":         "cc_set",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_set_name"].(string)).To(Equal("cc_set"))
			Expect(commonutil.GetStrByInterface(rsp.Data["bk_parent_id"])).To(Equal(strconv.FormatInt(childInstIdInt, 10)))
			Expect(int64(rsp.Data["bk_biz_id"].(float64))).To(Equal(bizIdInt))
			setId = strconv.FormatInt(int64(rsp.Data["bk_set_id"].(float64)), 10)
		})

		It(fmt.Sprintf("create set bk_biz_id=%s and bk_parent_id=%s", bizId, childInstId), func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_set_name"].(string)).To(Equal("test"))
			Expect(commonutil.GetStrByInterface(rsp.Data["bk_parent_id"])).To(Equal(strconv.FormatInt(childInstIdInt, 10)))
			Expect(int64(rsp.Data["bk_biz_id"].(float64))).To(Equal(bizIdInt))
			setId1 = strconv.FormatInt(int64(rsp.Data["bk_set_id"].(float64)), 10)
		})

		It("create set same bk_biz_id and bk_parent_id and bk_set_name", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create set invalid bk_biz_id", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test1",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           1000,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), "1000", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create set invalid bk_parent_id", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test2",
				"bk_parent_id":        1000,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create set less bk_parent_id", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test3",
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create set unmatch bk_biz_id and bk_parent_id", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test4",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           2,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), "2", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create set invalid bk_set_name", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "~!@#$%^&*()_+-=",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "2",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update set", func() {
			input := map[string]interface{}{
				"bk_set_name": "new_test",
			}
			rsp, err := instClient.UpdateSet(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update nonexist set", func() {
			input := map[string]interface{}{
				"bk_set_name": "test123",
			}
			rsp, err := instClient.UpdateSet(context.Background(), bizId, "10000", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update set same bk_set_name", func() {
			input := map[string]interface{}{
				"bk_set_name": "test",
			}
			rsp, err := instClient.UpdateSet(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update set invalid bk_set_name", func() {
			input := map[string]interface{}{
				"bk_set_name": "~!@#$%^&*()_+-=",
			}
			rsp, err := instClient.UpdateSet(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update set biz, parent", func() {
			input := map[string]interface{}{
				"bk_biz_id":    2,
				"bk_parent_id": 1,
			}
			rsp, err := instClient.UpdateSet(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete set", func() {
			rsp, err := instClient.DeleteSet(context.Background(), bizId, setId1, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search set", func() {
			input := &params.SearchParams{
				Condition: map[string]interface{}{},
				Page: map[string]interface{}{
					"sort": "id",
				},
			}
			rsp, err := instClient.SearchSet(context.Background(), "0", bizId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(map[string]interface{}(rsp.Data.Info[rsp.Data.Count-1])).To(HaveKeyWithValue("bk_set_name", "new_test"))
			Expect(map[string]interface{}(rsp.Data.Info[rsp.Data.Count-1])).To(HaveKeyWithValue("bk_service_status", "1"))
			Expect(map[string]interface{}(rsp.Data.Info[rsp.Data.Count-1])).To(HaveKeyWithValue("bk_set_env", "2"))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[rsp.Data.Count-1]["bk_parent_id"])).To(Equal(strconv.FormatInt(childInstIdInt, 10)))
			Expect(int64(rsp.Data.Info[rsp.Data.Count-1]["bk_biz_id"].(float64))).To(Equal(bizIdInt))
		})
	})

	Describe("module test", func() {
		var moduleId, moduleId1 string

		It(fmt.Sprintf("create module bk_biz_id=%s and bk_set_id=%s", bizId, setId), func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_module_name"].(string)).To(Equal("cc_module"))
			Expect(strconv.FormatInt(int64(rsp.Data["bk_set_id"].(float64)), 10)).To(Equal(setId))
			Expect(commonutil.GetStrByInterface(rsp.Data["bk_parent_id"])).To(Equal(setId))
			moduleId = strconv.FormatInt(int64(rsp.Data["bk_module_id"].(float64)), 10)
		})

		It(fmt.Sprintf("create module bk_biz_id=%s and bk_set_id=%s", bizId, setId), func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_module_name"].(string)).To(Equal("test_module"))
			Expect(strconv.FormatInt(int64(rsp.Data["bk_set_id"].(float64)), 10)).To(Equal(setId))
			Expect(commonutil.GetStrByInterface(rsp.Data["bk_parent_id"])).To(Equal(setId))
			moduleId1 = strconv.FormatInt(int64(rsp.Data["bk_module_id"].(float64)), 10)
		})

		It("create module same bk_module_name", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create module invalid bk_biz_id", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module1",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), "1000", setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create module invalid bk_set_id", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module2",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, "1000", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create module unmatch bk_biz_id and bk_set_id", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module3",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), "2", setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create module invalid bk_module_name", func() {
			input := map[string]interface{}{
				"bk_module_name":      "~!@#$%^&*()_+-=",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create module invalid bk_parent_id", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module4",
				"bk_parent_id":        1000,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create module less bk_parent_id", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test_module4",
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update module", func() {
			input := map[string]interface{}{
				"bk_module_name": "new_module",
			}
			rsp, err := instClient.UpdateModule(context.Background(), bizId, setId, moduleId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update nonexist module", func() {
			input := map[string]interface{}{
				"bk_module_name": "new_module",
			}
			rsp, err := instClient.UpdateModule(context.Background(), bizId, setId, "10000", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update module same bk_module_name", func() {
			input := map[string]interface{}{
				"bk_module_name": "new_module",
			}
			rsp, err := instClient.UpdateModule(context.Background(), bizId, setId, moduleId1, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update module invalid bk_module_name", func() {
			input := map[string]interface{}{
				"bk_module_name": "~!@#$%^&*()_+-=",
			}
			rsp, err := instClient.UpdateModule(context.Background(), bizId, setId, moduleId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update module set, biz, parent", func() {
			input := map[string]interface{}{
				"bk_set_id":    1,
				"bk_biz_id":    2,
				"bk_parent_id": 1,
			}
			rsp, err := instClient.UpdateModule(context.Background(), bizId, setId, moduleId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete module", func() {
			rsp, err := instClient.DeleteModule(context.Background(), bizId, setId, moduleId1, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module", func() {
			input := &params.SearchParams{
				Condition: map[string]interface{}{},
				Page: map[string]interface{}{
					"sort": "id",
				},
			}
			rsp, err := instClient.SearchModule(context.Background(), "0", bizId, setId, header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			Expect(map[string]interface{}(rsp.Data.Info[0])).To(HaveKeyWithValue("bk_module_name", "new_module"))
			Expect(strconv.FormatInt(int64(rsp.Data.Info[0]["bk_set_id"].(float64)), 10)).To(Equal(setId))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[0]["bk_parent_id"])).To(Equal(setId))
		})
	})
})
