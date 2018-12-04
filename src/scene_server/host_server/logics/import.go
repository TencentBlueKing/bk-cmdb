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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
	hutil "configcenter/src/scene_server/host_server/util"
	"configcenter/src/scene_server/validator"
)

func (lgc *Logics) AddHost(appID int64, moduleID []int64, ownerID string, pheader http.Header, hostInfos map[int64]map[string]interface{}, importType metadata.HostInputType) ([]string, []string, []string, error) {

	instance := NewImportInstance(ownerID, pheader, lgc.Engine)
	defLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	var err error
	instance.defaultFields, err = lgc.getHostFields(ownerID, pheader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get host fields failed, err: %v", err)
	}

	cond := hutil.NewOperation().WithOwnerID(ownerID).WithObjID(common.BKInnerObjIDHost).Data()
	assResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAssociations(context.Background(), pheader, cond)
	if err != nil || (err == nil && !assResult.Result) {
		return nil, nil, nil, fmt.Errorf("search host assosications failed, err: %v, result err: %s", err, assResult.ErrMsg)
	}
	instance.asstDes = assResult.Data

	hostMap, err := lgc.getAddHostIDMap(pheader, hostInfos)
	if err != nil {
		blog.Errorf("get hosts failed, err:%s", err.Error())
		return nil, nil, nil, fmt.Errorf("get hosts failed, err: %v", err)
	}

	var errMsg, updateErrMsg, succMsg []string
	logConents := make([]auditoplog.AuditLogExt, 0)
	auditHeaders, err := lgc.GetHostAttributes(ownerID, pheader)
	if err != nil {
		return nil, nil, nil, err
	}

	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, defLang.Languagef("host_import_innerip_empty", strconv.FormatInt(index, 10)))
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
		// host not db ,check params host info with host id
		iHostID, isOK = host[common.BKHostIDField]

		if false == isOK {
			key := lgc.getHostIPCloudKey(innerIP, iSubArea)
			iHost, isDBOK := hostMap[key]
			if isDBOK {
				isOK = isDBOK
				iHostID, _ = iHost[common.BKHostIDField]
			}

		}

		var err error
		var intHostID int64
		preData := make(map[string]interface{}, 0)
		if isOK {
			intHostID, err = util.GetInt64ByInterface(iHostID)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("invalid host id: %v", iHostID)
			}
			// delete system fields
			delete(host, common.BKHostIDField)
			preData, _, _ = lgc.GetHostInstanceDetails(pheader, ownerID, strconv.FormatInt(intHostID, 10))
			// update host instance.
			if err := instance.updateHostInstance(index, host, intHostID); err != nil {
				updateErrMsg = append(updateErrMsg, err.Error())
				continue
			}

		} else {
			intHostID, err = instance.addHostInstance(int64(common.BKDefaultDirSubArea), index, appID, moduleID, host)
			if err != nil {
				errMsg = append(errMsg, err.Error())
				continue
			}
			host[common.BKHostIDField] = intHostID
			hostMap[lgc.getHostIPCloudKey(innerIP, iSubArea)] = host
		}

		succMsg = append(succMsg, strconv.FormatInt(index, 10))
		curData, _, err := lgc.GetHostInstanceDetails(pheader, ownerID, strconv.FormatInt(intHostID, 10))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}

		logConents = append(logConents, auditoplog.AuditLogExt{
			ID: intHostID,
			Content: metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: auditHeaders,
			},
			ExtKey: innerIP,
		})
	}

	if len(logConents) > 0 {
		user := util.GetUser(pheader)
		log := map[string]interface{}{
			common.BKContentField: logConents,
			common.BKOpDescField:  "import host",
			common.BKOpTypeField:  auditoplog.AuditOpTypeAdd,
		}
		_, err := lgc.CoreAPI.AuditController().AddHostLogs(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, pheader, log)
		if err != nil {
			return succMsg, updateErrMsg, errMsg, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return succMsg, updateErrMsg, errMsg, errors.New(lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader)).Language("host_import_err"))
	}

	return succMsg, updateErrMsg, errMsg, nil
}

func (lgc *Logics) getHostFields(ownerID string, pheader http.Header) (map[string]*metadata.ObjAttDes, error) {
	page := metadata.BasePage{Start: 0, Limit: common.BKNoLimit}
	opt := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).WithPage(page).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("search host attributes failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	attributesDesc := make([]metadata.ObjAttDes, 0)
	for _, att := range result.Data {
		attributesDesc = append(attributesDesc, metadata.ObjAttDes{Attribute: att})
	}

	fields := make(map[string]*metadata.ObjAttDes)
	for index, f := range attributesDesc {
		fields[f.PropertyID] = &attributesDesc[index]
	}
	return fields, nil
}

func (lgc *Logics) getHostIPCloudKey(ip, cloudID interface{}) string {
	return fmt.Sprintf("%v-%v", ip, cloudID)
}

func (lgc *Logics) getAddHostIDMap(pheader http.Header, hostInfos map[int64]map[string]interface{}) (map[string]map[string]interface{}, error) {
	var ipArr []string
	for _, host := range hostInfos {
		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk && "" != innerIP {
			ipArr = append(ipArr, innerIP)
		}
	}

	if 0 == len(ipArr) {
		return nil, fmt.Errorf("not found host inner ip fields")
	}

	var conds map[string]interface{}
	if 0 < len(ipArr) {
		conds = map[string]interface{}{common.BKHostInnerIPField: common.KvMap{common.BKDBIN: ipArr}}

	}

	query := &metadata.QueryInput{
		Condition: conds,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}
	hResult, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
	if err != nil {
		return nil, errors.New(lgc.Language.Languagef("host_search_fail_with_errmsg", err.Error()))
	} else if err == nil && !hResult.Result {
		return nil, errors.New(lgc.Language.Languagef("host_search_fail_with_errmsg", hResult.ErrMsg))
	}

	hostMap := make(map[string]map[string]interface{})
	for _, h := range hResult.Data.Info {
		key := lgc.getHostIPCloudKey(h[common.BKHostInnerIPField], h[common.BKCloudIDField])
		hostMap[key] = h
	}

	return hostMap, nil
}

type importInstance struct {
	*backbone.Engine
	pheader   http.Header
	inputType metadata.HostInputType
	ownerID   string
	// cloudID       int64
	// hostInfos     map[int64]map[string]interface{}
	defaultFields map[string]*metadata.ObjAttDes
	asstDes       []metadata.Association
	rowErr        map[int64]error
}

func NewImportInstance(ownerID string, pheader http.Header, engine *backbone.Engine) *importInstance {
	return &importInstance{
		pheader: pheader,
		Engine:  engine,
		ownerID: ownerID,
	}
}

func (h *importInstance) updateHostInstance(index int64, host map[string]interface{}, hostID int64) error {
	delete(host, "import_from")
	delete(host, common.CreateTimeField)

	filterFields := []string{common.CreateTimeField}

	valid := validator.NewValidMapWithKeyFields(util.GetOwnerID(h.pheader), common.BKInnerObjIDHost, filterFields, h.pheader, h.Engine)
	err := valid.ValidMap(host, common.ValidUpdate, hostID)
	if nil != err {
		blog.Errorf("host valid error %v %v", index, err.Error())
		return fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("import_row_int_error_str", index, err.Error()))

	}

	input := make(map[string]interface{}, 2) //更新主机数据
	input["condition"] = map[string]interface{}{common.BKHostIDField: hostID}
	input["data"] = host

	uResult, err := h.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDHost, h.pheader, input)
	if err != nil || (err == nil && !uResult.Result) {
		ip, _ := host[common.BKHostInnerIPField].(string)
		return fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_update_fail", index, ip, fmt.Sprintf("%v, %v", err, uResult.ErrMsg)))
	}
	return nil
}

func (h *importInstance) addHostInstance(cloudID, index, appID int64, moduleID []int64, host map[string]interface{}) (int64, error) {
	ip, _ := host[common.BKHostInnerIPField].(string)
	_, ok := host[common.BKCloudIDField]
	if false == ok {
		host[common.BKCloudIDField] = cloudID
	}
	filterFields := []string{common.CreateTimeField}
	valid := validator.NewValidMapWithKeyFields(util.GetOwnerID(h.pheader), common.BKInnerObjIDHost, filterFields, h.pheader, h.Engine)
	err := valid.ValidMap(host, common.ValidCreate, 0)

	if nil != err {
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("import_row_int_error_str", index, err.Error()))
	}
	host[common.CreateTimeField] = time.Now().UTC()
	result, err := h.CoreAPI.HostController().Host().AddHost(context.Background(), h.pheader, host)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add host by ip:%s, err:%v, reply err:%s", ip, err, result.ErrMsg)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, result.ErrMsg))
	}

	hID, ok := result.Data.(map[string]interface{})[common.BKHostIDField]
	if !ok {
		blog.Errorf("add host by ip:%s reply not found hostID, err:%v, reply :%v", ip, err, result.Data)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip))
	}

	hostID, err := util.GetInt64ByInterface(hID)
	if err != nil {
		blog.Errorf("add host by ip:%s reply hostID not interger, err:%v, reply :%v", ip, err, result.Data)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, err.Error()))
	}

	opt := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		ModuleID:      moduleID,
		HostID:        hostID,
	}
	hResult, err := h.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), h.pheader, opt)
	if err != nil {
		blog.Errorf("add host module by ip:%s  err:%s", ip, err.Error())
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, err.Error()))
	} else if err == nil && !hResult.Result {
		blog.Errorf("add host module by ip:%s  err:%s", ip, hResult.ErrMsg)
		return 0, fmt.Errorf(h.Language.CreateDefaultCCLanguageIf(util.GetLanguage(h.pheader)).Languagef("host_import_add_fail", index, ip, hResult.ErrMsg))
	}

	return hostID, nil
}
