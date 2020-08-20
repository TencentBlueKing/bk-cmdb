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
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// audit provides common methods for all audit log utilities
type audit struct {
	clientSet apimachinery.ClientSetInterface
}

func (a *audit) SaveAuditLog(kit *rest.Kit, logs ...metadata.AuditLog) errors.CCErrorCoder {
	_, err := a.clientSet.CoreService().Audit().SaveAuditLog(kit.Ctx, kit.Header, logs...)
	return err
}

func (a *audit) getInstByCond(kit *rest.Kit, objID string, condition map[string]interface{}, fields []string) (
	[]mapstr.MapStr, errors.CCErrorCoder) {

	switch objID {
	case common.BKInnerObjIDHost:
		input := &metadata.QueryInput{Condition: condition}
		if len(fields) != 0 {
			input.Fields = strings.Join(fields, ",")
		}

		rsp, err := a.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, input)
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

		rsp, err := a.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
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

func (a *audit) getInstNameByID(kit *rest.Kit, objID string, instID int64) (string, errors.CCErrorCoder) {

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
