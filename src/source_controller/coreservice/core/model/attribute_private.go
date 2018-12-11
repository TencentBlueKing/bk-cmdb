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
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelAttribute) isExists(ctx core.ContextParams, propertyID string) (oneAttribute *metadata.Attribute, exists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: propertyID})
	cond.Element(&mongo.Eq{Key: metadata.AttributeFieldPropertyID, Val: propertyID})

	oneAttribute = &metadata.Attribute{}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).One(ctx, oneAttribute)
	if nil != err && m.dbProxy.IsNotFoundError(err) {
		return oneAttribute, exists, nil
	}
	return oneAttribute, m.dbProxy.IsNotFoundError(err), err
}
