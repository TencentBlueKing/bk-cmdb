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
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelAttribute) isExists(ctx core.ContextParams, propertyID string, meta metadata.Metadata) (oneAttribute *metadata.Attribute, exists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: propertyID})
	cond.Element(&mongo.Eq{Key: metadata.AttributeFieldPropertyID, Val: propertyID})

	// ATTETION: Currently only business dimension isolation is done,
	//           and there may be isolation requirements for other dimensions in the future.
	isExsit, bizID := meta.Label.Get(common.BKAppIDField)
	if isExsit {
		_, metaCond := cond.Embed(metadata.BKMetadata)
		_, lableCond := metaCond.Embed(metadata.BKLabel)
		lableCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}

	oneAttribute = &metadata.Attribute{}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).One(ctx, oneAttribute)
	if nil != err && !m.dbProxy.IsNotFoundError(err) {
		blog.Errorf("request(%s): database findone operation is failed, error info is %s", err.Error())
		return oneAttribute, false, err
	}
	return oneAttribute, !m.dbProxy.IsNotFoundError(err), nil
}
