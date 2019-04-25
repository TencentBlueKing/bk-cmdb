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
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPCreate, "/netdevice/import", nil, ImportNetDevice})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectPost, "/netdevice/export", nil, ExportNetDevice})
	wactions.RegisterNewAction(wactions.Action{
		common.HTTPSelectGet, "/netcollect/importtemplate/netdevice", nil, BuildDownLoadNetDeviceExcelTemplate})
}

func ImportNetDevice(c *gin.Context) {
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)
	logics.SetProxyHeader(c)

	// open uploaded file
	file, err, errMsg := openDeviceUploadedFile(c, defErr)
	if nil != err {
		blog.Errorf("[Import Net Device] open uploaded file error:%s", err.Error())
		c.String(http.StatusInternalServerError, string(errMsg))
		return
	}

	// get data from uploaded file
	apiSite, err := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	if nil != err {
		blog.Errorf("[Import Net Device] get api site error:%s", err.Error())
		c.String(http.StatusInternalServerError, string(errMsg))
		return
	}

	netDevice, err, errMsg := getNetDevicesFromFile(c.Request.Header, defLang, defErr, file, apiSite)
	if nil != err {
		blog.Errorf("[Import Net Device] http request id:%s, error:%s", util.GetHTTPCCRequestID(c.Request.Header), err.Error())
		c.String(http.StatusInternalServerError, string(errMsg))
		return
	}

	// http request get device
	url := apiSite + fmt.Sprintf("/api/%s/collector/netcollect/device/action/batch", webCommon.API_VERSION)
	blog.V(5).Infof("[Import Net Device] add device url: %v", url)

	data := make([]interface{}, 0)
	lineNumbers := []int{}
	for line, value := range netDevice {
		data = append(data, value)
		lineNumbers = append(lineNumbers, line)
	}
	params := map[string]interface{}{"Data": data}

	blog.V(5).Infof("[Import Net Device] import device content: %v", params)

	reply, err := httpRequest(url, params, c.Request.Header)
	blog.V(5).Infof("[Import Net Device] import device result: %v", reply)

	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// rebuild response body
	reply, err = rebuildDeviceReponseBody(reply, lineNumbers)
	if nil != err {
		c.String(http.StatusInternalServerError, getReturnStr(common.CCErrWebGetAddNetDeviceResultFail,
			defErr.Errorf(common.CCErrWebGetAddNetDeviceResultFail).Error(), nil))
		return
	}

	c.String(http.StatusOK, reply)
}

func ExportNetDevice(c *gin.Context) {
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)
	logics.SetProxyHeader(c)

	apiSite, err := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	if nil != err {
		blog.Errorf("[Export Net Device] get api site error:%s", err.Error())
		c.String(http.StatusInternalServerError, getReturnStr(common.CCErrWebGetNetDeviceFail,
			defErr.Errorf(common.CCErrWebGetNetDeviceFail, err.Error()).Error(), nil))
		return
	}

	deviceIDstr := c.PostForm(common.BKDeviceIDField)
	deviceInfo, err := logics.GetNetDeviceData(c.Request.Header, apiSite, deviceIDstr)
	if nil != err {
		blog.Errorf("[Export Net Device] get device data error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(common.BKNetDevice)
	if nil != err {
		blog.Errorf("[Export Net Device] create sheet error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	fields := logics.GetNetDevicefield(defLang)
	logics.AddNetDeviceExtFields(fields, defLang)

	if err = logics.BuildNetDeviceExcelFromData(defLang, fields, deviceInfo, sheet); nil != err {
		blog.Errorf("[Export Net Device] build net device excel data error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	if _, err = os.Stat(dirFileName); nil != err {
		if err = os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); nil != err {
			blog.Errorf("[Export Net Device] mkdir error:%s", err.Error())
			msg := getReturnStr(common.CCErrWebCreateEXCELFail,
				defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
			c.String(http.StatusInternalServerError, msg)
			return
		}
	}

	fileName := fmt.Sprintf("%dnetdevice.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(file, defLang)

	if err = file.Save(dirFileName); nil != err {
		blog.Error("[Export Net Device] save file error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	logics.AddDownExcelHttpHeader(c, "netdevice.xlsx")
	c.File(dirFileName)

	if err = os.Remove(dirFileName); nil != err {
		blog.Error("[Export Net Device] remove file error:%s", err.Error())
	}
}

func BuildDownLoadNetDeviceExcelTemplate(c *gin.Context) {
	logics.SetProxyHeader(c)
	cc := api.NewAPIResource()

	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	dir := webCommon.ResourcePath + "/template/"
	if _, err := os.Stat(dir); nil != err {
		if err = os.MkdirAll(dir, os.ModeDir|os.ModePerm); nil != err {
			blog.Errorf("[Build NetDevice Excel Template] mkdir error:%s", err.Error())
			msg := getReturnStr(common.CCErrCommExcelTemplateFailed,
				defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetDevice).Error(), nil)
			c.String(http.StatusInternalServerError, msg)
			return
		}
	}

	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, common.BKNetDevice, time.Now().UnixNano(), rand.Uint32())

	apiSite := cc.APIAddr()
	if err := logics.BuildNetDeviceExcelTemplate(c.Request.Header, defLang, apiSite, file); nil != err {
		blog.Errorf("Build NetDevice Excel Template fail, error:%s", err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetDevice).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("template_%s.xlsx", common.BKNetDevice))

	c.File(file)
	if err := os.Remove(file); nil != err {
		blog.Error("[Export Net Device] remove file error:%s", err.Error())
	}
	return
}

func openDeviceUploadedFile(c *gin.Context, defErr errors.DefaultCCErrorIf) (file *xlsx.File, err error, errMsg string) {
	fileHeader, err := c.FormFile("file")
	if nil != err {
		errMsg = getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		return nil, err, errMsg
	}

	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); nil != err {
		if err = os.MkdirAll(dir, os.ModeDir|os.ModePerm); nil != err {
			blog.Errorf("[Import Net Device] mkdir error:%s", err.Error())
			errMsg = getReturnStr(common.CCErrWebFileSaveFail,
				defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
			return nil, err, errMsg
		}
	}

	filePath := fmt.Sprintf("%s/importnetdevice-%d-%d.xlsx", dir, time.Now().UnixNano(), rand.Uint32())
	if err = c.SaveUploadedFile(fileHeader, filePath); nil != err {
		errMsg = getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		return nil, err, errMsg
	}

	defer os.Remove(filePath) // del file

	file, err = xlsx.OpenFile(filePath)
	if nil != err {
		errMsg = getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		return nil, err, errMsg
	}

	return file, nil, ""
}

func getNetDevicesFromFile(
	header http.Header, defLang lang.DefaultCCLanguageIf, defErr errors.DefaultCCErrorIf, file *xlsx.File, apiSite string) (
	netDevice map[int]map[string]interface{}, err error, errMsg string) {

	netDevice, errMsgs, err := logics.GetImportNetDevices(header, defLang, file, apiSite)
	if nil != err {
		blog.Errorf("[Import Net Device] http request id:%s, error:%s", util.GetHTTPCCRequestID(header), err.Error())
		errMsg = getReturnStr(common.CCErrWebFileContentFail,
			defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(),
			nil)
		return nil, err, errMsg
	}
	if 0 != len(errMsgs) {
		errMsg = getReturnStr(common.CCErrWebFileContentFail,
			defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(),
			common.KvMap{"err": errMsgs})
		return nil, err, errMsg
	}
	if 0 == len(netDevice) {
		errMsg = getReturnStr(common.CCErrWebFileContentEmpty,
			defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		return nil, err, errMsg
	}

	return netDevice, nil, ""
}

func rebuildDeviceReponseBody(reply string, lineNumbers []int) (string, error) {
	replyBody := new(meta.Response)
	if err := json.Unmarshal([]byte(reply), replyBody); nil != err {
		blog.Errorf("[Import Net Device] unmarshal response body err: %v", err)
		return "", err
	}

	addDeviceResult, ok := replyBody.Data.([]interface{})
	if !ok {
		blog.Errorf("[Import Net Device] 'Data' field of response body convert to []interface{} fail, replyBody.Data %#+v", replyBody.Data)
		return "", fmt.Errorf("convert response body fail")
	}

	var (
		errRow  []string
		succRow []string
	)
	for i, value := range addDeviceResult {
		data, ok := value.(map[string]interface{})
		if !ok {
			blog.Errorf("[Import Net Device] traverse replyBody.Data convert to map[string]interface{} fail, data %#+v", data)
			return "", fmt.Errorf("convert response body fail")
		}

		result, ok := data["result"].(bool)
		if !ok {
			blog.Errorf("[Import Net Device] data is not bool: %#+v", data["result"])
			return "", fmt.Errorf("convert response body fail")
		}

		switch result {
		case true:
			succRow = append(succRow, strconv.Itoa(lineNumbers[i]))
		case false:
			errMsg, ok := data["error_msg"].(string)
			if !ok {
				blog.Errorf("[Import Net Device] data is not string: %#+v", data["error_msg"])
				return "", fmt.Errorf("convert response body fail")
			}

			errRow = append(errRow, fmt.Sprintf("%d行%s", lineNumbers[i], errMsg))
		}
	}

	retData := make(map[string]interface{})
	if 0 < len(succRow) {
		retData["success"] = succRow
	}
	if 0 < len(errRow) {
		retData["error"] = errRow
	}

	replyBody.Data = retData

	replyByte, err := json.Marshal(replyBody)
	if nil != err {
		blog.Errorf("[Import Net Device] convert rebuilded response body fail, error: %v", err)
		return "", err
	}

	return string(replyByte), nil
}
