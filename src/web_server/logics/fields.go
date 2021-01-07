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
func (lgc *Logics) GetObjFieldIDs(objID string, filterFields []string, customFields []string, header http.Header, modelBizID int64) (map[string]Property, error) {
	fields, err := lgc.getObjFieldIDs(objID, header, modelBizID)
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

func (lgc *Logics) getObjectPrimaryFieldByObjID(objID string, header http.Header, modelBizID int64) ([]Property, error) {
	fields, err := lgc.getObjFieldIDsBySort(objID, common.BKPropertyIDField, header, nil, modelBizID)
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

func (lgc *Logics) getObjFieldIDs(objID string, header http.Header, modelBizID int64) ([]Property, error) {
	rid := util.GetHTTPCCRequestID(header)
	sort := fmt.Sprintf("%s", common.BKPropertyIndexField)

	// sortedFields 模型字段已经根据bk_property_index排序好了
	sortedFields, err := lgc.getObjFieldIDsBySort(objID, sort, header, nil, modelBizID)
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
	noUniqueFields := make([]Property, 0)
	noRequiredFields := make([]Property, 0)
	index := 1
	// 第一步，选出唯一校验字段；
	for _, field := range sortedFields {
		if field.IsOnly == true {
			field.ExcelColIndex = index
			index++
			fields = append(fields, field)
		} else {
			noUniqueFields = append(noUniqueFields, field)
		}
	}

	// 第二步，根据字段分组，对必填字段排序；并选出非必填字段
	for _, group := range groups {
		for _, field := range noUniqueFields {
			if field.Group != group.ID {
				continue
			}
			if field.IsRequire != true {
				noRequiredFields = append(noRequiredFields, field)
				continue
			}
			field.ExcelColIndex = index
			index++
			fields = append(fields, field)
		}
	}

	// 第三步，根据字段分组，用必填字段使用的index，继续对非必填字段进行排序
	for _, group := range groups {
		for _, field := range noRequiredFields {
			if field.Group != group.ID {
				continue
			}
			field.ExcelColIndex = index
			index++
			fields = append(fields, field)
		}
	}
	return fields, nil
}

func (lgc *Logics) getObjFieldIDsBySort(objID, sort string, header http.Header, conds mapstr.MapStr, modelBizID int64) ([]Property, error) {
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
		Page: metadata.BasePage{
			Start: 0,
			Limit: common.BKNoLimit,
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
	for _, attr := range result.Data {
		ret = append(ret, Property{
			ID:            attr.PropertyID,
			Name:          attr.PropertyName,
			PropertyType:  attr.PropertyType,
			IsRequire:     attr.IsRequired,
			IsPre:         attr.IsPre,
			Option:        attr.Option,
			Group:         attr.PropertyGroup,
			ExcelColIndex: int(attr.PropertyIndex),
			IsOnly:        keyIDs[uint64(attr.ID)],
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
		ExcelColIndex: 1, // why set ExcelColIndex=1? because ExcelColIndex=0 used by tip column
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
