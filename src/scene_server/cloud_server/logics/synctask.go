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
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

func (lgc *Logics) SearchVpc(kit *rest.Kit, accountID int64, vpcOpt *metadata.SearchVpcOption) (*metadata.VpcHostCntResult, error) {
	accountConf, err := lgc.GetCloudAccountConf(accountID)
	if err != nil {
		blog.Errorf("SearchVpc failed, rid:%s, accountID:%d, vpcOpt:%+v, err:%+v", kit.Rid, accountID, vpcOpt, err)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	result, err := lgc.GetVpcHostCnt(*accountConf, vpcOpt.Region)
	if err != nil {
		blog.Errorf("SearchVpc failed, rid:%s, accountConf:%+v, vpcOpt:%+v, err:%+v", kit.Rid, accountConf, vpcOpt, err)
		return nil, kit.CCError.CCError(common.CCErrCloudVpcGetFail)
	}
	return result, nil
}

func (lgc *Logics) CreateSyncTask(kit *rest.Kit, task *metadata.CloudSyncTask) (*metadata.CloudSyncTask, error) {
	result, err := lgc.CoreAPI.CoreService().Cloud().CreateSyncTask(kit.Ctx, kit.Header, task)
	if err != nil {
		blog.Errorf("CreateSyncTask failed, rid:%s, task:%+v, err:%+v", kit.Rid, task, err)
		return nil, err
	}

	// add auditLog
	auditLog := lgc.NewSyncTaskAuditLog(kit, kit.SupplierAccount)
	if err := auditLog.WithCurrent(kit, result.TaskID); err != nil {
		blog.Errorf("CreateSyncTask failed, rid:%s, task:%+v, err:%+v", kit.Rid, task, err)
		return nil, err
	}
	if err := auditLog.SaveAuditLog(kit, metadata.AuditCreate); err != nil {
		blog.Errorf("CreateSyncTask failed, rid:%s, task:%+v, err:%+v", kit.Rid, task, err)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) SearchSyncTask(kit *rest.Kit, option *metadata.SearchCloudOption) (*metadata.MultipleCloudSyncTask, error) {
	// set default limit
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}
	if option.Page.IsIllegal() {
		blog.Errorf("SearchSyncTask failed, Page is IsIllegal, rid:%s, page:%+v", kit.Rid, option.Page)
		return nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	// set default sort
	if option.Page.Sort == "" {
		option.Page.Sort = "-" + common.CreateTimeField
	}

	// if not exact search, change the string query to regexp
	if option.Exact != true {
		for k, v := range option.Condition {
			if reflect.TypeOf(v).Kind() == reflect.String {
				field := v.(string)
				option.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	result, err := lgc.CoreAPI.CoreService().Cloud().SearchSyncTask(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("SearchSyncTask failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) UpdateSyncTask(kit *rest.Kit, taskID int64, option map[string]interface{}) error {

	// add auditLog preData
	auditLog := lgc.NewSyncTaskAuditLog(kit, kit.SupplierAccount)
	if err := auditLog.WithPrevious(kit, taskID); err != nil {
		blog.Errorf("UpdateSyncTask failed, rid:%s, taskID:%d, option:%+v, err:%+v", kit.Rid, taskID, option, err)
		return err
	}

	err := lgc.CoreAPI.CoreService().Cloud().UpdateSyncTask(kit.Ctx, kit.Header, taskID, option)
	if err != nil {
		blog.Errorf("UpdateSyncTask failed, rid:%s, taskID:%d, option:%+v, err:%+v", kit.Rid, taskID, option, err)
		return err
	}

	// add auditLog
	if err := auditLog.WithCurrent(kit, taskID); err != nil {
		blog.Errorf("UpdateSyncTask failed, rid:%s, taskID:%d, option:%+v, err:%+v", kit.Rid, taskID, option, err)
		return err
	}
	if err := auditLog.SaveAuditLog(kit, metadata.AuditUpdate); err != nil {
		blog.Errorf("UpdateSyncTask failed, rid:%s, taskID:%d, option:%+v, err:%+v", kit.Rid, taskID, option, err)
		return err
	}

	return nil
}

func (lgc *Logics) DeleteSyncTask(kit *rest.Kit, taskID int64) error {
	// add auditLog preData
	auditLog := lgc.NewSyncTaskAuditLog(kit, kit.SupplierAccount)
	if err := auditLog.WithPrevious(kit, taskID); err != nil {
		blog.Errorf("DeleteSyncTask failed, rid:%s, taskID:%d, err:%+v", kit.Rid, taskID, err)
		return err
	}

	err := lgc.CoreAPI.CoreService().Cloud().DeleteSyncTask(kit.Ctx, kit.Header, taskID)
	if err != nil {
		blog.Errorf("DeleteSyncTask failed, rid:%s, taskID:%d, err:%+v", kit.Rid, taskID, err)
		return err
	}

	if err := auditLog.SaveAuditLog(kit, metadata.AuditDelete); err != nil {
		blog.Errorf("DeleteSyncTask failed, rid:%s, taskID:%d, err:%+v", kit.Rid, taskID, err)
		return err
	}

	return nil
}

func (lgc *Logics) CreateSyncHistory(kit *rest.Kit, history *metadata.SyncHistory) (*metadata.SyncHistory, error) {
	result, err := lgc.CoreAPI.CoreService().Cloud().CreateSyncHistory(kit.Ctx, kit.Header, history)
	if err != nil {
		blog.Errorf("CreateSyncHistory failed, rid:%s, history:%+v, err:%+v", kit.Rid, history, err)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) SearchSyncHistory(kit *rest.Kit, option *metadata.SearchSyncHistoryOption) (*metadata.MultipleSyncHistory, error) {
	// set default limit
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}
	if option.Page.IsIllegal() {
		blog.Errorf("SearchSyncHistory failed, Page is IsIllegal, rid:%s, page:%+v", kit.Rid, option.Page)
		return nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	// set default sort
	if option.Page.Sort == "" {
		option.Page.Sort = "-" + common.CreateTimeField
	}

	// if not exact search, change the string query to regexp
	if option.Exact != true {
		for k, v := range option.Condition {
			if reflect.TypeOf(v).Kind() == reflect.String {
				field := v.(string)
				option.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	result, err := lgc.CoreAPI.CoreService().Cloud().SearchSyncHistory(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("SearchSyncHistory failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) SearchSyncRegion(kit *rest.Kit, option *metadata.SearchSyncRegionOption) ([]metadata.SyncRegion, error) {
	accountConf, err := lgc.GetCloudAccountConf(option.AccountID)
	if err != nil {
		blog.Errorf("SearchSyncRegion failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result, err := lgc.GetRegionsInfo(*accountConf, option.WithHostCount)
	if err != nil {
		blog.Errorf("SearchSyncRegion failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, kit.CCError.CCError(common.CCErrCloudRegionGetFail)
	}

	return result, nil
}
