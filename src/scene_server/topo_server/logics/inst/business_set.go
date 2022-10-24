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
)

// BusinessSetOperationInterface business operation methods
type BusinessSetOperationInterface interface {

	// CreateBusinessSet create business set
	CreateBusinessSet(kit *rest.Kit, data *metadata.CreateBizSetRequest) (mapstr.MapStr, error)

	// SetProxy set business proxy
	SetProxy(inst InstOperationInterface)
}

// NewBusinessSetOperation create a business instance
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

	conditions := &metadata.Condition{
		Condition: map[string]interface{}{common.BKBizSetNameField: data.BizSetAttr[common.BKBizSetNameField]},
	}

	// count object instances num.
	resp, err := b.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, common.BKInnerObjIDBizSet,
		conditions)
	if err != nil {
		blog.Errorf("count object instances failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	if resp.Count > 0 {
		blog.Errorf("biz set name has been created, num: %d, err: %v, rid: %s", resp.Count, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrorBizSetNameDuplicated)
	}

	data.BizSetAttr[common.BKBizSetScopeField] = data.BizSetScope

	bizSetInst, err := b.inst.CreateInst(kit, common.BKInnerObjIDBizSet, data.BizSetAttr)
	if err != nil {
		blog.Errorf("create business failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	return bizSetInst, nil
}
