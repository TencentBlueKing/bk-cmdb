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
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	sencecommon "configcenter/src/scene_server/common"
	sourceAPI "configcenter/src/source_controller/api/object"
	"errors"
	"fmt"
	restful "github.com/emicklei/go-restful"
)

//getHostFields 获取主所有字段和默认值
func getHostFields(forward *sourceAPI.ForwardParam, ownerID, ObjAddr string) (map[string]*sourceAPI.ObjAttDes, error) {
	fields, err := GetObjectFields(forward, ownerID, common.BKInnerObjIDHost, ObjAddr, "")
	if nil != err {
		return nil, err
	}
	ret := make(map[string]*sourceAPI.ObjAttDes)
	for index, f := range fields {
		ret[f.PropertyID] = &fields[index]
	}
	return ret, nil
}

//convertHostInfo convert host info，InnerIP+SubArea key map[string]interface
func convertHostInfo(hosts []interface{}) map[string]map[string]interface{} {
	var hostMap map[string]map[string]interface{} = make(map[string]map[string]interface{})
	for _, host := range hosts {
		h := host.(map[string]interface{})

		key := fmt.Sprintf("%v-%v", h[common.BKHostInnerIPField], h[common.BKCloudIDField])
		hostMap[key] = h
	}
	return hostMap
}

//MoveHostToResourcePool move host to resource pool
func MoveHost2ResourcePool(CC *api.APIResource, req *restful.Request, appID int, hostID []int) (interface{}, error) {
	user := sencecommon.GetUserFromHeader(req)
	language := util.GetActionLanguage(req)
	//errHandle := CC.Error.CreateDefaultCCErrorIf(language)
	langHandle := CC.Lang.CreateDefaultCCLanguageIf(language)

	conds := make(map[string]interface{})
	conds[common.BKAppIDField] = appID
	appinfo, err := GetAppInfo(req, common.BKOwnerIDField, conds, CC.ObjCtrl(), langHandle)
	if err != nil {
		return nil, err
	}
	ownerID := appinfo[common.BKOwnerIDField].(string)
	if "" == ownerID {
		return nil, errors.New(langHandle.Language("host_resource_pool_not_exist")) // "未找到资源池")
	}
	//get default biz
	ownerAppID, err := GetDefaultAppID(req, ownerID, common.BKAppIDField, CC.ObjCtrl(), langHandle)
	if err != nil {
		return nil, errors.New(langHandle.Languagef("host_resource_pool_get_fail", err.Error()))

	}
	if 0 == appID {
		return nil, errors.New(langHandle.Language("host_resource_pool_not_exist")) // "未找到资源池")
		//return nil, errors.New("资源池不存在")
	}
	if ownerAppID == appID {
		return nil, errors.New(langHandle.Language("host_belong_resource_pool")) // "当前主机已经属于资源池，不需要转移")
	}

	//get resource set
	mconds := make(map[string]interface{})
	mconds[common.BKDefaultField] = common.DefaultResModuleFlag
	mconds[common.BKModuleNameField] = common.DefaultResModuleName
	mconds[common.BKAppIDField] = ownerAppID
	moduleID, err := GetSingleModuleID(req, mconds, CC.ObjCtrl())
	if nil != err {
		return nil, errors.New(langHandle.Languagef("host_resource_module_get_fail", err.Error())) //("获取资源池模块信息失败" + err.Error())
	}

	logClient, err := NewHostModuleConfigLog(req, hostID, CC.HostCtrl(), CC.ObjCtrl(), CC.AuditCtrl())

	conds[common.BKHostIDField] = hostID
	conds["bk_owner_module_id"] = moduleID
	conds["bk_owner_biz_id"] = ownerAppID
	url := CC.HostCtrl() + "/host/v1/meta/hosts/resource"
	isSucess, errmsg, data := GetHttpResult(req, url, common.HTTPUpdate, conds)
	if !isSucess {
		return data, errors.New(langHandle.Languagef("host_move_to_resource", errmsg)) //"更新主机关系失败;" + errmsg)
	}
	logClient.SetDesc("move host to resource pool")
	logErr := logClient.SaveLog(fmt.Sprintf("%d", appID), user)
	if nil != logErr {
		blog.Errorf("save host to resource pool error, hostID:%d, error:%s", hostID, logErr.Error())
	}

	return data, err
}
