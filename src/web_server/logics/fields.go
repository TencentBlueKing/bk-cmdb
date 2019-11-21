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

package logics

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
)

// Property object fields
type Property struct {
	ID            string
	Name          string
	PropertyType  string
	Option        interface{}
	IsPre         bool
	IsRequire     bool
	Group         string
	Index         int64
	ExcelColIndex int
	NotObjPropery bool //Not an attribute of the object, indicating that the field to be exported is needed for export,
	IsOnly        bool
	AsstObjID     string
	NotExport     bool
}

// PropertyGroup property group
type PropertyGroup struct {
	Name  string
	Index int64
	ID    string
}

type PropertyPrimaryVal struct {
	ID     string
	Name   string
	StrVal string
}

// GetObjFieldIDs get object fields
func (lgc *Logics) GetObjFieldIDs(objID string, filterFields []string, customFields []string, header http.Header, meta *metadata.Metadata) (map[string]Property, error) {

	fields, err := lgc.getObjFieldIDs(objID, header, meta)
	if nil != err {
		return nil, fmt.Errorf("get object fields failed, err: %+v", err)
	}
	groups, err := lgc.getObjectGroup(objID, header, meta)
	if nil != err {
		return nil, fmt.Errorf("get attribute group failed, err: %+v", err)
	}
	if len(groups) == 0 {
		return nil, errors.New("get attribute group by object not found")
	}

	ret := make(map[string]Property)

	// user specified export model field sort start index
	var ustomFieldStartIndex int
	//  sort the start of normal of the model
	var normalFieldStartIndex int
	// calculate the length of the user-psecified field
	for _, field := range fields {
		// filter fields cannot exported
		if util.InStrArr(filterFields, field.ID) {
			continue
		}
		// filter fields than do not exist
		if !util.InStrArr(customFields, field.ID) {
			continue
		}
		normalFieldStartIndex++
	}

	for _, group := range groups {
		for _, field := range fields {
			if field.Group == group.ID {
				if util.InStrArr(filterFields, field.ID) {
					field.NotExport = true
				} else if util.InStrArr(customFields, field.ID) {
					field.ExcelColIndex = ustomFieldStartIndex
					ustomFieldStartIndex++
				} else {
					field.ExcelColIndex = normalFieldStartIndex
					normalFieldStartIndex++
				}
				ret[field.ID] = field

			}
		}
	}
	return ret, nil
}

func (lgc *Logics) getObjectGroup(objID string, header http.Header, meta *metadata.Metadata) ([]PropertyGroup, error) {
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)
	condition := mapstr.MapStr{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: common.BKDefaultOwnerID,
		"page": mapstr.MapStr{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  common.BKPropertyGroupIndexField,
		},
		metadata.BKMetadata: meta,
	}
	result, err := lgc.Engine.CoreAPI.ApiServer().GetObjectGroup(context.Background(), header, ownerID, objID, condition)
	if nil != err {
		blog.Errorf("get %s fields group failed, err:%+v, rid: %s", objID, err, rid)
		return nil, fmt.Errorf("get attribute group failed, err: %+v", err)
	}
	if !result.Result {
		blog.Errorf("get %s fields group result failed. error code:%d, error message:%s, rid:%s", objID, result.Code, result.ErrMsg, rid)
		return nil, fmt.Errorf("get attribute group result false, result: %+v", result)
	}
	fields := result.Data
	ret := make([]PropertyGroup, 0)
	for _, mapField := range fields {
		propertyGroup := PropertyGroup{}
		propertyGroup.Index = mapField.GroupIndex
		propertyGroup.Name = mapField.GroupName
		propertyGroup.ID = mapField.GroupID
		ret = append(ret, propertyGroup)
	}
	blog.V(5).Infof("getObjectGroup count:%d, rid: %s", len(ret), rid)
	return ret, nil

}

func (lgc *Logics) getObjectPrimaryFieldByObjID(objID string, header http.Header, meta *metadata.Metadata) ([]Property, error) {
	fields, err := lgc.getObjFieldIDsBySort(objID, common.BKPropertyIDField, header, nil, meta)
	if nil != err {
		return nil, err
	}
	var ret []Property
	for _, field := range fields {
		if true == field.IsOnly {
			ret = append(ret, field)
		}
	}
	return ret, nil

}

func (lgc *Logics) getObjFieldIDs(objID string, header http.Header, meta *metadata.Metadata) ([]Property, error) {
	sort := fmt.Sprintf("-%s,bk_property_index", common.BKIsRequiredField)
	return lgc.getObjFieldIDsBySort(objID, sort, header, nil, meta)

}

func (lgc *Logics) getObjFieldIDsBySort(objID, sort string, header http.Header, conds mapstr.MapStr, meta *metadata.Metadata) ([]Property, error) {
	rid := util.GetHTTPCCRequestID(header)

	condition := mapstr.MapStr{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: util.GetOwnerID(header),
		metadata.PageName: mapstr.MapStr{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  sort,
		},
		metadata.BKMetadata: meta,
	}
	condition.Merge(conds)

	result, err := lgc.Engine.CoreAPI.ApiServer().GetObjectAttr(context.Background(), header, condition)
	if nil != err {
		blog.Errorf("getObjFieldIDsBySort get %s fields input:%s, error:%s ,rid:%s", objID, conds, err.Error(), rid)
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("getObjFieldIDsBySort get %s fields input:%s,  http reply info,error code:%d, error msg:%s ,rid:%s", objID, conds, result.Code, result.ErrMsg, rid)
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(result.Code, result.ErrMsg)
	}

	inputParam := metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Offset: 0,
			Limit:  common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKObjIDField: objID,
		}),
	}
	uniques, err := lgc.CoreAPI.CoreService().Model().ReadModelAttrUnique(context.Background(), header, inputParam)
	if nil != err {
		blog.Errorf("getObjectPrimaryFieldByObjID get unique for %s error: %v ,rid:%s", objID, err, rid)
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !uniques.Result {
		blog.Errorf("getObjectPrimaryFieldByObjID get unique for %s error: %v ,rid:%s", objID, uniques, rid)
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(uniques.Code, uniques.ErrMsg)
	}

	keyIDs := map[uint64]bool{}
	for _, unique := range uniques.Data.Info {
		if unique.MustCheck {
			for _, key := range unique.Keys {
				keyIDs[key.ID] = true
			}
			break
		}
	}
	if len(keyIDs) <= 0 {
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Error(common.CCErrTopoObjectUniqueSearchFailed)
	}

	ret := []Property{}
	for _, mapField := range result.Data {
		fieldIsOnly := keyIDs[uint64(mapField.ID)]
		fieldName := mapField.PropertyName
		fieldID := mapField.PropertyID
		fieldType := mapField.PropertyType
		fieldIsRequire := mapField.IsRequired
		fieldIsOption := mapField.Option
		fieldIsPre := mapField.IsPre
		fieldGroup := mapField.PropertyGroup
		fieldIndex := mapField.PropertyIndex

		ret = append(ret, Property{
			ID:           fieldID,
			Name:         fieldName,
			PropertyType: fieldType,
			IsRequire:    fieldIsRequire,
			IsPre:        fieldIsPre,
			Option:       fieldIsOption,
			Group:        fieldGroup,
			Index:        fieldIndex,
			IsOnly:       fieldIsOnly,
		})
	}
	blog.V(5).Infof("getObjFieldIDsBySort ret count:%d, rid: %s", len(ret), rid)
	return ret, nil
}

// getPropertyTypeAliasName  return propertyType name, whether to export,
func getPropertyTypeAliasName(propertyType string, defLang lang.DefaultCCLanguageIf) (string, bool) {
	var skip bool
	name := defLang.Language("field_type_" + propertyType)
	switch propertyType {
	case common.FieldTypeSingleChar:
	case common.FieldTypeLongChar:
	case common.FieldTypeInt:
	case common.FieldTypeFloat:
	case common.FieldTypeEnum:
	case common.FieldTypeDate:
	case common.FieldTypeTime:
	case common.FieldTypeUser:
	case common.FieldTypeBool:
	case common.FieldTypeTimeZone:

	}
	if "" == name {
		name = propertyType
	}
	return name, skip
}

// addSystemField add system field, get property not return property fields
func addSystemField(fields map[string]Property, objID string, defLang lang.DefaultCCLanguageIf) {
	for key, field := range fields {
		field.ExcelColIndex = field.ExcelColIndex + 1
		fields[key] = field
	}

	idProperty := Property{
		ID:            "",
		Name:          "",
		PropertyType:  common.FieldTypeInt,
		Group:         "defalut",
		ExcelColIndex: 0,
	}

	switch objID {
	case common.BKInnerObjIDHost:
		idProperty.ID = common.BKHostIDField
		idProperty.Name = defLang.Languagef("host_property_bk_host_id")
		fields[idProperty.ID] = idProperty
	case common.BKInnerObjIDObject:
		idProperty.ID = common.BKInstIDField
		idProperty.Name = defLang.Languagef("common_property_bk_inst_id")
		fields[idProperty.ID] = idProperty
	}

}
