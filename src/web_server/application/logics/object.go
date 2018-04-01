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
	webCommon "configcenter/src/web_server/common"
	_ "errors" //
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
)

// GetObjectData get object data
func GetObjectData(ownerID, objID, apiAddr string, header http.Header) ([]interface{}, error) {

	// construct the search condition
	searchCond := map[string]interface{}{
		"condition": []string{
			objID,
		},
	}

	// read objects
	url := fmt.Sprintf("%s/api/%s/object/search/batch", apiAddr, webCommon.API_VERSION)
	result, _ := httpRequest(url, searchCond, header)
	blog.Info("search object url:%s", url)
	blog.Info("search object return:%s", result)

	js, jsErr := simplejson.NewJson([]byte(result))
	if nil != jsErr {
		blog.Error("failed to unmarshal the json, error info is %s", jsErr.Error())
		return nil, jsErr
	}
	code, codeErr := js.Get("bk_error_code").Int()
	if nil != codeErr {
		blog.Errorf("failed to parse the code, error info is %s ", codeErr.Error())
		return nil, codeErr
	}

	if common.CCSuccess != code {
		msg, msgErr := js.Get("bk_error_msg").String()
		if nil != msgErr {
			blog.Error("failed to get the result, the reason is %s ", msgErr.Error())
			return nil, fmt.Errorf("failed to parse error info the reason is %s", msgErr.Error())
		}
		return nil, fmt.Errorf(msg)
	}

	// parse the result
	return js.GetPath("data", objID, "attr").Array()

}
