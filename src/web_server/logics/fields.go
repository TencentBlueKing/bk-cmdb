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
	"configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/tealeg/xlsx/v3"
)

// ExcelDataRange excel data range
type ExcelDataRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// HeaderTable excel header about table type
type HeaderTable struct {
	Start        int                 `json:"start"`
	End          int                 `json:"end"`
	NameIndexMap map[int]string      `json:"name_index_map"`
	Field        map[string]Property `json:"field"`
}

// Property object fields
type Property struct {
	ID            string
	Name          string
	PropertyType  string
	Option        interface{}
	IsPre         bool
	IsRequire     bool
	Group         string
	ExcelColIndex int
	NotObjPropery bool // Not an attribute of the object, indicating that the field to be exported is needed for export,
	AsstObjID     string
	NotExport     bool
	Length        int // 这个属性需要占用多少列excel
	NotEditable   bool
}

// ImportExcelPreData import excel pre data
type ImportExcelPreData struct {
	Fields       map[string]Property
	NameIndexMap map[int]string
	DataRange    []ExcelDataRange
	TableMap     map[string]HeaderTable
	Sheet        *xlsx.Sheet
}

// HandleFieldParam 处理Excel表格字段入参
type HandleFieldParam struct {
	Rid       string
	StyleCell *xlsx.Style
	Sheet     *xlsx.Sheet
	File      *xlsx.File
	Filter    []string
	DefLang   lang.DefaultCCLanguageIf
	CellStyle *xlsx.Style
	ColStyle  *xlsx.Style
}

// HandleHostDataParam 处理主机数据生成excel表格数据入参
type HandleHostDataParam struct {
	HostData          []mapstr.MapStr
	ExtFieldsTopoID   string
	ExtFieldsBizID    string
	ExtFieldsModuleID string
	ExtFieldsSetID    string
	CcErr             errors.DefaultCCErrorIf
	ExtFieldKey       []string
	UsernameMap       map[string]string
	PropertyList      []string
	Organization      []metadata.DepartmentItem
	OrgPropertyList   []string
	CcLang            lang.DefaultCCLanguageIf
	Sheet             *xlsx.Sheet
	File              *xlsx.File
	Rid               string
	ObjID             string
	ObjIDs            []string
	Fields            map[string]Property
	RowCount          int // 主机需要占用多少行excel
}

// HandleHostParam 处理主机数据入参
type HandleHostParam struct {
	RowIndex     int
	Data         []mapstr.MapStr
	CcErr        errors.DefaultCCErrorIf
	Fields       map[string]Property
	Rid          string
	ModelBizID   int64
	CustomLen    int
	ObjID        string
	UsernameMap  map[string]string
	PropertyList []string
	ObjName      []string
	CcLang       lang.DefaultCCLanguageIf
	Sheet        *xlsx.Sheet
}

// PropertyGroup property group
type PropertyGroup struct {
	Name  string
	Index int64
	ID    string
}

// PropertyPrimaryVal TODO
type PropertyPrimaryVal struct {
	ID     string
	Name   string
	StrVal string
}

// GetObjFieldIDs get object fields
func (lgc *Logics) GetObjFieldIDs(objID string, filterFields []string, customFields []string, header http.Header,
	modelBizID int64, index int) (map[string]Property, error) {
	fields, err := lgc.getObjFieldIDs(objID, header, modelBizID, customFields, index)
	if nil != err {
		return nil, fmt.Errorf("get object fields failed, err: %+v", err)
	}

	ret := make(map[string]Property)
	for _, field := range fields {
		if util.InStrArr(filterFields, field.ID) {
			field.NotExport = true
		}
		ret[field.ID] = field
	}

	return ret, nil
}

func (lgc *Logics) getObjectGroup(objID string, header http.Header, modelBizID int64) ([]PropertyGroup, error) {
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)
	condition := mapstr.MapStr{
		common.BKObjIDField: objID,
		"page": mapstr.MapStr{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  common.BKPropertyGroupIndexField,
		},
		common.BKAppIDField: modelBizID,
	}
	result, err := lgc.Engine.CoreAPI.ApiServer().GetObjectGroup(context.Background(), header, ownerID, objID,
		condition)
	if nil != err {
		blog.Errorf("get %s fields group failed, err:%+v, rid: %s", objID, err, rid)
		return nil, fmt.Errorf("get attribute group failed, err: %+v", err)
	}
	if !result.Result {
		blog.Errorf("get %s fields group result failed. error code:%d, error message:%s, rid:%s", objID, result.Code,
			result.ErrMsg, rid)
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

func (lgc *Logics) getObjFieldIDs(objID string, header http.Header, modelBizID int64, customFields []string,
	index int) ([]Property, error) {
	rid := util.GetHTTPCCRequestID(header)
	sort := fmt.Sprintf("%s", common.BKPropertyIndexField)

	customFieldsCond := make(map[string]interface{})
	if len(customFields) > 0 {
		fields := append(customFields, common.BKHostInnerIPField, common.BKCloudIDField)
		customFieldsCond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBIN: fields}
	}

	// sortedFields 模型字段已经根据bk_property_index排序好了
	sortedFields, err := lgc.getObjFieldIDsBySort(objID, sort, header, customFieldsCond, modelBizID)
	if err != nil {
		blog.Errorf("getObjFieldIDs, getObjFieldIDsBySort failed, sort: %s, rid: %s, err: %v", sort, rid, err)
		return nil, err
	}

	groups, err := lgc.getObjectGroup(objID, header, modelBizID)
	if nil != err {
		return nil, fmt.Errorf("getObjFieldIDs, get attribute group failed, err: %+v", err)
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("get attribute group by object not found")
	}

	fields := make([]Property, 0)
	requiredFieldMap := make(map[string][]Property)
	noRequiredFieldMap := make(map[string][]Property)

	// 构造必填字段和非必填字段所在分组的map
	for _, field := range sortedFields {
		if field.IsRequire {
			requiredFieldMap[field.Group] = append(requiredFieldMap[field.Group], field)
			continue
		}
		noRequiredFieldMap[field.Group] = append(noRequiredFieldMap[field.Group], field)
	}

	// 第二步，根据字段分组，对必填字段排序
	requiredFields, index, err := setFieldsIndex(groups, requiredFieldMap, index)
	if err != nil {
		return nil, err
	}
	fields = append(fields, requiredFields...)

	// 第三步，根据字段分组，用必填字段使用的index，继续对非必填字段进行排序
	noRequiredFields, index, err := setFieldsIndex(groups, noRequiredFieldMap, index)
	if err != nil {
		return nil, err
	}

	fields = append(fields, noRequiredFields...)
	return fields, nil
}

func setFieldsIndex(groups []PropertyGroup, fieldsGroupMap map[string][]Property, index int) ([]Property, int, error) {
	result := make([]Property, 0)
	for _, group := range groups {
		fields, ok := fieldsGroupMap[group.ID]
		if ok {
			for _, field := range fields {
				field.ExcelColIndex = index
				if field.PropertyType == common.FieldTypeInnerTable {
					option, err := metadata.ParseTableAttrOption(field.Option)
					if err != nil {
						return nil, 0, err
					}
					index += len(option.Header)
					field.Length = len(option.Header)
					result = append(result, field)
					continue
				}
				field.Length = 1
				result = append(result, field)
				index++
			}
		}
	}

	return result, index, nil
}

func (lgc *Logics) getObjFieldIDsBySort(objID, sort string, header http.Header, conds mapstr.MapStr, modelBizID int64) (
	[]Property, error) {

	rid := util.GetHTTPCCRequestID(header)

	condition := mapstr.MapStr{
		common.BKObjIDField: objID,
		metadata.PageName: mapstr.MapStr{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  sort,
		},
		common.BKAppIDField: modelBizID,
	}
	condition.Merge(conds)

	result, err := lgc.Engine.CoreAPI.ApiServer().ModelQuote().GetObjectAttrWithTable(context.Background(), header,
		condition)
	if err != nil {
		blog.Errorf("get object fields failed, objID: %s, input: %v, err: %v ,rid: %s", objID, conds, err, rid)
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).
			Error(common.CCErrCommHTTPDoRequestFailed)
	}

	ret := []Property{}
	for _, attr := range result {
		ret = append(ret, Property{
			ID:            attr.PropertyID,
			Name:          attr.PropertyName,
			PropertyType:  attr.PropertyType,
			IsRequire:     attr.IsRequired,
			IsPre:         attr.IsPre,
			Option:        attr.Option,
			Group:         attr.PropertyGroup,
			ExcelColIndex: int(attr.PropertyIndex),
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
	case common.FieldTypeEnumMulti:
	case common.FieldTypeEnumQuote:
	case common.FieldTypeDate:
	case common.FieldTypeTime:
	case common.FieldTypeUser:
	case common.FieldTypeOrganization:
	case common.FieldTypeBool:
	case common.FieldTypeTimeZone:

	}
	if "" == name {
		name = propertyType
	}
	return name, skip
}

// addSystemField add system field, get property not return property fields
func addSystemField(fields map[string]Property, objID string, defLang lang.DefaultCCLanguageIf, index int) {
	for key, field := range fields {
		field.ExcelColIndex = field.ExcelColIndex + 1
		fields[key] = field
	}

	idProperty := Property{
		ID:            "",
		Name:          "",
		PropertyType:  common.FieldTypeInt,
		Group:         "defalut",
		ExcelColIndex: index,
		Length:        1,
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
