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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	webCommon "configcenter/src/web_server/common"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/tealeg/xlsx"
)

//GetHostData get host data from excel
func GetHostData(appIDStr, hostIDStr, apiAddr string, header http.Header, kvMap map[string]string) ([]interface{}, error) {
	hostInfo := make([]interface{}, 0)
	sHostCond := make(map[string]interface{})
	appID, _ := strconv.Atoi(appIDStr)
	hostIDArr := strings.Split(hostIDStr, ",")
	iHostIDArr := make([]int, 0)
	for _, j := range hostIDArr {
		hostID, _ := strconv.Atoi(j)
		iHostIDArr = append(iHostIDArr, hostID)
	}
	if -1 != appID {
		sHostCond[common.BKAppIDField] = appID
		sHostCond["ip"] = make(map[string]interface{})
		sHostCond["condition"] = make([]interface{}, 0)
		sHostCond["page"] = make(map[string]interface{})
	} else {
		sHostCond[common.BKAppIDField] = -1
		sHostCond["ip"] = make(map[string]interface{})
		condArr := make([]interface{}, 0)
		condition := make(map[string]interface{})
		hostCondArr := make([]interface{}, 0)
		hostCond := make(map[string]interface{})
		hostCond["field"] = common.BKHostIDField
		hostCond["operator"] = "$in"
		hostCond["value"] = iHostIDArr
		hostCondArr = append(hostCondArr, hostCond)
		condition[common.BKObjIDField] = "host"
		condition["fields"] = make([]string, 0)
		condition["condition"] = hostCondArr
		condArr = append(condArr, condition)
		sHostCond["condition"] = condArr
		sHostCond["page"] = make(map[string]interface{})

	}
	url := apiAddr + fmt.Sprintf("/api/%s/hosts/search", webCommon.API_VERSION)
	result, _ := httpRequest(url, sHostCond, header)
	blog.Info("search host  url:%s", url)
	blog.Info("search host  return:%s", result)
	js, _ := simplejson.NewJson([]byte(result))
	hostData, _ := js.Map()
	hostResult := hostData["result"].(bool)
	if false == hostResult {
		return hostInfo, errors.New(hostData["bk_error_msg"].(string))
	}
	hostDataArr := hostData["data"].(map[string]interface{})
	hostInfo = hostDataArr["info"].([]interface{})
	hostCnt, _ := hostDataArr["count"].(json.Number).Int64()
	if !hostResult || 0 == hostCnt {
		return hostInfo, errors.New("no host")
	}

	url = apiAddr + fmt.Sprintf("/api/%s/object/attr/search", webCommon.API_VERSION)
	attrCond := make(map[string]interface{})
	attrCond[common.BKObjIDField] = common.BKInnerObjIDHost
	attrCond[common.BKOwnerIDField] = "0"
	result, _ = httpRequest(url, attrCond, header)
	blog.Info("get host attr  url:%s", url)
	blog.Info("get host attr return:%s", result)
	js, _ = simplejson.NewJson([]byte(result))
	hostAttr, _ := js.Map()
	attrData := hostAttr["data"].([]interface{})
	for _, j := range attrData {
		cell := j.(map[string]interface{})
		key := cell[common.BKPropertyIDField].(string)
		value, ok := cell[common.BKPropertyNameField].(string)
		if ok {
			kvMap[key] = value
		} else {
			kvMap[key] = ""
		}

	}
	return hostInfo, nil
}

//GetImportHosts get import hosts
func GetImportHosts(f *xlsx.File, url string, header http.Header) (map[int]map[string]interface{}, error) {

	if 0 == len(f.Sheets) {
		return nil, errors.New("文件内容不能为空,未找到工作簿")
	}
	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, errors.New("文件内容不能为空,工作内容不存在")
	}

	return GetExcelData(sheet, nil, common.KvMap{"import_from": common.HostAddMethodExcel}, false, 0)
}

//getPropertyTypeAliasName  return propertyType name, whether to export,
func getPropertyTypeAliasName(propertyType string) (string, bool) {
	var skip bool
	var name string
	switch propertyType {
	case common.FiledTypeSingleChar:
		name = common.FiledTypeSingleCharName
	case common.FiledTypeLongChar:
		name = common.FiledTypeLongCharName
	case common.FiledTypeInt:
		name = common.FiledTypeIntName
	case common.FiledTypeEnum:
		name = common.FiledTypeEnumName
	case common.FiledTypeDate:
		name = common.FiledTypeDateName
	case common.FiledTypeTime:
		name = common.FiledTypeDateName
	case common.FiledTypeUser:
		name = common.FiledTypeUserName
	case common.FiledTypeSingleAsst:
		name = common.FiledTypeSingleAsstName
	case common.FieldTypeMultiAsst:
		name = common.FiledTypeMultiAsstName
	case common.FiledTypeBool:
		name = common.FiledTypeBoolName
	case common.FieldTypeTimeZone:
		name = common.FiledTypeTimeZoneName
	default:
		name = "not found field type"
	}
	return name, skip
}

//httpRequest do http request
func httpRequest(url string, body interface{}, header http.Header) (string, error) {
	params, _ := json.Marshal(body)
	blog.Info("input:%s", string(params))
	httpClient := httpclient.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")

	reply, err := httpClient.POST(url, header, params)

	return string(reply), err
}

//httpRequestGet do http get request
func httpRequestGet(url string, body interface{}, header http.Header) (string, error) {
	params, _ := json.Marshal(body)
	blog.Info("input:%s", string(params))
	httpClient := httpclient.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")

	reply, err := httpClient.GET(url, header, params)

	return string(reply), err
}
