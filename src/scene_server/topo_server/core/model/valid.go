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
	"strings"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
)

// FieldValid field valid method
type FieldValid struct {
	lang language.DefaultCCLanguageIf
}

// Valid valid the field
func (f *FieldValid) Valid(kit *rest.Kit, data mapstr.MapStr, fieldID string) (string, error) {

	val, err := data.String(fieldID)
	if nil != err {
		return val, kit.CCError.New(common.CCErrCommParamsIsInvalid, fieldID+" "+err.Error())
	}
	if 0 == len(val) {
		return val, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, fieldID)
	}

	return val, nil
}

// ValidID check the property ID
func (f *FieldValid) ValidID(kit *rest.Kit, value string) error {
	if common.AttributeIDMaxLength < utf8.RuneCountInString(value) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			f.lang.Language("model_attr_bk_property_id"), common.AttributeIDMaxLength)
	}
	match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, value)
	if nil != err {
		return err
	}

	if !match {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, value)
	}

	return nil
}

// ValidName check the name
func (f *FieldValid) ValidName(kit *rest.Kit, value string) error {
	if common.AttributeNameMaxLength < utf8.RuneCountInString(value) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			f.lang.Language("model_attr_bk_property_name"), common.AttributeNameMaxLength)
	}
	value = strings.TrimSpace(value)
	return nil
}

// ValidPlaceHolder check the PlaceHolder
func (f *FieldValid) ValidPlaceHolder(kit *rest.Kit, value string) error {
	if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(value) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			f.lang.Language("model_attr_placeholder"), common.AttributePlaceHolderMaxLength)
	}
	return nil
}

// ValidNameWithRegex valid by regex
func (f *FieldValid) ValidNameWithRegex(kit *rest.Kit, value string) error {

	if err := f.ValidName(kit, value); nil != err {
		return err
	}

	/*
			match, err := regexp.MatchString(`^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[()+-《》_,，；;“”‘’。."\' \\/]){1,20}$`, value)
			if nil != err {
				fmt.Println("dd:", err.Error(), value)
				return err
			}

		if !match {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, value)
		}
	*/
	return nil
}
