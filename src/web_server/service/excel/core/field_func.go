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
	"fmt"
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
	"configcenter/src/common/resource/apigw"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/apigw/user"
	"configcenter/src/web_server/middleware/user/plugins"
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

		if len(enumQuoteIDList) == 0 {
			continue
		}

		enumQuoteIDs, err := util.SliceInterfaceToInt64(enumQuoteIDList)
		if err != nil {
			blog.Errorf("slice interface to int64 failed, data: %v, err: %v, rid: %s", enumQuoteIDList, err, kit.Rid)
			return err
		}

		quoteObjID, err := getEnumQuoteObjID(kit, property.Option)
		if err != nil {
			blog.Errorf("get enum quote option obj id failed, err: %s, rid: %s", err, kit.Rid)
			return err
		}

		input := &metadata.QueryCondition{
			Fields:         []string{common.GetInstNameField(quoteObjID)},
			Condition:      mapstr.MapStr{common.GetInstIDField(quoteObjID): mapstr.MapStr{common.BKDBIN: enumQuoteIDs}},
			DisableCounter: true,
		}
		resp, err := d.ApiClient.ReadInstance(kit.Ctx, kit.Header, quoteObjID, input)
		if err != nil {
			blog.Errorf("get quote inst name list failed, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
			return err
		}

		enumQuoteNames := make([]string, 0)
		for _, info := range resp.Data.Info {
			var ok bool
			var enumQuoteName string
			name, exist := info.Get(common.GetInstNameField(quoteObjID))
			if !exist {
				enumQuoteNames = append(enumQuoteNames, enumQuoteName)
				continue
			}

			enumQuoteName, ok = name.(string)
			if !ok {
				enumQuoteName = ""
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

	if len(names) == 0 {
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
	resp, err := d.ApiClient.ReadInstance(kit.Ctx, kit.Header, objID, input)
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

	quoteOption, err := metadata.ParseEnumQuoteOption(kit.Ctx, option)
	if err != nil {
		blog.Errorf("parse enum quote failed, data: %v, err: %v, rid: %s", option, err, kit.Rid)
		return "", err
	}

	if len(quoteOption) == 0 {
		return "", nil
	}

	return quoteOption[0].ObjID, nil
}

// GetInstWithOrgName get instance with organization name
func (d *Client) GetInstWithOrgName(kit *rest.Kit, ccLang language.DefaultCCLanguageIf, insts []mapstr.MapStr,
	colProps []ColProp) ([]mapstr.MapStr, error) {

	orgPropIDs := make([]string, 0)
	for _, property := range colProps {
		if property.PropertyType == common.FieldTypeOrganization {
			orgPropIDs = append(orgPropIDs, property.ID)
		}
	}

	orgIDList := make([]int64, 0)
	for _, inst := range insts {
		for _, propertyID := range orgPropIDs {
			if inst[propertyID] == nil {
				continue
			}

			orgIDs, ok := inst[propertyID].([]interface{})
			if !ok {
				return nil, fmt.Errorf("org id list type not []interface{}, real type is %T", inst[propertyID])
			}

			if len(orgIDs) == 0 {
				continue
			}

			ids, err := util.SliceInterfaceToInt64(orgIDs)
			if err != nil {
				blog.Errorf("slice interface to int64 failed, val: %v, err: %v, rid: %s", orgIDs, err, kit.Rid)
				return nil, err
			}
			orgIDList = append(orgIDList, ids...)
		}
	}

	orgIDs := util.IntArrayUnique(orgIDList)
	if len(orgIDs) == 0 {
		return insts, nil
	}

	organizations, err := getAllOrganization(kit, orgIDs)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	orgMap := make(map[int64]string)
	for _, item := range organizations {
		orgMap[item.ID] = item.Name
	}

	for idx, inst := range insts {
		for _, propertyID := range orgPropIDs {
			if inst[propertyID] == nil {
				continue
			}

			orgIDs, ok := inst[propertyID].([]interface{})
			if !ok {
				return nil, fmt.Errorf("org id list type not []interface{}, real type is %T", inst[propertyID])
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

			insts[idx][propertyID] = strings.Join(orgName, ",")
		}
	}

	return insts, nil
}

// getAllOrganization get organization info from paas
func getAllOrganization(kit *rest.Kit, orgIDs []int64) ([]user.DepartmentItem, errors.CCErrorCoder) {
	loginVersion, _ := cc.String("webServer.login.version")
	if loginVersion == common.BKOpenSourceLoginPluginVersion || loginVersion == common.BKSkipLoginPluginVersion {
		return make([]user.DepartmentItem, 0), nil
	}

	orgIDList := util.SplitArr(orgIDs, 100)
	departments := make([]user.DepartmentItem, 0)
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr error
	pipeline := make(chan bool, 10)

	for _, ids := range orgIDList {
		pipeline <- true
		wg.Add(1)
		go func(ids []int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			deptOpts := &user.BatchLookupDeptOpts{
				DeptIDs: ids,
			}
			result, err := apigw.Client().User().BatchLookupDept(kit.Ctx, kit.Header, deptOpts)
			if err != nil {
				firstErr = err
				blog.Errorf("get department by apigw failed, params: %+v, err: %v, rid: %s", deptOpts, err, kit.Rid)
				return
			}

			lock.Lock()
			departments = append(departments, result...)
			lock.Unlock()
		}(ids)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return departments, nil
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
		instances, err := d.ApiClient.ModelQuote().ListQuotedInstance(kit.Ctx, kit.Header, opt)
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

// GetInstWithUserFullName 将导出的实例的用户名转化为完整的用户名
func (d *Client) GetInstWithUserFullName(kit *rest.Kit, lang language.DefaultCCLanguageIf, objID string,
	insts []mapstr.MapStr, colProps []ColProp) ([]mapstr.MapStr, error) {

	propertyIDs := make([]string, 0)
	for _, property := range colProps {
		if property.PropertyType == common.FieldTypeUser {
			propertyIDs = append(propertyIDs, property.ID)
		}
	}

	names := make([]string, 0)
	for _, inst := range insts {
		for _, propertyID := range propertyIDs {
			if inst[propertyID] == nil {
				continue
			}

			username, ok := inst[propertyID].(string)
			if !ok {
				blog.Errorf("failed to cast %s instance from interface{} to string", objID, kit.Rid)
				return nil, fmt.Errorf("failed to cast %s instance from interface{} to string", objID)
			}

			names = append(names, strings.Split(username, ",")...)
		}
	}

	userList := util.RemoveDuplicatesAndEmpty(names)
	// get username from apigw
	fullNameMap, err := d.getUsernameFromApigw(kit, userList)
	if err != nil {
		blog.Errorf("get username map from apigw failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	for idx, inst := range insts {
		for _, propertyID := range propertyIDs {
			if inst[propertyID] == nil {
				continue
			}

			nameStr, ok := inst[propertyID].(string)
			if !ok {
				blog.Errorf("failed to cast %s instance from interface{} to string", objID, kit.Rid)
				return nil, fmt.Errorf("failed to cast %s instance from interface{} to string", objID)
			}

			if nameStr == "" {
				continue
			}

			oldNames := strings.Split(nameStr, ",")
			newNames := make([]string, 0)

			for _, name := range oldNames {
				fullName := fullNameMap[name]
				if fullName == "" {
					// return the original name and remind that the user is nonexistent in '()'
					fullName = fmt.Sprintf("%s(%s)", name, lang.Language("nonexistent_user"))
				}
				newNames = append(newNames, fullName)
			}
			insts[idx][propertyID] = strings.Join(newNames, ",")
		}
	}

	return insts, nil
}

func (d *Client) getUsernameFromApigw(kit *rest.Kit, userList []string) (map[string]string, error) {
	usernameMap := map[string]string{}

	if len(userList) == 0 {
		return usernameMap, nil
	}

	loginVersion, _ := cc.String("webServer.login.version")
	user := plugins.CurrentPlugin(loginVersion)

	userListArr := util.SplitArr(userList, 100)
	for _, users := range userListArr {
		userInfoList, errNew := user.GetUserList(d.GinCtx, &metadata.GetUserListOptions{Usernames: users})
		if errNew != nil {
			blog.Errorf("get users(%+v) list from apigw failed, err: %v, rid: %s", users, errNew, kit.Rid)
			return nil, errNew.ToCCError(kit.CCError)
		}

		for _, userInfo := range userInfoList {
			username := fmt.Sprintf("%s(%s)", userInfo.EnName, userInfo.CnName)
			usernameMap[userInfo.EnName] = username
		}
	}

	return usernameMap, nil
}
