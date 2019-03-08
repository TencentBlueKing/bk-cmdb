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

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) AddHost(ctx context.Context, appID int64, moduleID []int64, ownerID string, hostInfos map[int64]map[string]interface{}, importType metadata.HostInputType) ([]string, []string, []string, error) {

	instance := NewImportInstance(ctx, ownerID, lgc)
	var err error
	instance.defaultFields, err = lgc.getHostFields(ctx, ownerID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get host fields failed, err: %v", err)
	}

	hostMap, err := lgc.getAddHostIDMap(ctx, hostInfos)
	if err != nil {
		blog.Errorf("get hosts failed, err:%s", err.Error())
		return nil, nil, nil, fmt.Errorf("get hosts failed, err: %v", err)
	}

	var errMsg, updateErrMsg, succMsg []string
	logConents := make([]auditoplog.AuditLogExt, 0)
	auditHeaders, err := lgc.GetHostAttributes(ctx, ownerID, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, lgc.ccLang.Languagef("host_import_innerip_empty", strconv.FormatInt(index, 10)))
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
			preData, _, _ = lgc.GetHostInstanceDetails(ctx, ownerID, strconv.FormatInt(intHostID, 10))
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
		curData, _, err := lgc.GetHostInstanceDetails(ctx, ownerID, strconv.FormatInt(intHostID, 10))
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
		log := map[string]interface{}{
			common.BKContentField: logConents,
			common.BKOpDescField:  "import host",
			common.BKOpTypeField:  auditoplog.AuditOpTypeAdd,
		}
		_, err := lgc.CoreAPI.AuditController().AddHostLogs(ctx, ownerID, strconv.FormatInt(appID, 10), lgc.user, lgc.header, log)
		if err != nil {
			return succMsg, updateErrMsg, errMsg, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return succMsg, updateErrMsg, errMsg, errors.New(lgc.ccLang.Language("host_import_err"))
	}

	return succMsg, updateErrMsg, errMsg, nil
}

func (lgc *Logics) getHostFields(ctx context.Context, ownerID string) (map[string]*metadata.ObjAttDes, error) {
	opt := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(lgc.ownerID).WithAttrComm().MapStr()

	input := &metadata.QueryCondition{
		Condition: opt,
	}
	result, err := lgc.CoreAPI.CoreService().Model().
		ReadModelAttr(ctx, lgc.header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("getHostFields http do error, err:%s, input:%+v, rid:%s", err.Error(), input, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getHostFields http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, input, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	attributesDesc := make([]metadata.ObjAttDes, 0)
	for _, att := range result.Data.Info {
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

func (lgc *Logics) getAddHostIDMap(ctx context.Context, hostInfos map[int64]map[string]interface{}) (map[string]map[string]interface{}, error) {
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
	hResult, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, lgc.header, query)
	if err != nil {
		return nil, errors.New(lgc.ccLang.Languagef("host_search_fail_with_errmsg", err.Error()))
	}
	if !hResult.Result {
		return nil, errors.New(lgc.ccLang.Languagef("host_search_fail_with_errmsg", hResult.ErrMsg))
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
	rowErr        map[int64]error
	ctx           context.Context
	ccErr         ccErr.DefaultCCErrorIf
	ccLang        language.DefaultCCLanguageIf
	rid           string
}

func NewImportInstance(ctx context.Context, ownerID string, lgc *Logics) *importInstance {
	return &importInstance{
		pheader: lgc.header,
		Engine:  lgc.Engine,
		ownerID: ownerID,
		ctx:     ctx,
		ccErr:   lgc.ccErr,
		ccLang:  lgc.ccLang,
		rid:     lgc.rid,
	}
}

func (h *importInstance) updateHostInstance(index int64, host map[string]interface{}, hostID int64) error {
	delete(host, "import_from")
	delete(host, common.CreateTimeField)

	input := &metadata.UpdateOption{} //更新主机数据
	input.Condition = map[string]interface{}{common.BKHostIDField: hostID}
	input.Data = host
	uResult, err := h.CoreAPI.CoreService().Instance().UpdateInstance(h.ctx, h.pheader, common.BKInnerObjIDHost, input)
	if err != nil {
		ip, _ := host[common.BKHostInnerIPField].(string)
		blog.Errorf("updateHostInstance http do error,  err:%s,input:%+v,rid:%s", err.Error(), input, h.rid)
		return fmt.Errorf(h.ccLang.Languagef("host_import_update_fail", index, ip, err.Error()))
	}
	if !uResult.Result {
		ip, _ := host[common.BKHostInnerIPField].(string)
		blog.Errorf("updateHostInstance http response error,  err code:%d, err msg:%s,input:%+v,rid:%s", uResult.Code, uResult.ErrMsg, input, h.rid)
		return fmt.Errorf(h.ccLang.Languagef("host_import_update_fail", index, ip, uResult.ErrMsg))
	}
	return nil
}

func (h *importInstance) addHostInstance(cloudID, index, appID int64, moduleID []int64, host map[string]interface{}) (int64, error) {
	ip, _ := host[common.BKHostInnerIPField].(string)
	_, ok := host[common.BKCloudIDField]
	if false == ok {
		host[common.BKCloudIDField] = cloudID
	}

	input := &metadata.CreateModelInstance{
		Data: host,
	}

	result, err := h.CoreAPI.CoreService().Instance().CreateInstance(h.ctx, h.pheader, common.BKInnerObjIDHost, input) //(h.ctx, h.pheader, host)
	if err != nil {
		blog.Errorf("addHostInstance http do error,err:%s, input:%+v,rid:%s", err.Error(), host, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, err.Error()))
	}
	if !result.Result {
		blog.Errorf("addHostInstance http response error,err code:%d,err msg:%s, input:%+v,rid:%s", result.Code, result.ErrMsg, host, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, result.ErrMsg))
	}

	hostID := int64(result.Data.Created.ID)

	opt := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		ModuleID:      moduleID,
		HostID:        hostID,
	}
	hResult, err := h.CoreAPI.HostController().Module().AddModuleHostConfig(h.ctx, h.pheader, opt)
	if err != nil {
		blog.Errorf("add host module by ip:%s  err:%s,input:%+v,rid:%s", ip, err.Error(), opt, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, err.Error()))
	} else if err == nil && !hResult.Result {
		blog.Errorf("add host module by ip:%s  err code:%d,err msg:%s,input:%+v,rid:%s", ip, hResult.Code, hResult.ErrMsg, opt, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, hResult.ErrMsg))
	}

	return hostID, nil
}
