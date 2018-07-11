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

	"github.com/tidwall/gjson"
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

	// only search data not in diable status
	searchCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}

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
	rst, err := httpCli.POST(url, valid.forward.Header, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return false, err
	}
	count := gjson.Get(string(rst), "data.count").Int()

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
	if util.InArray(valid.objID, innerObject) {
		isInner = true
	} else {
		urlID = common.BKINnerObjIDObject
	}

	if 0 == len(valid.IsOnlyArr) {
		return true, nil
	}

	mapData := valid.getInstDataById(objID, instID)
	searchCond := make(map[string]interface{})

	for key, val := range mapData {
		if util.InArray(key, valid.IsOnlyArr) {
			searchCond[key] = val
		}
	}
	for key, val := range valData {
		if util.InArray(key, valid.IsOnlyArr) {
			searchCond[key] = val
		}
	}
	objIDName := util.GetObjIDByType(objID)
	searchCond[objIDName] = map[string]interface{}{common.BKDBNE: instID}

	// only search data not in diable status
	searchCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}

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
	blog.Infof("get insts by cond: %s, instID %v", string(info), instID)
	rst, err := httpCli.POST(fmt.Sprintf("%s/object/v1/insts/%s/search", valid.objCtrl, urlID), valid.forward.Header, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return false, valid.ccError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	count := gjson.Get(string(rst), "data.count").Int()
	if 0 != count {
		blog.Error("duplicate data ")
		return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
	}
	return true, nil
}

// getInstDataById get inst data by id
func (valid *ValidMap) getInstDataById(objID string, instID int) map[string]interface{} {
	isInner := false
	urlID := valid.objID

	if util.InArray(valid.objID, innerObject) {
		isInner = true
	} else {
		urlID = common.BKINnerObjIDObject
	}

	if 0 == len(valid.IsOnlyArr) {
		return nil
	}
	searchCond := make(map[string]interface{})

	if !isInner {
		searchCond[common.BKObjIDField] = objID
		searchCond[common.BKInstIDField] = instID
	} else {
		objIDName := util.GetObjIDByType(objID)
		searchCond[objIDName] = instID

	}
	condition := make(map[string]interface{})
	condition["condition"] = searchCond
	info, _ := json.Marshal(condition)
	httpCli := httpclient.NewHttpClient()
	httpCli.SetHeader("Content-Type", "application/json")
	httpCli.SetHeader("Accept", "application/json")
	blog.Infof("get insts by cond: %s instID: %v", string(info), instID)
	rst, err := httpCli.POST(fmt.Sprintf("%s/object/v1/insts/%s/search", valid.objCtrl, urlID), valid.forward.Header, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return nil
	}
	data := make(map[string]interface{})
	result := gjson.Get(string(rst), "data.info.0").Map()
	for key, val := range result {
		data[key] = val.Raw
	}

	return data

}
