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

func (m *modelAttrUnique) searchModelAttrUnique(ctx core.ContextParams, inputParam metadata.QueryCondition) (results []metadata.ObjectUnique, err error) {
	instHandler := m.dbProxy.Table(common.BKTableNameObjUnique).Find(inputParam.Condition)
	for _, sort := range inputParam.SortArr {
		fileld := sort.Field
		if sort.IsDsc {
			fileld = "-" + fileld
		}
		instHandler = instHandler.Sort(fileld)
	}
	err = instHandler.Start(uint64(inputParam.Limit.Offset)).Limit(uint64(inputParam.Limit.Limit)).All(ctx, &results)

	return results, err
}

func (m *modelAttrUnique) countModelAttrUnique(ctx core.ContextParams, cond mapstr.MapStr) (count uint64, err error) {

	count, err = m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond).Count(ctx)

	return count, err
}
