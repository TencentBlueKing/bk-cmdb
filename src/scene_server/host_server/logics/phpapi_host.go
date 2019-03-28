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
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
)

// helpers
func (phpapi *PHPAPI) UpdateHostMain(hostCondition, data map[string]interface{}, appID int64) (string, error) {
	//blog.V(5).Infof("updateHostMain start")
	blog.V(5).Infof("hostCondition:%v", hostCondition)

	_, hostIDArr, err := phpapi.GetHostMapByCond(hostCondition)

	blog.V(5).Infof("hostIDArr:%v", hostIDArr)
	if nil != err {
		return "", fmt.Errorf("GetHostIDByCond error:%v", err)
	}

	lenOfHostIDArr := len(hostIDArr)
	if lenOfHostIDArr != 1 {
		blog.V(5).Infof("GetHostMapByCond condition: %v, host:%v", hostCondition, hostIDArr)
		return "", errors.New("not find host info ")
	}

	ownerID := util.GetOwnerID(phpapi.header)
	valid := validator.NewValidMapWithKeyFields(ownerID, common.BKInnerObjIDHost, []string{common.CreateTimeField, common.LastTimeField, common.BKChildStr, common.BKOwnerIDField}, phpapi.header, phpapi.logic.Engine)
	validErr := valid.ValidMap(data, common.ValidUpdate, hostIDArr[0])
	if nil != validErr {
		blog.Errorf("updateHostMain error: %v", validErr)
		return "", validErr
	}

	configData, err := phpapi.logic.GetConfigByCond(phpapi.header, map[string][]int64{
		common.BKAppIDField:  []int64{appID},
		common.BKHostIDField: []int64{hostIDArr[0]},
	})

	if nil != err {
		return "", fmt.Errorf("GetConfigByCond error:%v", err)
	}

	lenOfConfigData := len(configData)

	if lenOfConfigData == 0 {
		blog.Errorf("not expected config lenth: appid:%d, hostid:%d", appID, hostIDArr[0])
		return "", fmt.Errorf("not expected config length: %d", lenOfConfigData)
	}

	hostID := configData[0][common.BKHostIDField]

	condition := make(map[string]interface{})
	condition[common.BKHostIDField] = hostID

	param := make(map[string]interface{})
	param["condition"] = condition
	param["data"] = data

	strHostID := strconv.FormatInt(hostID, 10)
	logContent := phpapi.logic.NewHostLog(phpapi.header, util.GetActionOnwerIDByHTTPHeader(phpapi.header))
	logContent.WithPrevious(strHostID, nil)
	res, err := phpapi.logic.Engine.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDHost, phpapi.header, param)
	if nil != err {
		return "", err
	}
	if false == res.Result {
		return "", errors.New(res.ErrMsg)
	}
	if nil == err && true == res.Result {
		//操作成功，新加操作日志日志resJs, err := simplejson.NewJson([]byte(res))
		if res.Result {
			user := util.GetUser(phpapi.header)
			ownerID := util.GetOwnerID(phpapi.header)
			logContent.WithCurrent(strHostID)
			content := logContent.GetContent(hostID)
			//(id interface{}, Content interface{}, OpDesc string, InnerIP, ownerID, appID, user string, OpType auditoplog.AuditOpType)
			phpapi.logic.CoreAPI.AuditController().AddHostLog(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, phpapi.header, content)
		}
	}

	return "", nil
}

func (phpapi *PHPAPI) AddHost(data map[string]interface{}) (int64, error) {
	hostID, err := phpapi.addObj(data, common.BKInnerObjIDHost)
	return hostID, err
}

func (phpapi *PHPAPI) AddModuleHostConfig(hostID, appID int64, moduleIDs []int64) error {
	data := &meta.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
		ModuleID:      moduleIDs,
	}
	blog.V(5).Infof("addModuleHostConfig start, data: %v", data)

	res, err := phpapi.logic.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), phpapi.header, data)
	if nil != err {
		blog.Errorf("AddModuleHostConfig http do error. cond:%#v, error:%s", data, err.Error())
		return err
	}

	if !res.Result {
		return errors.New(res.ErrMsg)
	}
	blog.V(5).Infof("addModuleHostConfig success, res: %v", res)
	return nil
}

func (phpapi *PHPAPI) addObj(data map[string]interface{}, objType string) (int64, error) {
	resMap := make(map[string]interface{})
	resp, err := phpapi.logic.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), objType, phpapi.header, data)
	if nil != err {
		return 0, err
	}

	if !resp.Result {
		return 0, errors.New(resp.ErrMsg)
	}

	blog.V(5).Infof("add object result : %v", resMap)

	objID, err := resp.Data.Int64(common.GetInstIDField(objType))
	if nil != err {
		blog.Errorf("addObj get id error, reply:%v, error:%s", resp, err.Error())
		return 0, fmt.Errorf("add object reply error, not found  id")
	}
	return objID, nil
}

//search host helpers

func (phpapi *PHPAPI) SetHostData(moduleHostConfig []map[string]int64, hostMap map[int64]map[string]interface{}) ([]mapstr.MapStr, error) {

	//total data
	hostData := make([]mapstr.MapStr, 0)

	appIDArr := make([]int64, 0)
	setIDArr := make([]int64, 0)
	moduleIDArr := make([]int64, 0)

	for _, config := range moduleHostConfig {
		setIDArr = append(setIDArr, config[common.BKSetIDField])
		moduleIDArr = append(moduleIDArr, config[common.BKModuleIDField])
		appIDArr = append(appIDArr, config[common.BKAppIDField])
	}

	moduleMap, err := phpapi.logic.GetModuleMapByCond(phpapi.header, "", map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDArr,
		},
	})
	if err != nil {
		return hostData, err
	}

	setMap, err := phpapi.logic.GetSetMapByCond(phpapi.header, "", map[string]interface{}{
		common.BKSetIDField: map[string]interface{}{
			common.BKDBIN: setIDArr,
		},
	})
	if err != nil {
		blog.Errorf("hostMap GetSetMapByCond  error, err:%s,rid:%s", err.Error(), phpapi.rid)
		return hostData, err
	}

	blog.V(5).Infof("GetAppMapByCond , appIDArr:%v, rid:%s", appIDArr, phpapi.rid)
	appMap, err := phpapi.logic.GetAppMapByCond(phpapi.header, "", map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: appIDArr,
		},
	})

	if err != nil {
		blog.Errorf("hostMap GetAppMapByCond  error, err:%s,rid:%s", err.Error(), phpapi.rid)
		return hostData, err
	}
	for _, config := range moduleHostConfig {
		hostItem, hasHost := hostMap[config[common.BKHostIDField]]
		if !hasHost {
			blog.Errorf("hostMap has not hostID: %d,rid:%s", config[common.BKHostIDField], phpapi.rid)
			continue
		}
		host := mapstr.New()
		host.Merge(hostItem)

		module := moduleMap[config[common.BKModuleIDField]]
		set := setMap[config[common.BKSetIDField]]
		app := appMap[config[common.BKAppIDField]]

		host[common.BKModuleIDField] = module[common.BKModuleIDField]
		host[common.BKModuleNameField] = module[common.BKModuleNameField]
		host[common.BKSetIDField], _ = set.Int64(common.BKSetIDField) //[common.BKSetIDField]
		host[common.BKSetNameField] = set[common.BKSetNameField]
		host[common.BKAppIDField], _ = app.Int64(common.BKAppIDField) //[common.BKAppIDField]
		host[common.BKAppNameField] = app[common.BKAppNameField]
		host[common.BKModuleTypeField] = module[common.BKModuleTypeField]
		host[common.BKOwnerIDField] = app[common.BKOwnerIDField]
		host[common.BKOperatorField] = module[common.BKOperatorField]
		host[common.BKBakOperatorField] = module[common.BKBakOperatorField]
		host[common.BKSupplierIDField] = app[common.BKSupplierIDField]

		hostData = append(hostData, host)
	}
	return hostData, nil
}
