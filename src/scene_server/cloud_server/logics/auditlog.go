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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) GetAddHostLog(kit *rest.Kit, curData map[string]interface{}) (*metadata.AuditLog, error) {

	auditLog := metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       metadata.AuditCreate,
		OperateFrom:  metadata.FromCloudSync,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: &metadata.BasicContent{
					PreData: nil,
					CurData: curData,
				},
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}

	return &auditLog, nil
}

// 获取资源池业务ID和名称
func (lgc *Logics) GetDefaultBizIDAndName(kit *rest.Kit) (int64, string, error) {
	condition := mapstr.MapStr{
		common.BKDefaultField: common.DefaultAppFlag,
	}
	cond := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField},
		Condition: condition,
	}
	res, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("GetDefaultBizIDAndName fail,err:%s, cond:%+v", err.Error(), *cond)
		return 0, "", err
	}
	if !res.Result {
		blog.Errorf("GetDefaultBizIDAndName fail,err:%s, cond:%+v", res.ErrMsg, *cond)
		return 0, "", fmt.Errorf("%s", res.ErrMsg)
	}

	if len(res.Data.Info) == 0 {
		blog.Errorf("GetDefaultBizIDAndName fail,err:%s, cond:%+v", "no default biz is found", *cond)
		return 0, "", fmt.Errorf("%s", "no default biz is found")
	}

	bizID, err := res.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("GetDefaultBizIDAndName fail,err:%s, cond:%+v", err.Error(), *cond)
		return 0, "", err
	}

	bizName, err := res.Data.Info[0].String(common.BKAppNameField)
	if err != nil {
		blog.Errorf("GetDefaultBizIDAndName fail,err:%s, cond:%+v", err.Error(), *cond)
		return 0, "", err
	}

	return bizID, bizName, nil
}

// 获取主机ID和内网IP
func getHostIDAndIP(hostInfo map[string]interface{}) (int64, string, error) {
	var hostID int64
	var innerIP string
	if hostIDI, ok := hostInfo[common.BKHostIDField]; ok {
		if hostIDVal, err := strconv.ParseInt(fmt.Sprintf("%v", hostIDI), 10, 64); err == nil {
			hostID = hostIDVal
		}
	}

	if innerIPI, ok := hostInfo[common.BKHostInnerIPField]; ok {
		innerIP = fmt.Sprintf("%s", innerIPI)
	}

	if hostID == 0 {
		blog.Errorf("getHostIDAndIP fail,hostID is 0, hostInfo:%+v", hostInfo)
		return 0, "", fmt.Errorf("%s", "hostID is 0")
	}

	return hostID, innerIP, nil
}

// 获取主机的业务ID和业务Name
func (lgc *Logics) GetBizIDAndName(kit *rest.Kit, hostID int64) (int64, error) {
	input := &metadata.HostModuleRelationRequest{HostIDArr: []int64{hostID}}
	moduleHost, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(context.Background(), kit.Header, input)
	if err != nil {
		blog.Errorf("GetBizIDAndName fail, err:%s, input:%+v", err.Error(), input)
		return 0, err
	}
	if !moduleHost.Result {
		blog.Errorf("GetBizIDAndName fail, err code:%d, err msg:%s, input:%+v", moduleHost.Code, moduleHost.ErrMsg, input)
		return 0, fmt.Errorf("%s", moduleHost.ErrMsg)
	}

	if len(moduleHost.Data.Info) == 0 {
		blog.Errorf("GetBizIDAndName fail, host biz is not found, input:%+v", input)
		return 0, fmt.Errorf("%s", "host biz is not found")
	}

	bizID := moduleHost.Data.Info[0].AppID

	return bizID, nil
}

func (lgc *Logics) GetUpdateHostLog(kit *rest.Kit, preData, curData map[string]interface{}) (*metadata.AuditLog, error) {
	// 获取主机ID和内网IP
	hostID, innerIP, err := getHostIDAndIP(preData)
	if err != nil {
		blog.Errorf("GetUpdateHostLog fail,err:%s, preData:%+v, curData:%+v", err.Error(), preData, curData)
		return nil, err
	}

	// 获取主机的业务ID和业务Name
	bizID, err := lgc.GetBizIDAndName(kit, hostID)
	if err != nil {
		blog.Errorf("GetUpdateHostLog fail,err:%s, preData:%+v, curData:%+v", err.Error(), preData, curData)
		return nil, err
	}

	auditLog := metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       metadata.AuditUpdate,
		OperateFrom:  metadata.FromCloudSync,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: innerIP,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: &metadata.BasicContent{
					PreData: preData,
					CurData: curData,
				},
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}

	return &auditLog, nil
}
