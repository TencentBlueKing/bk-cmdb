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
	"encoding/json"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/querybuilder"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host relation test", func() {
	ctx := context.Background()
	test.DeleteAllHosts()
	var hostID1, hostID2, hostID3, hostID4, hostID5, defaultCloudID, cloudID int64
	var idleHostID1, idleHostID2 int64
	var idleSetID, idleModuleID, recycleModuleID int64
	var categoryID, resDirID int64
	var hostTestBizID, hostTestSetID, hostTestModuleID int64
	var serviceTemplateID, setTemplateID int64

	It("host relation test", func() {
		By("add biz")
		func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "hostRelationTest",
				"time_zone":         "Asia/Shanghai",
			}
			rsp, err := apiServerClient.CreateBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			hostTestBizID, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
			recycleModuleID = test.GetBizRecycleModule(hostTestBizID)
			idleModuleID = test.GetBizIdleModule(hostTestBizID)
		}()

		By("create cloud area")
		func() {
			resp, err := hostServerClient.CreateCloudArea(context.Background(), header, map[string]interface{}{
				common.BKCloudNameField:     "hostRelationTestArea",
				common.BKProjectStatusField: "1",
				common.BKCloudVendor:        "1",
			})
			util.RegisterResponseWithRid(resp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			defaultCloudID = int64(resp.Data.Created.ID)
		}()

		By("add hosts to resource pool default module")
		func() {
			input := metadata.AddHostToResourcePoolHostList{
				HostInfo: []map[string]interface{}{
					{
						common.BKHostInnerIPField:   "20.0.0.1",
						common.BKHostNameField:      "host1",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           2,
						common.BKHostInnerIPv6Field: "::1",
					},
					{
						common.BKHostInnerIPField:   "20.0.0.2",
						common.BKHostNameField:      "host2",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           4,
						common.BKHostInnerIPv6Field: "::2",
					},
					{
						common.BKHostInnerIPField:   "20.0.0.3",
						common.BKHostNameField:      "host3",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           6,
						common.BKHostInnerIPv6Field: "::3",
					},
					{
						common.BKHostInnerIPField:   "20.0.0.4",
						common.BKHostNameField:      "host4",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           8,
						common.BKHostInnerIPv6Field: "::4",
					},
					{
						common.BKHostInnerIPField:   "20.0.0.5",
						common.BKHostNameField:      "host5",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           10,
						common.BKHostInnerIPv6Field: "::5",
					},
				},
			}
			hostRsp, err := hostServerClient.AddHostToResourcePool(ctx, header, input)
			util.RegisterResponseWithRid(hostRsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())
			js, err := json.Marshal(hostRsp.Data)
			Expect(err).NotTo(HaveOccurred())
			result := metadata.AddHostToResourcePoolResult{}
			err = json.Unmarshal(js, &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result.Success)).To(Equal(5))
			Expect(len(result.Error)).To(Equal(0))
			hostID1 = result.Success[0].HostID
			hostID2 = result.Success[1].HostID
			hostID3 = result.Success[2].HostID
			hostID4 = result.Success[3].HostID
			hostID5 = result.Success[4].HostID
		}()

		By("add hosts to business idle module")
		func() {
			input := metadata.HostListParam{
				ApplicationID: hostTestBizID,
				HostList: []mapstr.MapStr{
					{
						common.BKHostInnerIPField:   "20.0.0.6",
						common.BKHostNameField:      "host6",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           2,
						common.BKHostInnerIPv6Field: "::6",
					},
					{
						common.BKHostInnerIPField:   "20.0.0.7",
						common.BKHostNameField:      "host7",
						common.BKCloudIDField:       defaultCloudID,
						common.BKAddressingField:    common.BKAddressingStatic,
						common.BKCpuField:           4,
						common.BKHostInnerIPv6Field: "::7",
					},
				},
			}
			rsp, err := hostServerClient.AddHostToBizIdle(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.HostIDs)).To(Equal(2))
			idleHostID1 = rsp.HostIDs[0]
			idleHostID2 = rsp.HostIDs[1]
		}()

		By("find host module relation test")
		func() {
			input := metadata.HostModuleRelationParameter{
				AppID:  hostTestBizID,
				HostID: []int64{idleHostID1, idleHostID2},
			}
			rsp, err := hostServerClient.FindHostModules(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data)).To(Equal(2))

			Expect(rsp.Data[0].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Data[0].HostID).Should(BeElementOf(idleHostID1, idleHostID2))
			Expect(rsp.Data[0].SetID).To(Equal(rsp.Data[1].SetID))
			Expect(rsp.Data[0].ModuleID).To(Equal(idleModuleID))
			idleSetID = rsp.Data[0].SetID
		}()

		By("count biz host CPU test1")
		func() {
			input := metadata.CountHostCPUReq{
				BizID: hostTestBizID,
			}
			rsp, err := hostServerClient.CountBizHostCPU(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp)).To(Equal(1))
			Expect(rsp[0].BizID).To(Equal(hostTestBizID))
			Expect(rsp[0].CpuCount).To(Equal(int64(6)))
			Expect(rsp[0].HostCount).To(Equal(int64(2)))
			Expect(rsp[0].NoCpuHostCount).To(Equal(int64(0)))
		}()

		By("updatemany hosts all property test1")
		func() {
			input := metadata.UpdateHostOpt{
				Update: []metadata.UpdateHost{
					{
						HostIDs: []int64{idleHostID1, idleHostID2},
						Properties: map[string]interface{}{
							common.BKCpuField:       10,
							common.BKOSTypeField:    "1",
							common.BKOSNameField:    "cc_os1",
							common.BKOSVersionField: "v666",
							common.BKOSBitField:     "32",
							common.BKMemField:       2048,
							common.BKDiskField:      500,
							common.BKCpuArch:        "x86",
						},
					},
				},
			}
			err := hostServerClient.UpdateHostsAllProperty(ctx, header, &input)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("count biz host CPU test2")
		func() {
			input := metadata.CountHostCPUReq{
				BizID: hostTestBizID,
			}
			rsp, err := hostServerClient.CountBizHostCPU(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp)).To(Equal(1))
			Expect(rsp[0].BizID).To(Equal(hostTestBizID))
			Expect(rsp[0].CpuCount).To(Equal(int64(20)))
			Expect(rsp[0].HostCount).To(Equal(int64(2)))
			Expect(rsp[0].NoCpuHostCount).To(Equal(int64(0)))
		}()

		By("list hosts without biz test")
		func() {
			input := metadata.ListHostsWithNoBizParameter{
				HostPropertyFilter: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							querybuilder.AtomRule{
								Field:    common.BKCloudIDField,
								Operator: querybuilder.OperatorEqual,
								Value:    defaultCloudID,
							},
							querybuilder.AtomRule{
								Field:    common.BKHostNameField,
								Operator: querybuilder.OperatorIn,
								Value:    []string{"host6", "host7"},
							},
						},
					},
				},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKHostNameField, common.BKCloudIDField,
					common.BKHostInnerIPField, common.BKCpuField, common.BKOSTypeField, common.BKOSNameField,
					common.BKOSVersionField, common.BKOSBitField, common.BKMemField, common.BKDiskField, common.BKCpuArch},
			}
			rsp, err := hostServerClient.ListHostsWithoutApp(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.Info)).To(Equal(2))
			hostCloudID, err1 := commonutil.GetInt64ByInterface(rsp.Info[0][common.BKCloudIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostCloudID).To(Equal(defaultCloudID))
			hostName := commonutil.GetStrByInterface(rsp.Info[0][common.BKHostNameField])
			Expect(hostName).To(Equal("host6"))
			hostOSType := commonutil.GetStrByInterface(rsp.Info[0][common.BKOSTypeField])
			Expect(hostOSType).To(Equal("1"))
			hostOSName := commonutil.GetStrByInterface(rsp.Info[1][common.BKOSNameField])
			Expect(hostOSName).To(Equal("cc_os1"))
			hostOSVersion := commonutil.GetStrByInterface(rsp.Info[1][common.BKOSVersionField])
			Expect(hostOSVersion).To(Equal("v666"))
			hostOSBit := commonutil.GetStrByInterface(rsp.Info[0][common.BKOSBitField])
			Expect(hostOSBit).To(Equal("32"))
			hostCPUArch := commonutil.GetStrByInterface(rsp.Info[0][common.BKCpuArch])
			Expect(hostCPUArch).To(Equal("x86"))

			hostMen, err1 := commonutil.GetIntByInterface(rsp.Info[1][common.BKMemField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostMen).To(Equal(2048))
			hostDisk, err1 := commonutil.GetIntByInterface(rsp.Info[1][common.BKDiskField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostDisk).To(Equal(500))
		}()

		By("list resource pool hosts test1")
		func() {
			input := metadata.ListHostsParameter{
				HostPropertyFilter: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							querybuilder.AtomRule{
								Field:    common.BKCloudIDField,
								Operator: querybuilder.OperatorEqual,
								Value:    defaultCloudID,
							},
							querybuilder.AtomRule{
								Field:    common.BKHostInnerIPField,
								Operator: querybuilder.OperatorContains,
								Value:    "20.0.0",
							},
						},
					},
				},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKHostNameField, common.BKCloudIDField,
					common.BKHostInnerIPField, common.BKCpuField},
			}
			rsp, err := hostServerClient.ListResourcePoolHosts(ctx, header, &input)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.Info)).To(Equal(5))
			for i, host := range rsp.Info {
				cloudID, err1 := commonutil.GetInt64ByInterface(host[common.BKCloudIDField])
				Expect(err1).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(defaultCloudID))
				hostName := commonutil.GetStrByInterface(host[common.BKHostNameField])
				Expect(hostName).To(Equal("host" + strconv.Itoa(i+1)))
				hostIP := commonutil.GetStrByInterface(host[common.BKHostInnerIPField])
				Expect(hostIP).To(Equal("20.0.0." + strconv.Itoa(i+1)))
				hostCpu, err := commonutil.GetIntByInterface(host[common.BKCpuField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostCpu).To(Equal(2 * (i + 1)))
			}
		}()

		By("list resource pool hosts test2")
		func() {
			input := metadata.ListHostsParameter{
				HostPropertyFilter: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							querybuilder.AtomRule{
								Field:    common.BKHostInnerIPField,
								Operator: querybuilder.OperatorContains,
								Value:    "100.100.100",
							},
						},
					},
				},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKHostNameField, common.BKCloudIDField,
					common.BKHostInnerIPField, common.BKCpuField},
			}
			rsp, err := hostServerClient.ListResourcePoolHosts(ctx, header, &input)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.Info)).To(Equal(0))
		}()

		By("updatemany hosts all property test2")
		func() {
			input := metadata.UpdateHostOpt{
				Update: []metadata.UpdateHost{
					{
						HostIDs: []int64{hostID1, hostID2},
						Properties: map[string]interface{}{
							common.BKCpuField:       10,
							common.BKOSTypeField:    "2",
							common.BKOSNameField:    "cc_os2",
							common.BKOSVersionField: "v888",
							common.BKOSBitField:     "32",
							common.BKMemField:       2048,
							common.BKDiskField:      500,
							common.BKCpuArch:        "arm",
						},
					},
					{
						HostIDs: []int64{hostID3},
						Properties: map[string]interface{}{
							common.BKOSNameField: "cc_os2",
						},
					},
				},
			}
			err := hostServerClient.UpdateHostsAllProperty(ctx, header, &input)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("create cloud area")
		func() {
			resp, err := hostServerClient.CreateCloudArea(context.Background(), header, map[string]interface{}{
				common.BKCloudNameField:     "Area001",
				common.BKProjectStatusField: "1",
				common.BKCloudVendor:        "1",
			})
			util.RegisterResponseWithRid(resp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			cloudID = int64(resp.Data.Created.ID)
		}()

		By("update hosts cloud area test")
		func() {
			input := metadata.UpdateHostCloudAreaFieldOption{
				HostIDs: []int64{hostID2, hostID3},
				CloudID: cloudID,
			}
			err := hostServerClient.UpdateHostCloudArea(ctx, header, &input)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("list resource pool hosts test3")
		func() {
			input := metadata.ListHostsParameter{
				HostPropertyFilter: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							querybuilder.AtomRule{
								Field:    common.BKOSNameField,
								Operator: querybuilder.OperatorEqual,
								Value:    "cc_os2",
							},
						},
					},
				},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKHostNameField, common.BKCloudIDField, common.BKCpuArch},
			}
			rsp, err := hostServerClient.ListResourcePoolHosts(ctx, header, &input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.Info)).To(Equal(3))
			Expect(rsp.Count).To(Equal(3))
			hostName := commonutil.GetStrByInterface(rsp.Info[0][common.BKHostNameField])
			Expect(hostName).To(Equal("host1"))
			hostCPUArch := commonutil.GetStrByInterface(rsp.Info[1][common.BKCpuArch])
			Expect(hostCPUArch).To(Equal("arm"))

			hostID, err1 := commonutil.GetInt64ByInterface(rsp.Info[2][common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID3))
			hostCloudID, err1 := commonutil.GetInt64ByInterface(rsp.Info[2][common.BKCloudIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostCloudID).To(Equal(cloudID))
		}()

		By("create resource pool directory")
		func() {
			dir := map[string]interface{}{
				common.BKModuleNameField: "hostTestDir",
			}
			rsp, err := dirClient.CreateResourceDirectory(context.Background(), header, dir)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			resDirID = int64(rsp.Data.Created.ID)
		}()

		By("transfer resource pool directory test")
		func() {
			opt := metadata.TransferHostResourceDirectory{
				ModuleID: resDirID,
				HostID:   []int64{hostID4, hostID5},
			}
			err := hostServerClient.TransferHostResourceDirectory(context.Background(), header, &opt)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("find module host relation test")
		func() {
			opt := metadata.FindModuleHostRelationParameter{
				ModuleIDS:    []int64{resDirID},
				ModuleFields: []string{common.BKModuleIDField, common.BKModuleNameField},
				HostFields:   []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindModuleHostRelation(context.Background(), header, test.GetResBizID(), &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.Relation)).To(Equal(2))

			hostID, err1 := commonutil.GetInt64ByInterface(rsp.Relation[0].Host[common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostID).To(Equal(hostID4))
			hostName := commonutil.GetStrByInterface(rsp.Relation[0].Host[common.BKHostNameField])
			Expect(hostName).To(Equal("host4"))

			hostIDa, err1 := commonutil.GetInt64ByInterface(rsp.Relation[1].Host[common.BKHostIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(hostIDa).To(Equal(hostID5))
			hostNamea := commonutil.GetStrByInterface(rsp.Relation[1].Host[common.BKHostNameField])
			Expect(hostNamea).To(Equal("host5"))

			moduleID, err1 := commonutil.GetInt64ByInterface(rsp.Relation[0].Modules[0][common.BKModuleIDField])
			Expect(err1).NotTo(HaveOccurred())
			Expect(moduleID).To(Equal(resDirID))
			moduleName := commonutil.GetStrByInterface(rsp.Relation[0].Modules[0][common.BKModuleNameField])
			Expect(moduleName).To(Equal("hostTestDir"))
		}()

		By("transfer host to recycle module test")
		func() {
			opt := metadata.DefaultModuleHostConfigParams{
				ApplicationID: hostTestBizID,
				HostIDs:       []int64{idleHostID1, idleHostID2},
			}
			err := hostServerClient.UpdateHostToRecycle(context.Background(), header, &opt)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("transfer host to recycle module with wrong biz test")
		func() {
			opt := metadata.DefaultModuleHostConfigParams{
				ApplicationID: hostTestBizID,
				HostIDs:       []int64{hostID4, hostID5},
			}
			err := hostServerClient.UpdateHostToRecycle(context.Background(), header, &opt)
			Expect(err).To(HaveOccurred())
			Expect(err.GetCode()).To(Equal(common.CCErrCoreServiceHostNotBelongBusiness))
		}()

		By("find host topo relation test")
		func() {
			opt := metadata.HostModuleRelationRequest{
				ApplicationID: hostTestBizID,
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindHostTopoRelation(context.Background(), header, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(int64(2)))
			Expect(len(rsp.Info)).To(Equal(2))

			Expect(rsp.Info[0].HostID).To(Equal(idleHostID1))
			Expect(rsp.Info[0].ModuleID).To(Equal(recycleModuleID))

			Expect(rsp.Info[1].HostID).To(Equal(idleHostID2))
			Expect(rsp.Info[1].ModuleID).To(Equal(recycleModuleID))
		}()

		By("create service category")
		func() {
			input := map[string]interface{}{
				"bk_parent_id":      1,
				common.BKAppIDField: hostTestBizID,
				"name":              "hostTest",
			}
			rsp, err := procServerClient.CreateServiceCategory(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceCategory{}
			err = json.Unmarshal(j, &data)
			Expect(err).NotTo(HaveOccurred())
			categoryID = data.ID
		}()

		By("create service template")
		func() {
			input := map[string]interface{}{
				common.BKServiceCategoryIDField: categoryID,
				common.BKAppIDField:             hostTestBizID,
				common.BKFieldName:              "hostTestBizServTmp",
			}
			rsp, err := procServerClient.CreateServiceTemplate(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp.Data)
			Expect(err).NotTo(HaveOccurred())
			data := metadata.ServiceTemplate{}
			err = json.Unmarshal(j, &data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Name).To(Equal("hostTestBizServTmp"))
			Expect(data.ServiceCategoryID).To(Equal(categoryID))
			serviceTemplateID = data.ID
		}()

		By("create set template")
		func() {
			option := metadata.CreateSetTemplateOption{
				Name:               "hostTestBizSetTmp",
				ServiceTemplateIDs: []int64{serviceTemplateID},
			}
			rsp, err := setTmpClient.CreateSetTemplate(ctx, header, hostTestBizID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Name).To(Equal("hostTestBizSetTmp"))
			Expect(rsp.Data.ID).To(Not(Equal(int64(0))))
			setTemplateID = rsp.Data.ID
		}()

		By("create set")
		func() {
			input := mapstr.MapStr{
				common.BKSetNameField:       "hostRelationTest",
				common.BKParentIDField:      hostTestBizID,
				common.BKAppIDField:         hostTestBizID,
				common.BKSetStatusField:     "1",
				common.BKSetEnvField:        "3",
				common.BKSetTemplateIDField: setTemplateID,
			}
			rsp, e := instClient.CreateSet(context.Background(), hostTestBizID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			Expect(commonutil.GetStrByInterface(rsp[common.BKSetNameField])).To(Equal("hostRelationTest"))

			parentIdRes, err := commonutil.GetInt64ByInterface(rsp[common.BKParentIDField])
			Expect(err).NotTo(HaveOccurred())
			Expect(parentIdRes).To(Equal(hostTestBizID))
			bizIdRes, err := commonutil.GetInt64ByInterface(rsp[common.BKAppIDField])
			Expect(err).NotTo(HaveOccurred())
			Expect(bizIdRes).To(Equal(hostTestBizID))

			hostTestSetID, err = commonutil.GetInt64ByInterface(rsp[common.BKSetIDField])
			Expect(err).NotTo(HaveOccurred())
		}()

		By("transfer hosts to biz idle module")
		func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: hostTestBizID,
				HostIDs:       []int64{hostID1, hostID2, hostID3, hostID4, hostID5},
				ModuleID:      idleModuleID,
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("find host by topo inst test1")
		func() {
			opt := metadata.FindHostsByTopoOpt{
				ObjID:  common.BKInnerObjIDSet,
				InstID: hostTestSetID,
				Fields: []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
				},
			}
			rsp, err := hostServerClient.FindHostByTopoInst(context.Background(), header, hostTestBizID, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(0))
			Expect(len(rsp.Info)).To(Equal(0))
		}()

		By("find host by service template test1")
		func() {
			opt := metadata.FindHostsBySrvTplOpt{
				ServiceTemplateIDs: []int64{serviceTemplateID},
				Fields:             []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindHostByServiceTmpl(context.Background(), header, hostTestBizID, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(0))
			Expect(len(rsp.Info)).To(Equal(0))
		}()

		By("find host by set template test1")
		func() {
			opt := metadata.FindHostsBySetTplOpt{
				SetTemplateIDs: []int64{setTemplateID},
				Fields:         []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindHostBySetTmpl(context.Background(), header, hostTestBizID, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(0))
			Expect(len(rsp.Info)).To(Equal(0))
		}()

		By("search module")
		func() {
			input := &params.SearchParams{
				Condition: map[string]interface{}{},
				Page: map[string]interface{}{
					"sort": "id",
				},
			}
			rsp, err := instClient.SearchModule(context.Background(), hostTestBizID, hostTestSetID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())

			Expect(rsp.Count).To(Equal(1))
			bizID, err1 := rsp.Info[0].Int64(common.BKAppIDField)
			Expect(err1).NotTo(HaveOccurred())
			Expect(bizID).To(Equal(hostTestBizID))
			setID, err1 := rsp.Info[0].Int64(common.BKSetIDField)
			Expect(err1).NotTo(HaveOccurred())
			Expect(setID).To(Equal(hostTestSetID))

			Expect(rsp.Info[0].String(common.BKModuleNameField)).To(Equal("hostTestBizServTmp"))
			hostTestModuleID, err1 = rsp.Info[0].Int64(common.BKModuleIDField)
			Expect(err1).NotTo(HaveOccurred())
		}()

		By("transfer host module")
		func() {
			input := map[string]interface{}{
				common.BKAppIDField:       hostTestBizID,
				common.BKHostIDField:      []int64{hostID1, hostID2, hostID3, hostID4, hostID5},
				common.BKModuleIDField:    []int64{hostTestModuleID},
				common.BKIsIncrementField: false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()

		By("find host by topo inst test2")
		func() {
			opt := metadata.FindHostsByTopoOpt{
				ObjID:  common.BKInnerObjIDSet,
				InstID: hostTestSetID,
				Fields: []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindHostByTopoInst(context.Background(), header, hostTestBizID, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(5))
			Expect(len(rsp.Info)).To(Equal(5))
			for i, host := range rsp.Info {
				hostName := commonutil.GetStrByInterface(host[common.BKHostNameField])
				Expect(hostName).To(Equal("host" + strconv.Itoa(i+1)))
			}
		}()

		By("find host by service template test2")
		func() {
			opt := metadata.FindHostsBySrvTplOpt{
				ServiceTemplateIDs: []int64{serviceTemplateID},
				Fields:             []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindHostByServiceTmpl(context.Background(), header, hostTestBizID, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(5))
			Expect(len(rsp.Info)).To(Equal(5))
			for i, host := range rsp.Info {
				hostName := commonutil.GetStrByInterface(host[common.BKHostNameField])
				Expect(hostName).To(Equal("host" + strconv.Itoa(i+1)))
			}
		}()

		By("find host by set template test2")
		func() {
			opt := metadata.FindHostsBySetTplOpt{
				SetTemplateIDs: []int64{setTemplateID},
				Fields:         []string{common.BKHostIDField, common.BKHostNameField},
				Page: metadata.BasePage{
					Limit: 10,
					Start: 0,
					Sort:  common.BKHostIDField,
				},
			}
			rsp, err := hostServerClient.FindHostBySetTmpl(context.Background(), header, hostTestBizID, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(5))
			Expect(len(rsp.Info)).To(Equal(5))
			for i, host := range rsp.Info {
				hostName := commonutil.GetStrByInterface(host[common.BKHostNameField])
				Expect(hostName).To(Equal("host" + strconv.Itoa(i+1)))
			}
		}()

		// 该接口需调用catchservice相关逻辑，当前catchservice未改造，暂时注释，后续可放开
		//By("find host detail topo test")
		//func() {
		//	opt := metadata.ListHostsDetailAndTopoOption{
		//		HostPropertyFilter: &querybuilder.QueryFilter{
		//			Rule: querybuilder.CombinedRule{
		//				Condition: querybuilder.ConditionAnd,
		//				Rules: []querybuilder.Rule{
		//					querybuilder.AtomRule{
		//						Field:    common.BKHostIDField,
		//						Operator: querybuilder.OperatorEqual,
		//						Value:    hostID1,
		//					},
		//				},
		//			},
		//		},
		//		Page: metadata.BasePage{
		//			Sort:  common.BKHostIDField,
		//			Limit: 10,
		//			Start: 0,
		//		},
		//		Fields: []string{common.BKHostIDField, common.BKHostNameField},
		//	}
		//	rsp, err := hostServerClient.FindHostDetailTopo(context.Background(), header, &opt)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(rsp.Count).To(Equal(1))
		//	Expect(len(rsp.Info)).To(Equal(1))
		//
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Host[common.BKHostNameField])).To(Equal("host1"))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Instance.Object)).To(Equal(common.BKInnerObjIDSet))
		//	instID, err1 := commonutil.GetInt64ByInterface(rsp.Info[0].Topo[0].Instance.InstID)
		//	Expect(err1).NotTo(HaveOccurred())
		//	Expect(instID).To(Equal(hostTestSetID))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Instance.InstName)).
		//		To(Equal("hostRelationTest"))
		//
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Children[0].Instance.Object)).
		//		To(Equal(common.BKInnerObjIDModule))
		//	instID, err1 = commonutil.GetInt64ByInterface(rsp.Info[0].Topo[0].Children[0].Instance.InstID)
		//	Expect(err1).NotTo(HaveOccurred())
		//	Expect(instID).To(Equal(hostTestModuleID))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Children[0].Instance.InstName)).
		//		To(Equal("hostTestBizServTmp"))
		//}()

		By("find host relation with topo test1")
		func() {
			opt := metadata.FindHostRelationWtihTopoOpt{
				Business: hostTestBizID,
				ObjID:    common.BKInnerObjIDSet,
				InstIDs:  []int64{idleSetID},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKModuleIDField, common.BKSetIDField,
					common.BKAppIDField},
			}
			rsp, err := hostServerClient.FindHostRelationWithTopo(context.Background(), header, &opt)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(int64(2)))

			Expect(rsp.Info[0].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[0].SetID).To(Equal(idleSetID))
			Expect(rsp.Info[0].ModuleID).To(Equal(recycleModuleID))
			Expect(rsp.Info[0].HostID).To(Equal(idleHostID1))

			Expect(rsp.Info[1].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[1].SetID).To(Equal(idleSetID))
			Expect(rsp.Info[1].ModuleID).To(Equal(recycleModuleID))
			Expect(rsp.Info[1].HostID).To(Equal(idleHostID2))

		}()

		By("find host relation with topo without biz id test")
		func() {
			opt := metadata.FindHostRelationWtihTopoOpt{
				ObjID:   common.BKInnerObjIDSet,
				InstIDs: []int64{idleSetID},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKModuleIDField, common.BKSetIDField, common.BKAppIDField},
			}
			_, err := hostServerClient.FindHostRelationWithTopo(context.Background(), header, &opt)
			Expect(err).To(HaveOccurred())
			Expect(err.GetCode()).To(Equal(common.CCErrCommParamsInvalid))
		}()

		By("find host relation with topo test2")
		func() {
			opt := metadata.FindHostRelationWtihTopoOpt{
				Business: hostTestBizID,
				ObjID:    common.BKInnerObjIDModule,
				InstIDs:  []int64{hostTestModuleID},
				Page: metadata.BasePage{
					Sort:  common.BKHostIDField,
					Limit: 10,
					Start: 0,
				},
				Fields: []string{common.BKHostIDField, common.BKModuleIDField, common.BKSetIDField, common.BKAppIDField},
			}
			rsp, err := hostServerClient.FindHostRelationWithTopo(context.Background(), header, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(int64(5)))
			Expect(len(rsp.Info)).To(Equal(5))

			Expect(rsp.Info[0].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[0].SetID).To(Equal(hostTestSetID))
			Expect(rsp.Info[0].ModuleID).To(Equal(hostTestModuleID))
			Expect(rsp.Info[0].HostID).To(Equal(hostID1))

			Expect(rsp.Info[1].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[1].SetID).To(Equal(hostTestSetID))
			Expect(rsp.Info[1].ModuleID).To(Equal(hostTestModuleID))
			Expect(rsp.Info[1].HostID).To(Equal(hostID2))

			Expect(rsp.Info[2].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[2].SetID).To(Equal(hostTestSetID))
			Expect(rsp.Info[2].ModuleID).To(Equal(hostTestModuleID))
			Expect(rsp.Info[2].HostID).To(Equal(hostID3))

			Expect(rsp.Info[3].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[3].SetID).To(Equal(hostTestSetID))
			Expect(rsp.Info[3].ModuleID).To(Equal(hostTestModuleID))
			Expect(rsp.Info[3].HostID).To(Equal(hostID4))

			Expect(rsp.Info[4].AppID).To(Equal(hostTestBizID))
			Expect(rsp.Info[4].SetID).To(Equal(hostTestSetID))
			Expect(rsp.Info[4].ModuleID).To(Equal(hostTestModuleID))
			Expect(rsp.Info[4].HostID).To(Equal(hostID5))
		}()

		By("find host service template test")
		func() {
			opt := metadata.HostIDReq{
				HostIDs: []int64{hostID1, hostID2, idleHostID1},
			}
			rsp, err := hostServerClient.FindHostServiceTmpl(context.Background(), header, &opt)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data)).To(Equal(3))

			for _, hostSrvTmplRes := range rsp.Data {
				if hostSrvTmplRes.HostID == idleHostID1 {
					Expect(len(hostSrvTmplRes.SrvTmplIDs)).To(Equal(0))
					continue
				}
				Expect(len(hostSrvTmplRes.SrvTmplIDs)).To(Equal(1))
				Expect(hostSrvTmplRes.SrvTmplIDs[0]).To(Equal(serviceTemplateID))
			}
		}()

		// 该接口需调用catchservice相关逻辑，当前catchservice未改造，暂时注释，后续可放开
		//By("find host total mainline topo test1")
		//func() {
		//	opt := metadata.FindHostTotalTopo{
		//		HostPropertyFilter: &querybuilder.QueryFilter{
		//			Rule: querybuilder.CombinedRule{
		//				Condition: querybuilder.ConditionAnd,
		//				Rules: []querybuilder.Rule{
		//					querybuilder.AtomRule{
		//						Field:    common.BKHostIDField,
		//						Operator: querybuilder.OperatorEqual,
		//						Value:    idleHostID1,
		//					},
		//				},
		//			},
		//		},
		//		Page: metadata.BasePage{
		//			Sort:  common.BKHostIDField,
		//			Limit: 10,
		//			Start: 0,
		//		},
		//		Fields: []string{common.BKHostIDField, common.BKHostNameField},
		//	}
		//	rsp, err := hostServerClient.FindHostTotalMainlineTopo(context.Background(), header, hostTestBizID, &opt)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(rsp.Count).To(Equal(1))
		//	Expect(len(rsp.Info)).To(Equal(1))
		//
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Host[common.BKHostNameField])).To(Equal("host6"))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Instance.Object)).To(Equal(common.BKInnerObjIDSet))
		//	instID, err1 := commonutil.GetInt64ByInterface(rsp.Info[0].Topo[0].Instance.InstID)
		//	Expect(err1).NotTo(HaveOccurred())
		//	Expect(instID).To(Equal(idleSetID))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Instance.InstName)).
		//		To(Equal("空闲机池"))
		//
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Children[0].Instance.Object)).
		//		To(Equal(common.BKInnerObjIDModule))
		//	instID, err1 = commonutil.GetInt64ByInterface(rsp.Info[0].Topo[0].Children[0].Instance.InstID)
		//	Expect(err1).NotTo(HaveOccurred())
		//	Expect(instID).To(Equal(recycleModuleID))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Children[0].Instance.InstName)).
		//		To(Equal("待回收"))
		//}()

		// 该接口需调用catchservice相关逻辑，当前catchservice未改造，暂时注释，后续可放开
		//By("find host total mainline topo test2")
		//func() {
		//	opt := metadata.FindHostTotalTopo{
		//		SetPropertyFilter: &querybuilder.QueryFilter{
		//			Rule: querybuilder.CombinedRule{
		//				Condition: querybuilder.ConditionAnd,
		//				Rules: []querybuilder.Rule{
		//					querybuilder.AtomRule{
		//						Field:    common.BKSetNameField,
		//						Operator: querybuilder.OperatorEqual,
		//						Value:    "hostRelationTest",
		//					},
		//				},
		//			},
		//		},
		//		ModulePropertyFilter: &querybuilder.QueryFilter{
		//			Rule: querybuilder.CombinedRule{
		//				Condition: querybuilder.ConditionAnd,
		//				Rules: []querybuilder.Rule{
		//					querybuilder.AtomRule{
		//						Field:    common.BKModuleIDField,
		//						Operator: querybuilder.OperatorEqual,
		//						Value:    hostTestModuleID,
		//					},
		//				},
		//			},
		//		},
		//		Page: metadata.BasePage{
		//			Sort:  common.BKHostIDField,
		//			Limit: 10,
		//			Start: 0,
		//		},
		//		Fields: []string{common.BKHostIDField, common.BKHostNameField},
		//	}
		//	rsp, err := hostServerClient.FindHostTotalMainlineTopo(context.Background(), header, hostTestBizID, &opt)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(rsp.Count).To(Equal(5))
		//	Expect(len(rsp.Info)).To(Equal(5))
		//
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Host[common.BKHostNameField])).To(Equal("host1"))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Instance.Object)).To(Equal(common.BKInnerObjIDSet))
		//	instID, err1 := commonutil.GetInt64ByInterface(rsp.Info[0].Topo[0].Instance.InstID)
		//	Expect(err1).NotTo(HaveOccurred())
		//	Expect(instID).To(Equal(hostTestSetID))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Instance.InstName)).
		//		To(Equal("hostRelationTest"))
		//
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Children[0].Instance.Object)).
		//		To(Equal(common.BKInnerObjIDModule))
		//	instID, err1 = commonutil.GetInt64ByInterface(rsp.Info[0].Topo[0].Children[0].Instance.InstID)
		//	Expect(err1).NotTo(HaveOccurred())
		//	Expect(instID).To(Equal(hostTestModuleID))
		//	Expect(commonutil.GetStrByInterface(rsp.Info[0].Topo[0].Children[0].Instance.InstName)).
		//		To(Equal("hostTestBizServTmp"))
		//}()
	})
	test.DeleteAllHosts()
})
