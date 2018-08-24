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
	lang "configcenter/src/common/language"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/rentiansheng/xlsx"

	webCommon "configcenter/src/web_server/common"
)

//GetHostData get host data from excel
func GetHostData(appIDStr, hostIDStr, apiAddr string, header http.Header) ([]interface{}, error) {
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

		//host condition
		condition := make(map[string]interface{})
		hostCondArr := make([]interface{}, 0)
		hostCond := make(map[string]interface{})
		hostCond["field"] = common.BKHostIDField
		hostCond["operator"] = common.BKDBIN
		hostCond["value"] = iHostIDArr
		hostCondArr = append(hostCondArr, hostCond)
		condition[common.BKObjIDField] = common.BKInnerObjIDHost
		condition["fields"] = make([]string, 0)
		condition["condition"] = hostCondArr
		condArr = append(condArr, condition)

		//biz conditon
		condition = make(map[string]interface{})
		condition[common.BKObjIDField] = common.BKInnerObjIDApp
		condition["fields"] = make([]interface{}, 0)
		condition["condition"] = make([]interface{}, 0)
		condArr = append(condArr, condition)

		//set conditon
		condition = make(map[string]interface{})
		condition[common.BKObjIDField] = common.BKInnerObjIDSet
		condition["fields"] = make([]interface{}, 0)
		condition["condition"] = make([]interface{}, 0)
		condArr = append(condArr, condition)

		//module condition
		condition = make(map[string]interface{})
		condition[common.BKObjIDField] = common.BKInnerObjIDModule
		condition["fields"] = make([]interface{}, 0)
		condition["condition"] = make([]interface{}, 0)
		condArr = append(condArr, condition)

		sHostCond["condition"] = condArr
		sHostCond["page"] = make(map[string]interface{})

	}
	url := apiAddr + fmt.Sprintf("/api/%s/hosts/search/asstdetail", webCommon.API_VERSION)
	result, _ := httpRequest(url, sHostCond, header)
	blog.Infof("search host  url:%s", url)
	blog.Infof("search host  return:%s", result)
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

	return hostInfo, nil
}

// GetImportHosts get import hosts
// return inst array data, errmsg collection, error
func GetImportHosts(f *xlsx.File, url string, header http.Header, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	if 0 == len(f.Sheets) {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}
	fields, err := GetObjFieldIDs(common.BKInnerObjIDHost, url, nil, header)
	if nil != err {
		return nil, nil, errors.New(defLang.Languagef("web_get_object_field_failure", err.Error()))
	}

	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	return GetExcelData(sheet, fields, common.KvMap{"import_from": common.HostAddMethodExcel}, true, 0, defLang)
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
