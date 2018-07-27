/*
 * Tencent is pleased to support the open source community by making čé˛¸ available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

// FieldValid field valid method
type FieldValid struct {
}

// Valid valid the field
func (f *FieldValid) Valid(params types.ContextParams, data mapstr.MapStr, fieldID string) error {

	val, err := data.String(fieldID)
	if nil != err {
		return params.Err.Errorf(common.CCErrCommParamsIsInvalid, fieldID)
	}
	if 0 == len(val) {
		return params.Err.Errorf(common.CCErrCommParamsNeedSet, fieldID)
	}

	return nil
}
