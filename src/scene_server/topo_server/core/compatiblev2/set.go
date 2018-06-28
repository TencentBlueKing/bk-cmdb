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

package compatiblev2

import (
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

// SetInterface set interface
type SetInterface interface {
	UpdateMultiSet(bizID int64, data mapstr.MapStr, cond condition.Condition) error
	DeleteMultiSet(bizID int64, setIDS []int64, cond condition.Condition) error
	DeleteSetHost()
}

// NewSet ceate a new set instance
func NewSet() SetInterface {
	return &set{}
}

type set struct {
	params types.LogicParams
	client apimachinery.ClientSetInterface
}

func (s *set) UpdateMultiSet(bizID int64, data mapstr.MapStr, cond condition.Condition) error {

	input := mapstr.New()
	input.Set("data", data)
	input.Set("condition", cond.ToMapStr())

	rsp, err := s.client.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDSet, s.params.Header.ToHeader(), input)
	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, error info is %s", err.Error())
		return s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to update the set, error info is %s", err.Error())
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (s *set) DeleteMultiSet(bizID int64, setIDS []int64, cond condition.Condition) error {
	return nil
}
func (s *set) DeleteSetHost() {

}
