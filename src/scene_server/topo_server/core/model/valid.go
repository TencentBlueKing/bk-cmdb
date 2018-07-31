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
	"regexp"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

// FieldValid field valid method
type FieldValid struct {
}

// Valid valid the field
func (f *FieldValid) Valid(params types.ContextParams, data mapstr.MapStr, fieldID string) (string, error) {

	val, err := data.String(fieldID)
	if nil != err {
		return val, params.Err.New(common.CCErrCommParamsIsInvalid, fieldID+" "+err.Error())
	}
	if 0 == len(val) {
		return val, params.Err.Errorf(common.CCErrCommParamsNeedSet, fieldID)
	}

	return val, nil
}

// ValidID check the property ID
func (f *FieldValid) ValidID(params types.ContextParams, value string) error {

	match, err := regexp.MatchString(`[a-z\d_]+`, value)
	if nil != err {
		return err
	}

	if !match {
		return params.Err.Errorf(common.CCErrCommParamsIsInvalid, value)
	}

	return nil
}

// ValidName check the name
func (f *FieldValid) ValidName(params types.ContextParams, value string) error {
	if 20 < utf8.RuneCountInString(value) {
		return params.Err.Errorf(common.CCErrCommOverLimit, value)
	}
	return nil
}

// ValidNameWithRegex valid by regex
func (f *FieldValid) ValidNameWithRegex(params types.ContextParams, value string) error {

	if err := f.ValidName(params, value); nil != err {
		return err
	}

	/*
			match, err := regexp.MatchString(`^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[()+-《》_,，；;“”‘’。."\' \\/]){1,20}$`, value)
			if nil != err {
				fmt.Println("dd:", err.Error(), value)
				return err
			}

		if !match {
			return params.Err.Errorf(common.CCErrCommParamsIsInvalid, value)
		}
	*/
	return nil
}
