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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPCreate, "/netproperty/import", nil, ImportNetProperty})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectPost, "/netproperty/export", nil, ExportNetProperty})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/netcollect/importtemplate/netproperty", nil, BuildDownLoadNetPropertyExcelTemplate})
}

func ImportNetProperty(c *gin.Context) {
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("Import Net Device get file error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	logics.SetProxyHeader(c)

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	filePath := fmt.Sprintf("%s/importnetdevice-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	err = c.SaveUploadedFile(file, filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer os.Remove(filePath) //del file
	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	hosts, errMsg, err := logics.GetImportHosts(f, apiSite, c.Request.Header, defLang) //TODO

	if nil != err {
		blog.Errorf("ImportHost logID:%s, error:%s", util.GetHTTPCCRequestID(c.Request.Header), err.Error())
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(hosts) {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	url := apiSite + fmt.Sprintf("/api/%s/hosts/add", webCommon.API_VERSION)
	params := make(map[string]interface{})
	params["host_info"] = hosts
	params["bk_supplier_id"] = common.BKDefaultSupplierID
	params["input_type"] = common.InputTypeExcel

	blog.Infof("add host url: %v", url)
	blog.Infof("add host content: %v", params)
	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Infof("add host result: %v", reply)

	if nil != err {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, reply)
	}

}

func ExportNetProperty(c *gin.Context) {

}

func BuildDownLoadNetPropertyExcelTemplate(c *gin.Context) {

}
