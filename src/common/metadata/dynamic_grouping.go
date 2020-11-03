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

package metadata

import (
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"github.com/google/uuid"
)

// operators support for dynamic group conditions in system.
const (
	// DynamicGroupOperatorEQ eq operator.
	DynamicGroupOperatorEQ = "$eq"

	// DynamicGroupOperatorNE ne operator.
	DynamicGroupOperatorNE = "$ne"

	// DynamicGroupOperatorIN in operator.
	DynamicGroupOperatorIN = "$in"

	// DynamicGroupOperatorNIN nin operator.
	DynamicGroupOperatorNIN = "$nin"

	// DynamicGroupOperatorLTE lte operator.
	DynamicGroupOperatorLTE = "$lte"

	// DynamicGroupOperatorGTE gte operator.
	DynamicGroupOperatorGTE = "$gte"

	// DynamicGroupOperatorLIKE like operator.
	DynamicGroupOperatorLIKE = "$regex"
)

var (
	// DynamicGroupOperators all operators -> current newest operators.
	DynamicGroupOperators = map[string]string{
		DynamicGroupOperatorEQ:   DynamicGroupOperatorEQ,
		DynamicGroupOperatorNE:   DynamicGroupOperatorNIN,
		DynamicGroupOperatorIN:   DynamicGroupOperatorIN,
		DynamicGroupOperatorNIN:  DynamicGroupOperatorNIN,
		DynamicGroupOperatorLTE:  DynamicGroupOperatorLTE,
		DynamicGroupOperatorGTE:  DynamicGroupOperatorGTE,
		DynamicGroupOperatorLIKE: DynamicGroupOperatorLIKE,
	}

	// DynamicGroupConditionTypes all condition object types of dynamic group.
	DynamicGroupConditionTypes = map[string]map[string]string{
		// host dynamic group.
		common.BKInnerObjIDHost: map[string]string{
			common.BKInnerObjIDSet:    common.BKInnerObjIDSet,
			common.BKInnerObjIDModule: common.BKInnerObjIDModule,
			common.BKInnerObjIDHost:   common.BKInnerObjIDHost,
		},

		// set dynamic group.
		common.BKInnerObjIDSet: map[string]string{
			common.BKInnerObjIDSet: common.BKInnerObjIDSet,
		},
	}
)

// Validatefunc is func callback for validating.
type Validatefunc func(objectID string) ([]Attribute, error)

// DynamicGroupCondition is target resource search condition on fields level.
type DynamicGroupCondition struct {
	// Field is target field name for index resource.
	Field string `json:"field" bson:"field"`

	// Operator is index operator type, eg $ne/$eq/$in/$nin.
	Operator string `json:"operator" bson:"operator"`

	// Value is target field value for index resource(integer or string).
	Value interface{} `json:"value" bson:"value"`
}

// Validate validates dynamic group conditions format.
func (c *DynamicGroupCondition) Validate(attributeMap map[string]string) error {
	operator, isSupport := DynamicGroupOperators[c.Operator]
	if !isSupport {
		return fmt.Errorf("not support operator, %s", c.Operator)
	}

	if c.Field == common.BKDefaultField {
		return nil
	}

	attributeType, isSupport := attributeMap[c.Field]
	if !isSupport {
		return fmt.Errorf("not support condition filed, %+v", c.Field)
	}

	attrType, err := getAttributeType(attributeType)
	if err != nil {
		return err
	}

	if c.Field == common.BKServiceTemplateIDField {
		if c.Operator != DynamicGroupOperatorIN {
			return fmt.Errorf("service template field only support $in operator, not support operator, %s", c.Operator)
		}
	}

	switch operator {
	case DynamicGroupOperatorEQ, DynamicGroupOperatorNE, DynamicGroupOperatorLTE, DynamicGroupOperatorGTE:
		return validAttributeValueType(attrType, c.Value)
	case DynamicGroupOperatorIN, DynamicGroupOperatorNIN:
		valueArr, ok := c.Value.([]interface{})
		if !ok {
			return fmt.Errorf("operator %s only support array value, not support value, %+v", c.Operator, c.Value)
		}

		for _, value := range valueArr {
			if err := validAttributeValueType(attrType, value); err != nil {
				return err
			}
		}
	case DynamicGroupOperatorLIKE:
		if attrType != stringType {
			return fmt.Errorf("operator %s only support string value, not support attribute type, %s", c.Operator, attributeType)
		}

		return validAttributeValueType(attrType, c.Value)
	}

	return nil
}

func validAttributeValueType(attrType string, value interface{}) error {
	switch attrType {
	case stringType:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("attribute only support string value, not support value, %+v", value)
		}
	case numericType:
		if !util.IsNumeric(value) {
			return fmt.Errorf("attribute only support numeric value, not support value, %+v", value)
		}
	case boolType:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("attribute only support bool value, not support value, %+v", value)
		}
	}

	return nil
}

const (
	numericType = "numeric"
	boolType    = "bool"
	stringType  = "string"
)

func getAttributeType(attributeType string) (string, error) {
	switch attributeType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeEnum, common.FieldTypeDate, common.FieldTypeTime,
		common.FieldTypeTimeZone, common.FieldTypeUser, common.FieldTypeList, common.FieldTypeOrganization:
		return stringType, nil
	case common.FieldTypeInt, common.FieldTypeFloat:
		return numericType, nil
	case common.FieldTypeBool:
		return boolType, nil
	default:
		return "", fmt.Errorf("not support attribute type, %s", attributeType)
	}
}

// DynamicGroupInfoCondition is condition for dynamic grouping, user could search
// target source base on the conditions.
type DynamicGroupInfoCondition struct {
	// ObjID is cmdb object id, could be host/set now.
	ObjID string `json:"bk_obj_id" bson:"bk_obj_id"`

	// Condition is search condition on fields level.
	// Example: bk_host_name $eq my-host just index host which name is "my-host".
	Condition []DynamicGroupCondition `json:"condition" bson:"condition"`
}

// Validate validates dynamic group info conditions format.
func (c *DynamicGroupInfoCondition) Validate(validatefunc Validatefunc) error {
	attributes, err := validatefunc(c.ObjID)
	if err != nil {
		return fmt.Errorf("validate dynamic group failed, %+v", err)
	}

	attributeMap := make(map[string]string)
	for _, attribute := range attributes {
		attributeMap[attribute.PropertyID] = attribute.PropertyType
	}

	switch c.ObjID {
	case common.BKInnerObjIDSet:
		attributeMap[common.BKSetIDField] = common.FieldTypeInt

	case common.BKInnerObjIDModule:
		attributeMap[common.BKModuleIDField] = common.FieldTypeInt

	case common.BKInnerObjIDHost:
		attributeMap[common.BKHostIDField] = common.FieldTypeInt
	}

	blog.Infof("validate info conditions, object[%s] attributes[%+v]", c.ObjID, attributeMap)

	for _, cond := range c.Condition {
		if err := cond.Validate(attributeMap); err != nil {
			return err
		}
	}
	return nil
}

// DynamicGroupInfo is info field in DynamicGroup struct.
type DynamicGroupInfo struct {
	// Condition is dynamic group index conditions set.
	Condition []DynamicGroupInfoCondition `json:"condition" bson:"condition"`
}

// Validate validates dynamic group info format, it's OK if conditions empty in this level.
func (c *DynamicGroupInfo) Validate(objectID string, validatefunc Validatefunc) error {
	types, isSupport := DynamicGroupConditionTypes[objectID]
	if !isSupport {
		return fmt.Errorf("not support dynamic group type, %s", objectID)
	}

	for _, cond := range c.Condition {
		if _, isSupport = types[cond.ObjID]; !isSupport {
			return fmt.Errorf("not support condition type[%s] for %s dynamic group", cond.ObjID, objectID)
		}

		if err := cond.Validate(validatefunc); err != nil {
			return err
		}
	}
	return nil
}

// DynamicGroup is dynamic grouping of conditions for host/set data searching.
type DynamicGroup struct {
	// AppID is application id which dynamic group belongs to.
	AppID int64 `json:"bk_biz_id" bson:"bk_biz_id"`

	// ID is dynamic group instance unique id.
	ID string `json:"id" bson:"id"`

	// Name is dynamic group name.
	Name string `json:"name" bson:"name"`

	// ObjID is cmdb object id, could be host/set now.
	ObjID string `json:"bk_obj_id" bson:"bk_obj_id"`

	// Info is dynamic group core conditions information.
	Info DynamicGroupInfo `json:"info" bson:"info"`

	// CreateUser create user name.
	CreateUser string `json:"create_user" bson:"create_user"`

	// ModifyUser modify user name.
	ModifyUser string `json:"modify_user" bson:"modify_user"`

	// CreateTime create timestamp.
	CreateTime time.Time `json:"create_time" bson:"create_time"`

	// UpdateTime last update timestamp.
	UpdateTime time.Time `json:"last_time" bson:"last_time"`
}

// Validate validates dynamic group format.
func (g *DynamicGroup) Validate(validatefunc Validatefunc) error {
	if g.AppID <= 0 {
		return errors.New("empty bk_biz_id")
	}

	if len(g.Name) == 0 {
		return errors.New("empty name")
	}

	// check object id.
	if len(g.ObjID) == 0 {
		return errors.New("empty bk_obj_id")
	}

	// check conditions format.
	if len(g.Info.Condition) == 0 {
		// it's not OK if conditions empty in this level.
		return errors.New("empty info.condition")
	}
	return g.Info.Validate(g.ObjID, validatefunc)
}

// DynamicGroupResultBatch is batch result struct of dynamic group.
type DynamicGroupBatch struct {
	// Count batch count.
	Count uint64 `json:"count"`

	// Info batch data.
	Info []DynamicGroup `json:"info"`
}

// SearchDynamicGroupResult is result struct for dynamic group searching action.
type SearchDynamicGroupResult struct {
	BaseResp `json:",inline"`
	Data     DynamicGroupBatch `json:"data"`
}

// GetDynamicGroupResult is result struct for dynamic group detail query action.
type GetDynamicGroupResult struct {
	BaseResp `json:",inline"`
	Data     DynamicGroup `json:"data"`
}

// NewDynamicGroupID creates and returns a new dynamic group string unique ID.
func NewDynamicGroupID() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
