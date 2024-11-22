/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package core

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/pkg/errors"
)

// GetInst get instance
func (d *Client) GetInst(kit *rest.Kit, objID string, cond interface{}) ([]mapstr.MapStr, error) {
	instCond, ok := cond.(mapstr.MapStr)
	if !ok {
		blog.Errorf("get inst but condition parse failed, condition: %v, rid: %s", cond, kit.Rid)
		return nil, errors.New("get inst but condition parse failed")
	}

	result, err := d.ApiClient.GetInstDetail(kit.Ctx, kit.Header, objID, instCond)
	if err != nil {
		blog.Errorf("get inst data detail error: %v , search condition: %#v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	if result.Data.Count == 0 {
		blog.Errorf("get inst data detail, but got 0 instances, condition: %#v, rid: %s", cond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAPINoObjectInstancesIsFound)
	}

	return result.Data.Info, nil
}

// GetBiz get biz
func (d *Client) GetBiz(kit *rest.Kit, cond interface{}) ([]mapstr.MapStr, error) {
	bizCond, ok := cond.(*metadata.QueryBusinessRequest)
	if !ok {
		blog.Errorf("get biz but condition parse failed, condition: %v, rid: %s", cond, kit.Rid)
		return nil, errors.New("get biz but condition parse failed")
	}

	tenantID := httpheader.GetTenantID(kit.Header)
	result, err := d.ApiClient.SearchBiz(kit.Ctx, tenantID, kit.Header, bizCond)
	if err != nil {
		blog.Errorf("get biz data detail error: %v , search condition: %v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	return result.Data.Info, nil
}

// GetProject get project
func (d *Client) GetProject(kit *rest.Kit, cond interface{}) ([]mapstr.MapStr, error) {
	bizCond, ok := cond.(*metadata.SearchProjectOption)
	if !ok {
		blog.Errorf("get project but condition parse failed, condition: %v, rid: %s", cond, kit.Rid)
		return nil, errors.New("get project but condition parse failed")
	}

	result, err := d.ApiClient.SearchProject(kit.Ctx, kit.Header, bizCond)
	if err != nil {
		blog.Errorf("get project data detail error: %v , search condition: %v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	return result.Info, nil
}

// HandleImportedInst handle imported instance
func (d *Client) HandleImportedInst(kit *rest.Kit, param *ImportedParam) ([]int64, []string) {
	var result *metadata.ImportInstResp
	var err error

	switch param.HandleType {
	case AddHost:
		result, err = d.ApiClient.AddHostByExcel(kit.Ctx, kit.Header, param.Req)
	case UpdateHost:
		result, err = d.ApiClient.UpdateHost(kit.Ctx, kit.Header, param.Req)
	case AddInst:
		result, err = d.ApiClient.AddInstByImport(kit.Ctx, kit.Header, kit.TenantID,
			param.ObjID, param.Req)
	default:
		err = fmt.Errorf("handle type is invalid, type: %s", param.HandleType)
	}

	if err != nil {
		blog.Errorf("add instance failed, err: %v, rid: %s", err, kit.Rid)
		errMsg := make([]string, 0)
		defLang := param.Language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))
		for idx := range param.Instances {
			errMsg = append(errMsg, defLang.Languagef("import_data_fail", idx, err.Error()))
		}

		return nil, errMsg
	}

	return result.Data.Success, result.Data.Errors
}
