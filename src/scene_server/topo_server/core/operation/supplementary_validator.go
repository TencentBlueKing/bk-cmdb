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

package operation

import (
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
	"configcenter/src/scene_server/validator"
)

// ValidatorInterface the validator methods
type ValidatorInterface interface {
	ValidatorCreate(params types.ContextParams, obj model.Object, datas mapstr.MapStr) error
	ValidatorUpdate(params types.ContextParams, obj model.Object, datas mapstr.MapStr, instID int64, cond condition.Condition) error
}

type valid struct {
	inst InstOperationInterface
}

func (v *valid) ValidatorCreate(params types.ContextParams, obj model.Object, datas mapstr.MapStr) error {
	ignoreKeys := []string{
		common.BKOwnerIDField,
		common.BKDefaultField,
		common.BKInstParentStr,
		common.BKOwnerIDField,
		common.BKAppIDField,
		common.BKSupplierIDField,
		common.BKInstIDField,
	}
	validObj := validator.NewValidMapWithKeyFields(params.SupplierAccount, obj.GetID(), ignoreKeys, params.Header, params.Engin)
	return validObj.ValidMap(datas, common.ValidCreate, -1)
}
func (v *valid) ValidatorUpdate(params types.ContextParams, obj model.Object, datas mapstr.MapStr, instID int64, cond condition.Condition) error {

	ignoreKeys := []string{
		common.BKOwnerIDField,
		common.BKDefaultField,
		common.BKInstParentStr,
		common.BKOwnerIDField,
		common.BKAppIDField,
		common.BKDataStatusField,
		common.BKDataStatusField,
		common.BKSupplierIDField,
		common.BKInstIDField,
	}

	validObj := validator.NewValidMapWithKeyFields(params.SupplierAccount, obj.GetID(), ignoreKeys, params.Header, params.Engin)
	query := &metadata.QueryInput{}
	query.Fields = obj.GetInstIDFieldName()
	if instID < 0 {
		query.Condition = cond.ToMapStr()
		_, insts, err := v.inst.FindInst(params, obj, query, false)
		if nil != err {
			return err
		}

		for _, inst := range insts {
			id, err := inst.GetInstID()
			if nil != err {
				return err
			}
			if err = validObj.ValidMap(datas, common.ValidUpdate, id); nil != err {
				return err
			}
		}
		return nil
	}

	return validObj.ValidMap(datas, common.ValidUpdate, instID)
}
