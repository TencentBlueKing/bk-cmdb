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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ATTENTIONS: the dependent methods of the other module

// OperationDependences methods definition
type OperationDependences interface {
	// IsInstAsstExist used to check if the  instances  asst exist
	IsInstAsstExist(kit *rest.Kit, objID string, instID uint64) (exists bool, err error)

	// DeleteInstAsst used to delete inst asst
	DeleteInstAsst(kit *rest.Kit, objID string, instID uint64) error

	// SelectObjectAttWithParams select object att with params
	SelectObjectAttWithParams(kit *rest.Kit, objID string, bizIDs []int64) (attribute []metadata.Attribute, err error)

	// SelectObjectAttributes select object attributes
	SelectObjectAttributes(kit *rest.Kit, objID string, bizIDs []int64) ([]metadata.Attribute, error)

	// SearchUnique search unique attribute
	SearchUnique(kit *rest.Kit, objID string) (uniqueAttr []metadata.ObjectUnique, err error)

	// DeleteQuotedInst delete quoted instances by source instance ids
	DeleteQuotedInst(kit *rest.Kit, objID string, instIDs []int64) error

	// AttachQuotedInst attach quoted instances with source instance
	AttachQuotedInst(kit *rest.Kit, objID string, instID uint64, data mapstr.MapStr) error
}
