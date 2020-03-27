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
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (g *modelAttributeGroup) groupIDIsExists(ctx core.ContextParams, objID, groupID string, meta metadata.Metadata) (oneResult metadata.Group, isExists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldGroupID, Val: groupID})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: ctx.SupplierAccount})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})
	exist, bizID := meta.Label.Get(common.BKAppIDField)
	if exist {
		_, metaCond := cond.Embed(metadata.BKMetadata)
		_, labelCond := metaCond.Embed(metadata.BKLabel)
		labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}
	groups, err := g.search(ctx, cond)
	if nil != err {
		return oneResult, isExists, err
	}

	if 0 != len(groups) {
		return groups[0], true, nil
	}

	return oneResult, isExists, nil
}

func (g *modelAttributeGroup) groupNameIsExists(ctx core.ContextParams, objID, groupName string, meta metadata.Metadata) (oneResult metadata.Group, isExists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldGroupName, Val: groupName})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: ctx.SupplierAccount})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})
	exist, bizID := meta.Label.Get(common.BKAppIDField)
	if exist {
		_, metaCond := cond.Embed(metadata.BKMetadata)
		_, labelCond := metaCond.Embed(metadata.BKLabel)
		labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}
	groups, err := g.search(ctx, cond)
	if nil != err {
		return oneResult, isExists, err
	}

	if 0 != len(groups) {
		return groups[0], true, nil
	}

	return oneResult, isExists, nil
}

func (g *modelAttributeGroup) hasAttributes(ctx core.ContextParams, objID string, groupIDS []string) (isExists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: ctx.SupplierAccount})
	cond.Element(&mongo.In{Key: metadata.AttributeFieldPropertyGroup, Val: groupIDS})

	attrs, err := g.model.SearchModelAttributes(ctx, objID, metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	})

	if nil != err {
		return false, err
	}

	return 0 != attrs.Count, nil
}
