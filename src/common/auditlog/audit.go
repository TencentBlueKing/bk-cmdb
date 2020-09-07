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

package auditlog

import (
	"fmt"
	"strings"

	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// audit provides common methods for all audit log utilities
type audit struct {
	clientSet coreservice.CoreServiceClientInterface
}

func (a *audit) SaveAuditLog(kit *rest.Kit, logs ...metadata.AuditLog) error {
	_, err := a.clientSet.Audit().SaveAuditLog(kit.Ctx, kit.Header, logs...)
	return err
}

func (a *audit) getInstByCond(kit *rest.Kit, objID string, condition map[string]interface{}, fields []string) (
	[]mapstr.MapStr, error) {

	switch objID {
	case common.BKInnerObjIDHost:
		input := &metadata.QueryInput{Condition: condition}
		if len(fields) != 0 {
			input.Fields = strings.Join(fields, ",")
		}

		rsp, err := a.clientSet.Host().GetHosts(kit.Ctx, kit.Header, input)
		if err != nil {
			blog.ErrorfDepthf(1, "get host by cond %+v failed, err: %v, rid: %s", condition, err, kit.Rid)
			return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
		}

		if !rsp.Result {
			blog.ErrorfDepthf(1, "get host by cond %+v failed, err: %s, rid: %s", condition, rsp.ErrMsg, kit.Rid)
			return nil, rsp.CCError()
		}

		return rsp.Data.Info, nil
	default:
		input := &metadata.QueryCondition{Condition: condition}
		if len(fields) != 0 {
			input.Fields = fields
		}

		rsp, err := a.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.ErrorfDepthf(1, "get host by cond %+v failed, err: %v, rid: %s", condition, err, kit.Rid)
			return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
		}

		if !rsp.Result {
			blog.ErrorfDepthf(1, "get host by cond %+v failed, err: %s, rid: %s", condition, rsp.ErrMsg, kit.Rid)
			return nil, rsp.CCError()
		}

		return rsp.Data.Info, nil
	}
}

func (a *audit) getInstNameByID(kit *rest.Kit, objID string, instID int64) (string, error) {
	instIDField := common.GetInstIDField(objID)
	instNameField := common.GetInstNameField(objID)

	insts, err := a.getInstByCond(kit, objID, map[string]interface{}{instIDField: instID}, []string{instNameField})
	if err != nil {
		blog.ErrorfDepthf(1, "GetInstNameByID %d GetHosts failed, err: %v, rid: %s", instID, err, kit.Rid)
		return "", err
	}

	if len(insts) != 1 {
		blog.ErrorfDepthf(1, "GetInstNameByID %d GetHosts find %d insts, rid: %s", instID, len(insts), kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, instIDField)
	}

	instName, convErr := insts[0].String(instNameField)
	if convErr != nil {
		blog.ErrorfDepthf(1, "GetInstNameByID %d ReadInstance parse inst name failed, data: %+v, rid: %s", instID, insts[0], kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, objID, instNameField, "string", convErr.Error())
	}

	return instName, nil
}

// getBizIDByHostID get the bizID for host belong business.
func (a *audit) getBizIDByHostID(kit *rest.Kit, hostID int64) (int64, error) {
	input := &metadata.HostModuleRelationRequest{HostIDArr: []int64{hostID}, Fields: []string{common.BKAppIDField}}
	moduleHost, err := a.clientSet.Host().GetHostModuleRelation(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("", hostID, err, kit.Rid)
		return 0, err
	}
	if !moduleHost.Result {
		blog.Errorf("", hostID, moduleHost.ErrMsg, kit.Rid)
		return 0, fmt.Errorf("snapshot get moduleHostConfig failed, fail to create auditLog")
	}

	var bizID int64
	if len(moduleHost.Data.Info) > 0 {
		bizID = moduleHost.Data.Info[0].AppID
	}

	return bizID, nil
}

// getHostInstanceDetailByHostID get host data and hostIP by hostID.
func (a *audit) getHostInstanceDetailByHostID(kit *rest.Kit, hostID int64) (map[string]interface{}, string, error) {
	// get host details.
	result, err := a.clientSet.Host().GetHostByID(kit.Ctx, kit.Header, hostID)
	if err != nil {
		blog.Errorf("getHostInstanceDetailByHostID http do error, err: %s, input: %+v, rid: %s", err.Error(), hostID, kit.Rid)
		return nil, "", kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getHostInstanceDetailByHostID http response error, err code: %d, err msg: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, hostID, kit.Rid)
		return nil, "", kit.CCError.New(result.Code, result.ErrMsg)
	}

	hostInfo := result.Data
	if len(hostInfo) == 0 {
		return nil, "", nil
	}
	ip, ok := hostInfo[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("getHostInstanceDetailByHostID http response format error, convert bk_host_innerip to string error, inst: %#v  input:% #v, rid: %s",
			hostInfo, hostID, kit.Rid)
		return nil, "", kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostInnerIPField, "string", "not string")

	}

	return hostInfo, ip, nil
}

// getObjNameByObjID get the object name corresponding to object id.
func (a *audit) getObjNameByObjID(kit *rest.Kit, objID string) (string, error) {
	query := mapstr.MapStr{common.BKObjIDField: objID}
	// get objectName.
	resp, err := a.clientSet.Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		return "", err
	}
	if resp.Result != true {
		return "", kit.CCError.New(resp.Code, resp.ErrMsg)
	}
	if len(resp.Data.Info) <= 0 {
		return "", kit.CCError.CCError(common.CCErrorModelNotFound)
	}

	return resp.Data.Info[0].Spec.ObjectName, nil
}
