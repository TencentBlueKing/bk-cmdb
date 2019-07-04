/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) IsInstanceExist(ctx core.ContextParams, objID string, instID uint64) (exists bool, err error) {
	instIDFieldName := common.GetInstIDField(objID)
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: instIDFieldName, Val: instID})
	searchCond := metadata.QueryCondition{Condition: cond.ToMapStr()}
	result, err := s.core.InstanceOperation().SearchModelInstance(ctx, objID, searchCond)
	if nil != err {
		blog.Errorf("search model instance error: %v, rid: %s", err, ctx.ReqID)
		return false, err
	}
	if 0 == result.Count {
		return false, nil
	}
	return true, nil
}
