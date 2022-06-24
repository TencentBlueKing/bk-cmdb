/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package topo_server_test

import (
	"context"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/querybuilder"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("business set test", func() {
	ctx := context.Background()

	var sampleBizSetID, bizID3 int64
	It("prepare environment, create a biz set and biz in it with topo for searching biz and topo in biz set", func() {
		test.ClearDatabase()

		biz := map[string]interface{}{
			common.BKMaintainersField: "biz_set",
			common.BKAppNameField:     "biz_for_biz_set",
			"time_zone":               "Africa/Accra",
			"language":                "1",
		}
		bizResp, err := apiServerClient.CreateBiz(ctx, "0", header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID1, err := commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		biz[common.BKAppNameField] = "biz_for_biz_set1"
		bizResp, err = apiServerClient.CreateBiz(ctx, "0", header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID2, err := commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		biz[common.BKAppNameField] = "biz_not_for_biz_set"
		bizResp, err = apiServerClient.CreateBiz(ctx, "0", header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID3, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		mainlineObj := &metadata.MainLineObject{
			Object: metadata.Object{
				ObjCls:     "bk_biz_topo",
				ObjectID:   "mainline_obj_for_biz_set",
				ObjectName: "mainline_obj_for_biz_set",
				ObjIcon:    "icon-cc-business",
			},
			AssociationID: "biz",
		}
		mainlineObjResp, err := objectClient.CreateModel(ctx, header, mainlineObj)
		util.RegisterResponseWithRid(mainlineObjResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(mainlineObjResp.Result).To(Equal(true))

		mainlineInst := map[string]interface{}{
			common.BKInstNameField: "mainline_inst_for_biz_set",
			common.BKAppIDField:    bizID1,
			common.BKParentIDField: bizID1,
		}
		mainlineInstResp, err := instClient.CreateInst(ctx, "mainline_obj_for_biz_set", header, mainlineInst)
		util.RegisterResponseWithRid(mainlineInstResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(mainlineInstResp.Result).To(Equal(true))
		mainlineInstID, err := commonutil.GetInt64ByInterface(mainlineInstResp.Data[common.BKInstIDField])
		Expect(err).NotTo(HaveOccurred())

		set := map[string]interface{}{
			common.BKSetNameField:  "set_for_biz_set",
			common.BKAppIDField:    bizID1,
			common.BKParentIDField: mainlineInstID,
		}
		setResp, err := instClient.CreateSet(ctx, bizID1, header, set)
		util.RegisterResponseWithRid(setResp, header)
		Expect(err).NotTo(HaveOccurred())
		setID, err := commonutil.GetInt64ByInterface(setResp[common.BKSetIDField])
		Expect(err).NotTo(HaveOccurred())

		module := map[string]interface{}{
			common.BKModuleNameField: "module_for_biz_set",
			common.BKAppIDField:      bizID1,
			common.BKParentIDField:   setID,
		}
		moduleResp, err := instClient.CreateModule(ctx, bizID1, setID, header, module)
		util.RegisterResponseWithRid(moduleResp, header)
		Expect(err).NotTo(HaveOccurred())

		createBizSetOpt := metadata.CreateBizSetRequest{
			BizSetAttr: map[string]interface{}{
				common.BKBizSetNameField: "sample_biz_set",
			},
			BizSetScope: &metadata.BizSetScope{
				MatchAll: false,
				Filter: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							querybuilder.AtomRule{
								Field:    common.BKAppIDField,
								Operator: querybuilder.OperatorIn,
								Value:    []int64{bizID1, bizID2},
							},
						},
					}},
			},
		}

		sampleBizSetID, err = instClient.CreateBizSet(ctx, header, createBizSetOpt)
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("update business set test", func() {
		var bizSetID int64

		By("create business set")
		func() {
			createOpt := metadata.CreateBizSetRequest{
				BizSetAttr: map[string]interface{}{
					common.BKBizSetNameField: "biz_set_for_updating",
				},
				BizSetScope: &metadata.BizSetScope{
					MatchAll: true,
				},
			}

			var err error
			bizSetID, err = instClient.CreateBizSet(ctx, header, createOpt)
			util.RegisterResponseWithRid(bizSetID, header)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("update business set attributes with scope")
		searchBizSetOpt := &metadata.QueryBusinessSetRequest{
			BizSetPropertyFilter: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						querybuilder.AtomRule{
							Field:    common.BKBizSetIDField,
							Operator: querybuilder.OperatorEqual,
							Value:    bizSetID,
						},
					},
				}},
			Page: metadata.BasePage{EnableCount: false},
		}

		var updatedBizSetScopeJson []byte
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID},
				Data: &metadata.UpdateBizSetData{
					BizSetAttr: map[string]interface{}{
						common.BKBizSetNameField: "updated_biz_set",
					},
					Scope: &metadata.BizSetScope{
						MatchAll: false,
						Filter: &querybuilder.QueryFilter{
							Rule: querybuilder.CombinedRule{
								Condition: querybuilder.ConditionAnd,
								Rules: []querybuilder.Rule{
									querybuilder.AtomRule{
										Field:    common.BKAppIDField,
										Operator: querybuilder.OperatorEqual,
										Value:    bizID3,
									},
								},
							}}},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())

			bizSet, err := instClient.SearchBusinessSet(ctx, header, searchBizSetOpt)
			util.RegisterResponseWithRid(bizSet, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(bizSet.Count).To(Equal(0))
			Expect(len(bizSet.Info)).To(Equal(1))
			Expect(commonutil.GetStrByInterface(bizSet.Info[0][common.BKBizSetNameField])).To(Equal("updated_biz_set"))

			j, jErr := json.Marshal(bizSet.Info[0][common.BKBizSetScopeField])
			Expect(jErr).NotTo(HaveOccurred())
			scope := new(metadata.BizSetScope)
			jErr = json.Unmarshal(j, &scope)
			Expect(jErr).NotTo(HaveOccurred())
			Expect(scope.MatchAll).To(Equal(false))
			Expect(scope.Filter).NotTo(Equal(nil))
			combineRule, ok := scope.Filter.Rule.(querybuilder.CombinedRule)
			Expect(ok).To(Equal(true))
			Expect(combineRule.Condition).To(Equal(querybuilder.ConditionAnd))
			Expect(len(combineRule.Rules)).To(Equal(1))
			atomRule, ok := combineRule.Rules[0].(querybuilder.AtomRule)
			Expect(ok).To(Equal(true))
			Expect(atomRule.Field).To(Equal(common.BKAppIDField))
			Expect(atomRule.Operator).To(Equal(querybuilder.OperatorEqual))

			atomRuleVal, rawErr := commonutil.GetInt64ByInterface(atomRule.Value)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(atomRuleVal).To(Equal(bizID3))

			updatedBizSetScopeJson = j
		}()

		By("update business set attributes only")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID},
				Data: &metadata.UpdateBizSetData{
					BizSetAttr: map[string]interface{}{
						common.BKBizSetNameField: "updated_biz_set1",
						common.BKBizSetDescField: "updated_biz_set1",
					},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())

			bizSet, err := instClient.SearchBusinessSet(ctx, header, searchBizSetOpt)
			util.RegisterResponseWithRid(bizSet, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(bizSet.Count).To(Equal(0))
			Expect(len(bizSet.Info)).To(Equal(1))
			Expect(commonutil.GetStrByInterface(bizSet.Info[0][common.BKBizSetNameField])).To(Equal("updated_biz_set1"))
			Expect(commonutil.GetStrByInterface(bizSet.Info[0][common.BKBizSetDescField])).To(Equal("updated_biz_set1"))

			j, jErr := json.Marshal(bizSet.Info[0][common.BKBizSetScopeField])
			Expect(jErr).NotTo(HaveOccurred())
			Expect(j).To(Equal(updatedBizSetScopeJson))
		}()

		By("update business set scope only")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID},
				Data: &metadata.UpdateBizSetData{
					Scope: &metadata.BizSetScope{
						MatchAll: true,
					},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())

			bizSet, err := instClient.SearchBusinessSet(ctx, header, searchBizSetOpt)
			util.RegisterResponseWithRid(bizSet, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(bizSet.Count).To(Equal(0))
			Expect(len(bizSet.Info)).To(Equal(1))
			Expect(commonutil.GetStrByInterface(bizSet.Info[0][common.BKBizSetNameField])).To(Equal("updated_biz_set1"))
			Expect(commonutil.GetStrByInterface(bizSet.Info[0][common.BKBizSetDescField])).To(Equal("updated_biz_set1"))

			Expect(bizSet.Info[0][common.BKBizSetScopeField]).NotTo(Equal(nil))
			scope, ok := bizSet.Info[0][common.BKBizSetScopeField].(map[string]interface{})
			Expect(ok).To(Equal(true))
			matchAll, ok := scope[common.BKBizSetMatchField].(bool)
			Expect(ok).To(Equal(true))
			Expect(matchAll).To(Equal(true))
		}()

		By("update business set with no biz set id")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: make([]int64, 0),
				Data: &metadata.UpdateBizSetData{
					BizSetAttr: map[string]interface{}{
						common.BKBizSetDescField: "update_with_no_biz_id",
					},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()

		By("update business sets name attribute")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID, sampleBizSetID},
				Data: &metadata.UpdateBizSetData{
					BizSetAttr: map[string]interface{}{
						common.BKBizSetNameField: "update_multiple_names",
					},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()

		By("update business sets scope")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID, sampleBizSetID},
				Data: &metadata.UpdateBizSetData{
					Scope: &metadata.BizSetScope{
						MatchAll: true,
					},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()

		By("update business set with no attribute")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()

		By("update business set with invalid scope")
		func() {
			updateOpt := metadata.UpdateBizSetOption{
				BizSetIDs: []int64{bizSetID},
				Data: &metadata.UpdateBizSetData{
					Scope: &metadata.BizSetScope{
						MatchAll: false,
						Filter: &querybuilder.QueryFilter{
							Rule: querybuilder.AtomRule{
								Field:    common.BKBizSetIDField,
								Operator: querybuilder.OperatorEqual,
								Value:    bizID3,
							},
						},
					},
				},
			}

			err := instClient.UpdateBizSet(ctx, header, updateOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()
	})

	It("delete business set test", func() {
		var bizSetID int64

		By("create business set")
		func() {
			createOpt := metadata.CreateBizSetRequest{
				BizSetAttr: map[string]interface{}{
					common.BKBizSetNameField: "biz_set_for_deletion",
				},
				BizSetScope: &metadata.BizSetScope{
					MatchAll: true,
				},
			}

			var err error
			bizSetID, err = instClient.CreateBizSet(ctx, header, createOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("delete business set with no biz set ids")
		func() {
			noIDOpt := metadata.DeleteBizSetOption{
				BizSetIDs: make([]int64, 0),
			}

			err := instClient.DeleteBizSet(ctx, header, noIDOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()

		By("delete business set with biz set ids that exceeds max length")
		func() {
			tooManyIDOpt := metadata.DeleteBizSetOption{
				BizSetIDs: make([]int64, 101),
			}
			for i := 0; i < 101; i++ {
				tooManyIDOpt.BizSetIDs[i] = int64(i + 1)
			}

			err := instClient.DeleteBizSet(ctx, header, tooManyIDOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).To(HaveOccurred())
		}()

		By("delete business set")
		func() {
			delOpt := metadata.DeleteBizSetOption{
				BizSetIDs: []int64{bizSetID},
			}

			err := instClient.DeleteBizSet(ctx, header, delOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())

			count, err := clientSet.CoreService().Count().GetCountByFilter(ctx, header, common.BKTableNameBaseBizSet,
				[]map[string]interface{}{{common.BKBizSetIDField: bizSetID}})
			util.RegisterResponseWithRid(count, header)
			Expect(len(count)).To(Equal(1))
			Expect(count[0]).To(Equal(int64(0)))
		}()
	})

	It("find businesses in biz set test", func() {
		By("count businesses in biz set")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: sampleBizSetID,
				Page:     metadata.BasePage{EnableCount: true},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(biz.Count).To(Equal(2))
			Expect(len(biz.Info)).To(Equal(0))
		}()

		By("find businesses in biz set")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: sampleBizSetID,
				Page:     metadata.BasePage{Limit: 10, Sort: common.BKAppIDField},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(biz.Info)).To(Equal(2))
			Expect(commonutil.GetStrByInterface(biz.Info[0][common.BKAppNameField])).To(Equal("biz_for_biz_set"))
			Expect(commonutil.GetStrByInterface(biz.Info[1][common.BKAppNameField])).To(Equal("biz_for_biz_set1"))
		}()

		var bizSetID int64
		By("create business set that matches all biz")
		func() {
			createOpt := metadata.CreateBizSetRequest{
				BizSetAttr: map[string]interface{}{
					common.BKBizSetNameField: "biz_set",
				},
				BizSetScope: &metadata.BizSetScope{
					MatchAll: true,
				},
			}

			var err error
			bizSetID, err = instClient.CreateBizSet(ctx, header, createOpt)
			util.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("find businesses in biz set that matches all biz")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: bizSetID,
				Page:     metadata.BasePage{Limit: 10, Sort: common.BKAppIDField},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(biz.Info)).To(Equal(4))
			Expect(commonutil.GetStrByInterface(biz.Info[0][common.BKAppNameField])).To(Equal("蓝鲸"))
			Expect(commonutil.GetStrByInterface(biz.Info[1][common.BKAppNameField])).To(Equal("biz_for_biz_set"))
			Expect(commonutil.GetStrByInterface(biz.Info[2][common.BKAppNameField])).To(Equal("biz_for_biz_set1"))
			Expect(commonutil.GetStrByInterface(biz.Info[3][common.BKAppNameField])).To(Equal("biz_not_for_biz_set"))
		}()

		By("count businesses in biz set that matches all biz")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: bizSetID,
				Page:     metadata.BasePage{EnableCount: true},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(biz.Count).To(Equal(4))
			Expect(len(biz.Info)).To(Equal(0))
		}()

		By("find businesses in biz set with fields")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: sampleBizSetID,
				Fields:   []string{common.BKAppNameField},
				Page:     metadata.BasePage{Limit: 1},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(biz.Info)).To(Equal(1))
			Expect(biz.Info[0].Exists(common.BKAppNameField)).To(Equal(true))
			Expect(biz.Info[0].Exists(common.BKAppIDField)).To(Equal(false))
		}()

		By("find businesses in biz set with enable count")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: sampleBizSetID,
				Page:     metadata.BasePage{EnableCount: true, Limit: 10},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()

		By("find businesses in biz set by invalid biz set id")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: 1000000,
				Page:     metadata.BasePage{Limit: 10},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()

		By("find businesses in biz set with no limit")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: sampleBizSetID,
				Page:     metadata.BasePage{Start: 1},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()

		By("find businesses in biz set with a limit that exceeds maximum")
		func() {
			findBizOpt := &metadata.FindBizInBizSetOption{
				BizSetID: sampleBizSetID,
				Page:     metadata.BasePage{Limit: 10000},
			}

			biz, err := instClient.FindBizInBizSet(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()
	})

	It("find topo info in biz set test", func() {
		var bizID, mainlineInstID, setID int64

		By("find biz topo info in biz set")
		func() {
			findBizOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: common.BKInnerObjIDBizSet,
				ParentID:    sampleBizSetID,
			}

			biz, err := instClient.FindBizSetTopo(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(biz)).To(Equal(2))
			Expect(commonutil.GetStrByInterface(biz[0][common.BKObjIDField])).To(Equal(common.BKInnerObjIDApp))
			Expect(commonutil.GetStrByInterface(biz[0][common.BKInstNameField])).To(Equal("biz_for_biz_set"))
			Expect(commonutil.GetStrByInterface(biz[1][common.BKObjIDField])).To(Equal(common.BKInnerObjIDApp))
			Expect(commonutil.GetStrByInterface(biz[1][common.BKInstNameField])).To(Equal("biz_for_biz_set1"))

			var rawErr error
			bizID, rawErr = commonutil.GetInt64ByInterface(biz[0][common.BKInstIDField])
			Expect(rawErr).NotTo(HaveOccurred())
		}()

		By("find child topo info of biz in biz set, contains default set and the custom level under biz")
		func() {
			input := &params.SearchParams{
				Condition: map[string]interface{}{common.BKDefaultField: common.DefaultResSetFlag},
			}
			defaultSetResp, err := instClient.SearchSet(ctx, "0", strconv.FormatInt(bizID, 10), header, input)
			util.RegisterResponseWithRid(defaultSetResp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(defaultSetResp.Result).To(Equal(true))
			Expect(len(defaultSetResp.Data.Info)).To(Equal(1))

			findBizChildOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: common.BKInnerObjIDApp,
				ParentID:    bizID,
			}

			topo, err := instClient.FindBizSetTopo(ctx, header, findBizChildOpt)
			util.RegisterResponseWithRid(topo, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(topo)).To(Equal(3))
			Expect(commonutil.GetStrByInterface(topo[0][common.BKObjIDField])).To(Equal(common.BKInnerObjIDSet))
			Expect(topo[0][common.BKInstIDField]).To(Equal(defaultSetResp.Data.Info[0][common.BKSetIDField]))
			Expect(commonutil.GetStrByInterface(topo[1][common.BKObjIDField])).To(Equal("mainline_obj_for_biz_set"))
			Expect(commonutil.GetStrByInterface(topo[1][common.BKInstNameField])).To(Equal("mainline_obj_for_biz_set"))
			Expect(commonutil.GetStrByInterface(topo[2][common.BKObjIDField])).To(Equal("mainline_obj_for_biz_set"))
			Expect(commonutil.GetStrByInterface(topo[2][common.BKInstNameField])).To(Equal("mainline_inst_for_biz_set"))

			var rawErr error
			mainlineInstID, rawErr = commonutil.GetInt64ByInterface(topo[2][common.BKInstIDField])
			Expect(rawErr).NotTo(HaveOccurred())
		}()

		By("find set topo info in biz set under custom level")
		func() {
			findSetOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: "mainline_obj_for_biz_set",
				ParentID:    mainlineInstID,
			}

			set, err := instClient.FindBizSetTopo(ctx, header, findSetOpt)
			util.RegisterResponseWithRid(set, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(set)).To(Equal(1))
			Expect(commonutil.GetStrByInterface(set[0][common.BKObjIDField])).To(Equal(common.BKInnerObjIDSet))
			Expect(commonutil.GetStrByInterface(set[0][common.BKInstNameField])).To(Equal("set_for_biz_set"))

			var rawErr error
			setID, rawErr = commonutil.GetInt64ByInterface(set[0][common.BKInstIDField])
			Expect(rawErr).NotTo(HaveOccurred())
		}()

		By("find module topo info in biz set under custom level")
		func() {
			findModuleOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: common.BKInnerObjIDSet,
				ParentID:    setID,
			}

			module, err := instClient.FindBizSetTopo(ctx, header, findModuleOpt)
			util.RegisterResponseWithRid(module, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(module)).To(Equal(1))
			Expect(commonutil.GetStrByInterface(module[0][common.BKObjIDField])).To(Equal(common.BKInnerObjIDModule))
			Expect(commonutil.GetStrByInterface(module[0][common.BKInstNameField])).To(Equal("module_for_biz_set"))
		}()

		By("find biz topo info in biz set with invalid biz set id")
		func() {
			findBizOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    100000,
				ParentObjID: common.BKInnerObjIDBizSet,
				ParentID:    100000,
			}

			biz, err := instClient.FindBizSetTopo(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()

		By("find biz topo info in biz set with not matched biz set id and parent id")
		func() {
			findBizOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: common.BKInnerObjIDBizSet,
				ParentID:    bizID,
			}

			biz, err := instClient.FindBizSetTopo(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()

		By("find topo info in biz set with common object")
		func() {
			findBizOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: common.BKInnerObjIDHost,
				ParentID:    sampleBizSetID,
			}

			biz, err := instClient.FindBizSetTopo(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()

		By("find topo info in biz set with parent not in biz set")
		func() {
			findBizOpt := &metadata.FindBizSetTopoOption{
				BizSetID:    sampleBizSetID,
				ParentObjID: common.BKInnerObjIDApp,
				ParentID:    bizID3,
			}

			biz, err := instClient.FindBizSetTopo(ctx, header, findBizOpt)
			util.RegisterResponseWithRid(biz, header)
			Expect(err).To(HaveOccurred())
		}()
	})
})
