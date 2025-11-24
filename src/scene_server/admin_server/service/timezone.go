/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/tools"
	"configcenter/pkg/tenant/types"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful/v3"
)

// TimeZone struct for change timezone
type TimeZone struct {
	BizName  string `json:"bk_biz_name"`
	TimeZone string `json:"time_zone"`
}

// updateBizTimeZone 特殊需求，接口不对外暴露，修改业务对应时区
func (s *Service) updateBizTimeZone(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header

	kit := rest.NewKitFromHeader(rHeader, s.CCErr)
	if err := s.validTenant(kit); err != nil {
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err,
			ErrCode: common.CCErrAPICheckTenantInvalid})
		return
	}

	timeZone := new(TimeZone)
	if err := json.NewDecoder(req.Request.Body).Decode(timeZone); err != nil {
		blog.Errorf("failed with decode body err: %v, body: %s, rid:%s", err, req.Request.Body, kit.Rid)
		_ = resp.WriteError(http.StatusInternalServerError, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	if err := timeZone.validate(kit); err != nil {
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err,
			ErrCode: common.CCErrCommParamsInvalid})
		return
	}

	// update time zone
	cond := map[string]interface{}{
		common.BKAppNameField: timeZone.BizName,
	}

	updateData := map[string]interface{}{
		common.BKTimeZoneField: timeZone.TimeZone,
	}

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(s.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(
		updateData)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, common.BKInnerObjIDApp, cond)
	if ccErr != nil {
		blog.Errorf("update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: ccErr,
			ErrCode: common.CCErrAuditGenerateLogFailed})
		return
	}

	err := s.db.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Update(kit.Ctx, cond, updateData)
	if err != nil {
		blog.Errorf("update time zone failed, err: %v, rid: %s", err, kit.Rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err,
			ErrCode: common.CCErrCommDBUpdateFailed})
		return
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("update inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err,
			ErrCode: common.CCErrAuditSaveLogFailed})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
	return
}

func (s *Service) validTenant(kit *rest.Kit) error {
	// validate tenant mode
	tenantID, err := tools.ValidateDisableTenantMode(kit.TenantID, s.Config.EnableMultiTenantMode)
	if err != nil {
		blog.Errorf("tenant mode is not enabled, but tenant id is set")
		return fmt.Errorf("tenant mode is not enabled, but tenant id is set")
	}

	// check if tenant is validate
	tenantData, exist := tenant.GetTenant(tenantID)
	if !exist || tenantData.Status != types.EnabledStatus {
		blog.Errorf("invalid tenant: %s, rid: %s", tenantID, kit.Rid)
		return fmt.Errorf("invalid tenant: %s", tenantID)
	}

	return nil
}

func (t TimeZone) validate(kit *rest.Kit) error {

	if len(t.BizName) == 0 {
		blog.Errorf("biz name cannot be empty, rid: %s", kit.Rid)
		return fmt.Errorf("biz name cannot be empty")
	}

	if len(t.TimeZone) == 0 {
		blog.Errorf("time zone is empty, rid: %s", kit.Rid)
		return fmt.Errorf("time zone is empty")
	}

	return nil
}
