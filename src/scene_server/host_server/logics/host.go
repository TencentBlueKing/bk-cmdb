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
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
	"configcenter/src/scene_server/validator"
)

func (lgc *Logics) GetHostAttributes(ownerID string, header http.Header) ([]metadata.Header, error) {
	searchOp := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), header, searchOp)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	headers := make([]metadata.Header, 0)
	for _, p := range result.Data {
		if p.PropertyID == common.BKChildStr {
			continue
		}
		headers = append(headers, metadata.Header{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
		})
	}

	return headers, nil
}

func (lgc *Logics) GetHostInstanceDetails(pheader http.Header, ownerID, hostID string) (map[string]interface{}, string, error) {
	// get host details, pre data
	result, err := lgc.CoreAPI.HostController().Host().GetHostByID(context.Background(), hostID, pheader)
	if err != nil || (err == nil && !result.Result) {
		return nil, "", fmt.Errorf("get host  data failed, err, %v, %v", err, result.ErrMsg)
	}

	hostInfo := result.Data
	attributes, err := lgc.GetObjectAsst(ownerID, pheader)
	if err != nil {
		return nil, "", err
	}

	for key, val := range attributes {
		if item, ok := hostInfo[key]; ok {
			if item == nil {
				continue
			}

			strItem := util.GetStrByInterface(item)
			ids := make([]int64, 0)
			for _, strID := range strings.Split(strItem, ",") {
				if "" == strings.TrimSpace(strID) {
					continue
				}
				id, err := strconv.ParseInt(strID, 10, 64)
				if err != nil {
					return nil, "", err
				}
				ids = append(ids, id)
			}

			//cond := make(map[string]interface{})
			//cond[common.BKHostIDField] = map[string]interface{}{"$in": ids}
			q := &metadata.QueryInput{
				Condition: nil, //cond,
				Fields:    "",
				Start:     0,
				Limit:     common.BKNoLimit,
				Sort:      "",
			}

			asst, _, err := lgc.getInstAsst(ownerID, val, strings.Split(strItem, ","), pheader, q)
			if err != nil {
				return nil, "", fmt.Errorf("get instance asst failed, err: %v", err)
			}
			hostInfo[key] = asst
		}
	}

	ip := hostInfo[common.BKHostInnerIPField].(string)
	return hostInfo, ip, nil
}

func (lgc *Logics) GetConfigByCond(pheader http.Header, cond map[string][]int64) ([]map[string]int64, error) {
	configArr := make([]map[string]int64, 0)

	if 0 == len(cond) {
		return configArr, nil
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), pheader, cond)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get module host config failed, err: %v, %v", err, result.ErrMsg)
	}

	for _, info := range result.Data {
		data := make(map[string]int64)
		data[common.BKAppIDField] = info.AppID
		data[common.BKSetIDField] = info.SetID
		data[common.BKModuleIDField] = info.ModuleID
		data[common.BKHostIDField] = info.HostID
		configArr = append(configArr, data)
	}
	return configArr, nil
}

// EnterIP 将机器导入到制定模块或者空闲机器， 已经存在机器，不操作
func (lgc *Logics) EnterIP(pheader http.Header, ownerID string, appID, moduleID int64, ip string, cloudID int64, host map[string]interface{}, isIncrement bool) error {

	langType := util.GetLanguage(pheader)
	user := util.GetUser(pheader)
	lang := lgc.Language.CreateDefaultCCLanguageIf(langType)
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(langType)

	isExist, err := lgc.IsPlatExist(pheader, common.KvMap{common.BKCloudIDField: cloudID})
	if nil != err {
		return errors.New(lang.Languagef("plat_get_str_err", err.Error())) // "查询主机信息失败")
	}
	if !isExist {
		return errors.New(lang.Language("plat_id_not_exist"))
	}
	conds := map[string]interface{}{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}
	hostList, err := lgc.GetHostInfoByConds(pheader, conds)
	if nil != err {
		return errors.New(lang.Languagef("host_search_fail", err.Error())) // "查询主机信息失败")
	}

	hostID := int64(0)
	if len(hostList) == 0 {
		//host not exist, add host
		host[common.BKHostInnerIPField] = ip
		host[common.BKCloudIDField] = cloudID
		host["import_from"] = common.HostAddMethodAgent
		defaultFields, hasErr := lgc.getHostFields(ownerID, pheader)
		if nil != hasErr {
			blog.Errorf("get host property error; error:%s", hasErr.Error())
			return errors.New("get host property error")
		}
		//补充未填写字段的默认值
		for _, field := range defaultFields {
			_, ok := host[field.PropertyID]
			if !ok {
				if true == util.IsStrProperty(field.PropertyType) {
					host[field.PropertyID] = ""
				} else {
					host[field.PropertyID] = nil
				}
			}
		}
		valid := validator.NewValidMap(util.GetOwnerID(pheader), common.BKInnerObjIDHost, pheader, lgc.Engine)
		hasErr = valid.ValidMap(host, "create", 0)

		if nil != hasErr {
			return hasErr
		}

		cond := hutil.NewOperation().WithOwnerID(ownerID).WithObjID(common.BKInnerObjIDHost).Data()
		assResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), pheader, cond)
		if err != nil || (err == nil && !assResult.Result) {
			return fmt.Errorf("search host assosications failed, err: %v, result err: %s", err, assResult.ErrMsg)
		}

		result, err := lgc.CoreAPI.HostController().Host().AddHost(context.Background(), pheader, host)
		if err != nil {
			return errors.New(lang.Languagef("host_agent_add_host_fail", err.Error()))
		} else if err == nil && !result.Result {
			return errors.New(lang.Languagef("host_agent_add_host_fail", result.ErrMsg))
		}

		retHost := result.Data.(map[string]interface{})
		hostID, err = util.GetInt64ByInterface(retHost[common.BKHostIDField])
		if err != nil {
			return errors.New(lang.Languagef("host_agent_add_host_fail", err.Error()))
		}
	} else if false == isIncrement {
		//Not an additional relationship model
		return nil
	} else {

		hostID, err = util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if err != nil {
			return errors.New(lang.Languagef("host_search_fail", err.Error())) // "查询主机信息失败"
		}
		if 0 == hostID {
			return errors.New(lang.Languagef("host_search_fail", err.Error()))
		}
		bl, hasErr := lgc.IsHostExistInApp(appID, hostID, pheader)
		if nil != hasErr {
			blog.Errorf("check host is exist in app error, params:{appid:%d, hostid:%d}, error:%s", appID, hostID, hasErr.Error())
			return ccErr.Errorf(common.CCErrHostNotINAPPFail, hostID)

		}
		if false == bl {
			blog.Errorf("Host does not belong to the current application; error, params:{appid:%d, hostid:%d}", appID, hostID)
			return ccErr.Errorf(common.CCErrHostNotINAPP, hostID)
		}

	}

	//del host relation from default  module
	conf := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
	}
	result, err := lgc.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(context.Background(), pheader, conf)
	if err != nil || (err == nil && !result.Result) {
		return ccErr.Errorf(common.CCErrHostDELResourcePool, hostID)
	}

	cfg := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		ModuleID:      []int64{moduleID},
		HostID:        hostID,
	}
	result, err = lgc.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), pheader, cfg)
	if err != nil {
		blog.Errorf("enter ip, add module host config failed, err: %v", err)
		return errors.New(lang.Languagef("host_agent_add_host_module_fail", err.Error()))
	} else if err == nil && !result.Result {
		blog.Errorf("enter ip, add module host config failed, err: %v", result.ErrMsg)
		return errors.New(lang.Languagef("host_agent_add_host_module_fail", result.ErrMsg))
	}

	audit := lgc.NewHostLog(pheader, ownerID)
	if err := audit.WithPrevious(strconv.FormatInt(hostID, 10), nil); err != nil {
		return fmt.Errorf("audit host log, but get pre data failed, err: %v", err)
	}
	content := audit.GetContent(hostID)
	log := common.KvMap{common.BKContentField: content, common.BKOpDescField: "enter ip host", common.BKHostInnerIPField: audit.ip, common.BKOpTypeField: auditoplog.AuditOpTypeAdd, "inst_id": hostID}
	aResult, err := lgc.CoreAPI.AuditController().AddHostLog(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, pheader, log)
	if err != nil || (err == nil && !aResult.Result) {
		return fmt.Errorf("audit host module log failed, err: %v, %v", err, aResult.ErrMsg)
	}

	hmAudit := lgc.NewHostModuleLog(pheader, []int64{hostID})
	if err := hmAudit.WithPrevious(); err != nil {
		return fmt.Errorf("audit host module log, but get pre data failed, err: %v", err)
	}
	if err := hmAudit.SaveAudit(strconv.FormatInt(appID, 10), user, "host module change"); err != nil {
		return fmt.Errorf("audit host module log, but get pre data failed, err: %v", err)
	}
	return nil
}

func (lgc *Logics) GetHostInfoByConds(pheader http.Header, cond map[string]interface{}) ([]mapstr.MapStr, error) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}

	result, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get hosts info failed, err: %v, %v", err, result.ErrMsg)
	}

	return result.Data.Info, nil
}

// HostSearch search host by mutiple condition
const (
	SplitFlag      = "##"
	TopoSetName    = "TopSetName"
	TopoModuleName = "TopModuleName"
)

func (lgc *Logics) GetHostIDByCond(pheader http.Header, cond map[string][]int64) ([]int64, error) {
	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), pheader, cond)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data {
		hostIDs = append(hostIDs, val.HostID)
	}

	return hostIDs, nil
}
