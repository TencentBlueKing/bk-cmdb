/*
* Tencent is pleased to support the open source community by making
* 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
* Copyright (C) 2017 THL A29 Limited,
* a Tencent company. All rights reserved.
* Licensed under the MIT License (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at http://opensource.org/licenses/MIT
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on
* an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
* either express or implied. See the License for the
* specific language governing permissions and limitations under the License.
* We undertake not to change the open source license (MIT license) applicable
* to the current version of the project delivered to anyone in the future.
 */

package topo_server_test

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var (
	// 字段组合模板id
	fieldTemplateID, fieldTemplateID1, fieldTemplateID2 int64
	// 克隆的字段组合模板id
	cloneFieldTmplID int64
	// 字段组合模板属性字段数组
	fieldTmplAttrArr []metadata.FieldTemplateAttr
	// 字段组合模板唯一校验数组
	fieldTmplUniqueArr []metadata.FieldTemplateUnique
	// 字段组合模板测试模型
	fieldTmplObject, fieldTmplObject1, fieldTmplObject2 metadata.Object
	// 字段组合模板同步信息到模型的任务id
	taskID, failureTaskID string
	// 模型唯一校验信息
	objUnique metadata.ObjectUnique
	// 模型属性字段信息
	objAttr metadata.Attribute
)

var _ = Describe("field template test", func() {
	ctx := context.Background()

	It("normal field template test", func() {
		By("create field template")
		func() {
			option := metadata.CreateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					Name:        "my_tmp",
					Description: "my_tmp",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop",
						PropertyName: "prop",
						PropertyType: "int",
						Option: map[string]interface{}{
							"min": 1,
							"max": 1000,
						},
						Default: 500,
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{{
					Keys: []string{"prop"},
				}},
			}
			rsp, err := fieldTemplateClient.CreateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.ID).To(Not(Equal(int64(0))))
			fieldTemplateID = rsp.Data.ID
		}()

		By("create field template 'my_tmp1'")
		func() {
			option := metadata.CreateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					Name:        "my_tmp1",
					Description: "my_tmp1",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop1",
						PropertyName: "prop1",
						PropertyType: "singlechar",
						Option:       "",
						Default:      "aaa",
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{{
					Keys: []string{"prop1"},
				}},
			}
			rsp, err := fieldTemplateClient.CreateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.ID).To(Not(Equal(int64(0))))
			fieldTemplateID1 = rsp.Data.ID
		}()

		By("create field template 'my_tmp2'")
		func() {
			option := metadata.CreateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					Name:        "my_tmp2",
					Description: "my_tmp2",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop2",
						PropertyName: "prop2",
						PropertyType: "enummulti",
						Option: []map[string]interface{}{
							{
								"id":         "1",
								"is_default": true,
								"name":       "aaa",
								"type":       "text",
							},
							{
								"id":         "2",
								"is_default": false,
								"name":       "bbb",
								"type":       "text",
							},
							{
								"id":         "3",
								"is_default": false,
								"name":       "ccc",
								"type":       "text",
							},
						},
					},
				},
			}
			rsp, err := fieldTemplateClient.CreateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.ID).To(Not(Equal(int64(0))))
			fieldTemplateID2 = rsp.Data.ID
		}()

		By("find field template by id")
		func() {
			rsp, err := fieldTemplateClient.FindFieldTemplateByID(ctx, header, fieldTemplateID)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.ID).To(Equal(fieldTemplateID))
			Expect(rsp.Data.Name).To(Equal("my_tmp"))
		}()

		By("list field template")
		func() {
			option := metadata.CommonQueryOption{
				Page: metadata.BasePage{
					Limit: 100,
				},
			}
			rsp, err := fieldTemplateClient.ListFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Data.Info[0].ID).Should(BeElementOf(fieldTemplateID, fieldTemplateID1, fieldTemplateID2))
			Expect(rsp.Data.Info[0].Name).Should(BeElementOf("my_tmp", "my_tmp1", "my_tmp2"))
			Expect(rsp.Data.Info[0].Description).Should(BeElementOf("my_tmp", "my_tmp1", "my_tmp2"))
		}()

		By("again create field template 'my_tmp'")
		func() {
			option := metadata.CreateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					Name:        "my_tmp",
					Description: "my_tmp",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop",
						PropertyName: "prop",
						PropertyType: "int",
						Option: map[string]interface{}{
							"min": 1,
							"max": 1000,
						},
						Default: 500,
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{{
					Keys: []string{"prop"},
				}},
			}
			rsp, err := fieldTemplateClient.CreateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("create field template duplicate property")
		func() {
			option := metadata.CreateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					Name:        "field_tmp_test",
					Description: "field_tmp_test",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop_test",
						PropertyName: "prop_test",
						PropertyType: "int",
						Option: map[string]interface{}{
							"min": 1,
							"max": 1000,
						},
						Default: 500,
					},
					{
						PropertyID:   "prop_test",
						PropertyName: "prop_test",
						PropertyType: "int",
						Option: map[string]interface{}{
							"min": 1,
							"max": 1000,
						},
						Default: 500,
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{{
					Keys: []string{"prop_test"},
				}},
			}
			rsp, err := fieldTemplateClient.CreateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("create field template no attributes")
		func() {
			option := metadata.CreateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					Name:        "field_tmp_test",
					Description: "field_tmp_test",
				},
			}
			rsp, err := fieldTemplateClient.CreateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("update field template info")
		func() {
			option := metadata.FieldTemplate{
				ID:          fieldTemplateID,
				Name:        "my_tmp_updated",
				Description: "my_tmp_updated",
			}
			rsp, err := fieldTemplateClient.UpdateFieldTemplateInfo(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
		}()

		By("find field template by id after the update")
		func() {
			rsp, err := fieldTemplateClient.FindFieldTemplateByID(ctx, header, fieldTemplateID)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.ID).To(Equal(fieldTemplateID))
			Expect(rsp.Data.Name).To(Equal("my_tmp_updated"))
			Expect(rsp.Data.Description).To(Equal("my_tmp_updated"))
		}()

		By("update field template unsupported type")
		func() {
			option := metadata.UpdateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					ID:   fieldTemplateID,
					Name: "my_tmp",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop",
						PropertyName: "prop",
						PropertyType: "int",
						Option: map[string]interface{}{
							"min": 1,
							"max": 1000,
						},
						Default: 600,
					},
					{
						PropertyID:   "prop2",
						PropertyName: "prop2",
						PropertyType: "innertable",
						Option:       "",
						Default:      "",
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{
					{
						Keys: []string{"prop2"},
					},
				},
			}
			rsp, err := fieldTemplateClient.UpdateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("update field template")
		func() {
			option := metadata.UpdateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					ID:   fieldTemplateID,
					Name: "my_tmp",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop",
						PropertyName: "prop",
						PropertyType: "int",
						Option: map[string]interface{}{
							"min": 1,
							"max": 1000,
						},
						Default: 600,
					},
					{
						PropertyID:   "prop2",
						PropertyName: "prop2",
						PropertyType: "singlechar",
						Option:       "",
						Default:      "test",
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{
					{
						Keys: []string{"prop2"},
					},
				},
			}
			rsp, err := fieldTemplateClient.UpdateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
		}()

		By("count field template attribute")
		func() {
			option := metadata.CountFieldTmplResOption{
				TemplateIDs: []int64{fieldTemplateID},
			}
			rsp, err := fieldTemplateClient.CountFieldTemplateAttr(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Data[0].TemplateID).To(Equal(fieldTemplateID))
			Expect(rsp.Data[0].Count).To(Equal(2))
		}()

		By("list field template attribute")
		func() {
			option := metadata.ListFieldTmplAttrOption{
				TemplateID: fieldTemplateID,
			}
			rsp, err := fieldTemplateClient.ListFieldTemplateAttr(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Data.Info).NotTo(BeEmpty())
			fieldTmplAttrArr = rsp.Data.Info
		}()

		By("list field template unique")
		func() {
			option := metadata.ListFieldTmplUniqueOption{
				TemplateID: fieldTemplateID,
			}
			rsp, err := fieldTemplateClient.ListFieldTemplateUnique(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Data.Info).NotTo(BeNil())
			fieldTmplUniqueArr = rsp.Data.Info
		}()

		By("clone field template")
		func() {
			option := metadata.CloneFieldTmplOption{
				ID: fieldTemplateID,
				FieldTemplate: metadata.FieldTemplate{
					Name:        "my_tmp-copy",
					Description: "my_tmp-copy",
				},
			}
			rsp, err := fieldTemplateClient.CloneFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Data.ID).NotTo(Equal(0))
			cloneFieldTmplID = rsp.Data.ID
		}()

		By("delete field template")
		func() {
			option := metadata.DeleteFieldTmplOption{
				ID: cloneFieldTmplID,
			}
			rsp, err := fieldTemplateClient.DeleteFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(BeNil())
		}()

		By("find field template by id")
		func() {
			rsp, err := fieldTemplateClient.FindFieldTemplateByID(ctx, header, cloneFieldTmplID)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.ID).To(Equal(int64(0)))
		}()

		By("create object classification")
		func() {
			option := metadata.Classification{
				ClassificationID:   "field_tmpl_test_class",
				ClassificationName: "字段组合模板测试模型分组",
			}
			rsp, err := objectClient.CreateClassification(ctx, header, &option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("create object")
		func() {
			option := metadata.Object{
				ObjCls:     "field_tmpl_test_class",
				ObjectID:   "field_tmpl_test_obj",
				ObjectName: "字段组合模板测试模型",
				ObjIcon:    "icon-cc-default",
				TenantID:   "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ObjCls).To(Equal(option.ObjCls))
			Expect(rsp.Data.ObjectID).To(Equal(option.ObjectID))
			Expect(rsp.Data.ObjectName).To(Equal(option.ObjectName))
			Expect(rsp.Data.TenantID).To(Equal(option.TenantID))
			Expect(rsp.Data.Creator).To(Equal(option.Creator))
			fieldTmplObject = rsp.Data
		}()

		By("create object 'field_tmpl_test_obj1'")
		func() {
			option := metadata.Object{
				ObjCls:     "field_tmpl_test_class",
				ObjectID:   "field_tmpl_test_obj1",
				ObjectName: "字段组合模板测试模型1",
				ObjIcon:    "icon-cc-default",
				TenantID:   "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ObjCls).To(Equal(option.ObjCls))
			Expect(rsp.Data.ObjectID).To(Equal(option.ObjectID))
			Expect(rsp.Data.ObjectName).To(Equal(option.ObjectName))
			Expect(rsp.Data.TenantID).To(Equal(option.TenantID))
			Expect(rsp.Data.Creator).To(Equal(option.Creator))
			fieldTmplObject1 = rsp.Data
		}()

		By("create object 'field_tmpl_test_obj2'")
		func() {
			option := metadata.Object{
				ObjCls:     "field_tmpl_test_class",
				ObjectID:   "field_tmpl_test_obj2",
				ObjectName: "字段组合模板测试模型2",
				ObjIcon:    "icon-cc-default",
				TenantID:   "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ObjCls).To(Equal(option.ObjCls))
			Expect(rsp.Data.ObjectID).To(Equal(option.ObjectID))
			Expect(rsp.Data.ObjectName).To(Equal(option.ObjectName))
			Expect(rsp.Data.TenantID).To(Equal(option.TenantID))
			Expect(rsp.Data.Creator).To(Equal(option.Creator))
			fieldTmplObject2 = rsp.Data
		}()

		By("field template bind object")
		func() {
			option := metadata.FieldTemplateBindObjOpt{
				ID:        fieldTemplateID,
				ObjectIDs: []int64{fieldTmplObject.ID, fieldTmplObject1.ID, fieldTmplObject2.ID},
			}
			rsp, err := fieldTemplateClient.FieldTemplateBindObject(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data).To(BeNil())
		}()

		By("field template 'fieldTemplateID2' bind object 'fieldTmplObject2'")
		func() {
			option := metadata.FieldTemplateBindObjOpt{
				ID:        fieldTemplateID2,
				ObjectIDs: []int64{fieldTmplObject2.ID},
			}
			rsp, err := fieldTemplateClient.FieldTemplateBindObject(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data).To(BeNil())
		}()

		By("list object and field template relationship")
		func() {
			option := metadata.ListObjFieldTmplRelOption{
				TemplateIDs: []int64{fieldTemplateID},
				ObjectIDs:   []int64{fieldTmplObject.ID},
			}
			rsp, err := fieldTemplateClient.ListObjFieldTmplRel(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Info[0].ObjectID).Should(BeElementOf(fieldTmplObject.ID, fieldTmplObject1.ID,
				fieldTmplObject2.ID))
			Expect(rsp.Data.Info[0].TemplateID).To(Equal(fieldTemplateID))
		}()

		By("list field template by object")
		func() {
			option := metadata.ListFieldTmplByObjOption{
				ObjectID: fieldTmplObject2.ID,
			}
			rsp, err := fieldTemplateClient.ListFieldTmplByObj(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Info[0].ID).Should(BeElementOf(fieldTemplateID, fieldTemplateID2))
		}()

		By("list object by field template")
		func() {
			option := metadata.ListObjByFieldTmplOption{
				TemplateID: fieldTemplateID,
				CommonQueryOption: metadata.CommonQueryOption{
					Page: metadata.BasePage{Limit: 100},
				},
			}
			rsp, err := fieldTemplateClient.ListObjByFieldTmpl(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Info[0].ID).Should(BeElementOf(fieldTmplObject.ID, fieldTmplObject1.ID,
				fieldTmplObject2.ID))
		}()

		By("sync field template info to object")
		func() {
			option := metadata.FieldTemplateSyncOption{
				TemplateID: fieldTemplateID,
				ObjectIDs:  []int64{fieldTmplObject.ID, fieldTmplObject2.ID},
			}
			rsp, err := fieldTemplateClient.SyncFieldTemplateInfoToObjects(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := rsp.Data.([]interface{})
			taskID = data[0].(string)
		}()

		By("list field template tasks status")
		func() {
			for {
				option := metadata.ListFieldTmplTaskStatusOption{
					TaskIDs: []string{taskID},
				}
				rsp, err := fieldTemplateClient.ListFieldTemplateTasksStatus(ctx, header, option)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data[0].TaskID).To(Equal(taskID))
				if rsp.Data[0].Status == "finished" {
					break
				}
				time.Sleep(3 * time.Second)
			}
		}()

		By("list field template sync status")
		func() {
			option := metadata.ListFieldTmpltSyncStatusOption{
				ID:        fieldTemplateID,
				ObjectIDs: []int64{fieldTmplObject.ID},
			}
			rsp, err := fieldTemplateClient.ListFieldTemplateSyncStatus(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data[0].ObjectID).To(Equal(fieldTmplObject.ID))
			Expect(rsp.Data[0].NeedSync).To(Equal(false))
		}()

		By("update field template")
		func() {
			option := metadata.UpdateFieldTmplOption{
				FieldTemplate: metadata.FieldTemplate{
					ID:   fieldTemplateID2,
					Name: "my_tmp2",
				},
				Attributes: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop2",
						PropertyName: "prop2",
						PropertyType: "enummulti",
						Option: []map[string]interface{}{
							{
								"id":         "1",
								"is_default": true,
								"name":       "aaa",
								"type":       "text",
							},
							{
								"id":         "2",
								"is_default": false,
								"name":       "bbb",
								"type":       "text",
							},
							{
								"id":         "3",
								"is_default": false,
								"name":       "ccc",
								"type":       "text",
							},
						},
					},
					{
						PropertyID:   "prop3",
						PropertyName: "prop3",
						PropertyType: "singlechar",
						Option:       "",
						Default:      "test",
					},
					{
						PropertyID:   "prop4",
						PropertyName: "prop4",
						PropertyType: "singlechar",
						Option:       "",
						Default:      "test",
					},
				},
				Uniques: []metadata.FieldTmplUniqueOption{
					{
						Keys: []string{"prop3"},
					},
				},
			}
			rsp, err := fieldTemplateClient.UpdateFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
		}()

		By("compare object attribute and field template attribute")
		func() {
			option := metadata.CompareFieldTmplAttrOption{
				TemplateID: fieldTemplateID2,
				ObjectID:   fieldTmplObject2.ID,
				Attrs: []metadata.FieldTemplateAttr{
					{
						PropertyID:   "prop2",
						PropertyName: "prop2",
						PropertyType: "enummulti",
						Option: []map[string]interface{}{
							{
								"id":         "1",
								"is_default": true,
								"name":       "aaa",
								"type":       "text",
							},
							{
								"id":         "2",
								"is_default": false,
								"name":       "bbb",
								"type":       "text",
							},
							{
								"id":         "3",
								"is_default": false,
								"name":       "ccc",
								"type":       "text",
							},
						},
					},
					{
						PropertyID:   "prop3",
						PropertyName: "prop3",
						PropertyType: "singlechar",
						Option:       "",
						Default:      "test",
					},
					{
						PropertyID:   "prop4",
						PropertyName: "prop4",
						PropertyType: "singlechar",
						Option:       "",
						Default:      "test",
					},
				},
			}
			rsp, err := fieldTemplateClient.CompareFieldTemplateAttr(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).NotTo(BeNil())
			Expect(rsp.Data.Create[0].PropertyID).Should(BeElementOf("prop3", "prop4"))
			Expect(rsp.Data.Conflict[0].PropertyID).Should(BeElementOf("prop2"))
			Expect(rsp.Data.Unchanged[0].PropertyID).Should(BeElementOf("prop", "bk_inst_name"))
		}()

		By("compare object unique and field template unique")
		func() {
			option := metadata.CompareFieldTmplUniqueOption{
				TemplateID: fieldTemplateID,
				ObjectID:   fieldTmplObject.ID,
				Uniques: []metadata.FieldTmplUniqueForUpdate{
					{
						Keys: []string{"prop3"},
					},
				},
			}
			rsp, err := fieldTemplateClient.CompareFieldTemplateUnique(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).NotTo(BeNil())
			Expect(rsp.Data.Create[0].Index).To(Equal(0))
		}()

		By("list field template sync status")
		func() {
			option := metadata.ListFieldTmpltSyncStatusOption{
				ID:        fieldTemplateID2,
				ObjectIDs: []int64{fieldTmplObject2.ID},
			}
			rsp, err := fieldTemplateClient.ListFieldTemplateSyncStatus(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data[0].ObjectID).To(Equal(fieldTmplObject2.ID))
			Expect(rsp.Data[0].NeedSync).To(Equal(true))
		}()

		By("sync field template info to object failure")
		func() {
			option := metadata.FieldTemplateSyncOption{
				TemplateID: fieldTemplateID2,
				ObjectIDs:  []int64{fieldTmplObject2.ID},
			}
			rsp, err := fieldTemplateClient.SyncFieldTemplateInfoToObjects(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := rsp.Data.([]interface{})
			failureTaskID = data[0].(string)
		}()

		By("list field template tasks status")
		func() {
			for {
				option := metadata.ListFieldTmplTaskStatusOption{
					TaskIDs: []string{failureTaskID},
				}
				rsp, err := fieldTemplateClient.ListFieldTemplateTasksStatus(ctx, header, option)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data[0].TaskID).To(Equal(failureTaskID))
				if rsp.Data[0].Status == "failure" {
					break
				}
				time.Sleep(3 * time.Second)
			}
		}()

		By("delete field template failure")
		func() {
			option := metadata.DeleteFieldTmplOption{
				ID: fieldTemplateID2,
			}
			rsp, err := fieldTemplateClient.DeleteFieldTemplate(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("search object unique")
		func() {
			rsp, err := objectClient.SearchObjectUnique(ctx, fieldTmplObject.ObjectID, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			objUniqueArr := make([]metadata.ObjectUnique, 0)
			rspData, err := json.Marshal(rsp.Data)
			err = json.Unmarshal(rspData, &objUniqueArr)
			Expect(err).ShouldNot(HaveOccurred())
			for _, ou := range objUniqueArr {
				if ou.TemplateID == 0 {
					continue
				}
				objUnique = ou
			}
			Expect(objUnique.ObjID).To(Equal(fieldTmplObject.ObjectID))
		}()

		By("list field template simple data by unique option")
		func() {
			option := metadata.ListTmplSimpleByUniqueOption{
				TemplateID: fieldTmplUniqueArr[0].ID,
				UniqueID:   int64(objUnique.ID),
			}
			rsp, err := fieldTemplateClient.ListFieldTmplByUniqueTmplIDForUI(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ID).To(Equal(fieldTemplateID))
		}()

		By("search object attribute")
		func() {
			data := map[string]interface{}{
				common.BKObjIDField:      fieldTmplObject.ObjectID,
				common.BKPropertyIDField: fieldTmplAttrArr[0].PropertyID,
			}
			rsp, err := objectClient.SelectObjectAttWithParams(ctx, header, data)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			objAttrResp := make([]metadata.ObjAttDes, 0)
			rspData, err := json.Marshal(rsp.Data)
			err = json.Unmarshal(rspData, &objAttrResp)
			Expect(err).ShouldNot(HaveOccurred())
			objAttr = objAttrResp[0].Attribute
			Expect(objAttr.ObjectID).To(Equal(fieldTmplObject.ObjectID))
		}()

		By("list field template simple data by object attribute option")
		func() {
			option := metadata.ListTmplSimpleByAttrOption{
				TemplateID: fieldTmplAttrArr[0].ID,
				AttrID:     objAttr.ID,
			}
			rsp, err := fieldTemplateClient.ListFieldTmplByObjectTmplIDForUI(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ID).To(Equal(fieldTemplateID))
		}()

		By("list field template simple data by object attribute option")
		func() {
			option := metadata.ListFieldTmplModelStatusOption{
				ID:        fieldTemplateID,
				ObjectIDs: []int64{fieldTmplObject.ID},
			}
			rsp, err := fieldTemplateClient.ListFieldTemplateModelStatus(ctx, header, option)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Info[0].ObjectID).To(Equal(fieldTmplObject.ID))
		}()
	})
})
