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

package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/validator"
	"configcenter/src/web_server/application/middleware/types"
	webCommon "configcenter/src/web_server/common"
)

type OwnerManager struct {
	httpCli  *httpclient.HttpClient
	OwnerID  string
	UserName string
}

func NewOwnerManager(userName, ownerID, language string) *OwnerManager {
	ownerManager := new(OwnerManager)
	ownerManager.UserName = userName
	ownerManager.OwnerID = ownerID
	ownerManager.httpCli = httpclient.NewHttpClient()
	ownerManager.httpCli.SetHeader(common.BKHTTPHeaderUser, userName)
	ownerManager.httpCli.SetHeader(common.BKHTTPLanguage, language)
	ownerManager.httpCli.SetHeader(common.BKHTTPOwnerID, ownerID)
	return ownerManager
}

func (m *OwnerManager) InitOwner() error {
	blog.Infof("init owner %s", m.OwnerID)

	exist, err := m.defaultAppIsExist()
	if err != nil {
		return err
	}
	if !exist {
		rediscli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
		for {
			ok, err := rediscli.SetNX(common.BKCacheKeyV3Prefix+"owner_init_lock"+m.OwnerID, m.OwnerID, 60*time.Second).Result()
			if nil != err {
				blog.Errorf("owner_init_lock error %s", err.Error())
				return err
			}
			if ok {
				defer rediscli.Del(common.BKCacheKeyV3Prefix + "owner_init_lock" + m.OwnerID)
				break
			}
			time.Sleep(time.Second)
		}
		exist, err = m.defaultAppIsExist()
		if err != nil {
			return err
		}
		if !exist {
			err = m.addDefaultApp()
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func (m *OwnerManager) addDefaultApp() error {
	blog.Info("addDefaultApp %s", m.OwnerID)
	params, err := m.getObjectFields(common.BKInnerObjIDApp)
	if err != nil {
		return err
	}
	params[common.BKAppNameField] = common.DefaultAppName
	params[common.BKMaintainersField] = "admin"
	params[common.BKProductPMField] = "admin"
	params[common.BKTimeZoneField] = "Asia/Shanghai"
	params[common.BKLanguageField] = "1" //中文
	params[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	byteParams, _ := json.Marshal(params)
	url := fmt.Sprintf("%s/api/%s/biz/default/%s", api.GetAPIResource().APIAddr(), webCommon.API_VERSION, m.OwnerID)
	reply, err := m.httpCli.POST(url, nil, byteParams)
	if err != nil {
		return err
	}

	result := CreateAppResult{}
	err = json.Unmarshal(reply, &result)
	if nil != err {
		return err
	}

	if result.Code != common.CCSuccess && result.Code != common.CCErrCommDuplicateItem {
		return fmt.Errorf("create app faild %s", result.Message)
	}
	return nil
}

func (m *OwnerManager) defaultAppIsExist() (bool, error) {
	params := make(map[string]interface{})
	params["condition"] = make(map[string]interface{})
	params["fields"] = []string{common.BKAppIDField}
	params["start"] = 0
	params["limit"] = 20

	byteParams, _ := json.Marshal(params)
	url := fmt.Sprintf("%s/api/%s/biz/default/%s/search", api.GetAPIResource().APIAddr(), webCommon.API_VERSION, m.OwnerID)

	reply, err := m.httpCli.POST(url, nil, byteParams)
	if err != nil {
		return false, err
	}

	result := types.SearchAppResult{}
	err = json.Unmarshal(reply, &result)
	if nil != err {
		return false, err
	}

	if result.Code != common.CCSuccess {
		return false, fmt.Errorf("search default app err: %s", result.Message)
	}

	return 0 < result.Data.Count, nil
}

func (m *OwnerManager) getObjectFields(objID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s/object/attr/search", api.GetAPIResource().APIAddr(), webCommon.API_VERSION)
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	byteParams, _ := json.Marshal(conds)
	reply, err := m.httpCli.POST(url, nil, byteParams)
	if err != nil {
		return nil, err
	}

	replyVal := gjson.ParseBytes(reply)
	if !replyVal.Get("result").Bool() {
		return nil, fmt.Errorf("get object fields faile: %s", replyVal.Get(common.HTTPBKAPIErrorMessage))
	}

	fields := []metadata.Attribute{}
	err = json.Unmarshal([]byte(replyVal.Get("data").String()), &fields)
	if nil != err {
		return nil, err
	}

	ret := map[string]interface{}{}
	validator.FillLostedFieldValue(ret, fields, nil)
	return ret, nil
}

type CreateAppResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"bk_error_code"`
	Message interface{} `json:"bk_error_msg"`
}
