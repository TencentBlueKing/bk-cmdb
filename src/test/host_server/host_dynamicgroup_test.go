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

package host_server_test

import (
	"context"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dynamic group test", func() {
	ctx := context.Background()
	var dynamicGroupTestBizID int64
	var hostID1, hostID2, hostID3, hostID4, hostID5, defaultCloudID int64
	var setID1, setID3 int64
	var groupID1, groupID2 string

	test.DeleteAllHosts()

	It("dynamic group object is the host test", func() {
		By("add biz")
		func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "dynamicGroupTest",
				"time_zone":         "Asia/Shanghai",
			}
			rsp, err := apiServerClient.CreateBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			dynamicGroupTestBizID, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		}()

		By("create cloud area")
		func() {
			resp, err := hostServerClient.CreateCloudArea(context.Background(), header, map[string]interface{}{
				common.BKCloudNameField:     "dynamicGroupTestArea",
				common.BKProjectStatusField: "1",
				common.BKCloudVendor:        "1",
			})
			util.RegisterResponseWithRid(resp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			defaultCloudID = int64(resp.Data.Created.ID)
		}()

		By("add hosts to business idle module")
		func() {
			input := metadata.HostListParam{
				ApplicationID: dynamicGroupTestBizID,
				HostList: []mapstr.MapStr{
					{
						common.BKHostInnerIPField:   "21.0.0.1",
						common.BKHostNameField:      "host1",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           2,
						common.BKOSTypeField:        "1",
						common.BKOSNameField:        "cc_os1",
						common.BKOSVersionField:     "v666",
						common.BKOSBitField:         "32",
						common.BKMemField:           1024,
						common.BKDiskField:          500,
						common.BKCpuArch:            "x86",
						common.BKHostInnerIPv6Field: "2001::1",
					},
					{
						common.BKHostInnerIPField:   "21.0.0.2",
						common.BKHostNameField:      "host2",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           4,
						common.BKOSTypeField:        "2",
						common.BKOSNameField:        "cc_os1",
						common.BKOSVersionField:     "v888",
						common.BKOSBitField:         "32",
						common.BKMemField:           2048,
						common.BKDiskField:          500,
						common.BKCpuArch:            "arm",
						common.BKHostInnerIPv6Field: "2001::2",
					},
					{
						common.BKHostInnerIPField:   "21.0.0.3",
						common.BKHostNameField:      "host3",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           6,
						common.BKOSTypeField:        "3",
						common.BKOSNameField:        "cc_os2",
						common.BKOSVersionField:     "v666",
						common.BKOSBitField:         "64",
						common.BKMemField:           2048,
						common.BKDiskField:          500,
						common.BKCpuArch:            "x86",
						common.BKHostInnerIPv6Field: "2001::3",
					},
					{
						common.BKHostInnerIPField:   "21.0.0.4",
						common.BKHostNameField:      "host4",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           8,
						common.BKOSTypeField:        "4",
						common.BKOSNameField:        "cc_os2",
						common.BKOSVersionField:     "v888",
						common.BKOSBitField:         "64",
						common.BKMemField:           4096,
						common.BKDiskField:          1000,
						common.BKCpuArch:            "arm",
						common.BKHostInnerIPv6Field: "2001::4",
					},
					{
						common.BKHostInnerIPField:   "21.0.0.5",
						common.BKHostNameField:      "host5",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           10,
						common.BKOSTypeField:        "5",
						common.BKOSNameField:        "cc_os2",
						common.BKOSVersionField:     "v666",
						common.BKOSBitField:         "64",
						common.BKMemField:           4096,
						common.BKDiskField:          1000,
						common.BKCpuArch:            "x86",
						common.BKHostInnerIPv6Field: "2001::5",
					},
				},
			}
			rsp, err := hostServerClient.AddHostToBizIdle(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.HostIDs)).To(Equal(5))
			hostID1 = rsp.HostIDs[0]
			hostID2 = rsp.HostIDs[1]
			hostID3 = rsp.HostIDs[2]
			hostID4 = rsp.HostIDs[3]
			hostID5 = rsp.HostIDs[4]
		}()

		By("create dynamic group test1")
		func() {
			input := metadata.DynamicGroup{
				AppID: dynamicGroupTestBizID,
				ObjID: common.BKInnerObjIDHost,
				Name:  "cc_group1",
				Info: metadata.DynamicGroupInfo{
					Condition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDHost,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKCloudIDField,
									Operator: common.BKDBEQ,
									Value:    defaultCloudID,
								},
							},
						},
					},
					VariableCondition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDHost,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKOSNameField,
									Operator: common.BKDBEQ,
									Value:    "cc_os2",
								},
								{
									Field:    common.BKCpuArch,
									Operator: common.BKDBIN,
									Value:    []string{"x86", "arm"},
								},
							},
						},
					},
				},
			}
			rsp, err := hostServerClient.CreateDynamicGroup(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			groupID1 = rsp.Data.ID
		}()

		By("execute dynamic group 'groupID1' test1")
		func() {
			input := metadata.ExecuteOption{
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKHostNameField, common.BKHostInnerIPField},
			}
			rsp, err := hostServerClient.ExecuteDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID1, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
			Expect(len(rsp.Data.Info)).To(Equal(3))

			hostID, err1 := commonutil.GetInt64ByInterface(rsp.Data.Info[0][common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID3))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[0][common.BKHostNameField])).To(Equal("host3"))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[0][common.BKHostInnerIPField])).To(Equal("21.0.0.3"))

			hostID, err1 = commonutil.GetInt64ByInterface(rsp.Data.Info[1][common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID4))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[1][common.BKHostNameField])).To(Equal("host4"))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[1][common.BKHostInnerIPField])).To(Equal("21.0.0.4"))

			hostID, err1 = commonutil.GetInt64ByInterface(rsp.Data.Info[2][common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID5))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[2][common.BKHostNameField])).To(Equal("host5"))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[2][common.BKHostInnerIPField])).To(Equal("21.0.0.5"))
		}()

		By("execute dynamic group 'groupID1' test2")
		func() {
			input := metadata.ExecuteOption{
				VariableCondition: []metadata.DynamicGroupInfoCondition{
					{
						ObjID: common.BKInnerObjIDHost,
						Condition: []metadata.DynamicGroupCondition{
							{
								Field:    common.BKOSNameField,
								Operator: "contains",
								Value:    "os1",
							},
						},
					},
				},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKHostNameField, common.BKHostInnerIPField},
			}
			rsp, err := hostServerClient.ExecuteDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID1, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			Expect(len(rsp.Data.Info)).To(Equal(2))

			hostID, err1 := commonutil.GetInt64ByInterface(rsp.Data.Info[0][common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID1))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[0][common.BKHostNameField])).To(Equal("host1"))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[0][common.BKHostInnerIPField])).To(Equal("21.0.0.1"))

			hostID, err1 = commonutil.GetInt64ByInterface(rsp.Data.Info[1][common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID2))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[1][common.BKHostNameField])).To(Equal("host2"))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[1][common.BKHostInnerIPField])).To(Equal("21.0.0.2"))
		}()

		By("get dynamic group test1")
		func() {
			rsp, err := hostServerClient.GetDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID1, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			Expect(rsp.Data.Name).To(Equal("cc_group1"))
			Expect(rsp.Data.ObjID).To(Equal(common.BKInnerObjIDHost))
			Expect(len(rsp.Data.Info.Condition)).To(Equal(1))
			Expect(len(rsp.Data.Info.Condition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info.Condition[0].Condition[0].Field).To(Equal(common.BKCloudIDField))

			Expect(len(rsp.Data.Info.VariableCondition)).To(Equal(1))
			Expect(len(rsp.Data.Info.VariableCondition[0].Condition)).To(Equal(2))
			Expect(rsp.Data.Info.VariableCondition[0].Condition[0].Field).Should(BeElementOf(common.BKOSNameField,
				common.BKCpuArch))
		}()

		By("update dynamic group test")
		func() {
			input := metadata.DynamicGroup{
				AppID: dynamicGroupTestBizID,
				ObjID: common.BKInnerObjIDHost,
				Name:  "cc_group1",
				Info: metadata.DynamicGroupInfo{
					Condition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDModule,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKModuleNameField,
									Operator: common.BKDBIN,
									Value:    []string{"空闲机", "待回收", "故障机"},
								},
							},
						},
					},
					VariableCondition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDHost,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKOSVersionField,
									Operator: common.BKDBEQ,
									Value:    "v666",
								},
							},
						},
					},
				},
			}
			rsp, err := hostServerClient.UpdateDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID1, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("get dynamic group test2")
		func() {
			rsp, err := hostServerClient.GetDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID1, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			Expect(rsp.Data.AppID).To(Equal(dynamicGroupTestBizID))
			Expect(rsp.Data.Name).To(Equal("cc_group1"))
			Expect(rsp.Data.ObjID).To(Equal(common.BKInnerObjIDHost))
			Expect(len(rsp.Data.Info.Condition)).To(Equal(1))
			Expect(len(rsp.Data.Info.Condition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info.Condition[0].ObjID).To(Equal(common.BKInnerObjIDModule))
			Expect(rsp.Data.Info.Condition[0].Condition[0].Field).To(Equal(common.BKModuleNameField))

			Expect(len(rsp.Data.Info.VariableCondition)).To(Equal(1))
			Expect(len(rsp.Data.Info.VariableCondition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info.VariableCondition[0].ObjID).To(Equal(common.BKInnerObjIDHost))
			Expect(rsp.Data.Info.VariableCondition[0].Condition[0].Field).To(Equal(common.BKOSVersionField))
		}()

		By("search dynamic group test")
		func() {
			input := metadata.QueryCondition{
				Condition: map[string]interface{}{
					common.BKFieldName: "cc_group1",
				},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
				},
			}
			rsp, err := hostServerClient.SearchDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(uint64(1)))
			Expect(len(rsp.Data.Info)).To(Equal(1))

			Expect(rsp.Data.Info[0].Name).To(Equal("cc_group1"))
			Expect(rsp.Data.Info[0].ObjID).To(Equal(common.BKInnerObjIDHost))
			Expect(len(rsp.Data.Info[0].Info.Condition)).To(Equal(1))
			Expect(len(rsp.Data.Info[0].Info.Condition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info[0].Info.Condition[0].ObjID).To(Equal(common.BKInnerObjIDModule))
			Expect(rsp.Data.Info[0].Info.Condition[0].Condition[0].Field).To(Equal(common.BKModuleNameField))

			Expect(len(rsp.Data.Info[0].Info.VariableCondition)).To(Equal(1))
			Expect(len(rsp.Data.Info[0].Info.VariableCondition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info[0].Info.VariableCondition[0].ObjID).To(Equal(common.BKInnerObjIDHost))
			Expect(rsp.Data.Info[0].Info.VariableCondition[0].Condition[0].Field).To(Equal(common.BKOSVersionField))
		}()
	})

	It("dynamic group object is the set test", func() {
		By("create set1")
		func() {
			input := mapstr.MapStr{
				common.BKSetNameField:       "set1",
				common.BKParentIDField:      dynamicGroupTestBizID,
				common.BKAppIDField:         dynamicGroupTestBizID,
				common.BKSetStatusField:     "1",
				common.BKSetEnvField:        "1",
				common.BKSetTemplateIDField: 0,
				common.BKSetDescField:       "dynamicGroupTest",
			}
			rsp, err := instClient.CreateSet(context.Background(), dynamicGroupTestBizID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			setID, err1 := commonutil.GetInt64ByInterface(rsp[common.BKSetIDField])
			Expect(err1).NotTo(HaveOccurred())
			setID1 = setID
		}()

		By("create set2")
		func() {
			input := mapstr.MapStr{
				common.BKSetNameField:       "set2",
				common.BKParentIDField:      dynamicGroupTestBizID,
				common.BKAppIDField:         dynamicGroupTestBizID,
				common.BKSetStatusField:     "2",
				common.BKSetEnvField:        "2",
				common.BKSetTemplateIDField: 0,
				common.BKSetDescField:       "dynamicGroupTest",
			}
			rsp, err := instClient.CreateSet(context.Background(), dynamicGroupTestBizID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("create set3")
		func() {
			input := mapstr.MapStr{
				common.BKSetNameField:       "set3",
				common.BKParentIDField:      dynamicGroupTestBizID,
				common.BKAppIDField:         dynamicGroupTestBizID,
				common.BKSetStatusField:     "1",
				common.BKSetEnvField:        "3",
				common.BKSetTemplateIDField: 0,
				common.BKSetDescField:       "dynamicGroupTest",
			}
			rsp, err := instClient.CreateSet(context.Background(), dynamicGroupTestBizID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			setID, err1 := commonutil.GetInt64ByInterface(rsp[common.BKSetIDField])
			Expect(err1).NotTo(HaveOccurred())
			setID3 = setID
		}()

		By("create dynamic group test2")
		func() {
			input := metadata.DynamicGroup{
				AppID: dynamicGroupTestBizID,
				ObjID: common.BKInnerObjIDSet,
				Name:  "cc_group2",
				Info: metadata.DynamicGroupInfo{
					Condition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDSet,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKSetStatusField,
									Operator: common.BKDBEQ,
									Value:    "1",
								},
								{
									Field:    common.BKSetDescField,
									Operator: common.BKDBEQ,
									Value:    "dynamicGroupTest",
								},
							},
						},
					},
					VariableCondition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDSet,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKSetEnvField,
									Operator: common.BKDBIN,
									Value:    []string{"1", "2", "3"},
								},
							},
						},
					},
				},
			}
			rsp, err := hostServerClient.CreateDynamicGroup(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			groupID2 = rsp.Data.ID
		}()

		By("execute dynamic group 'groupID2' test")
		func() {
			input := metadata.ExecuteOption{
				Page: metadata.BasePage{
					Sort:  common.BKSetIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKSetIDField, common.BKSetNameField},
			}
			rsp, err := hostServerClient.ExecuteDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID2, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			Expect(len(rsp.Data.Info)).To(Equal(2))

			setID, err1 := commonutil.GetInt64ByInterface(rsp.Data.Info[0][common.BKSetIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(setID).To(Equal(setID1))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[0][common.BKSetNameField])).To(Equal("set1"))

			setID, err1 = commonutil.GetInt64ByInterface(rsp.Data.Info[1][common.BKSetIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(setID).To(Equal(setID3))
			Expect(commonutil.GetStrByInterface(rsp.Data.Info[1][common.BKSetNameField])).To(Equal("set3"))
		}()

		By("get dynamic group test1")
		func() {
			rsp, err := hostServerClient.GetDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID2, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			Expect(rsp.Data.Name).To(Equal("cc_group2"))
			Expect(rsp.Data.ObjID).To(Equal(common.BKInnerObjIDSet))
			Expect(len(rsp.Data.Info.Condition)).To(Equal(1))
			Expect(len(rsp.Data.Info.Condition[0].Condition)).To(Equal(2))
			Expect(rsp.Data.Info.Condition[0].Condition[0].Field).To(Equal(common.BKSetStatusField))

			Expect(len(rsp.Data.Info.VariableCondition)).To(Equal(1))
			Expect(len(rsp.Data.Info.VariableCondition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info.VariableCondition[0].Condition[0].Field).To(Equal(common.BKSetEnvField))
		}()

		By("update dynamic group test")
		func() {
			input := metadata.DynamicGroup{
				AppID: dynamicGroupTestBizID,
				ObjID: common.BKInnerObjIDSet,
				Name:  "cc_group2",
				Info: metadata.DynamicGroupInfo{
					Condition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDSet,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKSetEnvField,
									Operator: common.BKDBEQ,
									Value:    "3",
								},
							},
						},
					},
					VariableCondition: []metadata.DynamicGroupInfoCondition{
						{
							ObjID: common.BKInnerObjIDSet,
							Condition: []metadata.DynamicGroupCondition{
								{
									Field:    common.BKSetStatusField,
									Operator: common.BKDBEQ,
									Value:    "2",
								},
							},
						},
					},
				},
			}
			rsp, err := hostServerClient.UpdateDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID2, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("get dynamic group test2")
		func() {
			rsp, err := hostServerClient.GetDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID2, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			Expect(rsp.Data.Name).To(Equal("cc_group2"))
			Expect(rsp.Data.ObjID).To(Equal(common.BKInnerObjIDSet))
			Expect(len(rsp.Data.Info.Condition)).To(Equal(1))
			Expect(len(rsp.Data.Info.Condition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info.Condition[0].Condition[0].Field).To(Equal(common.BKSetEnvField))

			Expect(len(rsp.Data.Info.VariableCondition)).To(Equal(1))
			Expect(len(rsp.Data.Info.VariableCondition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info.VariableCondition[0].Condition[0].Field).To(Equal(common.BKSetStatusField))
		}()

		By("search dynamic group test1")
		func() {
			input := metadata.QueryCondition{
				Condition: map[string]interface{}{
					common.BKFieldName: "cc_group2",
				},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
				},
			}
			rsp, err := hostServerClient.SearchDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			Expect(rsp.Data.Info[0].Name).To(Equal("cc_group2"))
			Expect(rsp.Data.Info[0].ObjID).To(Equal(common.BKInnerObjIDSet))
			Expect(len(rsp.Data.Info[0].Info.Condition)).To(Equal(1))
			Expect(len(rsp.Data.Info[0].Info.Condition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info[0].Info.Condition[0].ObjID).To(Equal(common.BKInnerObjIDSet))
			Expect(rsp.Data.Info[0].Info.Condition[0].Condition[0].Field).To(Equal(common.BKSetEnvField))

			Expect(len(rsp.Data.Info[0].Info.VariableCondition)).To(Equal(1))
			Expect(len(rsp.Data.Info[0].Info.VariableCondition[0].Condition)).To(Equal(1))
			Expect(rsp.Data.Info[0].Info.VariableCondition[0].ObjID).To(Equal(common.BKInnerObjIDSet))
			Expect(rsp.Data.Info[0].Info.VariableCondition[0].Condition[0].Field).To(Equal(common.BKSetStatusField))
		}()

		By("search dynamic group test2")
		func() {
			input := metadata.QueryCondition{
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKFieldName,
				},
			}
			rsp, err := hostServerClient.SearchDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data.Info)).To(Equal(2))

			Expect(rsp.Data.Info[0].Name).To(Equal("cc_group1"))
			Expect(rsp.Data.Info[1].Name).To(Equal("cc_group2"))
		}()

		By("delete dynamic group test1")
		func() {
			rsp, err := hostServerClient.DeleteDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID1, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("delete dynamic group test2")
		func() {
			rsp, err := hostServerClient.DeleteDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				groupID2, header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("search dynamic group test3")
		func() {
			input := metadata.QueryCondition{
				Page: metadata.BasePage{
					Sort:  "-name",
					Limit: 10,
					Start: 0,
				},
			}
			rsp, err := hostServerClient.SearchDynamicGroup(ctx, strconv.FormatInt(dynamicGroupTestBizID, 10),
				header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data.Info)).To(Equal(0))
		}()
	})
	test.DeleteAllHosts()
})
