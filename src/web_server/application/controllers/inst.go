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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"configcenter/src/common/types"

	"configcenter/src/web_server/application/logics"

	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{Verb: common.HTTPCreate, Path: "/insts/owner/:bk_supplier_account/object/:bk_obj_id/import", Params: nil, Handler: ImportInst})
	wactions.RegisterNewAction(wactions.Action{Verb: common.HTTPSelectPost, Path: "/insts/owner/:bk_supplier_account/object/:bk_obj_id/export", Params: nil, Handler: ExportInst})
}

// ImportInst import inst
func ImportInst(c *gin.Context) {
	logics.SetProxyHeader(c)

	cc := api.NewAPIResource()
	language := util.GetActionLanguageByHTTPHeader(c.Request.Header)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)
	c.Header(common.BKHTTPLanguage, language)

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

	apiAddr, err := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	url := apiAddr
	insts, err := logics.GetImportInsts(f, url, c.Request.Header, 0, defLang)
	if 0 == len(insts) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebOpenFileFail, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	blog.Debug("insts data from file:%+v", insts)
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	url = apiSite + "/api/" + webCommon.API_VERSION + "/inst/" + c.Param("bk_supplier_account") + "/" + c.Param("bk_obj_id")
	blog.Debug("batch insert insts, the url is %s", url)
	params := make(map[string]interface{})
	params["BatchInfo"] = insts
	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Debug("return the result:", reply)
	if nil != err {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, reply)
	}

}

// ExportInst export inst
func ExportInst(c *gin.Context) {
	logics.SetProxyHeader(c)
	cc := api.NewAPIResource()
	language := util.GetActionLanguageByHTTPHeader(c.Request.Header)
	//defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)
	c.Header(common.BKHTTPLanguage, language)


	ownerID := c.Param(common.BKOwnerIDField)
	objID := c.Param(common.BKObjIDField)
	instIDStr := c.PostForm(common.BKInstIDField)

	kvMap := make(map[string]string)
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	instInfo, err := logics.GetInstData(ownerID, objID, instIDStr, apiSite, c.Request.Header, kvMap)
	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)

		c.String(http.StatusBadGateway, msg, nil)
		return
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("inst")
	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return

	}
	row = sheet.AddRow()
	kArray := make([]string, 0)
	for i, k := range kvMap {
		cell = row.AddCell()
		cell.Value = k
		kArray = append(kArray, i)
	}
	kLength := len(kArray)
	for _, j := range instInfo {
		instcell := j.(map[string]interface{})
		blog.Debug("inst data :%#v", instcell)
		//instcell := instData["Host"].(map[string]interface{})
		row = sheet.AddRow()
		for i := 0; i != kLength; i++ {
			cell = row.AddCell()
			kName := kArray[i]

			n, ok := instcell[kName]

			if ok {
				switch dval := n.(type) {
				case string:
					cell.Value = dval
				case int, int8, int16, int32, int64:
					cell.SetInt64(dval.(int64))
				case []interface{}:
					cell.Value = ""
					for _, idxVal := range dval {
						if instVal, ok := idxVal.(map[string]interface{}); ok {
							if id, idOk := instVal["id"]; idOk {
								if 0 != len(cell.Value) {
									cell.Value += ","
								}
								if name, nameOk := instVal["name"]; nameOk {
									cell.Value += fmt.Sprintf("%+v:%+v", id, name)
								} else {
									cell.Value += fmt.Sprintf("%+v:", id)
								}
							} else {
								cell.Value += ""
							}
						} else {
							cell.Value += ""
						}
					}
				case nil:
					cell.Value = ""
				default:
					cell.SetValue(n)
					blog.Debug("the %s kind is %s,value %#v", kName, reflect.TypeOf(n), n)
				}

			} else {
				cell.Value = ""
			}
		}
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)
	//fileName := fmt.Sprintf("tmp/%s_inst.xls", time.Now().UnixNano())
	err = file.Save(dirFileName)
	if err != nil {
		blog.Error("ExportInst save file error:%s", err.Error())
		fmt.Printf(err.Error())
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("inst_%s.xlsx", objID))
	c.File(dirFileName)

	os.Remove(dirFileName)
}
