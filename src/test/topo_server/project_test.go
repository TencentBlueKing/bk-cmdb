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
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("project test", func() {
	ctx := context.Background()

	var id1 int64
	var id2 int64
	It("create project", func() {
		project1 := map[string]interface{}{
			common.BKProjectIDField:     "99bc6e8c406b48c6b12e3a9327172999",
			common.BKProjectNameField:   "project1",
			common.BKProjectCodeField:   "project1",
			common.BKProjectDescField:   "project1 desc",
			common.BKProjectTypeField:   "mobile_game",
			common.BKProjectSecLvlField: "public",
			common.BKProjectOwnerField:  "admin",
			common.BKProjectTeamField:   []int64{1, 2},
			common.BKProjectIconField:   "https://127.0.0.1/file/png/11111",
		}

		project2 := map[string]interface{}{
			common.BKProjectNameField:   "project2",
			common.BKProjectCodeField:   "project2",
			common.BKProjectDescField:   "project2 desc",
			common.BKProjectTypeField:   "mobile_game",
			common.BKProjectSecLvlField: "public",
			common.BKProjectOwnerField:  "admin",
			common.BKProjectTeamField:   []int64{1, 2},
			common.BKProjectIconField:   "https://127.0.0.1/file/png/11111",
		}

		opt := metadata.CreateProjectOption{
			Data: []mapstr.MapStr{project1, project2},
		}
		result, err := topoServerClient.Instance().CreateProject(ctx, header, &opt)
		util.RegisterResponseWithRid(result, header)
		Expect(err).NotTo(HaveOccurred())
		id1 = result.IDs[0]
		id2 = result.IDs[1]
	})

	It("update project", func() {
		opt := metadata.UpdateProjectOption{
			IDs:  []int64{id1},
			Data: mapstr.MapStr{common.BKProjectStatusField: "disabled"},
		}
		err := topoServerClient.Instance().UpdateProject(ctx, header, &opt)
		Expect(err).NotTo(HaveOccurred())
	})

	It("find project", func() {
		filter := &querybuilder.QueryFilter{
			Rule: &querybuilder.CombinedRule{
				Condition: querybuilder.ConditionAnd,
				Rules: []querybuilder.Rule{
					querybuilder.AtomRule{
						Field:    common.BKProjectStatusField,
						Operator: querybuilder.OperatorEqual,
						Value:    "enable",
					},
				},
			},
		}

		// get project data
		page := metadata.BasePage{
			Start: 0,
			Limit: 10,
		}
		fields := []string{common.BKFieldID}
		opt := metadata.SearchProjectOption{
			Filter: filter,
			Page:   page,
			Fields: fields,
		}
		result, err := topoServerClient.Instance().SearchProject(ctx, header, &opt)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0][common.BKFieldID].(json.Number).Int64()).To(Equal(id2))

		// get project count
		page = metadata.BasePage{
			EnableCount: true,
		}
		opt = metadata.SearchProjectOption{
			Filter: filter,
			Page:   page,
		}
		opt.Page.EnableCount = true
		result, err = topoServerClient.Instance().SearchProject(ctx, header, &opt)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Count).To(Equal(1))
	})

	It("update project id", func() {
		opt := metadata.UpdateProjectIDOption{
			ID:        id1,
			ProjectID: "99bc6e8c406b48c6b12e3a9327172666",
		}
		err := topoServerClient.Instance().UpdateProjectID(ctx, header, &opt)
		Expect(err).NotTo(HaveOccurred())
	})

	It("delete project", func() {
		deleteOpt := metadata.DeleteProjectOption{
			IDs: []int64{id1, id2},
		}

		err := topoServerClient.Instance().DeleteProject(ctx, header, &deleteOpt)
		Expect(err).NotTo(HaveOccurred())
	})
})
