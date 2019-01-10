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

package types

import (
	"time"

	"configcenter/src/framework/core/types"
)

type CreateModelCtx struct {
	BaseCtx
	ModelInfo CreateModelInfo
}

type CreateModelInfo struct {
	// name of creator of this model
	Creator string `json:"creator"`
	// the class id that this model belongs to.
	// should not be empty.
	ClassID string `json:"bk_classification_id"`
	// object id
	ID string `json:"bk_obj_id"`
	// object name
	Name    string `json:"bk_obj_name"`
	Tenancy string `json:"bk_supplier_account"`
	Icon    string `json:"bk_obj_icon"`
}

type CreateModelResponse struct {
	BaseResp `json:",inline"`
	Data     struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

type DeleteModelCtx struct {
	BaseCtx
	ModelID int64
}

type UpdateModelCtx struct {
	BaseCtx
	ModelID   int64
	ModelInfo UpdateModelInfo
}

type UpdateModelInfo struct {
	ID       int64  `json:"id,omitempty"`
	Modifier string `json:"modifier,omitempty"`
	// required field, can not be empty.
	ClassID string `json:"bk_classification_id"`
	Name    string `json:"bk_obj_name,omitempty"`
	// required field, can not be empty.
	Tenancy  string `json:"bk_supplier_account"`
	Icon     string `json:"bk_obj_icon,omitempty"`
	Position string `json:"position,omitempty"`
}

type GetModelsCtx struct {
	BaseCtx
	Filters types.MapStr
}

type GetModelsResult struct {
	BaseResp `json:",inline"`
	Data     []ModelInfo `json:"data"`
}

type ModelInfo struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"bk_obj_name"`
	ClassID     string    `json:"bk_classification_id"`
	ObjectID    string    `json:"bk_object_id"`
	Tenancy     string    `json:"bk_supplier_account"`
	Creator     string    `json:"creator"`
	Modifier    string    `json:"modifier"`
	Description string    `json:"description"`
	IsPaused    bool      `json:"bk_ispaused"`
	IsPre       bool      `json:"ispre"`
	Icon        string    `json:"bk_obj_icon"`
	Position    string    `json:"position"`
	LastTime    time.Time `json:"last_time"`
	CreateTime  time.Time `json:"create_time"`
}
