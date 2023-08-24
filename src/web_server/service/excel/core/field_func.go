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

package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/util"
)

// TransEnumQuoteIDToName transfer enum quote field id to name
func (d *Client) TransEnumQuoteIDToName(kit *rest.Kit, infos []mapstr.MapStr, colProps []ColProp) ([]mapstr.MapStr,
	error) {

	for _, rowMap := range infos {
		if err := d.getEnumQuoteInstNames(kit, colProps, rowMap); err != nil {
			blog.Errorf("get enum quote inst name list failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
	}

	return infos, nil
}

// getEnumQuoteInstNames search inst detail and return a bk_inst_name
func (d *Client) getEnumQuoteInstNames(kit *rest.Kit, colProps []ColProp, rowMap mapstr.MapStr) error {
	for _, property := range colProps {
		if property.PropertyType != common.FieldTypeEnumQuote {
			continue
		}

		enumQuoteIDInterface, exist := rowMap[property.ID]
		if !exist || enumQuoteIDInterface == nil {
			continue
		}

		enumQuoteIDList, ok := enumQuoteIDInterface.([]interface{})
		if !ok {
			blog.Errorf("rowMap[%s] type to array failed, rowMap: %v, rowMap type: %T, rid: %s", property,
				rowMap[property.ID], rowMap[property.ID], kit.Rid)
			return fmt.Errorf("convert variable rowMap[%v] type to int array failed", property)
		}

		enumQuoteIDMap := make(map[int64]interface{}, 0)
		for _, enumQuoteID := range enumQuoteIDList {
			id, err := util.GetInt64ByInterface(enumQuoteID)
			if err != nil {
				blog.Errorf("convert enumQuoteID[%d] to int64 failed, type: %T, err: %v, rid: %s", enumQuoteID,
					enumQuoteID, kit.Rid)
				return err
			}

			if id == 0 {
				return fmt.Errorf("enum quote instID is %d, it is illegal", id)
			}

			enumQuoteIDMap[id] = struct{}{}
		}

		if len(enumQuoteIDMap) == 0 {
			continue
		}

		quoteObjID, err := getEnumQuoteObjID(kit, property.Option)
		if err != nil {
			blog.Errorf("get enum quote option obj id failed, err: %s, rid: %s", err, kit.Rid)
			return err
		}

		enumQuoteIDs := make([]int64, 0)
		for enumQuoteID := range enumQuoteIDMap {
			enumQuoteIDs = append(enumQuoteIDs, enumQuoteID)
		}
		input := &metadata.QueryCondition{
			Fields:         []string{common.GetInstNameField(quoteObjID)},
			Condition:      mapstr.MapStr{common.GetInstIDField(quoteObjID): mapstr.MapStr{common.BKDBIN: enumQuoteIDs}},
			DisableCounter: true,
		}
		resp, err := d.Engine.CoreAPI.ApiServer().ReadInstance(kit.Ctx, kit.Header, quoteObjID, input)
		if err != nil {
			blog.Errorf("get quote inst name list failed, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
			return err
		}

		enumQuoteNames := make([]string, 0)
		for _, info := range resp.Data.Info {
			var ok bool
			var enumQuoteName string
			if name, exist := info.Get(common.GetInstNameField(quoteObjID)); exist {
				enumQuoteName, ok = name.(string)
				if !ok {
					enumQuoteName = ""
				}
			}
			enumQuoteNames = append(enumQuoteNames, enumQuoteName)
		}

		rowMap[property.ID] = strings.Join(enumQuoteNames, "\n")
	}

	return nil
}

// TransEnumQuoteNameToID transfer enum quote field name to id
func (d *Client) TransEnumQuoteNameToID(kit *rest.Kit, names []string, prop *ColProp) ([]int64, error) {
	if prop == nil {
		blog.Errorf("property is nil, rid: %s", kit.Rid)
		return nil, fmt.Errorf("property is nil")
	}

	if names == nil || len(names) == 0 {
		return nil, nil
	}

	objID, err := getEnumQuoteObjID(kit, prop.Option)
	if err != nil {
		blog.Errorf("get enum quote option objID failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	input := &metadata.QueryCondition{
		Fields:         []string{common.GetInstIDField(objID)},
		Condition:      mapstr.MapStr{common.GetInstNameField(objID): mapstr.MapStr{common.BKDBIN: names}},
		DisableCounter: true,
	}
	resp, err := d.Engine.CoreAPI.ApiServer().ReadInstance(kit.Ctx, kit.Header, objID, input)
	if err != nil {
		blog.Errorf("get instance id list failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	ids := make([]int64, 0)
	for _, info := range resp.Data.Info {
		id, err := info.Int64(common.GetInstIDField(objID))
		if err != nil {
			blog.Errorf("get enum quote id failed, err: %v, rid: %s", err, kit.Rid)
			continue
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// getEnumQuoteObjID get enum quote field option bk_obj_id and bk_inst_id value
func getEnumQuoteObjID(kit *rest.Kit, option interface{}) (string, error) {
	if option == nil {
		return "", fmt.Errorf("enum quote option is nil")
	}

	arrOption, ok := option.([]interface{})
	if !ok {
		blog.Errorf("option %v not enum quote option, rid: %s", option, kit.Rid)
		return "", fmt.Errorf("enum quote option is unvalid")
	}

	for _, o := range arrOption {
		mapOption, ok := o.(map[string]interface{})
		if !ok || mapOption == nil {
			blog.Errorf("option %v not enum quote option, enum quote option item must bk_obj_id, rid: %s", option,
				kit.Rid)
			return "", fmt.Errorf("convert option map[string]interface{} failed")
		}

		objIDVal, objIDOk := mapOption[common.BKObjIDField]
		if !objIDOk || objIDVal == "" {
			blog.Errorf("enum quote option bk_obj_id can't be empty, rid: %s", option, kit.Rid)
			return "", fmt.Errorf("enum quote option bk_obj_id can't be empty")
		}

		objID, ok := objIDVal.(string)
		if !ok {
			blog.Errorf("objIDVal %v not string, rid: %s", objIDVal, kit.Rid)
			return "", fmt.Errorf("enum quote option bk_obj_id is not string")
		}

		return objID, nil
	}

	return "", nil
}

// GetInstWithOrgName get instance with organization name
func (d *Client) GetInstWithOrgName(kit *rest.Kit, ccLang language.DefaultCCLanguageIf, insts []mapstr.MapStr,
	colProps []ColProp) ([]mapstr.MapStr, error) {

	orgIDList := make([]int64, 0)
	for _, inst := range insts {
		for _, property := range colProps {
			if property.PropertyType != common.FieldTypeOrganization || inst[property.ID] == nil {
				continue
			}

			orgIDs, ok := inst[property.ID].([]interface{})
			if !ok {
				return nil, fmt.Errorf("org id list type not []interface{}, real type is %T", inst[property.ID])
			}

			if len(orgIDs) == 0 {
				continue
			}

			for _, orgID := range orgIDs {
				id, err := util.GetInt64ByInterface(orgID)
				if err != nil {
					return nil, err
				}
				orgIDList = append(orgIDList, id)
			}
		}
	}

	orgIDs := util.IntArrayUnique(orgIDList)
	organizations, err := getAllOrganization(kit, orgIDs)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	orgMap := make(map[int64]string)
	for _, item := range organizations.Results {
		orgMap[item.ID] = item.FullName
	}

	for idx, inst := range insts {
		for _, property := range colProps {
			if property.PropertyType != common.FieldTypeOrganization || inst[property.ID] == nil {
				continue
			}

			orgIDs, ok := inst[property.ID].([]interface{})
			if !ok {
				return nil, fmt.Errorf("org id list type not []interface{}, real type is %T", inst[property.ID])
			}

			if len(orgIDs) == 0 {
				continue
			}

			orgName := make([]string, 0)
			for _, orgID := range orgIDs {
				id, err := util.GetInt64ByInterface(orgID)
				if err != nil {
					return nil, err
				}

				name, exist := orgMap[id]
				if !exist {
					orgName = append(orgName, fmt.Sprintf("[%d]%s", id, ccLang.Language("nonexistent_org")))
					continue
				}

				orgName = append(orgName, fmt.Sprintf("[%d]%s", id, name))
			}

			insts[idx][property.ID] = strings.Join(orgName, ",")
		}
	}

	return insts, nil
}

// getAllOrganization get organization info from paas
func getAllOrganization(kit *rest.Kit, orgIDs []int64) (*metadata.DepartmentData, errors.CCErrorCoder) {

	loginVersion, _ := cc.String("webServer.login.version")
	if loginVersion == common.BKOpenSourceLoginPluginVersion || loginVersion == common.BKSkipLoginPluginVersion {
		return &metadata.DepartmentData{}, nil
	}

	orgIDList := getOrgListStr(orgIDs)
	departments := &metadata.DepartmentData{}
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr error
	pipeline := make(chan bool, 10)

	for _, subStr := range orgIDList {
		pipeline <- true
		wg.Add(1)
		go func(subStr string) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			params := make(map[string]string)
			params["exact_lookups"] = subStr
			result, esbErr := esb.EsbClient().User().GetAllDepartment(kit.Ctx, kit.Header, params)
			if esbErr != nil {
				firstErr = esbErr
				blog.Errorf("get department by esb client failed, params: %+v, err: %v, rid: %s", params, esbErr,
					kit.Rid)
				return
			}
			if !result.Result {
				blog.Errorf("get department by esb client failed, params: %+v, rid: %s", params, kit.Rid)
				firstErr = fmt.Errorf("get department by esb failed, params: %v", params)
				return
			}

			lock.Lock()
			departments.Count += result.Data.Count
			departments.Results = append(departments.Results, result.Data.Results...)
			lock.Unlock()
		}(subStr)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return departments, nil
}

const organizationMaxLength = 500

// getOrgListStr get org list str
func getOrgListStr(orgIDList []int64) []string {
	orgListStr := make([]string, 0)

	orgBuffer := bytes.Buffer{}
	for _, orgID := range orgIDList {
		if orgBuffer.Len()+len(strconv.FormatInt(orgID, 10)) > organizationMaxLength {
			orgBuffer.WriteString(strconv.FormatInt(orgID, 10))
			orgStr := orgBuffer.String()
			orgListStr = append(orgListStr, orgStr)
			orgBuffer.Reset()
			continue
		}

		orgBuffer.WriteString(strconv.FormatInt(orgID, 10))
		orgBuffer.WriteByte(',')
	}

	if orgBuffer.Len() == 0 {
		return []string{}
	}

	orgStr := orgBuffer.String()
	orgListStr = append(orgListStr, orgStr[:len(orgStr)-1])

	return orgListStr
}

// GetInstWithTable 第一个返回值是返回带有表格数据的实例信息，第二个返回值返回的每个实例数据所占用的excel行数
func (d *Client) GetInstWithTable(kit *rest.Kit, objID string, insts []mapstr.MapStr, colProps []ColProp) (
	[]mapstr.MapStr, []int, error) {

	instHeights := make([]int, len(insts))
	for i := range instHeights {
		instHeights[i] = 1
	}

	// 1. 找出表格字段
	tableProperty := make([]ColProp, 0)
	for _, property := range colProps {
		if property.PropertyType == common.FieldTypeInnerTable {
			tableProperty = append(tableProperty, property)
		}
	}
	if len(tableProperty) == 0 {
		return insts, instHeights, nil
	}

	ids := make([]int64, 0)
	dataMap := make(map[int64]map[string][]mapstr.MapStr)
	for _, info := range insts {
		id, err := info.Int64(common.GetInstIDField(objID))
		if err != nil {
			return nil, nil, fmt.Errorf("data is invalid, err: %v", err)
		}
		ids = append(ids, id)
		dataMap[id] = make(map[string][]mapstr.MapStr)
	}

	// 2. 查询数据对应的表格字段的值
	queryOpt := metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(
			common.BKInstIDField, filter.In, ids)},
		Page: metadata.BasePage{Limit: common.BKMaxPageSize},
	}
	for _, property := range tableProperty {
		opt := &metadata.ListQuotedInstOption{ObjID: objID, PropertyID: property.ID, CommonQueryOption: queryOpt}
		instances, err := d.Engine.CoreAPI.ApiServer().ModelQuote().ListQuotedInstance(kit.Ctx, kit.Header, opt)
		if err != nil {
			return nil, nil, err
		}
		for _, inst := range instances.Info {
			instID, err := inst.Int64(common.BKInstIDField)
			if err != nil {
				return nil, nil, err
			}
			dataMap[instID][property.ID] = append(dataMap[instID][property.ID], inst)
		}
	}

	// 3. 整理返回带表格数据的结果, 以及每条数据需要占用excel多少行
	for idx := range insts {
		id, err := insts[idx].Int64(common.GetInstIDField(objID))
		if err != nil {
			return nil, nil, fmt.Errorf("data is invalid, err: %v", err)
		}

		for propertyID, data := range dataMap[id] {
			insts[idx][propertyID] = data

			if len(data) > instHeights[idx] {
				instHeights[idx] = len(data)
			}
		}
	}

	return insts, instHeights, nil
}
