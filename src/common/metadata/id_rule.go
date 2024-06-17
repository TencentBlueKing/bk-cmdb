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
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/util"
)

// AsstIDOption asst id option
type AsstIDOption struct {
	Rule string `json:"rule"`
}

// ParseAsstIDOption parse asst id option
func ParseAsstIDOption(val interface{}) (*AsstIDOption, error) {
	if val == nil || val == "" {
		return nil, fmt.Errorf("option val is invalid")
	}

	switch option := val.(type) {
	case map[string]interface{}:
		rule, ok := option["rule"]
		if !ok {
			return nil, fmt.Errorf("option val is invalid")
		}
		return &AsstIDOption{Rule: util.GetStrByInterface(rule)}, nil
	case AsstIDOption:
		return &option, nil
	case string:
		res := new(AsstIDOption)
		err := json.Unmarshal([]byte(option), res)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("unknow val type: %T for asstid option", val)
	}
}

// RuleKind rule kind
type RuleKind string

const (
	// Const rules for constant types
	Const RuleKind = "const"
	// Attr rules for attribute types
	Attr RuleKind = "attr"
	// GlobalID rules for global id types
	GlobalID RuleKind = "globalID"
	// LocalID rules for local model id types
	LocalID RuleKind = "localID"
	// RandomID rules for random id types
	RandomID RuleKind = "randomID"
)

// SubAssetRule sub asset rule
type SubAssetRule struct {
	Val  string
	Kind RuleKind
	Len  int64
}

const (
	idRuleVarLimit = 4
	idLenLimit     = 32
	varTypeLimit   = 1

	// IDRuleFieldLimit object id rule field limit
	IDRuleFieldLimit = 1
)

// ParseSubIDRules parse sub asset rule
func ParseSubIDRules(val interface{}) ([]SubAssetRule, error) {
	option, err := ParseAsstIDOption(val)
	if err != nil {
		return nil, err
	}
	rule := option.Rule

	pattern := `\{\{[^\}]+\}\}` // 正则表达式，用于匹配 {{xxx}} 数据
	re := regexp.MustCompile(pattern)
	indexes := re.FindAllStringIndex(rule, -1) // 查找匹配的变量在原字符串中的位置

	// id规则变量个数限制
	if len(indexes) > idRuleVarLimit {
		return nil, fmt.Errorf("option.rule var count:%d exceed max count:%d", indexes, idRuleVarLimit)
	}

	lastIdx := 0
	result := make([]SubAssetRule, 0)
	globalVarCount, localVarCount, randVarCount := 0, 0, 0

	for _, idx := range indexes {
		// 处理普通的字符串
		varConst := rule[lastIdx:idx[0]]
		result = append(result, SubAssetRule{Val: varConst, Kind: Const, Len: int64(len(varConst))})
		lastIdx = idx[1]

		// 处理变量, 去掉左右括号
		varStr := rule[idx[0]+2 : idx[1]-2]
		split := strings.Split(varStr, ".")
		// 引用模型的其他字段
		if len(split) == 1 {
			result = append(result, SubAssetRule{Val: varStr, Kind: Attr})
			continue
		}

		varVal := strings.Join(split[0:len(split)-1], ".")
		var kind RuleKind
		switch varVal {
		case common.GlobalIncrIDVar:
			kind = GlobalID
			globalVarCount++
		case common.LocalIncrIDVar:
			kind = LocalID
			localVarCount++
		case common.RandomIDVar:
			kind = RandomID
			randVarCount++
		default:
			return nil, fmt.Errorf("option var %s is invalid", varVal)
		}

		// 随机ID、自增id等长度限制
		length, err := strconv.ParseInt(split[len(split)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		if length > idLenLimit {
			return nil, fmt.Errorf("option.rule var %s length exceed max length %d", varVal, idLenLimit)
		}

		result = append(result, SubAssetRule{Val: varVal, Kind: kind, Len: length})
	}

	if globalVarCount == 0 && localVarCount == 0 && randVarCount == 0 {
		return nil, fmt.Errorf("option.rule has no %s, %s and %s", GlobalID, LocalID, RandomID)
	}

	if globalVarCount > varTypeLimit || localVarCount > varTypeLimit || randVarCount > varTypeLimit {
		return nil, fmt.Errorf("option.rule var type exceed max count %d", idLenLimit)
	}

	lastVal := rule[lastIdx:]
	if lastVal != "" {
		result = append(result, SubAssetRule{Val: lastVal})
	}

	return result, nil
}

// GetIDRuleRandomID get id rule random id
func GetIDRuleRandomID(length int64) string {
	rand.Seed(time.Now().UnixNano())
	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var sb strings.Builder
	for i := int64(0); i < length; i++ {
		digit := digits[rand.Intn(len(digits))]
		sb.WriteString(digit)
	}

	return sb.String()
}

// MakeUpDigit make up the number of digits
func MakeUpDigit(id uint64, length int64) (string, error) {
	idStr := strconv.FormatUint(id, 10)
	idLen := int64(len(idStr))
	if idLen == length {
		return idStr, nil
	}

	if idLen > length {
		return "", fmt.Errorf("id exceed length limit: %d", length)
	}

	return strings.Repeat("0", int(length-idLen)) + idStr, nil
}

// UpdateIDGenOption update id generator option
type UpdateIDGenOption struct {
	Type       string `json:"type"`
	SequenceID int64  `json:"sequence_id"`
}

// Validate validate UpdateIDGenOption
func (u *UpdateIDGenOption) Validate() ccErr.RawErrorInfo {
	if u.Type == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	if u.SequenceID <= 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsIsInvalid, Args: []interface{}{common.BKFieldSeqID}}
	}

	return ccErr.RawErrorInfo{}
}

// SyncIDRuleOption sync id rule option
type SyncIDRuleOption struct {
	ObjID string `json:"bk_obj_id"`
}

// Validate validate SyncIDRuleOption
func (s *SyncIDRuleOption) Validate() ccErr.RawErrorInfo {
	if s.ObjID == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKObjIDField}}
	}

	return ccErr.RawErrorInfo{}
}

// SyncIDRuleRes sync id rule result
type SyncIDRuleRes struct {
	TaskID string `json:"task_id"`
}

// UpdateInstIDRuleOption update instance id rule field option
type UpdateInstIDRuleOption struct {
	ObjID      string  `json:"bk_obj_id"`
	IDs        []int64 `json:"ids"`
	PropertyID string  `json:"bk_property_id"`
}

// Validate validate UpdateInstIDRuleOption
func (u *UpdateInstIDRuleOption) Validate() ccErr.RawErrorInfo {
	if u.ObjID == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKObjIDField}}
	}

	if u.PropertyID == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKPropertyIDField}}
	}

	if len(u.IDs) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(u.IDs) > common.BKMaxPageSize {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit,
			Args: []interface{}{"ids", common.BKMaxPageSize}}
	}

	return ccErr.RawErrorInfo{}
}

// IDRuleTaskOption id rule task option
type IDRuleTaskOption struct {
	TaskID string `json:"task_id"`
}

// Validate validate IDRuleTaskOption
func (a *IDRuleTaskOption) Validate() ccErr.RawErrorInfo {
	if a.TaskID == "" {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKTaskIDField}}
	}

	return ccErr.RawErrorInfo{}
}

// GetIDRule 获取对应id rule自增id的唯一标识，目前bk_obj_id唯一，后续涉及到多租户，可能需要调整
func GetIDRule(flag string) string {
	return fmt.Sprintf("%s%s", common.IDRulePrefix, flag)
}

// IsValidAttrRuleType check if attribute rule type is valid
func IsValidAttrRuleType(typ string) bool {
	switch typ {
	case common.FieldTypeSingleChar, common.FieldTypeEnum, common.FieldTypeList:
		return true
	}
	return false
}
