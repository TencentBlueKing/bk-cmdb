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

package instances

import (
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

// ATTENTIONS: the dependent methods of the other module

// OperationDependences methods definition
type OperationDependences interface {

	// IsInstanceExist used to check if the  instances  asst exist
	IsInstAsstExist(ctx core.ContextParams, objID string, instID uint64) (exists bool, err error)

	// DeleteInstAsst used to delete inst asst
	DeleteInstAsst(ctx core.ContextParams, objID string, instID uint64) error

	// SelectObjectAttWithParams select object att with params
	SelectObjectAttWithParams(ctx core.ContextParams, objID string, bizID int64) (attribute []metadata.Attribute, err error)

	// SearchUnique search unique attribute
	SearchUnique(ctx core.ContextParams, objID string) (uniqueAttr []metadata.ObjectUnique, err error)
}
