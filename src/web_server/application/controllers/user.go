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
	"configcenter/src/common/http/httpclient"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/gin-gonic/gin"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/user/list", nil, GetUserList})
}

//GetUserList get user list
func GetUserList(c *gin.Context) {
	session := sessions.Default(c)
	skiplogin := session.Get("skiplogin")
	skiplogins, ok := skiplogin.(string)
	if ok && "1" == skiplogins {
		admindata := make([]interface{}, 0)
		admincell := make(map[string]interface{})
		admincell["chinese_name"] = "admin"
		admincell["english_name"] = "admin"
		admindata = append(admindata, admincell)
		c.JSON(200, gin.H{
			"result":        true,
			"bk_error_msg":  "get user list ok",
			"bk_error_code": "00",
			"data":          admindata,
		})
		return
	}

	a := api.NewAPIResource()
	config, _ := a.ParseConfig()
	accountURL := config["site.bk_account_url"]

	token := session.Get("bk_token")
	getURL := fmt.Sprintf(accountURL, token)
	blog.Info("get user list url：%s", getURL)
	httpClient := httpclient.NewHttpClient()
	type userResult struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Code    string      `json:"code"`
		Result  bool        `json:"result"`
	}
	reply, err := httpClient.GET(getURL, nil, nil)
	if nil != err {
		blog.Error("get user list error：%v", err)
		c.JSON(200, gin.H{
			"result":        false,
			"bk_error_msg":  "get user list false",
			"bk_error_code": "",
			"data":          nil,
		})
		return
	}
	blog.Info("get user list return：%s", reply)
	var result userResult
	err = json.Unmarshal([]byte(reply), &result)
	if nil != err || false == result.Result {
		c.JSON(200, gin.H{
			"result":        false,
			"bk_error_msg":  "get user list false",
			"bk_error_code": "",
			"data":          nil,
		})
		return
	}
	data, ok := result.Data.([]interface{})
	if false == ok {
		c.JSON(200, gin.H{
			"result":        false,
			"bk_error_msg":  "get user list false",
			"bk_error_code": "",
			"data":          nil,
		})
		return
	}
	info := make([]interface{}, 0)
	for _, i := range data {
		j, ok := i.(map[string]interface{})
		if false == ok {
			continue
		}
		name, ok := j["username"]
		if false == ok {
			continue
		}
		chname, ok := j["chname"]
		if false == ok {
			continue
		}
		cellData := make(map[string]interface{})
		cellData["chinese_name"] = chname
		cellData["english_name"] = name
		info = append(info, cellData)
	}
	c.JSON(200, gin.H{
		"result":        true,
		"bk_error_msg":  "get user list ok",
		"bk_error_code": "00",
		"data":          info,
	})
	return
}
