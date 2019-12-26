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

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

// isExists 需要支持的情况
// 1. 公有模型加入业务私有字段：私有字段不能与当前业务私有字段重复，且不能与公有字段重复
// 2. 公有模型加入业务公有字段：公有字段不能与其它公有字段重复，且不能与任何业务的私有字段重复(即忽略业务参数)
func (m *modelAttribute) isExists(ctx core.ContextParams, objID, propertyID string, meta metadata.Metadata) (oneAttribute *metadata.Attribute, exists bool, err error) {
	filter := map[string]interface{}{
		metadata.AttributeFieldPropertyID: propertyID,
		common.BKObjIDField:               objID,
	}

	bizID, err := meta.ParseBizID()
	if err != nil {
		blog.Errorf("request(%s): database findOne operation is failed, parse biz id failed, error info is %s", ctx.ReqID, err.Error())
		return oneAttribute, false, err
	}
	if bizID != 0 {
		oc := metadata.NewPublicOrBizConditionByBizID(bizID)
		if _, ok := oc[common.BKDBOR]; ok == true {
			filter[common.BKDBOR] = oc[common.BKDBOR]
		}
	}

	condMap := util.SetModOwner(filter, ctx.SupplierAccount)
	oneAttribute = &metadata.Attribute{}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(condMap).One(ctx, oneAttribute)
	blog.V(5).Infof("isExists cond:%#v, rid:%s", condMap, ctx.ReqID)
	if nil != err && !m.dbProxy.IsNotFoundError(err) {
		blog.Errorf("request(%s): database findOne operation is failed, error info is %s", ctx.ReqID, err.Error())
		return oneAttribute, false, err
	}
	return oneAttribute, !m.dbProxy.IsNotFoundError(err), nil
}
