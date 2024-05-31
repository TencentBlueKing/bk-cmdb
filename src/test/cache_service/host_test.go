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

package cache_service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"
	"github.com/tidwall/gjson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/xid"
)

var _ = Describe("host cache test", func() {
	hosts := make([]mapstr.MapStr, 5)
	for i := 0; i < 5; i++ {
		hosts[i] = mapstr.MapStr{
			common.BKCloudIDField:     int64(0),
			common.BKHostInnerIPField: fmt.Sprintf("127.0.0.%d", i+1),
			common.BKAgentIDField:     xid.New().String(),
			common.BKAddressingField:  common.BKAddressingStatic,
			"bk_comment":              fmt.Sprintf("test%d", i+1),
		}
	}
	hostIDs := make([]int64, 5)

	It("list host cache with no host data", func() {
		By("clear data", func() {
			test.DeleteAllHosts()
			deleteAllFullSyncCond()
			time.Sleep(10 * time.Second)
		})

		By("list host", func() {
			opt := &metadata.ListHostWithPage{
				Page: metadata.BasePage{Limit: 5},
			}
			cnt, jsonStr, err := hostCli.ListHostWithPage(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(Equal(int64(0)))
			Expect(jsonStr).To(Equal("[]"))
		})
	})

	It("list host cache after adding host data", func() {
		By("add hosts", func() {
			addHostInfo := make(mapstr.MapStr)
			for i, host := range hosts {
				addHostInfo[strconv.Itoa(i)] = host
			}
			input := mapstr.MapStr{
				"host_info": addHostInfo,
			}
			rsp, err := hostSvrCli.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		By("list host with page", func() {
			opt := &metadata.ListHostWithPage{
				Page: metadata.BasePage{Limit: 3},
			}
			cnt, jsonStr, err := hostCli.ListHostWithPage(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(Equal(int64(len(hosts))))
			hostInfos := make([]mapstr.MapStr, 0)
			err = json.Unmarshal([]byte(jsonStr), &hostInfos)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(3))
			for i, info := range hostInfos {
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[i][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[i][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[i][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[i][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[i]["bk_comment"]))
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				hostIDs[i] = hostID
			}
		})

		By("list host with start and fields", func() {
			opt := &metadata.ListHostWithPage{
				Fields: []string{common.BKHostIDField, common.BKAgentIDField},
				Page:   metadata.BasePage{Start: 3, Limit: 2},
			}
			cnt, jsonStr, err := hostCli.ListHostWithPage(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(Equal(int64(len(hosts))))
			hostInfos := make([]mapstr.MapStr, 0)
			err = json.Unmarshal([]byte(jsonStr), &hostInfos)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(2))
			for i, info := range hostInfos {
				Expect(len(info)).To(Equal(2))
				Expect(info[common.BKAgentIDField] == hosts[i+3][common.BKAgentIDField]).To(Equal(true))
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				hostIDs[i+3] = hostID
			}
		})

		By("list host with ids", func() {
			opt := &metadata.ListHostWithPage{
				HostIDs: hostIDs,
			}
			cnt, jsonStr, err := hostCli.ListHostWithPage(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(Equal(int64(len(hosts))))
			hostInfos := make([]mapstr.MapStr, 0)
			err = json.Unmarshal([]byte(jsonStr), &hostInfos)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts)))
			for i, info := range hostInfos {
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[i][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[i][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[i][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[i][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[i]["bk_comment"]))
			}
		})
	})

	fullSyncCondIDs := make([]int64, 0)
	It("create full sync cond with no condition and list cache", func() {
		By("create full sync cond with no condition", func() {
			opt := &fullsynccond.CreateFullSyncCondOpt{
				Resource: general.Host,
				IsAll:    true,
				Interval: 6,
			}
			fullSyncCondID, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(fullSyncCondID, header)
			Expect(err).NotTo(HaveOccurred())
			fullSyncCondIDs = append(fullSyncCondIDs, fullSyncCondID)
		})

		By("list cache using full sync cond with no condition", func() {
			opt := &fullsynccond.ListCacheByFullSyncCondOpt{
				CondID: fullSyncCondIDs[0],
				Cursor: 0,
				Limit:  10,
			}
			res, err := generalResCli.ListCacheByFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts)))
			for i, info := range hostInfos {
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[i]))
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[i][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[i][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[i][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[i][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[i]["bk_comment"]))
			}
		})
	})

	It("create full sync cond with condition and list cache", func() {
		By("create full sync cond with condition", func() {
			opt := &fullsynccond.CreateFullSyncCondOpt{
				Resource: general.Host,
				IsAll:    false,
				Condition: &filter.Expression{RuleFactory: &filter.CombinedRule{
					Condition: filter.And,
					Rules: []filter.RuleFactory{
						filtertools.GenAtomFilter(common.BKHostIDField, filter.NotEqual, hostIDs[2]),
						filtertools.GenAtomFilter("bk_comment", filter.NotIn,
							[]string{"test1", "test1-3", "comment333"}),
					},
				}},
				Interval: 6,
			}
			fullSyncCondID, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(fullSyncCondID, header)
			Expect(err).NotTo(HaveOccurred())
			fullSyncCondIDs = append(fullSyncCondIDs, fullSyncCondID)
		})

		By("list cache using full sync cond with condition", func() {
			opt := &fullsynccond.ListCacheByFullSyncCondOpt{
				CondID: fullSyncCondIDs[len(fullSyncCondIDs)-1],
				Cursor: 0,
				Limit:  10,
			}
			res, err := generalResCli.ListCacheByFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts) - 2))
			j := 1
			for _, info := range hostInfos {
				if j == 2 {
					j++
				}
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[j]))
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[j][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[j][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[j][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[j][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[j]["bk_comment"]))
				j++
			}
		})
	})

	It("list cache after changing host data", func() {
		By("add hosts", func() {
			addHostInfo := make(mapstr.MapStr)
			for i := 0; i < 4; i++ {
				host := mapstr.MapStr{
					common.BKCloudIDField:     int64(0),
					common.BKHostInnerIPField: fmt.Sprintf("127.1.0.%d", i+1),
					common.BKAgentIDField:     xid.New().String(),
					common.BKAddressingField:  common.BKAddressingDynamic,
					"bk_comment":              fmt.Sprintf("test1-%d", i+1),
				}
				hosts = append(hosts, host)
				addHostInfo[strconv.Itoa(i)] = host
			}
			input := mapstr.MapStr{
				"host_info": addHostInfo,
			}
			rsp, err := hostSvrCli.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		By("update hosts", func() {
			updateInput := map[string]interface{}{
				"update": []map[string]interface{}{
					{
						common.BKHostIDField: hostIDs[3],
						"properties": map[string]interface{}{
							"bk_comment": "comment333",
						},
					},
					{
						common.BKHostIDField: hostIDs[0],
						"properties": map[string]interface{}{
							"bk_comment": "comment000",
						},
					},
				},
			}
			err := hostSvrCli.UpdateHostPropertyBatch(context.Background(), header, updateInput)
			util.RegisterResponseWithRid(err, header)
			Expect(err).NotTo(HaveOccurred())

			hosts[3]["bk_comment"] = "comment333"
			hosts[0]["bk_comment"] = "comment000"
		})

		By("delete hosts", func() {
			input := map[string]interface{}{
				common.BKHostIDField: fmt.Sprintf("%d,%d", hostIDs[1], hostIDs[4]),
			}
			rsp, err := hostSvrCli.DeleteHostBatch(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())

			hosts = append([]mapstr.MapStr{hosts[0], hosts[2], hosts[3]}, hosts[5:]...)
			hostIDs = append([]int64{hostIDs[0], hostIDs[2], hostIDs[3]}, hostIDs[5:]...)
		})

		time.Sleep(10 * time.Second)

		By("list host with page", func() {
			opt := &metadata.ListHostWithPage{
				Page: metadata.BasePage{Limit: 10},
			}
			cnt, jsonStr, err := hostCli.ListHostWithPage(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(cnt).To(Equal(int64(len(hosts))))
			hostInfos := make([]mapstr.MapStr, 0)
			err = json.Unmarshal([]byte(jsonStr), &hostInfos)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts)))
			for i, info := range hostInfos {
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[i][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[i][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[i][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[i][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[i]["bk_comment"]))
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				if i <= 2 {
					Expect(hostIDs[i]).To(Equal(hostID))
					continue
				}
				hostIDs = append(hostIDs, hostID)
			}
		})

		By("list cache using full sync cond with no condition", func() {
			opt := &fullsynccond.ListCacheByFullSyncCondOpt{
				CondID: fullSyncCondIDs[0],
				Cursor: hostIDs[1],
				Limit:  10,
			}
			res, err := generalResCli.ListCacheByFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts) - 2))
			for i, info := range hostInfos {
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[i+2]))
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[i+2][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[i+2][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[i+2][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[i+2][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[i+2]["bk_comment"]))
			}
		})

		By("list cache using full sync cond with condition", func() {
			opt := &fullsynccond.ListCacheByFullSyncCondOpt{
				CondID: fullSyncCondIDs[len(fullSyncCondIDs)-1],
				Cursor: 0,
				Limit:  10,
			}
			res, err := generalResCli.ListCacheByFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts) - 3))
			j := 0
			for _, info := range hostInfos {
				if j == 1 {
					j += 2
				} else if j == 5 {
					j++
				}
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[j]))
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[j][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[j][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[j][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[j][common.BKAddressingField]))
				j++
			}
		})
	})

	It("search host cache", func() {
		By("search host with inner ip", func() {
			opt := &metadata.SearchHostWithInnerIPOption{
				InnerIP: commonutil.GetStrByInterface(hosts[1][common.BKHostInnerIPField]),
				CloudID: 0,
				Fields:  []string{common.BKHostIDField, common.BKAgentIDField},
			}
			jsonStr, err := hostCli.SearchHostWithInnerIPForStatic(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(gjson.Get(jsonStr, common.BKHostIDField).Int()).To(Equal(hostIDs[1]))
			Expect(gjson.Get(jsonStr, common.BKAgentIDField).String()).To(Equal(hosts[1][common.BKAgentIDField]))
		})

		By("search host with agent id", func() {
			opt := &metadata.SearchHostWithAgentID{
				AgentID: commonutil.GetStrByInterface(hosts[5][common.BKAgentIDField]),
				Fields:  []string{common.BKHostIDField, common.BKCloudIDField, common.BKHostInnerIPField},
			}
			jsonStr, err := hostCli.SearchHostWithAgentID(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(gjson.Get(jsonStr, common.BKHostIDField).Int()).To(Equal(hostIDs[5]))
			Expect(gjson.Get(jsonStr, common.BKCloudIDField).Int()).To(Equal(hosts[5][common.BKCloudIDField]))
			Expect(gjson.Get(jsonStr, common.BKHostInnerIPField).String()).
				To(Equal(hosts[5][common.BKHostInnerIPField]))
		})

		By("search host with host id", func() {
			opt := &metadata.SearchHostWithIDOption{
				HostID: hostIDs[0],
				Fields: []string{common.BKAgentIDField, common.BKCloudIDField, common.BKHostInnerIPField},
			}
			jsonStr, err := hostCli.SearchHostWithHostID(context.Background(), header, opt)
			util.RegisterResponseWithRid(jsonStr, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(gjson.Get(jsonStr, common.BKAgentIDField).String()).To(Equal(hosts[0][common.BKAgentIDField]))
			Expect(gjson.Get(jsonStr, common.BKCloudIDField).Int()).To(Equal(hosts[0][common.BKCloudIDField]))
			Expect(gjson.Get(jsonStr, common.BKHostInnerIPField).String()).
				To(Equal(hosts[0][common.BKHostInnerIPField]))
		})
	})

	It("list host detail cache", func() {
		By("list host detail cache with inner ip", func() {
			opt := &general.ListDetailByUniqueKeyOpt{
				Resource: general.Host,
				Type:     general.IPCloudIDType,
				Keys: []string{
					general.IPCloudIDKey(commonutil.GetStrByInterface(hosts[2][common.BKHostInnerIPField]), 0),
					general.IPCloudIDKey(commonutil.GetStrByInterface(hosts[0][common.BKHostInnerIPField]), 0),
				},
				Fields: []string{common.BKHostIDField, common.BKAgentIDField},
			}
			res, err := generalResCli.ListGeneralCacheByUniqueKey(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(2))
			for i, info := range hostInfos {
				j := 2
				if i == 1 {
					j = 0
				}
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[j]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[j][common.BKAgentIDField]))
			}
		})

		By("list host detail cache with agent id", func() {
			opt := &general.ListDetailByUniqueKeyOpt{
				Resource: general.Host,
				Type:     general.AgentIDType,
				Keys: []string{
					general.AgentIDKey(commonutil.GetStrByInterface(hosts[1][common.BKAgentIDField])),
					general.AgentIDKey(commonutil.GetStrByInterface(hosts[5][common.BKAgentIDField])),
				},
				Fields: []string{common.BKHostIDField, common.BKHostInnerIPField},
			}
			res, err := generalResCli.ListGeneralCacheByUniqueKey(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(2))
			for i, info := range hostInfos {
				j := 1
				if i == 1 {
					j = 5
				}
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[j]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[j][common.BKHostInnerIPField]))
			}
		})

		By("list host detail cache with host id", func() {
			opt := &general.ListDetailByIDsOpt{
				Resource: general.Host,
				IDs:      hostIDs,
			}
			res, err := generalResCli.ListGeneralCacheByIDs(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			hostInfos, rawErr := general.DecodeStringArrRes[mapstr.MapStr](res.Info)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(len(hostInfos)).To(Equal(len(hosts)))
			for i, info := range hostInfos {
				hostID, err := commonutil.GetInt64ByInterface(info[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(hostID).To(Equal(hostIDs[i]))
				cloudID, err := commonutil.GetInt64ByInterface(info[common.BKCloudIDField])
				Expect(err).NotTo(HaveOccurred())
				Expect(cloudID).To(Equal(hosts[i][common.BKCloudIDField]))
				Expect(info[common.BKHostInnerIPField]).To(Equal(hosts[i][common.BKHostInnerIPField]))
				Expect(info[common.BKAgentIDField]).To(Equal(hosts[i][common.BKAgentIDField]))
				Expect(info[common.BKAddressingField]).To(Equal(hosts[i][common.BKAddressingField]))
				Expect(info["bk_comment"]).To(Equal(hosts[i]["bk_comment"]))
			}
		})
	})
})
