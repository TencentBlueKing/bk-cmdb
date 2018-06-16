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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	ccError "configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/util"
	scenecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/validator"
	sourceAuditAPI "configcenter/src/source_controller/api/auditlog"
	sourceAPI "configcenter/src/source_controller/api/object"
)

//AddHost, return error info
func AddHost(req *restful.Request, ownerID string, appID int, hostInfos map[int]map[string]interface{}, inputType string, moduleID []int, cc *api.APIResource) (error, []string, []string, []string) {

	hostsInst, err := NewHostsInstance(req, ownerID, inputType, common.BKDefaultDirSubArea, cc)
	if nil != err {
		blog.Errorf("get host map by hostinfos error, errror:%s", err.Error())

		return err, nil, nil, nil
	}

	hostMap, err := hostsInst.GetAddHostIDMap(hostInfos)
	if nil != err {
		blog.Errorf("get host map by hostinfos error, errror:%s", err.Error())

		return err, nil, nil, nil
	}
	_, _, err = hostsInst.GetHostAsstHande(hostInfos)
	if nil != err {
		blog.Errorf("get host assocate info  error, errror:%s", err.Error())
		return err, nil, nil, nil
	}

	// parase host excel assocate fields
	for index, host := range hostInfos {
		hostsInst.ParseHostInstanceAssocate(index, host)
	}

	var errMsg, updateErrMsg, succMsg []string

	ts := time.Now().UTC()
	//operator log
	var logConents []auditoplog.AuditLogExt
	hostLogFields, _ := GetHostLogFields(req, ownerID, cc.ObjCtrl())
	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, hostsInst.langHandle.Languagef("host_import_innerip_empty", index))
			continue
		}

		var iSubArea interface{}
		iSubArea, ok := host[common.BKCloudIDField]
		if false == ok {
			iSubArea = host[common.BKCloudIDField]
		}
		if nil == iSubArea {
			iSubArea = common.BKDefaultDirSubArea
		}

		var iHostID interface{}
		var isOK bool
		if inputType == common.InputTypeExcel {
			// host not db ,check params host info with host id
			iHostID, isOK = host[common.BKHostIDField]

		}

		if false == isOK {
			key := fmt.Sprintf("%s-%v", innerIP, iSubArea)
			iHost, isDBOK := hostMap[key]
			if isDBOK {
				isOK = isDBOK
				iHostID, _ = iHost[common.BKHostIDField]
			}

		}

		// delete system fields
		delete(host, common.BKHostIDField)

		//生产日志
		if isOK {

			hostID, _ := util.GetInt64ByInterface(iHostID)
			//prepare the log
			strHostID := fmt.Sprintf("%d", hostID)
			logObj := NewHostLog(req, common.BKDefaultOwnerID, strHostID, cc.HostCtrl(), cc.ObjCtrl(), hostLogFields)

			err := hostsInst.UpdateHostInstance(index, host, hostID)
			if nil != err {
				updateErrMsg = append(updateErrMsg, err.Error())
				continue
			}
			logContent, _ := logObj.GetHostLog(strHostID, false)
			logConents = append(logConents, auditoplog.AuditLogExt{ID: int(hostID), Content: logContent, ExtKey: innerIP})

		} else {
			//prepare the log
			logObj := NewHostLog(req, common.BKDefaultOwnerID, "", cc.HostCtrl(), cc.ObjCtrl(), hostLogFields)
			hostID, err := hostsInst.AddHostInstance(index, appID, moduleID, host, ts)
			if nil != err {
				errMsg = append(errMsg, err.Error())
				continue
			}
			strHostID := fmt.Sprintf("%d", hostID)
			logContent, _ := logObj.GetHostLog(strHostID, false)

			logConents = append(logConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: innerIP})

		}

		succMsg = append(succMsg, fmt.Sprintf("%d", index))
	}

	if 0 < len(logConents) {

		logAPIClient := sourceAuditAPI.NewClient(cc.AuditCtrl())
		_, err := logAPIClient.AuditHostsLog(logConents, "import host", ownerID, fmt.Sprintf("%d", appID), hostsInst.user, auditoplog.AuditOpTypeAdd)
		//addAuditLogs(req, logAdd, "新加主机", ownerID, appID, user, auditAddr)
		if nil != err {
			blog.Errorf("add audit log error %s", err.Error())
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return errors.New(hostsInst.langHandle.Language("host_import_err")), succMsg, updateErrMsg, errMsg
	}

	return nil, succMsg, updateErrMsg, errMsg
}

// EnterIP 将机器导入到制定模块或者空闲机器， 已经存在机器，不操作
func EnterIP(req *restful.Request, ownerID string, appID, moduleID int, ip string, cloudID int64, host map[string]interface{}, isIncrement bool, cc *api.APIResource) error {

	user := scenecommon.GetUserFromHeader(req)

	hostAddr := cc.HostCtrl()
	ObjAddr := cc.ObjCtrl()
	auditAddr := cc.AuditCtrl()

	language := util.GetActionLanguage(req)
	errHandle := cc.Error.CreateDefaultCCErrorIf(language)
	langHandle := cc.Lang.CreateDefaultCCLanguageIf(language)

	addHostURL := hostAddr + "/host/v1/insts/"

	addParams := make(map[string]interface{})
	addParams[common.BKAppIDField] = appID
	addParams[common.BKModuleIDField] = []int{moduleID}
	addModulesURL := hostAddr + "/host/v1/meta/hosts/modules/"

	isExist, err := IsExistPlat(req, ObjAddr, common.KvMap{common.BKCloudIDField: cloudID})
	if nil != err {
		return errors.New(langHandle.Languagef("plat_get_str_err", err.Error())) // "查询主机信息失败")
	}
	if !isExist {
		return errors.New(langHandle.Language("plat_id_not_exist"))
	}
	conds := map[string]interface{}{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}
	hostList, err := GetHostInfoByConds(req, hostAddr, conds, langHandle)
	if nil != err {
		return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
	}

	hostID := 0
	if len(hostList) == 0 {
		//host not exist, add host
		host[common.BKHostInnerIPField] = ip
		host[common.BKCloudIDField] = cloudID
		host["import_from"] = common.HostAddMethodAgent
		forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
		defaultFields, hasErr := getHostFields(forward, ownerID, ObjAddr)
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
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDHost, ObjAddr, forward, errHandle)
		_, hasErr = valid.ValidMap(host, "create", 0)

		if nil != hasErr {
			return hasErr
		}

		isSuccess, message, retData := GetHttpResult(req, addHostURL, common.HTTPCreate, host)
		if !isSuccess {
			return errors.New(langHandle.Languagef("host_agent_add_host_fail", message))
		}

		retHost := retData.(map[string]interface{})
		hostID, _ = util.GetIntByInterface(retHost[common.BKHostIDField])
	} else if false == isIncrement {
		//Not an additional relationship model
		return nil
	} else {
		hostMap, ok := hostList[0].(map[string]interface{})
		if false == ok {
			return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
		}
		hostID, _ = util.GetIntByInterface(hostMap[common.BKHostIDField])
		if 0 == hostID {
			return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
		}
		//func IsExistHostIDInApp(CC *api.APIResource, req *restful.Request, appID int, hostID int, defLang language.DefaultCCLanguageIf) (bool, error) {
		bl, hasErr := IsExistHostIDInApp(cc, req, appID, hostID, langHandle)
		if nil != hasErr {
			blog.Error("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", appID, hostID, hasErr.Error())
			return errHandle.Errorf(common.CCErrHostNotINAPPFail, hostID)

		}
		if false == bl {
			blog.Error("Host does not belong to the current application; error, params:{appid:%d, hostid:%s}", appID, hostID)
			return errHandle.Errorf(common.CCErrHostNotINAPP, fmt.Sprintf("%d", hostID))
		}

	}

	//del host relation from default  module
	params := make(map[string]interface{})
	params[common.BKAppIDField] = appID
	params[common.BKHostIDField] = hostID
	delModulesURL := cc.HostCtrl() + "/host/v1/meta/hosts/defaultmodules"
	isSuccess, errMsg, _ := GetHttpResult(req, delModulesURL, common.HTTPDelete, params)
	if !isSuccess {
		blog.Error("remove hosthostconfig error, params:%v, error:%s", params, errMsg)
		return errHandle.Errorf(common.CCErrHostDELResourcePool, hostID)
	}

	addParams[common.BKHostIDField] = hostID

	isSuccess, message, _ := GetHttpResult(req, addModulesURL, common.HTTPCreate, addParams)
	if !isSuccess {
		blog.Error("enterip add hosthostconfig error, params:%v, error:%s", addParams, message)
		return errors.New(langHandle.Languagef("host_agent_add_host_module_fail", message))
	}

	//prepare the log
	hostLogFields, _ := GetHostLogFields(req, ownerID, ObjAddr)
	logObj := NewHostLog(req, common.BKDefaultOwnerID, "", hostAddr, ObjAddr, hostLogFields)
	content, _ := logObj.GetHostLog(fmt.Sprintf("%d", hostID), false)
	logAPIClient := sourceAuditAPI.NewClient(auditAddr)
	logAPIClient.AuditHostLog(hostID, content, "enter IP HOST", ip, ownerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeAdd)
	logClient, err := NewHostModuleConfigLog(req, nil, hostAddr, ObjAddr, auditAddr)
	logClient.SetHostID([]int{hostID})
	logClient.SetDesc("host module change")
	logClient.SaveLog(fmt.Sprintf("%d", appID), user)
	return nil

}

type hostsInstance struct {
	forward       *sourceAPI.ForwardParam
	user          string
	hostAddr      string
	objAddr       string
	auditAddr     string
	inputType     string
	ownerID       string
	cloudID       int
	rowErr        map[int]error
	defaultFields map[string]*sourceAPI.ObjAttDes
	langHandle    lang.DefaultCCLanguageIf
	errHandle     ccError.DefaultCCErrorIf
	req           *restful.Request
	cc            *api.APIResource
	assObjectInt  *scenecommon.AsstObjectInst
	asstDes       []sourceAPI.ObjAsstDes
}

func NewHostsInstance(req *restful.Request, ownerID, inputType string, cloudID int, cc *api.APIResource) (*hostsInstance, error) {
	language := util.GetActionLanguage(req)

	h := &hostsInstance{
		req:        req,
		inputType:  inputType,
		ownerID:    ownerID,
		forward:    &sourceAPI.ForwardParam{Header: req.Request.Header},
		user:       scenecommon.GetUserFromHeader(req),
		hostAddr:   cc.HostCtrl(),
		objAddr:    cc.ObjCtrl(),
		auditAddr:  cc.AuditCtrl(),
		errHandle:  cc.Error.CreateDefaultCCErrorIf(language),
		langHandle: cc.Lang.CreateDefaultCCLanguageIf(language),
		cc:         cc,
		cloudID:    cloudID,
	}
	var err error
	h.defaultFields, err = getHostFields(h.forward, ownerID, h.objAddr)
	if nil != err {
		return nil, errors.New("get host property failure")
	}

	//get asst field
	objCli := sourceAPI.NewClient("")
	objCli.SetAddress(h.objAddr)
	asst := map[string]interface{}{}
	asst[common.BKOwnerIDField] = ownerID
	asst[common.BKObjIDField] = common.BKInnerObjIDHost
	searchData, _ := json.Marshal(asst)
	objCli.SetAddress(h.objAddr)
	h.asstDes, err = objCli.SearchMetaObjectAsst(h.forward, searchData)
	if nil != err {
		return nil, errors.New(h.langHandle.Language("host_search_fail"))
	}

	return h, nil
}

func (h *hostsInstance) ParseHostInstanceAssocate(index int, host map[string]interface{}) {
	if common.InputTypeExcel == h.inputType {
		if _, ok := h.rowErr[index]; true == ok {
			return
		}
		err := h.assObjectInt.SetObjAsstPropertyVal(host)
		if nil != err {
			blog.Error("host assocate property error %d %s", index, err.Error())
			h.rowErr[index] = err
		}

	}
}

func (h *hostsInstance) UpdateHostInstance(index int, host map[string]interface{}, hostID int64) error {
	// InputTypeApiNewHostSync method used only for synchronizing new hosts, not allowed to update
	if common.InputTypeApiNewHostSync == h.inputType {
		return h.errHandle.Error(common.CCErrCommDuplicateItem)
	}

	//delete(host, common.BKCloudIDField)
	delete(host, "import_from")
	delete(host, common.CreateTimeField)

	filterFields := []string{common.CreateTimeField}

	valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDHost, h.objAddr, filterFields, h.forward, h.errHandle)
	_, err := valid.ValidMap(host, common.ValidUpdate, int(hostID))
	if nil != err {
		blog.Error("host valid error %v %v", index, err)
		return fmt.Errorf(h.langHandle.Languagef("import_row_int_error_str", index, err.Error()))

	}
	//update host asst attr
	err = scenecommon.UpdateInstAssociation(h.objAddr, h.req, int(hostID), h.ownerID, common.BKInnerObjIDHost, host) //hostAsstData, ownerID, host)
	if nil != err {
		blog.Error("update host asst attr error : %v", err)
		return fmt.Errorf(h.langHandle.Languagef("import_row_int_error_str", index, err.Error()))
	}

	uHostURL := h.objAddr + "/object/v1/insts/host"

	condInput := make(map[string]interface{}, 1) //更新主机条件
	input := make(map[string]interface{}, 2)     //更新主机数据

	condInput[common.BKHostIDField] = hostID
	input["condition"] = condInput
	input["data"] = host
	isSuccess, message, _ := GetHttpResult(h.req, uHostURL, common.HTTPUpdate, input)
	innerIP := host[common.BKHostInnerIPField].(string)
	if !isSuccess {
		blog.Error("host update error %v %v", index, message)
		return fmt.Errorf(h.langHandle.Languagef("host_import_update_fail", index, innerIP, message))
	}
	return nil
}

func (h *hostsInstance) AddHostInstance(index, appID int, moduleID []int, host map[string]interface{}, ts time.Time) (int, error) {

	switch h.inputType {
	case common.InputTypeExcel:
		host["import_from"] = common.HostAddMethodExcel
	case common.InputTypeApiNewHostSync:
		host["import_from"] = common.HostAddMethodAPI
	}
	_, ok := host[common.BKCloudIDField]
	if false == ok {
		host[common.BKCloudIDField] = h.cloudID
	}
	filterFields := []string{common.CreateTimeField}
	valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDHost, h.objAddr, filterFields, h.forward, h.errHandle)
	_, err := valid.ValidMap(host, common.ValidCreate, 0)

	if nil != err {
		return 0, fmt.Errorf(h.langHandle.Languagef("import_row_int_error_str", index, err.Error()))
	}
	host[common.CreateTimeField] = ts

	addHostURL := h.hostAddr + "/host/v1/insts/"
	isSuccess, message, retData := GetHttpResult(h.req, addHostURL, common.HTTPCreate, host)
	if !isSuccess {
		ip, _ := host["InnerIP"].(string)
		return 0, fmt.Errorf(h.langHandle.Languagef("host_import_add_fail", index, ip, message))
	}

	retHost := retData.(map[string]interface{})
	hostID, _ := util.GetIntByInterface(retHost[common.BKHostIDField])

	//add host asst attr
	hostAsstData := scenecommon.ExtractDataFromAssociationField(int64(hostID), host, h.asstDes)
	err = scenecommon.CreateInstAssociation(h.objAddr, h.req, hostAsstData)
	if nil != err {
		blog.Error("add host asst attr error : %v", err)
		return 0, fmt.Errorf(h.langHandle.Languagef("import_row_int_error_str", index, err.Error()))
	}

	addParams := make(map[string]interface{})
	addParams[common.BKAppIDField] = appID
	addParams[common.BKModuleIDField] = moduleID
	addModulesURL := h.hostAddr + "/host/v1/meta/hosts/modules/"
	addParams[common.BKHostIDField] = hostID
	innerIP := host[common.BKHostInnerIPField].(string)

	isSuccess, message, _ = GetHttpResult(h.req, addModulesURL, common.HTTPCreate, addParams)
	if !isSuccess {
		blog.Error("add hosthostconfig error, params:%v, error:%s", addParams, message)
		return 0, fmt.Errorf(h.langHandle.Languagef("host_import_add_host_module", index, innerIP))
	}

	return hostID, nil
}

// getAddHostIDMap   InnerIP+SubArea key map[string]interface
func (h *hostsInstance) GetAddHostIDMap(hostInfos map[int]map[string]interface{}) (map[string]map[string]interface{}, error) {
	var ipArr []string
	for _, host := range hostInfos {
		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk && "" != innerIP {
			ipArr = append(ipArr, innerIP)
		}
	}

	var conds map[string]interface{}
	if 0 < len(ipArr) {
		conds = map[string]interface{}{common.BKHostInnerIPField: common.KvMap{common.BKDBIN: ipArr}}

	}

	allHostList, err := GetHostInfoByConds(h.req, h.hostAddr, conds, h.langHandle)
	if nil != err {
		return nil, errors.New(h.langHandle.Language("host_search_fail"))
	}

	hostMap := convertHostInfo(allHostList)

	return hostMap, nil
}

// getHostAsstHande get assocate object handle interface
func (h *hostsInstance) GetHostAsstHande(hostInfos map[int]map[string]interface{}) (*scenecommon.AsstObjectInst, map[int]error, error) {

	if common.InputTypeExcel == h.inputType {
		h.assObjectInt = scenecommon.NewAsstObjectInst(h.req, h.ownerID, h.objAddr, h.defaultFields, h.langHandle)
		err := h.assObjectInt.GetObjAsstObjectPrimaryKey()
		if nil != err {
			return nil, nil, fmt.Errorf("get host assocate object  property failure, error:%s", err.Error())
		}
		h.rowErr, err = h.assObjectInt.InitInstFromData(hostInfos)
		if nil != err {
			return nil, nil, fmt.Errorf("get host assocate object instance data failure, error:%s", err.Error())
		}

	}
	return h.assObjectInt, h.rowErr, nil
}
