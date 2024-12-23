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

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("full sync cond crud test", func() {
	var id int64

	var _ = Describe("full sync cond crud normal test", func() {
		It("create full sync cond", func() {
			By("test preparation", func() {
				deleteAllFullSyncCond()
			})

			opt := &fullsynccond.CreateFullSyncCondOpt{
				Resource: general.Biz,
				IsAll:    true,
				Interval: 6,
			}
			var err error
			id, err = generalResCli.CreateFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(id, header)
			Expect(err).NotTo(HaveOccurred())
		})

		It("list full sync cond by resource", func() {
			opt := &fullsynccond.ListFullSyncCondOpt{
				Resource: general.Biz,
			}
			res, err := generalResCli.ListFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(res.Info)).To(Equal(1))
			Expect(res.Info[0].ID).To(Equal(id))
			Expect(res.Info[0].Resource).To(Equal(general.Biz))
			Expect(res.Info[0].SubResource).To(Equal(""))
			Expect(res.Info[0].IsAll).To(Equal(true))
			Expect(res.Info[0].Interval).To(Equal(6))
			Expect(res.Info[0].Condition == nil).To(Equal(true))
			Expect(res.Info[0].SupplierAccount).To(Equal(httpheader.GetSupplierAccount(header)))
		})

		It("update full sync cond", func() {
			opt := &fullsynccond.UpdateFullSyncCondOpt{
				ID:   id,
				Data: &fullsynccond.UpdateFullSyncCondData{Interval: 7},
			}
			err := generalResCli.UpdateFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(err, header)
			Expect(err).NotTo(HaveOccurred())
		})

		It("list full sync cond by ids", func() {
			opt := &fullsynccond.ListFullSyncCondOpt{
				Resource: general.Biz,
				IDs:      []int64{id},
			}
			res, err := generalResCli.ListFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(res.Info)).To(Equal(1))
			Expect(res.Info[0].ID).To(Equal(id))
			Expect(res.Info[0].Resource).To(Equal(general.Biz))
			Expect(res.Info[0].SubResource).To(Equal(""))
			Expect(res.Info[0].IsAll).To(Equal(true))
			Expect(res.Info[0].Interval).To(Equal(7))
			Expect(res.Info[0].Condition == nil).To(Equal(true))
			Expect(res.Info[0].SupplierAccount).To(Equal(httpheader.GetSupplierAccount(header)))
		})

		It("delete full sync cond", func() {
			opt := &fullsynccond.DeleteFullSyncCondOpt{
				ID: id,
			}
			err := generalResCli.DeleteFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(err, header)
			Expect(err).NotTo(HaveOccurred())

			listOpt := &fullsynccond.ListFullSyncCondOpt{
				Resource: general.Biz,
			}
			res, err := generalResCli.ListFullSyncCond(context.Background(), header, listOpt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(res.Info)).To(Equal(0))
		})

		It("create full sync cond with sub res and cond", func() {
			opt := &fullsynccond.CreateFullSyncCondOpt{
				Resource:    general.ObjectInstance,
				SubResource: "bk_switch",
				IsAll:       false,
				Condition:   filtertools.GenAtomFilter(common.BKInstNameField, filter.NotEqual, "aaa"),
				Interval:    6,
			}
			var err error
			id, err = generalResCli.CreateFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(id, header)
			Expect(err).NotTo(HaveOccurred())

			listOpt := &fullsynccond.ListFullSyncCondOpt{
				Resource:    general.ObjectInstance,
				SubResource: "bk_switch",
			}
			res, err := generalResCli.ListFullSyncCond(context.Background(), header, listOpt)
			util.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(res.Info)).To(Equal(1))
			Expect(res.Info[0].ID).To(Equal(id))
			Expect(res.Info[0].Resource).To(Equal(opt.Resource))
			Expect(res.Info[0].SubResource).To(Equal(opt.SubResource))
			Expect(res.Info[0].IsAll).To(Equal(opt.IsAll))
			Expect(res.Info[0].Interval).To(Equal(opt.Interval))
			Expect(res.Info[0].Condition.String()).To(Equal(opt.Condition.String()))
			Expect(res.Info[0].SupplierAccount).To(Equal(httpheader.GetSupplierAccount(header)))
		})
	})

	var _ = Describe("full sync cond crud abnormal test", func() {
		It("create full sync cond abnormal test", func() {
			By("create full sync cond with invalid resource", func() {
				opt := &fullsynccond.CreateFullSyncCondOpt{
					Resource: "xxxx",
					IsAll:    true,
					Interval: 6,
				}
				id, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())
			})

			By("create obj inst full sync cond with no sub resource", func() {
				opt := &fullsynccond.CreateFullSyncCondOpt{
					Resource: general.ObjectInstance,
					IsAll:    true,
					Interval: 6,
				}
				id, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())
			})

			By("create full sync cond with no cond when is_all==false", func() {
				opt := &fullsynccond.CreateFullSyncCondOpt{
					Resource: general.Biz,
					IsAll:    false,
					Interval: 6,
				}
				id, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())
			})

			By("create full sync cond with no cond when is_all==false", func() {
				opt := &fullsynccond.CreateFullSyncCondOpt{
					Resource: general.Biz,
					IsAll:    false,
					Interval: 6,
				}
				id, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())
			})

			By("create full sync cond with invalid cond", func() {
				opt := &fullsynccond.CreateFullSyncCondOpt{
					Resource: general.Biz,
					IsAll:    false,
					Condition: &filter.Expression{
						RuleFactory: &filter.AtomRule{
							Field:    "",
							Operator: "xxx",
							Value:    "a",
						},
					},
					Interval: 6,
				}
				id, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())
			})

			By("create full sync cond with invalid interval", func() {
				opt := &fullsynccond.CreateFullSyncCondOpt{
					Resource: general.Biz,
					IsAll:    true,
					Interval: 2,
				}
				id, err := generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())

				opt.Interval = 8 * 24
				id, err = generalResCli.CreateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(id, header)
				Expect(err).To(HaveOccurred())
			})
		})

		It("update full sync cond abnormal test", func() {
			By("update full sync cond with invalid id", func() {
				opt := &fullsynccond.UpdateFullSyncCondOpt{
					Data: &fullsynccond.UpdateFullSyncCondData{Interval: 7},
				}
				err := generalResCli.UpdateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(err, header)
				Expect(err).To(HaveOccurred())
			})

			By("update full sync cond with invalid interval", func() {
				opt := &fullsynccond.UpdateFullSyncCondOpt{
					ID:   id,
					Data: &fullsynccond.UpdateFullSyncCondData{Interval: 1},
				}
				err := generalResCli.UpdateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(err, header)
				Expect(err).To(HaveOccurred())

				opt.Data.Interval = 8 * 24
				err = generalResCli.UpdateFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(err, header)
				Expect(err).To(HaveOccurred())
			})
		})

		It("delete full sync cond abnormal test", func() {
			opt := &fullsynccond.DeleteFullSyncCondOpt{
				ID: 0,
			}
			err := generalResCli.DeleteFullSyncCond(context.Background(), header, opt)
			util.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		It("list full sync cond abnormal test", func() {
			By("list full sync cond with invalid resource", func() {
				opt := &fullsynccond.ListFullSyncCondOpt{
					Resource: "xxx",
				}
				res, err := generalResCli.ListFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(res, header)
				Expect(err).To(HaveOccurred())
			})

			By("list obj inst full sync cond with no sub resource", func() {
				opt := &fullsynccond.ListFullSyncCondOpt{
					Resource: general.ObjectInstance,
				}
				res, err := generalResCli.ListFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(res, header)
				Expect(err).To(HaveOccurred())
			})

			By("list obj inst full sync cond with no info", func() {
				opt := new(fullsynccond.ListFullSyncCondOpt)
				res, err := generalResCli.ListFullSyncCond(context.Background(), header, opt)
				util.RegisterResponseWithRid(res, header)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
