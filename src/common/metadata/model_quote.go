/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
)

// ModelQuoteRelation model relationship table information.
// For example, the disk table field on the host is a table type,
// then DestModel is bk_host_disk, SrcModel is host, PropertyID is
// disk, Type is table.
type ModelQuoteRelation struct {
	// DestModel derived model
	DestModel string `json:"dest_model" bson:"dest_model"`
	// SrcModel source model, such as a tabular model
	SrcModel string `json:"src_model" bson:"src_model"`
	// PropertyID model attribute id
	PropertyID string `json:"bk_property_id" bson:"bk_property_id"`
	// Type the specific type of model such as the table
	Type common.ModelQuoteType `json:"type" bson:"type"`
	// SupplierAccount supplier account
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// ListModelQuoteRelRes list model quote relationship table result.
type ListModelQuoteRelRes struct {
	Count uint64               `json:"count"`
	Info  []ModelQuoteRelation `json:"info"`
}

// ListModelQuoteRelResp list model quote relationship table response.
type ListModelQuoteRelResp struct {
	BaseResp `json:",inline"`
	Data     *ListModelQuoteRelRes `json:"data"`
}

// BatchCreateQuotedInstOption batch create quoted instance option
type BatchCreateQuotedInstOption struct {
	ObjID      string          `json:"bk_obj_id"`
	PropertyID string          `json:"bk_property_id"`
	Data       []mapstr.MapStr `json:"data"`
}

// Validate batch create quoted instance option
func (c *BatchCreateQuotedInstOption) Validate() errors.RawErrorInfo {
	if c.ObjID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKObjIDField}}
	}

	if c.PropertyID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyIDField}}
	}

	if len(c.Data) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"data"}}
	}

	if len(c.Data) > common.BKMaxWriteOpLimit {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"data"}}
	}

	return errors.RawErrorInfo{}
}

// BatchUpdateQuotedInstOption batch update quoted instance option
type BatchUpdateQuotedInstOption struct {
	ObjID      string        `json:"bk_obj_id"`
	PropertyID string        `json:"bk_property_id"`
	IDs        []uint64      `json:"ids"`
	Data       mapstr.MapStr `json:"data"`
}

// Validate batch update quoted instance option
func (c *BatchUpdateQuotedInstOption) Validate() errors.RawErrorInfo {
	if c.ObjID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKObjIDField}}
	}

	if c.PropertyID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyIDField}}
	}

	if len(c.IDs) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(c.IDs) > common.BKMaxWriteOpLimit {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(c.Data) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"data"}}
	}

	return errors.RawErrorInfo{}
}

// BatchDeleteQuotedInstOption batch delete quoted instance option
type BatchDeleteQuotedInstOption struct {
	ObjID      string   `json:"bk_obj_id"`
	PropertyID string   `json:"bk_property_id"`
	IDs        []uint64 `json:"ids"`
}

// Validate batch delete quoted instance option
func (c *BatchDeleteQuotedInstOption) Validate() errors.RawErrorInfo {
	if c.ObjID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKObjIDField}}
	}

	if c.PropertyID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyIDField}}
	}

	if len(c.IDs) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(c.IDs) > common.BKMaxDeletePageSize {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	return errors.RawErrorInfo{}
}

// ListQuotedInstOption list quoted instance option
type ListQuotedInstOption struct {
	ObjID             string `json:"bk_obj_id"`
	PropertyID        string `json:"bk_property_id"`
	CommonQueryOption `json:",inline"`
}

// Validate list quoted instance option
func (c *ListQuotedInstOption) Validate() errors.RawErrorInfo {
	if c.ObjID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKObjIDField}}
	}

	if c.PropertyID == "" {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyIDField}}
	}

	if err := c.CommonQueryOption.Validate(); err.ErrCode != 0 {
		return err
	}

	return errors.RawErrorInfo{}
}

// GenerateModelQuoteObjID generate the obj id referenced by the model.
func GenerateModelQuoteObjID(srcModel, propertyID string) string {
	return "bk_" + srcModel + "#" + propertyID
}

// GenerateModelQuoteObjName generate the obj name referenced by the model.
func GenerateModelQuoteObjName(srcModel, propertyName string) string {
	return "bk_" + srcModel + "#" + propertyName
}
