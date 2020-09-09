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
		blog.Errorf("failed to getting instance name, instID: %d, objID: %d, err: %v, rid: %s",
			instID, objID, err, kit.Rid)
		return "", err
	}

	if len(insts) != 1 {
		blog.Errorf("failed to getting instance name, instance not one, instID: %d, objID: %d, rid: %s", instID, objID, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, instIDField)
	}

	instName, convErr := insts[0].String(instNameField)
	if convErr != nil {
		blog.Errorf("getting instance name failed, failed to %s fields convert to string, instID: %d, data: %+v, rid: %s",
			instNameField, instID, insts[0], kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, objID, instNameField, "string", convErr.Error())
	}

	return instName, nil
}

// getObjNameByObjID get the object name by object id.
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

// getDefaultAppID get default businessID under designated supplier account.
func (a *audit) getDefaultAppID(kit *rest.Kit) (int64, error) {
	cond := mapstr.MapStr{
		common.BKDefaultField: common.DefaultAppFlag,
	}
	fields := []string{common.BKAppIDField, common.BkSupplierAccount}

	results, err := a.getInstByCond(kit, common.BKInnerObjIDApp, cond, fields)
	if err != nil {
		blog.Errorf("get default appID failed, err: %v, rid :%s", err, kit.Rid)
		return 0, err
	}

	for _, data := range results {
		ownID, err := data.String(common.BkSupplierAccount)
		if err != nil {
			return 0, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp,
				common.BkSupplierAccount, "string", err.Error())
		}

		if kit.SupplierAccount == ownID {
			bizID, err := data.Int64(common.BKAppIDField)
			if err != nil {
				return 0, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp,
					common.BKAppIDField, "int64", err.Error())
			}

			return bizID, nil
		}

	}

	blog.Errorf("no such default business when supplier account is %s, rid: %s", kit.SupplierAccount, kit.Rid)
	return 0, fmt.Errorf("no such default business when supplier account is %s", kit.SupplierAccount)
}
