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

package metadata

import (
	"configcenter/src/common/mapstr"
)

// CreateManyModelClassifiaction create many input params
type CreateManyModelClassifiaction struct {
	Data []Classification `json:"datas"`
}

// CreateOneModelClassification create one model classification
type CreateOneModelClassification struct {
	Data Classification `json:"data"`
}

// SetManyModelClassification set many input params
type SetManyModelClassification CreateManyModelClassifiaction

// SetOneModelClassification set one input params
type SetOneModelClassification CreateOneModelClassification

// DeleteModelClassificationResult delete the model classification result
type DeleteModelClassificationResult struct {
	BaseResp `json:",inline"`
	Data     DeletedCount `json:"data"`
}

// CreateModel create model params
type CreateModel struct {
	Spec       ObjectDes   `json:"spec"`
	Attributes []Attribute `json:"attributes"`
}

// SetModel define SetMode method input params
type SetModel CreateModel

// SearchModelInfo search  model params
type SearchModelInfo struct {
	Spec       mapstr.MapStr   `json:"spec"`
	Attributes []mapstr.MapStr `json:"attributes"`
}

// CreateModelAttributes create model attributes
type CreateModelAttributes struct {
	Attributes []Attribute `json:"attributes"`
}

type SetModelAttributes CreateModelAttributes

type CreateModelInstance struct {
	Data mapstr.MapStr `json:"data"`
}

type CreateManyModelInstance struct {
	Datas []mapstr.MapStr `json:"datas"`
}

type SetModelInstance CreateModelInstance
type SetManyModelInstance CreateManyModelInstance

type CreateAssociationKind struct {
	Data AssociationKind `json:"data"`
}

type CreateManyAssociationKind struct {
	Data []AssociationKind `json:"datas"`
}
type SetAssociationKind CreateAssociationKind
type SetManyAssociationKind CreateManyAssociationKind

type CreateModelAssociation struct {
	Spec Association `json:"spec"`
}

type SetModelAssociation CreateModelAssociation

type CreateOneInstanceAssociation struct {
	Datas InstAsst `json:"data"`
}
type CreateManyInstanceAssociation struct {
	Datas InstAsst `json:"datas"`
}

type SetOneInstanceAssociation CreateOneInstanceAssociation
type SetManyInstanceAssociation CreateManyInstanceAssociation
