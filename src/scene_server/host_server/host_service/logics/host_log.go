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
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/instapi"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"errors"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

type HostLog struct {
	curData  interface{}
	preData  interface{}
	headers  []metadata.Header
	innerIP  interface{}
	req      *restful.Request
	ownerID  string
	hostCtrl string
	objCtrl  string
	instID   string
}

type HostModuleConfigLog struct {
	cur       []interface{}
	pre       []interface{}
	req       *restful.Request
	ownerID   string
	hostCtrl  string
	objCtrl   string
	auditCtrl string
	hostInfos []interface{}
	instID    []int
	desc      string
}

func NewHostLog(req *restful.Request, ownerID, instID, hostCtl, objCtrl string, headers []metadata.Header) *HostLog {
	h := &HostLog{
		req:      req,
		ownerID:  ownerID,
		instID:   instID,
		objCtrl:  objCtrl,
		hostCtrl: hostCtl,
	}
	if nil == h.headers {
		h.headers = headers
	} else {
		h.headers, _ = h.getHostFields()
	}
	if "" != instID {
		h.preData, _ = h.getHostDetail(instID)
	}
	return h
}

//getHostFields operate log content header
func (h *HostLog) getHostFields() ([]metadata.Header, int) {
	return GetHostLogFields(h.req, h.ownerID, h.objCtrl)
}

//GetInnerIP  return innerip for host detail
func (h *HostLog) GetInnerIP() string {
	if nil == h.innerIP {
		return ""
	}
	return h.innerIP.(string)
}

func (h *HostLog) getHostDetail(instID string) (interface{}, int) {

	h.instID = instID
	gHostURL := h.hostCtrl + "/host/v1/host/" + instID

	gHostRe, err := httpcli.ReqHttp(h.req, gHostURL, common.HTTPSelectGet, nil)
	if nil != err {
		blog.Error("GetHostDetail info error :%v, url:%s", err, gHostURL)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}

	// deal the association id
	instapi.Inst.InitInstHelper(h.hostCtrl, h.objCtrl)
	gHostRe, retStrErr := instapi.Inst.GetInstDetails(h.req, common.BKInnerObjIDHost, h.ownerID, gHostRe, map[string]interface{}{
		"start": 0,
		"limit": common.BKNoLimit,
		"sort":  "",
	})

	if common.CCSuccess != retStrErr {
		blog.Error("failed to replace association object, error code is %d", retStrErr)
	}
	//
	js, err := simplejson.NewJson([]byte(gHostRe))
	gHostData, _ := js.Map()
	gResult := gHostData["result"].(bool)
	if false == gResult {
		blog.Error("GetHostDetail  info error :%v", err)
		return nil, common.CC_Err_Comm_Host_Get_FAIL
	}

	hostData := gHostData["data"].(map[string]interface{})
	if nil != hostData {
		h.innerIP, _ = hostData[common.BKHostInnerIPField]
	}
	return hostData, common.CCSuccess
}

func (h *HostLog) GetPreHostData() *metadata.Content {
	logContent := &metadata.Content{}
	logContent.CurData = h.curData
	logContent.PreData = h.preData
	logContent.Headers = h.headers

	return logContent
}

func (h *HostLog) GetHostLog(instID string, isDel bool) (*metadata.Content, int) {
	//gHostURL := "http://" + cli.CC.HostCtrl + "/host/v1/host/" + hostID
	if false == isDel {
		if "" != h.instID && instID != h.instID {
			errString := fmt.Sprintf("instID error: instId not equal， source:%s, curent:%s", h.instID, instID)
			blog.Errorf(errString)
			return nil, common.CC_Err_Comm_APP_CHECK_HOST_FAIL
		}
		h.instID = instID

		h.curData, _ = h.getHostDetail(h.instID)
	}

	logContent := &metadata.Content{}
	logContent.CurData = h.curData
	logContent.PreData = h.preData
	logContent.Headers = h.headers
	return logContent, common.CCSuccess
}

//GetHostLogFields  get host fields
func GetHostLogFields(req *restful.Request, ownerID, objCtrl string) ([]metadata.Header, int) {
	gHostAttrURL := objCtrl + "/object/v1/meta/objectatts"
	searchBody := make(map[string]interface{})
	searchBody[common.BKObjIDField] = common.BKInnerObjIDHost
	searchBody[common.BKOwnerIDField] = ownerID
	searchJson, _ := json.Marshal(searchBody)
	gHostAttrRe, err := httpcli.ReqHttp(req, gHostAttrURL, common.HTTPSelectPost, []byte(searchJson))
	if nil != err {
		blog.Error("GetHostDetailByID  attr error :%v", err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}

	js, err := simplejson.NewJson([]byte(gHostAttrRe))
	gHostAttr, _ := js.Map()
	gAttrResult := gHostAttr["result"].(bool)
	if false == gAttrResult {
		blog.Error("GetHostDetailByID  attr error :%v", err)
		//cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return nil, common.CCErrCommHTTPReadBodyFailed
	}
	hostAttrArr := gHostAttr["data"].([]interface{})
	reResult := make([]metadata.Header, 0)
	for _, i := range hostAttrArr {
		attr := i.(map[string]interface{})
		data := metadata.Header{}
		propertyID := attr[common.BKPropertyIDField].(string)
		if propertyID == common.BKChildStr {
			continue
		}
		data.PropertyID = propertyID
		data.PropertyName = attr[common.BKPropertyNameField].(string)

		reResult = append(reResult, data)
	}
	return reResult, common.CCSuccess
}

func NewHostModuleConfigLog(req *restful.Request, instID []int, hostCtl, objCtrl, auditCtrl string) (*HostModuleConfigLog, error) {
	h := &HostModuleConfigLog{
		req:       req,
		instID:    instID,
		objCtrl:   objCtrl,
		hostCtrl:  hostCtl,
		auditCtrl: auditCtrl,
	}
	if nil != instID {
		h.pre = h.getHostModuleConfig()
		h.hostInfos = h.getInnerIP()
		if len(h.hostInfos) != len(h.instID) {
			blog.Infof("NewHostModuleConfigLog get hostinfo error, hostID:%v, replay:%v", instID, h.hostInfos)
		}
	}

	return h, nil
}

func (h *HostModuleConfigLog) getHostModuleConfig() []interface{} {

	conds := common.KvMap{common.BKHostIDField: h.instID}
	inputJson, _ := json.Marshal(conds)
	gHostURL := h.hostCtrl + "/host/v1/meta/hosts/module/config/search"

	gHostRe, err := httpcli.ReqHttp(h.req, gHostURL, common.HTTPSelectPost, inputJson)
	blog.Infof("GetHostModuleConfig, input:%s, return:%s", string(inputJson), gHostRe)
	if nil != err {
		blog.Error("getHostModuleConfig info error :%v, url:%s", err, gHostURL)
		return nil
	}
	//
	js, err := simplejson.NewJson([]byte(gHostRe))

	gResult, _ := js.Get("result").Bool()
	if false == gResult {
		blog.Error("getHostModuleConfig  info error :%v", err)
		return nil
	}

	//
	hostData, _ := js.Get("data").Array()
	return hostData
}

func (h *HostModuleConfigLog) getInnerIP() []interface{} {

	var dat commondata.ObjQueryInput

	dat.Fields = fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField)
	dat.Condition = common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.instID}}
	dat.Start = 0
	dat.Limit = common.BKNoLimit
	inputJson, _ := json.Marshal(dat)
	gHostURL := h.hostCtrl + "/host/v1/hosts/search"

	gHostRe, err := httpcli.ReqHttp(h.req, gHostURL, common.HTTPSelectPost, inputJson)
	blog.Infof("getInnerIP, input:%s, replay:%s", string(inputJson), gHostRe)
	if nil != err {
		blog.Error("GetInnerIP info error :%v, url:%s", err, gHostURL)
		return nil
	}
	//
	js, err := simplejson.NewJson([]byte(gHostRe))

	gResult, _ := js.Get("result").Bool()
	if false == gResult {
		blog.Error("GetHostDetail  info error :%v", err)
		return nil
	}

	//
	hostData, _ := js.Get("data").Get("info").Array()
	return hostData
}

func (h *HostModuleConfigLog) getModules(moduleIds []int) ([]interface{}, error) {
	if 0 == len(moduleIds) {
		return nil, nil
	}

	var dat commondata.ObjQueryInput

	dat.Fields = fmt.Sprintf("%s,%s,%s,%s,%s", common.BKModuleIDField, common.BKSetIDField, common.BKModuleNameField, common.BKAppIDField, common.BKOwnerIDField)
	dat.Limit = common.BKNoLimit
	dat.Start = 0
	dat.Condition = common.KvMap{common.BKModuleIDField: common.KvMap{common.BKDBIN: moduleIds}}
	bodyContent, _ := json.Marshal(dat)
	url := h.objCtrl + "/object/v1/insts/module/search"
	blog.Info("getModules url :%s", url)
	blog.Info("getModules content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(h.req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("getModules return :%s", string(reply))
	if err != nil {
		blog.Errorf("getModules url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Errorf("getModules url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())

		return nil, err
	}
	return js.Get("data").Get("info").Array()
}

func (h *HostModuleConfigLog) getSets(setIds []int) ([]interface{}, error) {
	var dat commondata.ObjQueryInput

	dat.Fields = fmt.Sprintf("%s,%s,%s", common.BKSetNameField, common.BKSetIDField, common.BKOwnerIDField)
	dat.Limit = common.BKNoLimit
	dat.Start = 0
	dat.Condition = common.KvMap{common.BKSetIDField: common.KvMap{common.BKDBIN: setIds}}
	bodyContent, _ := json.Marshal(dat)
	url := h.objCtrl + "/object/v1/insts/set/search"
	blog.Info("getSets url :%s", url)
	blog.Info("getSets content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(h.req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("getSets return :%s", string(reply))
	if err != nil {
		blog.Errorf("getSets url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Errorf("getSets url:%s, input:%s error:%s, ", url, string(bodyContent), err.Error())

		return nil, err
	}
	return js.Get("data").Get("info").Array()

}

//set host id, host id must be nil
func (h *HostModuleConfigLog) SetHostID(hostID []int) error {
	if nil == h.instID {
		h.instID = hostID
		h.hostInfos = h.getInnerIP()
		return nil
	}
	return errors.New("hostID not empty")
}

func (h *HostModuleConfigLog) SetDesc(desc string) {
	h.desc = desc
}

func (h *HostModuleConfigLog) SaveLog(appID, user string) error {
	//gHostURL := "http://" + cli.CC.HostCtrl + "/host/v1/host/" + hostID
	h.cur = h.getHostModuleConfig()

	var setIDs []int
	var moduleIDs []int
	preMap := make(map[int]map[int]interface{})
	curMap := make(map[int]map[int]interface{})

	for _, val := range h.pre {
		valMap, _ := val.(map[string]interface{})
		hostID, _ := util.GetIntByInterface(valMap[common.BKHostIDField])
		mID, _ := util.GetIntByInterface(valMap[common.BKModuleIDField])
		sID, _ := util.GetIntByInterface(valMap[common.BKSetIDField])
		if _, ok := preMap[hostID]; false == ok {
			preMap[hostID] = make(map[int]interface{}, 0)
		}
		preMap[hostID][mID] = valMap
		setIDs = append(setIDs, sID)
		moduleIDs = append(moduleIDs, mID)
	}
	for _, val := range h.cur {
		valMap, _ := val.(map[string]interface{})
		hostID, _ := util.GetIntByInterface(valMap[common.BKHostIDField])
		mID, _ := util.GetIntByInterface(valMap[common.BKModuleIDField])
		sID, _ := util.GetIntByInterface(valMap[common.BKSetIDField])
		if _, ok := curMap[hostID]; false == ok {
			curMap[hostID] = make(map[int]interface{}, 0)
		}
		curMap[hostID][mID] = valMap
		setIDs = append(setIDs, sID)
		moduleIDs = append(moduleIDs, mID)
	}
	moduels, err := h.getModules(moduleIDs)
	if nil != err {
		return fmt.Errorf("HostModuleConfigLog get module error:%s", err.Error())
	}
	sets, err := h.getSets(setIDs)
	if nil != err {
		return fmt.Errorf("HostModuleConfigLog get set error:%s", err.Error())
	}

	setMap := make(map[int]metadata.Ref, 0)
	for _, set := range sets {
		setInfo := set.(map[string]interface{})
		instID, _ := util.GetIntByInterface(setInfo[common.BKSetIDField])
		setMap[instID] = metadata.Ref{
			RefID:   instID,
			RefName: setInfo[common.BKSetNameField].(string),
		}
	}
	type ModuleRef struct {
		metadata.Ref
		Set     []interface{} `json:"set"`
		appID   interface{}
		ownerID string
	}
	moduleMap := make(map[int]ModuleRef, 0)
	for _, module := range moduels {
		moduleInfo := module.(map[string]interface{})
		mID, _ := util.GetIntByInterface(moduleInfo[common.BKModuleIDField])
		sID, _ := util.GetIntByInterface(moduleInfo[common.BKSetIDField])
		moduleRef := ModuleRef{}
		moduleRef.Set = append(moduleRef.Set, setMap[sID])
		moduleRef.RefID = mID
		moduleRef.RefName = moduleInfo[common.BKModuleNameField].(string)
		moduleRef.appID = moduleInfo[common.BKAppIDField]
		moduleRef.ownerID = moduleInfo[common.BKOwnerIDField].(string)
		moduleMap[mID] = moduleRef
	}
	moduleReName := "module"
	setRefName := "set"
	headers := []metadata.Header{
		metadata.Header{PropertyID: moduleReName, PropertyName: "module"},
		metadata.Header{PropertyID: setRefName, PropertyName: "app"},
		metadata.Header{PropertyID: common.BKAppIDField, PropertyName: "business ID"},
	}
	logs := []auditoplog.AuditLogExt{}

	for _, host := range h.hostInfos {
		host := host.(map[string]interface{})
		instID, _ := util.GetIntByInterface(host[common.BKHostIDField])
		log := auditoplog.AuditLogExt{ID: instID}
		log.ExtKey = host[common.BKHostInnerIPField].(string)

		preModule := make([]interface{}, 0)
		var preApp interface{}
		for moduleID, _ := range preMap[instID] {
			preModule = append(preModule, moduleMap[moduleID])
			preApp = moduleMap[moduleID].appID
			h.ownerID = moduleMap[moduleID].ownerID
		}

		curModule := make([]interface{}, 0)
		var curApp interface{}

		for moduleID, _ := range curMap[instID] {
			curModule = append(curModule, moduleMap[moduleID])
			curApp = moduleMap[moduleID].appID
			h.ownerID = moduleMap[moduleID].ownerID
		}

		log.Content = metadata.Content{
			PreData: common.KvMap{moduleReName: preModule, common.BKAppIDField: preApp},
			CurData: common.KvMap{moduleReName: curModule, common.BKAppIDField: curApp},
			Headers: headers,
		}
		logs = append(logs, log)

	}
	if "" == h.desc {
		h.desc = "host module change"
	}
	opClient := auditlog.NewClient(h.auditCtrl, h.req.Request.Header)
	_, err = opClient.AuditHostsLog(logs, h.desc, h.ownerID, appID, user, auditoplog.AuditOpTypeHostModule)

	return err
}
