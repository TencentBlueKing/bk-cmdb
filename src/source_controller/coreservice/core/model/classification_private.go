/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelClassification) hasModel(ctx core.ContextParams, cond mapstr.MapStr) (cnt uint64, exists bool, err error) {

	cnt, err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond).Count(ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *modelClassification) searchClassification(ctx core.ContextParams, cond mapstr.MapStr) ([]metadata.Classification, error) {

	results := []metadata.Classification{}
	err := m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond).All(ctx, &results)

	return results, err
}
