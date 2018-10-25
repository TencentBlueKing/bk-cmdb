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
	"context"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/validator"

	redis "gopkg.in/redis.v5"
)

type OwnerManager struct {
	Engine   *backbone.Engine
	CacheCli *redis.Client
	OwnerID  string
	UserName string
	header   http.Header
}

func NewOwnerManager(userName, ownerID, language string) *OwnerManager {
	ownerManager := new(OwnerManager)
	ownerManager.UserName = userName
	ownerManager.OwnerID = ownerID

	header := make(http.Header)
	header.Add(common.BKHTTPHeaderUser, userName)
	header.Add(common.BKHTTPLanguage, language)
	header.Add(common.BKHTTPOwnerID, ownerID)
	ownerManager.header = header
	return ownerManager
}

func (m *OwnerManager) InitOwner() error {
	blog.Infof("init owner %s", m.OwnerID)

	exist, err := m.defaultAppIsExist()
	if err != nil {
		return err
	}
	if !exist {
		rediscli := m.CacheCli
		for {
			ok, err := rediscli.SetNX(common.BKCacheKeyV3Prefix+"owner_init_lock:"+m.OwnerID, m.OwnerID, 60*time.Second).Result()
			if nil != err {
				blog.Errorf("owner_init_lock error %s", err.Error())
				return err
			}
			if ok {
				defer func() {
					if err := rediscli.Del(common.BKCacheKeyV3Prefix + "owner_init_lock:" + m.OwnerID).Err(); err != nil {
						blog.Errorf("owner_init_lock error %s", err.Error())
					}
				}()
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

	result, err := m.Engine.CoreAPI.ApiServer().AddDefaultApp(context.Background(), m.header, m.OwnerID, params)
	if err != nil {
		return err
	}

	if result.Code != common.CCSuccess && result.Code != common.CCErrCommDuplicateItem {
		return fmt.Errorf("create app faild %s", result.ErrMsg)
	}
	return nil
}

func (m *OwnerManager) defaultAppIsExist() (bool, error) {

	params := make(map[string]interface{})
	params["condition"] = make(map[string]interface{})
	params["fields"] = []string{common.BKAppIDField}
	params["start"] = 0
	params["limit"] = 20
	result, err := m.Engine.CoreAPI.ApiServer().SearchDefaultApp(context.Background(), m.header, m.OwnerID, params)
	if err != nil {
		return false, err
	}

	if result.Code != common.CCSuccess {
		return false, fmt.Errorf("search default app err: %s", result.ErrMsg)
	}

	return 0 < result.Data.Count, nil
}

func (m *OwnerManager) getObjectFields(objID string) (map[string]interface{}, error) {

	conds := mapstr.MapStr{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"skip": 0, "limit": common.BKNoLimit}}
	result, err := m.Engine.CoreAPI.ApiServer().GetObjectAttr(context.Background(), m.header, conds)
	if err != nil {
		return nil, err
	}

	fields := result.Data

	ret := map[string]interface{}{}
	validator.FillLostedFieldValue(ret, fields, nil)
	return ret, nil
}
