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
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/rentiansheng/xlsx"
)

// GetHostData get host data from excel
func (lgc *Logics) GetHostData(appIDStr, hostIDStr string, header http.Header) ([]mapstr.MapStr, error) {
	hostInfo := make([]mapstr.MapStr, 0)
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
	result, err := lgc.Engine.CoreAPI.ApiServer().GetHostData(context.Background(), header, sHostCond)
	if nil != err || false == result.Result {
		return hostInfo, errors.New("no host")
	}

	return result.Data.Info, nil
}

// GetImportHosts get import hosts
// return inst array data, errmsg collection, error
func (lgc *Logics) GetImportHosts(f *xlsx.File, header http.Header, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	if 0 == len(f.Sheets) {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}
	fields, err := lgc.GetObjFieldIDs(common.BKInnerObjIDHost, nil, header)
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

// ImportHosts import host info
func (lgc *Logics) ImportHosts(ctx context.Context, f *xlsx.File, header http.Header, defLang lang.DefaultCCLanguageIf) (resultData mapstr.MapStr, errCode int, err error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	hosts, errMsg, err := lgc.GetImportHosts(f, header, defLang)
	resultData = mapstr.New()

	if nil != err {
		blog.Errorf("ImportHost  get import hosts from excel err, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
	}
	if 0 != len(errMsg) {
		resultData.Set("err", errMsg)
		return resultData, common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, " file empty")
	}
	if 0 == len(hosts) {
		return nil, common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "")
	}

	params := mapstr.MapStr{}
	params["host_info"] = hosts
	params["bk_supplier_id"] = common.BKDefaultSupplierID
	params["input_type"] = common.InputTypeExcel

	result, resultErr := lgc.CoreAPI.ApiServer().AddHost(context.Background(), header, params)
	if nil != resultErr {
		blog.Errorf("ImportHosts add host info  http request  error:%s, rid:%s", resultErr.Error(), util.GetHTTPCCRequestID(header))
		return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	resultData.Merge(result.Data)
	errCode = result.Code
	err = defErr.New(result.Code, result.ErrMsg)

	if len(f.Sheets) > 2 {
		asstInfoMap := GetAssociationExcelData(f.Sheets[1], common.HostAddMethodExcelAssociationIndexOffset)
		if len(asstInfoMap) > 0 {
			asstInfoMapInput := &metadata.RequestImportAssociation{
				AssociationInfoMap: asstInfoMap,
			}
			asstResult, asstResultErr := lgc.CoreAPI.ApiServer().ImportAssociation(ctx, header, common.BKInnerObjIDHost, asstInfoMapInput)
			if nil != asstResultErr {
				blog.Errorf("ImportHosts logics http request import association error:%s, rid:%s", asstResultErr.Error(), util.GetHTTPCCRequestID(header))
				return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
			}

			resultData.Set("asst_error", asstResult.Data.ErrMsgMap)

			if result.Result && !asstResult.Result {
				errCode = asstResult.Code
				err = defErr.New(asstResult.Code, asstResult.ErrMsg)
			}

		}

	}

	return

}
