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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) AddHost(ctx context.Context, appID int64, moduleID []int64, ownerID string, hostInfos map[int64]map[string]interface{}, importType metadata.HostInputType) ([]int64, []string, []string, []string, error) {

	hostIDs := make([]int64, 0)
	instance := NewImportInstance(ctx, ownerID, lgc)
	var err error
	instance.defaultFields, err = lgc.getHostFields(ctx, ownerID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	hostIDMap, err := instance.GetHostIDByHostInfoArr(ctx, hostInfos)
	if err != nil {
		blog.Errorf("get hosts failed, err:%s, rid:%s", err.Error(), lgc.rid)
		return nil, nil, nil, nil, err
	}

	var errMsg, updateErrMsg, succMsg []string
	logConents := make([]metadata.SaveAuditLogParams, 0)
	auditHeaders, err := lgc.GetHostAttributes(ctx, ownerID, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, lgc.ccLang.Languagef("host_import_innerip_empty", index))
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

		iSubAreaVal, err := util.GetInt64ByInterface(iSubArea)
		if err != nil || iSubAreaVal < 0 {
			errMsg = append(errMsg, lgc.ccLang.Language("import_host_cloudID_invaild"))
			continue
		}

		var intHostID int64
		var existInDB bool

		// we support update host info both base on hostID and innerIP, hostID has higher priority then innerIP
		hostIDFromInput, bHostIDInInput := host[common.BKHostIDField]
		if bHostIDInInput == true {
			intHostID, err = util.GetInt64ByInterface(hostIDFromInput)
			if err != nil {
				errMsg = append(errMsg, lgc.ccLang.Language("import_host_hostID_not_int"))
				continue
			}
			existInDB = true
		} else {
			// try to get hostID from db
			key := generateHostCloudKey(innerIP, iSubArea)
			intHostID, existInDB = hostIDMap[key]
		}
		var preData mapstr.MapStr
		// remove unchangeable fields
		delete(host, common.BKHostIDField)
		if existInDB {
			// remove unchangeable fields
			delete(host, common.BKHostInnerIPField)

			// get host info before really change it
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
			hostIDMap[generateHostCloudKey(innerIP, iSubArea)] = intHostID
		}
		// add current host operate result to  batch add result
		succMsg = append(succMsg, strconv.FormatInt(index, 10))

		// host info after it changed
		curData, _, err := lgc.GetHostInstanceDetails(ctx, ownerID, strconv.FormatInt(intHostID, 10))
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}

		// add audit log
		logConents = append(logConents, metadata.SaveAuditLogParams{
			ID: intHostID,
			Content: metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: auditHeaders,
			},
			ExtKey: innerIP,
			OpType: auditoplog.AuditOpTypeAdd,
			BizID:  appID,
			OpDesc: "import host",
			Model:  common.BKInnerObjIDHost,
		})
		hostIDs = append(hostIDs, intHostID)
	}

	if len(logConents) > 0 {
		_, err := lgc.CoreAPI.CoreService().Audit().SaveAuditLog(context.Background(), lgc.header, logConents...)
		if err != nil {
			return hostIDs, succMsg, updateErrMsg, errMsg, fmt.Errorf("generate audit log, but get host instance defail failed, err: %v", err)
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return hostIDs, succMsg, updateErrMsg, errMsg, errors.New(lgc.ccLang.Language("host_import_err"))
	}

	return hostIDs, succMsg, updateErrMsg, errMsg, nil
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

// generateHostCloudKey generate a cloudKey for host that is unique among clouds by appending the cloudID.
func generateHostCloudKey(ip, cloudID interface{}) string {
	return fmt.Sprintf("%v-%v", ip, cloudID)
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

	// 更新主机数据
	input := &metadata.UpdateOption{}
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

	// (h.ctx, h.pheader, host)
	result, err := h.CoreAPI.CoreService().Instance().CreateInstance(h.ctx, h.pheader, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Errorf("addHostInstance http do error,err:%s, input:%+v,rid:%s", err.Error(), host, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, err.Error()))
	}
	if !result.Result {
		blog.Errorf("addHostInstance http response error,err code:%d,err msg:%s, input:%+v,rid:%s", result.Code, result.ErrMsg, host, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, result.ErrMsg))
	}

	hostID := int64(result.Data.Created.ID)

	opt := &metadata.HostsModuleRelation{
		ApplicationID: appID,
		ModuleID:      moduleID,
		HostID:        []int64{hostID},
	}
	hResult, err := h.CoreAPI.CoreService().Host().TransferToNormalModule(h.ctx, h.pheader, opt)
	if err != nil {
		blog.Errorf("add host module by ip:%s  err:%s,input:%+v,rid:%s", ip, err.Error(), opt, h.rid)
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, err.Error()))
	}
	if !hResult.Result {
		blog.Errorf("add host module by ip:%s , result:%#v, input:%#v, rid:%s", ip, hResult.Code, opt, h.rid)
		if len(hResult.Data) > 0 {
			return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, hResult.Data[0].Message))
		}
		return 0, fmt.Errorf(h.ccLang.Languagef("host_import_add_fail", index, ip, hResult.ErrMsg))
	}

	return hostID, nil
}

// GetHostIDByHostInfoArr get host id map at it best, ip not exist in db will be ignored.
func (h *importInstance) GetHostIDByHostInfoArr(ctx context.Context, hostInfos map[int64]map[string]interface{}) (map[string]int64, error) {
	// TODO why don't it just return a data structure of cloudKey: hostID map ?
	// step1. extract all innerIP from hostInfos
	var ipArr []string
	for _, host := range hostInfos {
		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk && "" != innerIP {
			ipArr = append(ipArr, innerIP)
		}
	}

	// step2. query host info by innerIPs
	var conds map[string]interface{}
	if 0 < len(ipArr) {
		conds = map[string]interface{}{common.BKHostInnerIPField: common.KvMap{common.BKDBIN: ipArr}}

	}
	query := &metadata.QueryCondition{
		Condition: conds,
		Limit: metadata.SearchLimit{
			Offset: 0,
			Limit:  common.BKNoLimit,
		},
	}
	hResult, err := h.CoreAPI.CoreService().Instance().ReadInstance(ctx, h.pheader, common.BKInnerObjIDHost, query)
	if err != nil {
		blog.Errorf("GetHostIDByHostInfoArr ReadInstance http do err. error:%s, input:%#v, rid:%s", err.Error(), query, h.rid)
		return nil, h.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hResult.Result {
		blog.Errorf("GetHostIDByHostInfoArr ReadInstance http reply err. reply:%#v, input:%#v, rid:%s", hResult, query, h.rid)
		return nil, h.ccErr.New(hResult.Code, hResult.ErrMsg)
	}

	// step3. arrange data as a map, cloudKey: hostID
	hostMap := make(map[string]int64, 0)
	for _, host := range hResult.Data.Info {
		key := generateHostCloudKey(host[common.BKHostInnerIPField], host[common.BKCloudIDField])
		hostID, err := host.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("GetHostIDByHostInfoArr get hostID error. err:%s, hostInfo:%#v, rid:%s", err.Error(), host, h.rid)
			// convert %s  field %s to %s error %s
			return hostMap, h.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
		}
		hostMap[key] = hostID
	}

	return hostMap, nil
}
