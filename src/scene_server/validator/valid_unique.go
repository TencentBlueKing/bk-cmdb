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

package validator

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"strings"
)

// validCreateUnique  valid create unique
func (valid *ValidMap) validCreateUnique(valData map[string]interface{}) (bool, error) {
	isInner := false
	objID := valid.objID
	if util.InArray(valid.objID, innerObject) {
		isInner = true
	} else {
		objID = common.BKINnerObjIDObject
	}

	if 0 == len(valid.IsOnlyArr) {
		blog.Debug("is only array is zero %+v", valid.IsOnlyArr)
		return true, nil
	}
	searchCond := make(map[string]interface{})
	for key, val := range valData {
		if util.InArray(key, valid.IsOnlyArr) {
			searchCond[key] = val
		}
	}

	if !isInner {
		searchCond[common.BKObjIDField] = valid.objID
	}

	if 0 == len(searchCond) {
		return true, nil
	}
	condition := make(map[string]interface{})
	condition["condition"] = searchCond
	info, _ := json.Marshal(condition)
	httpCli := httpclient.NewHttpClient()
	httpCli.SetHeader("Content-Type", "application/json")
	httpCli.SetHeader("Accept", "application/json")
	blog.Info("get insts by cond: %s", string(info))
	url := fmt.Sprintf("%s/object/v1/insts/%s/search", valid.objCtrl, objID)
	if !strings.HasPrefix(url, "http://") {
		url = fmt.Sprintf("http://%s", url)
	}
	blog.Info("get insts by url : %s", url)
	rst, err := httpCli.POST(url, nil, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return false, err
	}

	var rstRes InstRst
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return false, jserr
	}
	if false == rstRes.Result {
		blog.Error("get rst res error :%v", rstRes)
		return false, valid.ccError.Error(common.CCErrCommUniqueCheckFailed)
	}

	data := rstRes.Data.(map[string]interface{})
	count, err := util.GetIntByInterface(data["count"])
	if nil != err {
		blog.Error("get data error :%v", data)
		return false, valid.ccError.Error(common.CCErrCommParseDataFailed)
	}
	if 0 != count {
		blog.Error("duplicate data ")
		return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
	}
	return true, nil
}

// validUpdateUnique valid update unique
func (valid *ValidMap) validUpdateUnique(valData map[string]interface{}, objID string, instID int) (bool, error) {
	isInner := false
	urlID := valid.objID
	searchOnlykeyNum := 0
	if util.InArray(valid.objID, innerObject) {
		isInner = true
	} else {
		urlID = common.BKINnerObjIDObject
	}

	if 0 == len(valid.IsOnlyArr) {
		return true, nil
	}
	searchCond := make(map[string]interface{})
	for key, val := range valData {
		if util.InArray(key, valid.IsOnlyArr) {
			searchCond[key] = val
			searchOnlykeyNum++
		}
	}

	//if no only key in params return true, if part of it return false
	if 0 == searchOnlykeyNum {
		return true, nil
	} else {
		if searchOnlykeyNum != len(valid.IsOnlyArr) {
			return false, valid.ccError.Error(common.CCErrCommUniqueCheckFailed)
		}
	}

	if 1 == len(searchCond) {
		for key := range searchCond {
			if key == common.BKAppIDField {
				return true, nil
			}
		}
	}

	if !isInner {
		searchCond[common.BKObjIDField] = valid.objID
	}
	if 0 == len(searchCond) {
		return true, nil
	}
	condition := make(map[string]interface{})
	condition["condition"] = searchCond
	info, _ := json.Marshal(condition)
	httpCli := httpclient.NewHttpClient()
	httpCli.SetHeader("Content-Type", "application/json")
	httpCli.SetHeader("Accept", "application/json")
	blog.Info("get insts by cond: %s", string(info))
	blog.Info("get insts by cond instID: %v", instID)
	rst, err := httpCli.POST(fmt.Sprintf("%s/object/v1/insts/%s/search", valid.objCtrl, urlID), nil, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return false, valid.ccError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	var rstRes InstRst
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return false, valid.ccError.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	if false == rstRes.Result {
		blog.Error("valid update unique false: %v", rstRes)
		return false, valid.ccError.Error(common.CCErrCommUniqueCheckFailed)
	}
	data := rstRes.Data.(map[string]interface{})
	count, err := util.GetIntByInterface(data["count"])
	if nil != err {
		err := "data false"
		blog.Error("data struct false %v", err)
		return false, valid.ccError.Error(common.CCErrCommParseDataFailed)
	}
	if 0 == count {
		return true, nil
	} else if 1 == count {
		info, ok := data["info"]
		if false == ok {
			blog.Error("data struct false lack info %v", data)
			return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
		}
		infoMap, ok := info.([]interface{})
		if false == ok {
			blog.Error("data struct false lack info is not array%v", data)
			return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
		}
		for _, j := range infoMap {
			i := j.(map[string]interface{})
			objIDName := util.GetObjIDByType(objID)
			instIDc, ok := i[objIDName]
			if false == ok {
				blog.Error("data struct false no objID%v", objIDName)
				return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
			}
			instIDci, err := util.GetIntByInterface(instIDc)

			if nil != err {
				blog.Error("instID not int , error info is %s", err.Error())
				return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
			}
			if instIDci == instID {
				return true, nil
			}
			blog.Error("duplicate data ")
			return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
		}
	} else {
		blog.Error("duplicate data ")
		return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
	}
	return true, nil
}
