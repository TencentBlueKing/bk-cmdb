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

package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	lang "configcenter/src/common/language"
	"configcenter/src/common/types"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{Verb: common.HTTPCreate, Path: "/object/owner/:bk_supplier_account/object/:bk_obj_id/import", Params: nil, Handler: ImportObject})
	wactions.RegisterNewAction(wactions.Action{Verb: common.HTTPSelectPost, Path: "/object/owner/:bk_supplier_account/object/:bk_obj_id/export", Params: nil, Handler: ExportObject})
}

var sortFields = []string{
	"bk_property_id",
	"bk_property_name",
	"bk_property_type",
	"bk_property_group_name",
	"option",
	"unit",
	"description",
	"placeholder",
	"editable",
	"isrequired",
	"isreadonly",
	"isonly",
}

// ImportObject import object attribute
func ImportObject(c *gin.Context) {
	logics.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)

	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	err = c.SaveUploadedFile(file, filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer os.Remove(filePath) //delete file
	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)

	attrItems, errMsg, err := logics.GetImportInsts(f, objID, apiSite, c.Request.Header, 3, false, defLang)
	if 0 == len(attrItems) {
		msg := ""
		if nil != err {
			msg = getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		} else {
			msg = getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, "").Error(), nil)
		}
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, strings.Join(errMsg, ",")).Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}

	logics.ConvAttrOption(attrItems)

	blog.Debug("the object file content:%+v", attrItems)

	url := fmt.Sprintf("%s/api/%s/object/batch", apiSite, webCommon.API_VERSION)
	blog.Debug("batch insert insts, the url is %s", url)

	params := map[string]interface{}{
		objID: map[string]interface{}{
			"meta": nil,
			"attr": attrItems,
		},
	}

	blog.Debug("import the params(%+v)", params)

	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Debug("return the result:", reply)

	if nil != err {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, reply)
	}

}

func setExcelSubTitle(row *xlsx.Row) *xlsx.Row {
	for _, key := range sortFields {
		cell := row.AddCell()
		cell.Value = key
	}
	return row
}

func setExcelTitle(row *xlsx.Row, defLang lang.DefaultCCLanguageIf) *xlsx.Row {
	fields := logics.GetPropertyFieldDesc(defLang)
	for _, key := range sortFields {
		cell := row.AddCell()
		cell.Value = fields[key]
		blog.Debug("key:%s value:%v", key, fields[key])
	}
	return row
}

func setExcelTitleType(row *xlsx.Row, defLang lang.DefaultCCLanguageIf) *xlsx.Row {
	fieldType := logics.GetPropertyFieldType(defLang)
	for _, key := range sortFields {
		cell := row.AddCell()
		cell.Value = fieldType[key]
		blog.Debug("key:%s value:%v", key, fieldType[key])
	}
	return row
}

func setExcelRow(row *xlsx.Row, item interface{}) *xlsx.Row {

	itemMap, ok := item.(map[string]interface{})
	if !ok {
		blog.Debug("failed to convert to map")
		return row
	}

	// key is the object filed, value is the object filed value
	for _, key := range sortFields {

		cell := row.AddCell()
		//cell.SetValue([]string{"v1", "v2"})
		keyVal, ok := itemMap[key]
		if !ok {
			blog.Warn("not fount the key(%s), skip it", key)
			continue
		}
		blog.Debug("key:%s value:%v", key, keyVal)
		if nil == keyVal {
			cell.SetString("")
			continue
		}
		switch t := keyVal.(type) {
		case bool:
			cell.SetBool(t)
		case string:
			if "\"\"" == t {
				cell.SetValue("")
			} else {
				cell.SetValue(t)
			}
		default:
			switch key {
			case common.BKOptionField:

				bOptions, err := json.Marshal(t)
				if nil != err {
					blog.Errorf("option format error:%v", t)
					cell.SetValue("error info:" + err.Error())
				} else {
					cell.SetString(string(bOptions))
				}

			default:
				if nil != keyVal {
					cell.SetValue(t)
				}
			}
		}
	}

	return row
}

// ExportObject export object
func ExportObject(c *gin.Context) {

	logics.SetProxyHeader(c)
	cc := api.NewAPIResource()

	ownerID := c.Param(common.BKOwnerIDField)
	objID := c.Param(common.BKObjIDField)

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)

	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	// get the all attribute of the object
	arrItems, err := logics.GetObjectData(ownerID, objID, apiSite, c.Request.Header)
	if nil != err {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg)
		return
	}

	blog.Debug("the result:%+v", arrItems)

	// construct the excel file
	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()

	sheet, err = file.AddSheet(objID)

	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}

	// set the title
	setExcelTitle(sheet.AddRow(), defLang)
	setExcelTitleType(sheet.AddRow(), defLang)
	setExcelSubTitle(sheet.AddRow())

	/*
		dd := xlsx.NewXlsxCellDataValidation(true, true, true)
		dd.SetDropList([]string{})
		sheet.Col(2).SetDataValidationWithStart(dd, 3)
		sheet.Cell(1,1).SetString()
	*/

	// add the value
	for _, item := range arrItems {

		innerRow := item.(map[string]interface{})
		blog.Debug("object attribute data :%+v", innerRow)

		// set row value
		setExcelRow(sheet.AddRow(), innerRow)

	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%d_%s.xlsx", time.Now().UnixNano(), objID)
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Error("ExportInst save file error:%s", err.Error())
		fmt.Printf(err.Error())
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("inst_%s.xlsx", objID))
	c.File(dirFileName)

	os.Remove(dirFileName)

}
