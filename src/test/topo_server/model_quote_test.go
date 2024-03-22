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

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	testutil "configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	// 模型引用类型测试模型
	modelQuoteObj metadata.Object
	// 引用类型字段
	quoteProp metadata.Attribute
	// 模型实例id
	instanceID int64
	// 引用类型实例id数组
	quoteInstIDArr []uint64
)

var _ = Describe("model quote type test", func() {
	ctx := context.Background()
	It("model quote type test", func() {
		By("create object classification")
		func() {
			option := metadata.Classification{
				ClassificationID:   "model_quote_test_class",
				ClassificationName: "模型引用类型测试模型分组",
			}
			rsp, err := objectClient.CreateClassification(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("create object")
		func() {
			option := metadata.Object{
				ObjCls:     "model_quote_test_class",
				ObjectID:   "model_quote_test_obj",
				ObjectName: "模型引用类型测试模型",
				ObjIcon:    "icon-cc-default",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ObjCls).To(Equal(option.ObjCls))
			Expect(rsp.Data.ObjectID).To(Equal(option.ObjectID))
			Expect(rsp.Data.ObjectName).To(Equal(option.ObjectName))
			Expect(rsp.Data.OwnerID).To(Equal(option.OwnerID))
			Expect(rsp.Data.Creator).To(Equal(option.Creator))
			modelQuoteObj = rsp.Data
		}()

		By("create model attribute")
		func() {
			option := metadata.ObjAttDes{
				Attribute: metadata.Attribute{
					ObjectID:      modelQuoteObj.ObjectID,
					OwnerID:       "0",
					Creator:       "admin",
					PropertyID:    "model_quote_prop",
					PropertyName:  "model_quote_prop",
					PropertyGroup: common.BKDefaultField,
					Unit:          "",
					Placeholder:   "",
					IsEditable:    true,
					IsRequired:    false,
					PropertyType:  common.FieldTypeInnerTable,
					Option: map[string]interface{}{
						"header": []map[string]interface{}{
							{
								common.BKPropertyIDField:   "cc1",
								common.BKPropertyNameField: "cc1",
								common.BKPropertyTypeField: "int",
								"unit":                     "",
								"placeholder":              "",
								"editable":                 true,
								"isrequired":               false,
								"ismultiple":               false,
								"option": map[string]interface{}{
									"min": 1,
									"max": 1000,
								},
								"default": 500,
							},
							{
								common.BKPropertyIDField:   "cc2",
								common.BKPropertyNameField: "cc2",
								common.BKPropertyTypeField: "singlechar",
								"unit":                     "",
								"placeholder":              "",
								"editable":                 true,
								"isrequired":               false,
								"ismultiple":               false,
								"option":                   "",
								"default":                  "abc",
							},
						},
					},
				},
			}
			rsp, err := objectClient.CreateObjectAtt(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			quoteProp = rsp.Data.Attribute
			Expect(quoteProp.ObjectID).To(Equal(modelQuoteObj.ObjectID))
			Expect(quoteProp.PropertyID).To(Equal(option.PropertyID))
			Expect(quoteProp.PropertyGroup).To(Equal(common.BKDefaultField))
			Expect(quoteProp.PropertyType).To(Equal(common.FieldTypeInnerTable))
			Expect(quoteProp.Option).NotTo(BeNil())
		}()

		By("create model attribute model no exist")
		func() {
			option := metadata.ObjAttDes{
				Attribute: metadata.Attribute{
					ObjectID:      "model_quote_test_obj_noexist",
					OwnerID:       "0",
					Creator:       "admin",
					PropertyID:    "model_quote_prop",
					PropertyName:  "model_quote_prop",
					PropertyGroup: common.BKDefaultField,
					Unit:          "",
					Placeholder:   "",
					IsEditable:    true,
					IsRequired:    false,
					PropertyType:  common.FieldTypeInnerTable,
					Option: map[string]interface{}{
						"header": []map[string]interface{}{
							{
								common.BKPropertyIDField:   "cc1",
								common.BKPropertyNameField: "cc1",
								common.BKPropertyTypeField: "int",
								"unit":                     "",
								"placeholder":              "",
								"editable":                 true,
								"isrequired":               false,
								"ismultiple":               false,
								"option": map[string]interface{}{
									"min": 1,
									"max": 1000,
								},
								"default": 500,
							},
						},
					},
				},
			}
			rsp, err := objectClient.CreateObjectAtt(ctx, header, &option)
			Expect(err).NotTo(HaveOccurred())
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(rsp.Code).To(Equal(common.CCErrCommParamsIsInvalid))
		}()

		By("create model attribute 'model_quote_prop' again")
		func() {
			option := metadata.ObjAttDes{
				Attribute: metadata.Attribute{
					ObjectID:      modelQuoteObj.ObjectID,
					OwnerID:       "0",
					Creator:       "admin",
					PropertyID:    "model_quote_prop",
					PropertyName:  "model_quote_prop",
					PropertyGroup: common.BKDefaultField,
					Unit:          "",
					Placeholder:   "",
					IsEditable:    true,
					IsRequired:    false,
					PropertyType:  common.FieldTypeInnerTable,
					Option: map[string]interface{}{
						"header": []map[string]interface{}{
							{
								common.BKPropertyIDField:   "cc1",
								common.BKPropertyNameField: "cc1",
								common.BKPropertyTypeField: "int",
								"unit":                     "",
								"placeholder":              "",
								"editable":                 true,
								"isrequired":               false,
								"ismultiple":               false,
								"option": map[string]interface{}{
									"min": 1,
									"max": 1000,
								},
								"default": 500,
							},
						},
					},
				},
			}
			rsp, err := objectClient.CreateObjectAtt(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Code).To(Equal(common.CCErrCommDuplicateItem))
		}()

		By("create object instance")
		func() {
			option := mapstr.MapStr{
				"bk_inst_name": "obj_quote_test",
			}
			rsp, err := instClient.CreateInst(ctx, modelQuoteObj.ObjectID, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			instID, idExist := rsp.Data.Get(common.BKInstIDField)
			Expect(idExist).To(Equal(true))
			instanceID, err = util.GetInt64ByInterface(instID)
			Expect(err).ShouldNot(HaveOccurred())
			instName, nameExist := rsp.Data.Get(common.BKInstNameField)
			Expect(nameExist).To(Equal(true))
			Expect(util.GetStrByInterface(instName)).To(Equal("obj_quote_test"))
		}()

		By("create model quote instance")
		func() {
			option := metadata.BatchCreateQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				Data: []mapstr.MapStr{
					{
						"cc1":        10,
						"cc2":        "aaa",
						"bk_inst_id": instanceID,
					},
					{
						"cc1":        20,
						"cc2":        "bbb",
						"bk_inst_id": instanceID,
					},
				},
			}
			rsp, err := modelQuoteClient.BatchCreateQuotedInstance(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp).To(HaveLen(2))
			quoteInstIDArr = rsp
		}()

		By("create model quote instance but model instance no exist")
		func() {
			option := metadata.BatchCreateQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				Data: []mapstr.MapStr{
					{
						"cc1":        10,
						"cc2":        "aaa",
						"bk_inst_id": instanceID + 3256,
					},
					{
						"cc1":        20,
						"cc2":        "bbb",
						"bk_inst_id": instanceID + 3256,
					},
				},
			}
			rsp, err := modelQuoteClient.BatchCreateQuotedInstance(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("create model quote instance but model property no exist")
		func() {
			option := metadata.BatchCreateQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: "prop_no_exist",
				Data: []mapstr.MapStr{
					{
						"cc1":        10,
						"cc2":        "aaa",
						"bk_inst_id": instanceID + 3256,
					},
					{
						"cc1":        20,
						"cc2":        "bbb",
						"bk_inst_id": instanceID + 3256,
					},
				},
			}
			rsp, err := modelQuoteClient.BatchCreateQuotedInstance(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		}()

		By("list model quote instance")
		func() {
			option := metadata.ListQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				CommonQueryOption: metadata.CommonQueryOption{
					Page: metadata.BasePage{
						Limit: 100,
					},
				},
			}
			rsp, err := modelQuoteClient.ListQuotedInstance(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Info).To(HaveLen(2))
			for _, inst := range rsp.Info {
				instID, _ := inst.Get(common.BKInstIDField)
				instID64, err := util.GetInt64ByInterface(instID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(instID64).To(Equal(instanceID))
			}
		}()

		By("update model quote instance")
		func() {
			option := metadata.BatchUpdateQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				IDs:        quoteInstIDArr,
				Data: mapstr.MapStr{
					"cc1":        30,
					"cc2":        "ccc",
					"bk_inst_id": instanceID,
				},
			}
			err := modelQuoteClient.BatchUpdateQuotedInstance(ctx, header, &option)
			Expect(err).ShouldNot(HaveOccurred())
		}()

		By("update model quote instance but model property no exist")
		func() {
			option := metadata.BatchUpdateQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: "prop_no_exist",
				IDs:        quoteInstIDArr,
				Data: mapstr.MapStr{
					"cc1":        30,
					"cc2":        "ccc",
					"bk_inst_id": instanceID,
				},
			}
			err := modelQuoteClient.BatchUpdateQuotedInstance(ctx, header, &option)
			Expect(err).To(HaveOccurred())
		}()

		By("list model quote instance")
		func() {
			option := metadata.ListQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				CommonQueryOption: metadata.CommonQueryOption{
					Page: metadata.BasePage{
						Limit: 100,
					},
				},
			}
			rsp, err := modelQuoteClient.ListQuotedInstance(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Info).To(HaveLen(2))
			for _, inst := range rsp.Info {
				instID, _ := inst.Get(common.BKInstIDField)
				instIDInt64, err := util.GetInt64ByInterface(instID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(instIDInt64).To(Equal(instanceID))
				cc1Field, _ := inst.Get("cc1")
				cc1FieldInt, err := util.GetIntByInterface(cc1Field)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(cc1FieldInt).To(Equal(30))
				cc2Field, _ := inst.Get("cc2")
				Expect(cc2Field.(string)).To(Equal("ccc"))
			}
		}()

		By("delete model quote instance")
		func() {
			option := metadata.BatchDeleteQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				IDs:        quoteInstIDArr,
			}
			err := modelQuoteClient.BatchDeleteQuotedInstance(ctx, header, &option)
			Expect(err).ShouldNot(HaveOccurred())
		}()

		By("list model quote instance")
		func() {
			option := metadata.ListQuotedInstOption{
				ObjID:      modelQuoteObj.ObjectID,
				PropertyID: quoteProp.PropertyID,
				CommonQueryOption: metadata.CommonQueryOption{
					Page: metadata.BasePage{
						Limit: 100,
					},
				},
			}
			rsp, err := modelQuoteClient.ListQuotedInstance(ctx, header, &option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(rsp.Info).To(HaveLen(0))
		}()
	})
})
