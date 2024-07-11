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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetInst get instance
func (d *Client) GetInst(kit *rest.Kit, objID string, cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	result, err := d.ApiClient.GetInstDetail(kit.Ctx, kit.Header, objID, cond)
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
func (d *Client) GetBiz(kit *rest.Kit, cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	ownerID := kit.Header.Get(common.BKHTTPOwnerID)
	result, err := d.ApiClient.SearchBiz(kit.Ctx, ownerID, kit.Header, cond)
	if err != nil {
		blog.Errorf("get biz data detail error: %v , search condition: %v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	if len(result.Data.Info) == 0 {
		return nil, nil
	}

	return result.Data.Info, nil
}

// GetProject get project
func (d *Client) GetProject(kit *rest.Kit, cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	result, err := d.ApiClient.SearchProject(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("get project data detail error: %v , search condition: %v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	if len(result.Data.Info) == 0 {
		return nil, nil
	}

	return result.Data.Info, nil
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
		result, err = d.ApiClient.AddInstByImport(kit.Ctx, kit.Header, kit.SupplierAccount,
			param.ObjID, param.Req)
	default:
		err = fmt.Errorf("handle type is invalid, type: %s", param.HandleType)
	}

	if err != nil {
		blog.Errorf("add instance failed, err: %v, rid: %s", err, kit.Rid)
		errMsg := make([]string, 0)
		defLang := param.Language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header))
		for idx := range param.Instances {
			errMsg = append(errMsg, defLang.Languagef("import_data_fail", idx, err.Error()))
		}

		return nil, errMsg
	}

	return result.Data.Success, result.Data.Errors
}
