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

package inst

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"fmt"
)

// BusinessOperationInterface business operation methods
type BusinessSetOperationInterface interface {

	// CreateBusinessSet create business set
	CreateBusinessSet(kit *rest.Kit, data *metadata.CreateBizSetRequest) (mapstr.MapStr, error)

	// FindBizSet  find biz set find biz set by cond
	FindBizSet(kit *rest.Kit, cond *metadata.CommonSearchFilter) (*metadata.CommonSearchResult, error)

	// CountBizSet  count biz set find biz set by cond
	CountBizSet(kit *rest.Kit, cond *metadata.CommonCountFilter) (*metadata.CommonCountResult, error)
	// SetProxy set business proxy
	SetProxy(inst InstOperationInterface)
}

// NewBusinessOperation create a business instance
func NewBusinessSetOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) BusinessSetOperationInterface {
	return &businessSet{
		clientSet:   client,
		authManager: authManager,
	}
}

type businessSet struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
}

// SetProxy SetProxy
func (b *businessSet) SetProxy(inst InstOperationInterface) {
	b.inst = inst
}

// CreateBusinessSet create business set
func (b *businessSet) CreateBusinessSet(kit *rest.Kit, data *metadata.CreateBizSetRequest) (mapstr.MapStr, error) {

	// 进行biz set name 的唯一性校验
	input := &metadata.CommonCountFilter{
		Conditions: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{
				Condition: querybuilder.ConditionAnd,
				Rules: []querybuilder.Rule{
					&querybuilder.AtomRule{
						Field:    common.BKAppSetNameField,
						Operator: querybuilder.OperatorEqual,
						Value:    data.BizSetAttr[common.BKAppSetNameField].(string)},
				},
			},
		},
	}
	result, err := b.inst.CountObjectInstances(kit, common.BKInnerObjIDAppSet, input)
	if err != nil {
		blog.Errorf("get biz_set name %s fail, err: %v, rid: %s", data.BizSetAttr[common.BKAppSetNameField].(string),
			err, kit.Rid)
		return nil, err
	}
	if result.Count > 0 {
		blog.Errorf("biz set name %s has been created, num: %d,err: %v, rid: %s",
			data.BizSetAttr[common.BKAppSetNameField].(string), result.Count, err, kit.Rid)
		return nil, fmt.Errorf("biz set name %s has been created", data.BizSetAttr[common.BKAppSetNameField].(string))
	}

	bizSetInfo := mapstr.New()

	for key, value := range data.BizSetAttr {
		bizSetInfo[key] = value
	}

	cond, errKey, err := data.BizSetScope.Filter.ToMgo()
	if err != nil {
		blog.Errorf(" biz set scope convert to mongo condition fail,scope: %+v, errKey: %s, err: %v, rid: %s",
			data.BizSetScope, errKey, err, kit.Rid)
		return mapstr.MapStr{}, err
	}

	bizSetInfo[common.BKAppSetScopeField] = cond
	bizInst, err := b.inst.CreateInst(kit, common.BKInnerObjIDAppSet, bizSetInfo)
	if err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}

	return bizInst, nil
}

// FindBizSet find biz set by condition.
func (b *businessSet) FindBizSet(kit *rest.Kit, cond *metadata.CommonSearchFilter) (
	*metadata.CommonSearchResult, error) {

	result, err := b.inst.SearchObjectInstances(kit, cond.ObjectID, cond)
	if err != nil {
		blog.Errorf("search biz set instances failed, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	return result, err
}

// CountBizSet count biz set by condition.
func (b *businessSet) CountBizSet(kit *rest.Kit, cond *metadata.CommonCountFilter) (
	*metadata.CommonCountResult, error) {

	result, err := b.inst.CountObjectInstances(kit, cond.ObjectID, cond)
	if err != nil {
		blog.Errorf("count biz set num failed, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	return result, err
}
