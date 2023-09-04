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
	"configcenter/pkg/filter"
	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
)

// CommonFilterOption common filter option.
type CommonFilterOption struct {
	Filter *filter.Expression `json:"filter"`
}

// Validate common filter option.
func (c *CommonFilterOption) Validate(opt ...*filter.ExprOption) ccErr.RawErrorInfo {
	if c.Filter == nil {
		return ccErr.RawErrorInfo{}
	}

	var op *filter.ExprOption
	if len(opt) != 0 {
		op = opt[0]
	} else {
		op = filter.NewDefaultExprOpt(nil)
		op.IgnoreRuleFields = true
	}

	if err := c.Filter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return ccErr.RawErrorInfo{}
}

// ToMgo convert common filter option to mongodb filter.
func (c *CommonFilterOption) ToMgo() (map[string]interface{}, error) {
	if c.Filter == nil {
		return make(map[string]interface{}), nil
	}

	mgo, err := c.Filter.ToMgo()
	if err != nil {
		return nil, err
	}

	return mgo, nil
}

// CommonQueryOption common query option.
type CommonQueryOption struct {
	CommonFilterOption `json:",inline"`
	Page               BasePage `json:"page"`
	Fields             []string `json:"fields"`
}

// Validate common query option.
func (c *CommonQueryOption) Validate(opt ...*filter.ExprOption) ccErr.RawErrorInfo {
	if err := c.Page.ValidateWithEnableCount(false, common.BKMaxPageSize); err.ErrCode != 0 {
		return err
	}

	if err := c.CommonFilterOption.Validate(opt...); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// CommonUpdateOption common update option.
type CommonUpdateOption struct {
	CommonFilterOption `json:",inline"`
	Data               mapstr.MapStr `json:"data"`
}

// Validate common update option.
func (c *CommonUpdateOption) Validate(opt ...*filter.ExprOption) ccErr.RawErrorInfo {
	if err := c.CommonFilterOption.Validate(opt...); err.ErrCode != 0 {
		return err
	}

	if len(c.Data) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"data"}}
	}

	return ccErr.RawErrorInfo{}
}
