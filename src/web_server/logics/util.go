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
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
)

const (
	fieldTypeBoolTrue  = "true"
	fieldTypeBoolFalse = "false"
)

// getFieldsIDIndexMap get field property index
func getFieldsIDIndexMap(fields map[string]Property) map[string]int {
	index := 0
	IDNameMap := make(map[string]int)
	for id := range fields {
		IDNameMap[id] = index
		index++
	}
	return IDNameMap
}

// getAssociatePrimaryKey  get getAssociate object name
func getAssociatePrimaryKey(a []interface{}, primaryField []Property) []string {
	vals := []string{}
	for _, valRow := range a {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			instMap, ok := mapVal["inst_info"].(map[string]interface{})
			if true == ok {
				var itemVals []string
				for _, field := range primaryField {
					val, _ := instMap[field.ID]
					if nil == val {
						val = ""
					}
					itemVals = append(itemVals, fmt.Sprintf("%v", val))
				}
				vals = append(vals, strings.Join(itemVals, common.ExcelAsstPrimaryKeySplitChar))
			}
		}
	}

	return vals
}

// getEnumNameByID get enum name from option
func getEnumNameByID(id string, items []interface{}) string {
	var name string
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			enumID, ok := mapVal["id"].(string)
			if true == ok {
				if enumID == id {
					name = mapVal["name"].(string)
				}
			}
		}
	}

	return name
}

// getEnumIDByName get enum name from option
func getEnumIDByName(name string, items []interface{}) string {
	id := name
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			enumName, ok := mapVal["name"].(string)
			if true == ok {
				if enumName == name {
					id = mapVal["id"].(string)
				}
			}
		}
	}

	return id
}

// getEnumNames get enum name from option
func getEnumNames(items []interface{}) []string {
	var names []string
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {

			name, ok := mapVal["name"].(string)
			if ok {
				names = append(names, name)
			}

		}
	}

	return names
}

// getHeaderCellGeneralStyle get excel header general style by C6EFCE,000000
func getHeaderCellGeneralStyle() *xlsx.Style {
	return getCellStyle(common.ExcelHeaderOtherRowColor, common.ExcelHeaderOtherRowFontColor)
}

// getHeaderFirstRowCellStyle TODO
func getHeaderFirstRowCellStyle(isRequire bool) *xlsx.Style {
	if isRequire {
		return getCellStyle(common.ExcelHeaderFirstRowColor, common.ExcelHeaderFirstRowRequireFontColor)
	}

	return getCellStyle(common.ExcelHeaderFirstRowColor, common.ExcelHeaderFirstRowFontColor)
}

// getCellStyle get cell style from fgColor and fontcolor
func getCellStyle(fgColor, fontColor string) *xlsx.Style {
	style := xlsx.NewStyle()
	style.Fill = *xlsx.DefaultFill()
	style.Font = *xlsx.DefaultFont()
	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = true

	style.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	style.Border.BottomColor = common.ExcelCellDefaultBorderColor
	style.Border.TopColor = common.ExcelCellDefaultBorderColor
	style.Border.LeftColor = common.ExcelCellDefaultBorderColor
	style.Border.RightColor = common.ExcelCellDefaultBorderColor

	style.Fill.FgColor = fgColor
	style.Fill.PatternType = "solid"

	style.Font.Color = fontColor

	return style
}

// addExtFields  add extra fields,
func addExtFields(fields map[string]Property, extFields map[string]string, extFieldKey []string) map[string]Property {
	excelColIndex := 0
	for _, extFieldID := range extFieldKey {
		fields[extFieldID] = Property{
			ID:            "",
			Name:          extFields[extFieldID],
			NotObjPropery: true,
			ExcelColIndex: excelColIndex,
			NotEditable:   true,
		}
		excelColIndex++
	}

	for _, field := range fields {
		if excelColIndex < field.ExcelColIndex {
			excelColIndex = field.ExcelColIndex
		}
	}

	return fields
}

func replaceEnName(rid string, rowMap mapstr.MapStr, usernameMap map[string]string, propertyList []string,
	defLang language.DefaultCCLanguageIf) (mapstr.MapStr, error) {
	// propertyList是用户自定义的objuser型的attr名列表
	for _, property := range propertyList {
		if rowMap[property] == nil {
			continue
		}

		userListString, ok := rowMap[property].(string)
		if !ok {
			blog.Errorf("convert variable rowMap[%s] type to string field , rowMap: %v, rowMap type: %T, rid: %s", property, rowMap[property], rowMap[property], rid)
			return nil, fmt.Errorf("convert variable rowMap[%s] type to string field", property)
		}
		userListString = strings.TrimSpace(userListString)
		if userListString == "" {
			continue
		}

		newUserList := []string{}
		enNameList := strings.Split(userListString, ",")
		for _, enName := range enNameList {
			username := usernameMap[enName]
			if username == "" {
				// return the original user name and remind that the user is nonexistent in '()'
				username = fmt.Sprintf("%s(%s)", enName, defLang.Language("nonexistent_user"))
			}
			newUserList = append(newUserList, username)
		}
		rowMap[property] = strings.Join(newUserList, ",")
	}

	return rowMap, nil
}

// setExcelCellIgnored set the excel cell to be ignored
func setExcelCellIgnored(sheet *xlsx.Sheet, style *xlsx.Style, row int, col int) error {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return err
	}
	cell.Value = common.ExcelCellIgnoreValue
	cell.SetStyle(style)
	return nil
}

// replaceDepartmentFullName replace attribute organization's id by fullname in export excel
func replaceDepartmentFullName(rid string, rowMap mapstr.MapStr, org []metadata.DepartmentItem, propertyList []string,
	defLang language.DefaultCCLanguageIf) (mapstr.MapStr, error) {
	orgMap := make(map[int64]string)
	for _, item := range org {
		orgMap[item.ID] = item.FullName
	}

	for _, property := range propertyList {
		orgIDInterface, exist := rowMap[property]
		if !exist || orgIDInterface == nil {
			continue
		}

		orgIDList, ok := orgIDInterface.([]interface{})
		if !ok {
			blog.Errorf("rowMap[%s] type to array failed, rowMap: %v, rowMap type: %T, rid: %s", property,
				rowMap[property], rowMap[property], rid)
			return nil, fmt.Errorf("convert variable rowMap[%s] type to int array failed", property)
		}

		orgName := make([]string, 0)
		for _, orgID := range orgIDList {
			id, err := util.GetInt64ByInterface(orgID)
			if err != nil {
				blog.Errorf("convert orgID to int64 failed, type: %T, err: %v, rid: %s", orgID, err, rid)
				return nil, fmt.Errorf("convert variable orgID[%v] type to int64 failed", orgID)
			}

			name, exist := orgMap[id]
			if !exist {
				blog.Errorf("organization[%d] does no exist, rid: %s", id, rid)
				orgName = append(orgName, fmt.Sprintf("[%d]%s", id, defLang.Language("nonexistent_org")))
				continue
			}

			orgName = append(orgName, fmt.Sprintf("[%d]%s", id, name))
		}
		rowMap[property] = strings.Join(orgName, ",")
	}

	return rowMap, nil
}

// replaceEnumMultiName replace attribute enummulti's id by name in export excel
func replaceEnumMultiName(rid string, rowMap mapstr.MapStr, fields map[string]Property) (mapstr.MapStr, error) {

	for id, property := range fields {
		enumMultiIDInterface, exist := rowMap[id]
		if !exist || enumMultiIDInterface == nil {
			continue
		}

		switch property.PropertyType {
		case common.FieldTypeEnumMulti:
			enumMultiIDList, ok := enumMultiIDInterface.([]interface{})
			if !ok {
				blog.Errorf("rowMap[%s] type to array failed, rowMap: %v, rowMap type: %T, rid: %s", property,
					rowMap[id], rowMap[id], rid)
				return nil, fmt.Errorf("convert variable rowMap[%s] type to int array failed", property)
			}
			enumMultiName := make([]string, 0)
			for _, enumMultiID := range enumMultiIDList {
				id, ok := enumMultiID.(string)
				if !ok {
					blog.Errorf("convert enumMultiID[%s] to string failed, type: %T, rid: %s", enumMultiID,
						enumMultiID, rid)
					return nil, fmt.Errorf("convert variable enumMultiID[%s] type to string failed", enumMultiID)
				}
				items, ok := property.Option.([]interface{})
				if !ok {
					blog.Errorf("convert option to []interface{} failed, type: %T, rid: %s", property.Option, rid)
					return nil, fmt.Errorf("enum multi option param is invalid, option: %v", property.Option)
				}
				name := getEnumNameByID(id, items)
				enumMultiName = append(enumMultiName, name)
			}
			rowMap[id] = strings.Join(enumMultiName, "\n")
		}
	}

	return rowMap, nil
}

// GetUsernameMapWithPropertyList 依照"bk_obj_id"和"bk_property_type":"objuser"查询"cc_ObjAttDes"集合,得到"bk_property_id"
// 的值; 然后以它的值为key,取得Info中的value,然后以value作为param访问ESB,得到其中文名。
func (lgc *Logics) GetUsernameMapWithPropertyList(c *gin.Context, objID string, infoList []mapstr.MapStr,
	config *options.Config) (map[string]string, []string, error) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	cond := metadata.QueryCondition{
		Fields: []string{metadata.AttributeFieldPropertyID},
		Condition: map[string]interface{}{
			metadata.AttributeFieldObjectID:     objID,
			metadata.AttributeFieldPropertyType: common.FieldTypeUser,
		},
	}
	attrRsp, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(c, c.Request.Header, objID, &cond)
	if err != nil {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), rid)
		return nil, nil, err
	}

	usernameList := []string{}
	propertyList := []string{}
	ok := true
	for _, info := range infoList {
		// 主机模型的info内容比inst模型的info内容多封装了一层，需要将内容提取出来。
		if objID == common.BKInnerObjIDHost {
			info, ok = info[common.BKInnerObjIDHost].(map[string]interface{})
			if !ok {
				err = fmt.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, "+
					"rid: %s", objID, rid)
				blog.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, rid: %s",
					objID, rid)
				return nil, nil, err
			}
		}
		for _, item := range attrRsp.Info {
			propertyList = append(propertyList, item.PropertyID)
			if info[item.PropertyID] != nil {
				username, ok := info[item.PropertyID].(string)
				if !ok {
					err = fmt.Errorf("failed to cast %s instance info from interface{} to string, rid: %s", objID, rid)
					blog.Errorf("failed to cast %s instance info from interface{} to string, rid: %s", objID, rid)
					return nil, nil, err
				}
				usernameList = append(usernameList, strings.Split(username, ",")...)
			}
		}
	}
	propertyList = util.RemoveDuplicatesAndEmpty(propertyList)
	userList := util.RemoveDuplicatesAndEmpty(usernameList)
	// get username from esb
	usernameMap, err := lgc.getUsernameFromEsb(c, config, userList)
	if err != nil {
		blog.ErrorJSON("get username map from ESB failed, err: %s, rid: %s", err.Error(), rid)
		return nil, nil, err
	}

	return usernameMap, propertyList, nil
}

func (lgc *Logics) getUsernameFromEsb(c *gin.Context, config *options.Config, userList []string) (map[string]string,
	error) {

	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(c.Request.Header))
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	usernameMap := map[string]string{}

	if len(userList) == 0 {
		return usernameMap, nil
	}

	user := plugins.CurrentPlugin(config.LoginVersion)

	// 处理请求的用户数据，将用户拼接成不超过500字节的字符串进行用户数据的获取
	userListStr := lgc.getUserListStr(userList)

	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	userListEsb := make([]*metadata.LoginSystemUserInfo, 0)

	for _, subStr := range userListStr {
		pipeline <- true
		wg.Add(1)
		go func(subStr string) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			lock.Lock()
			params := make(map[string]string)
			params["fields"] = "username,display_name"
			params["exact_lookups"] = subStr
			c.Request.Header = c.Request.Header.Clone()
			lock.Unlock()

			userListEsbSub, errNew := user.GetUserList(c, params)
			if errNew != nil {
				firstErr = errNew.ToCCError(defErr)
				blog.Errorf("get users(%s) list from ESB failed, err: %v, rid: %s", subStr, firstErr, rid)
				return
			}

			lock.Lock()
			userListEsb = append(userListEsb, userListEsbSub...)
			lock.Unlock()
		}(subStr)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	for _, userInfo := range userListEsb {
		username := fmt.Sprintf("%s(%s)", userInfo.EnName, userInfo.CnName)
		usernameMap[userInfo.EnName] = username
	}
	return usernameMap, nil
}

const getUserMaxLength = 500

// getUserListStr get user list str
func (lgc *Logics) getUserListStr(userList []string) []string {
	userListStr := make([]string, 0)

	userBuffer := bytes.Buffer{}
	for _, user := range userList {
		if userBuffer.Len()+len(user) > getUserMaxLength {
			userBuffer.WriteString(user)
			userStr := userBuffer.String()
			userListStr = append(userListStr, userStr)
			userBuffer.Reset()
			continue
		}

		userBuffer.WriteString(user)
		userBuffer.WriteByte(',')
	}

	if userBuffer.Len() == 0 {
		return userList
	}

	userStr := userBuffer.String()
	userListStr = append(userListStr, userStr[:len(userStr)-1])

	return userListStr
}

// GetDepartmentDetail search department detail and return a id-fullname map
func (lgc *Logics) GetDepartmentDetail(c *gin.Context, objID string, config *options.Config, infoList []mapstr.MapStr) (
	[]metadata.DepartmentItem, []string, error) {

	rid := util.GetHTTPCCRequestID(c.Request.Header)
	cond := metadata.QueryCondition{
		Fields: []string{metadata.AttributeFieldPropertyID},
		Condition: map[string]interface{}{
			metadata.AttributeFieldObjectID:     objID,
			metadata.AttributeFieldPropertyType: common.FieldTypeOrganization,
		},
	}
	attrRsp, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(c, c.Request.Header, objID, &cond)
	if err != nil {
		blog.Errorf("search object[%s] attribute failed, err: %v, rid: %s", objID, err, rid)
		return nil, nil, err
	}

	if len(attrRsp.Info) == 0 {
		return make([]metadata.DepartmentItem, 0), make([]string, 0), nil
	}

	propertyList := make([]string, 0)
	for _, item := range attrRsp.Info {
		propertyList = append(propertyList, item.PropertyID)
	}

	ok := true
	orgIDList := make([]int64, 0)
	for _, info := range infoList {
		// 主机模型的info内容比inst模型的info内容多封装了一层，需要将内容提取出来。
		if objID == common.BKInnerObjIDHost {
			info, ok = info[common.BKInnerObjIDHost].(map[string]interface{})
			if !ok {
				blog.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, rid: %s",
					objID, rid)
				return nil, nil, fmt.Errorf("failed to cast %s instance info convert to map[string]interface{}",
					objID)
			}
		}
		for _, item := range attrRsp.Info {
			if info[item.PropertyID] != nil {
				orgIDs, ok := info[item.PropertyID].([]interface{})
				if !ok {
					return nil, nil, fmt.Errorf("org id list type not []interface{}, real type is %T",
						info[item.PropertyID])
				}
				if len(orgIDs) == 0 {
					continue
				}
				for _, orgID := range orgIDs {
					id, err := util.GetInt64ByInterface(orgID)
					if err != nil {
						return nil, nil, err
					}
					orgIDList = append(orgIDList, id)
				}
			}
		}
	}

	orgIDs := util.IntArrayUnique(orgIDList)
	department, err := lgc.GetAllDepartment(c, config, orgIDs)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	return department.Results, propertyList, nil
}

// HandleExportEnumQuoteInst search inst detail and return a id-bk_inst_name map
func (lgc *Logics) HandleExportEnumQuoteInst(c *gin.Context, h http.Header, data []mapstr.MapStr, objID string,
	fields map[string]Property, rid string) ([]mapstr.MapStr, error) {

	for _, rowMap := range data {
		if objID == common.BKInnerObjIDHost {
			hostData, err := mapstr.NewFromInterface(rowMap[common.BKInnerObjIDHost])
			if err != nil {
				blog.Errorf("get host data failed, hostData: %#v, err: %v, rid: %s", hostData, err, rid)
				return nil, err
			}
			if err := lgc.getEnumQuoteInstNames(c, h, fields, rid, hostData); err != nil {
				blog.Errorf("get enum quote inst name list failed, err: %v, rid: %s", err, rid)
				return nil, err
			}
		}

		if err := lgc.getEnumQuoteInstNames(c, h, fields, rid, rowMap); err != nil {
			blog.Errorf("get enum quote inst name list failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	}

	return data, nil
}

// getEnumQuoteInstNames search inst detail and return a bk_inst_name
func (lgc *Logics) getEnumQuoteInstNames(c *gin.Context, h http.Header, fields map[string]Property, rid string,
	rowMap mapstr.MapStr) error {

	for id, property := range fields {
		switch property.PropertyType {
		case common.FieldTypeEnumQuote:
			enumQuoteIDInterface, exist := rowMap[id]
			if !exist || enumQuoteIDInterface == nil {
				continue
			}
			enumQuoteIDList, ok := enumQuoteIDInterface.([]interface{})
			if !ok {
				blog.Errorf("rowMap[%s] type to array failed, rowMap: %v, rowMap type: %T, rid: %s", property,
					rowMap[id], rowMap[id], rid)
				return fmt.Errorf("convert variable rowMap[%s] type to int array failed", property)
			}

			enumQuoteIDMap := make(map[int64]interface{}, 0)
			for _, enumQuoteID := range enumQuoteIDList {
				id, err := util.GetInt64ByInterface(enumQuoteID)
				if err != nil {
					blog.Errorf("convert enumQuoteID[%d] to int64 failed, type: %T, err: %v, rid: %s",
						enumQuoteID, enumQuoteID, rid)
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

			quoteObjID, err := GetEnumQuoteObjID(property.Option, rid)
			if err != nil {
				blog.Errorf("get enum quote option obj id failed, err: %s, rid: %s", err, rid)
				return fmt.Errorf("get enum quote option obj id failed, option: %v", property.Option)
			}

			enumQuoteIDs := make([]int64, 0)
			for enumQuoteID := range enumQuoteIDMap {
				enumQuoteIDs = append(enumQuoteIDs, enumQuoteID)
			}
			input := &metadata.QueryCondition{
				Fields: []string{common.GetInstNameField(quoteObjID)},
				Condition: mapstr.MapStr{
					common.GetInstIDField(quoteObjID): mapstr.MapStr{common.BKDBIN: enumQuoteIDs},
				},
				DisableCounter: true,
			}
			resp, err := lgc.Engine.CoreAPI.ApiServer().ReadInstance(c, h, quoteObjID, input)
			if err != nil {
				blog.Errorf("get quote inst name list failed, input: %+v, err: %v, rid: %s", input, err, rid)
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

			rowMap[id] = strings.Join(enumQuoteNames, "\n")
		}
	}

	return nil
}
