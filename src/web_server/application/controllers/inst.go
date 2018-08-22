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
	"configcenter/src/common/types"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{Verb: common.HTTPCreate, Path: "/insts/owner/:bk_supplier_account/object/:bk_obj_id/import", Params: nil, Handler: ImportInst})
	wactions.RegisterNewAction(wactions.Action{Verb: common.HTTPSelectPost, Path: "/insts/owner/:bk_supplier_account/object/:bk_obj_id/export", Params: nil, Handler: ExportInst})
}

// ImportInst import inst
func ImportInst(c *gin.Context) {
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

	apiAddr, err := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	url := apiAddr
	insts, errMsg, err := logics.GetImportInsts(f, objID, url, c.Request.Header, 0, true, defLang)
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, strings.Join(errMsg, ",")).Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebOpenFileFail, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(insts) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	blog.Debug("insts data from file:%+v", insts)
	apiSite := cc.APIAddr()
	url = apiSite + "/api/" + webCommon.API_VERSION + "/inst/" + c.Param("bk_supplier_account") + "/" + objID
	blog.Debug("batch insert insts, the url is %s", url)
	params := make(map[string]interface{})
	params["input_type"] = common.InputTypeExcel
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
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

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

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("inst")
	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return

	}

	fields, err := logics.GetObjFieldIDs(objID, apiSite, nil, c.Request.Header)
	err = logics.BuildExcelFromData(objID, fields, nil, instInfo, sheet, defLang)
	if nil != err {
		blog.Errorf("ExportHost object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)
	//fileName := fmt.Sprintf("tmp/%s_inst.xls", time.Now().UnixNano())
	logics.ProductExcelCommentSheet(file, defLang)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Error("ExportInst save file error:%s", err.Error())
		if err != nil {
			blog.Error("ExportInst save file error:%s", err.Error())
			reply := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
			c.Writer.Write([]byte(reply))
			return
		}
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("inst_%s.xlsx", objID))
	c.File(dirFileName)

	os.Remove(dirFileName)
}
