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
	ID            int64
	PropertyID    string
	PropertyName  string
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

type AttrUnique struct {
	ID        int64
	IsPre     bool
	MustCheck bool
	KeyID     int64
	KeyKind   string
}

// GetObjFieldIDs get object fields
func (lgc *Logics) GetObjFieldIDs(objID string, filterFields []string, customFields []string, header http.Header, meta *metadata.Metadata) (map[string]Property, error) {
	fields, err := lgc.getObjFieldIDs(objID, header, meta)
	if nil != err {
		return nil, fmt.Errorf("get object fields failed, err: %+v", err)
	}

	ret := make(map[string]Property)
	for _, field := range fields {
		if util.InStrArr(filterFields, field.PropertyID) {
			field.NotExport = true
		}
		ret[field.PropertyID] = field
	}

	return ret, nil
}

func (lgc *Logics) getObjectGroup(objID string, header http.Header, meta *metadata.Metadata) ([]PropertyGroup, error) {
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)
	condition := mapstr.MapStr{
		common.BKObjIDField: objID,
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

func (lgc *Logics) getObjectUnique(objID string, header http.Header, meta *metadata.Metadata) ([]AttrUnique, error) {
	rid := util.GetHTTPCCRequestID(header)
	condition := mapstr.MapStr{
		"condition": []string{
			objID,
		},
		"metadata": meta,
	}

	result, err := lgc.Engine.CoreAPI.ApiServer().GetObjectUnique(context.Background(), header, objID, condition)
	if nil != err {
		return nil, fmt.Errorf("get object unique data failed, err: %v, rid: %s", err, rid)
	}

	if !result.Result {
		return nil, fmt.Errorf("get object unique data failed, but got err: %s, rid: %s", result.ErrMsg, rid)
	}

	attrUniqueList := make([]AttrUnique, 0)
	for _, objUnique := range result.Data {
		for _, key := range objUnique.Keys {
			attrUnique := AttrUnique{}
			attrUnique.ID = int64(objUnique.ID)
			attrUnique.IsPre = objUnique.Ispre
			attrUnique.MustCheck = objUnique.MustCheck
			attrUnique.KeyID = int64(key.ID)
			attrUnique.KeyKind = key.Kind
			attrUniqueList = append(attrUniqueList, attrUnique)
		}
	}

	return attrUniqueList, nil

}

func (lgc *Logics) getObjFieldIDs(objID string, header http.Header, meta *metadata.Metadata) ([]Property, error) {
	rid := util.GetHTTPCCRequestID(header)
	sort := fmt.Sprintf("%s", common.BKPropertyIndexField)

	// sortedFields 模型字段已经根据bk_property_index排序好了
	sortedFields, err := lgc.getObjFieldIDsBySort(objID, sort, header, nil, meta)
	if err != nil {
		blog.Errorf("getObjFieldIDs, getObjFieldIDsBySort failed, sort: %s, rid: %s, err: %v", sort, rid, err)
		return nil, err
	}

	groups, err := lgc.getObjectGroup(objID, header, meta)
	if nil != err {
		return nil, fmt.Errorf("getObjFieldIDs, get attribute group failed, err: %+v", err)
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("get attribute group by object not found")
	}

	uniques, err := lgc.getObjectUnique(objID, header, meta)
	if nil != err {
		return nil, fmt.Errorf("getObjectUnique, get object unique data failed, err: %v, rid: %s", err)
	}
	if len(uniques) == 0 {
		return nil, fmt.Errorf("get object unique data,data not found")
	}

	fields := make([]Property, 0)
	noRequiredFields := make([]Property, 0)
	noUniqueFields := sortedFields
	index := 1
	// 第一步，根据ObjectUnique表，从sortedFields拉取与之相同id的字段；
	for _, unique := range uniques {
		for i, field := range sortedFields {
			//不考虑“属性空值时不校验”的唯一性校验
			if unique.KeyID != field.ID {
				continue
			}
			if unique.MustCheck != true {
				continue
			}
			noUniqueFields[i] = Property{}
			field.ExcelColIndex = index
			index++
			fields = append(fields, field)
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

func (lgc *Logics) getObjFieldIDsBySort(objID, sort string, header http.Header, conds mapstr.MapStr, meta *metadata.Metadata) ([]Property, error) {
	rid := util.GetHTTPCCRequestID(header)

	condition := mapstr.MapStr{
		common.BKObjIDField: objID,
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

	//inputParam := metadata.QueryCondition{
	//	Page: metadata.BasePage{
	//		Start: 0,
	//		Limit: common.BKNoLimit,
	//	},
	//	Condition: mapstr.MapStr(map[string]interface{}{
	//		common.BKObjIDField: objID,
	//	}),
	//}
	//uniques, err := lgc.CoreAPI.CoreService().Model().ReadModelAttrUnique(context.Background(), header, inputParam)
	//if nil != err {
	//	blog.Errorf("getObjectPrimaryFieldByObjID get unique for %s error: %v ,rid:%s", objID, err, rid)
	//	return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Error(common.CCErrCommHTTPDoRequestFailed)
	//}
	//if !uniques.Result {
	//	blog.Errorf("getObjectPrimaryFieldByObjID get unique for %s error: %v ,rid:%s", objID, uniques, rid)
	//	return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(uniques.Code, uniques.ErrMsg)
	//}

	//keyIDs := map[uint64]bool{}
	//for _, unique := range uniques.Data.Info {
	//	if unique.MustCheck {
	//		for _, key := range unique.Keys {
	//			keyIDs[key.ID] = true
	//		}
	//		break
	//	}
	//}
	//if len(keyIDs) <= 0 {
	//	return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Error(common.CCErrTopoObjectUniqueSearchFailed)
	//}

	ret := []Property{}
	for _, attr := range result.Data {
		ret = append(ret, Property{
			ID:            attr.ID,
			PropertyID:    attr.PropertyID,
			PropertyName:  attr.PropertyName,
			PropertyType:  attr.PropertyType,
			IsRequire:     attr.IsRequired,
			IsPre:         attr.IsPre,
			Option:        attr.Option,
			Group:         attr.PropertyGroup,
			ExcelColIndex: int(attr.PropertyIndex),
			IsOnly:        attr.IsOnly,
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
		PropertyID:    "",
		PropertyName:  "",
		PropertyType:  common.FieldTypeInt,
		Group:         "defalut",
		ExcelColIndex: 1, // why set ExcelColIndex=1? because ExcelColIndex=0 used by tip column
	}

	switch objID {
	case common.BKInnerObjIDHost:
		idProperty.PropertyID = common.BKHostIDField
		idProperty.PropertyName = defLang.Languagef("host_property_bk_host_id")
		fields[idProperty.PropertyID] = idProperty
	case common.BKInnerObjIDObject:
		idProperty.PropertyID = common.BKInstIDField
		idProperty.PropertyName = defLang.Languagef("common_property_bk_inst_id")
		fields[idProperty.PropertyID] = idProperty
	}

}
