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
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type CloudAuditLog interface {
	WithPrevious(*rest.Kit, int64) errors.CCError
	WithCurrent(*rest.Kit, int64) errors.CCError
	SaveAuditLog(*rest.Kit, metadata.ActionType) errors.CCError
}

type SyncTaskAuditLog struct {
	logic    *Logics
	header   http.Header
	ownerID  string
	taskName string
	taskID   int64
	content  metadata.CloudSyncTaskOpContent
}

func (lgc *Logics) NewSyncTaskAuditLog(kit *rest.Kit, ownerID string) *SyncTaskAuditLog {
	return &SyncTaskAuditLog{
		logic:   lgc,
		header:  kit.Header,
		ownerID: ownerID,
	}
}

func (log *SyncTaskAuditLog) WithPrevious(kit *rest.Kit, taskID int64) errors.CCError {
	preData, err := log.buildLogData(kit, taskID)
	if err != nil {
		return err
	}
	log.content.PreData = &preData

	return nil
}

func (log *SyncTaskAuditLog) WithCurrent(kit *rest.Kit, taskID int64) errors.CCError {
	curData, err := log.buildLogData(kit, taskID)
	if err != nil {
		return err
	}
	log.content.CurData = &curData

	return nil
}

func (log *SyncTaskAuditLog) buildLogData(kit *rest.Kit, taskID int64) (metadata.CloudSyncTask, errors.CCError) {
	option := metadata.SearchCloudOption{
		Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID},
	}
	res, err := log.logic.CoreAPI.CoreService().Cloud().SearchSyncTask(kit.Ctx, kit.Header, &option)
	if err != nil {
		return metadata.CloudSyncTask{}, err
	}
	if len(res.Info) <= 0 {
		return metadata.CloudSyncTask{}, kit.CCError.CCErrorf(common.CCErrCloudValidSyncTaskParamFail, common.BKCloudAccountID)
	}

	log.taskID = taskID
	log.taskName = res.Info[0].TaskName

	return res.Info[0], nil
}

func (log *SyncTaskAuditLog) SaveAuditLog(kit *rest.Kit, action metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudSyncTaskRes,
		Action:       action,
		OperationDetail: &metadata.CloudSyncTaskOpDetail{
			TaskName: log.taskName,
			TaskID:   log.taskID,
			Details:  &log.content,
		},
	}

	auditResult, err := log.logic.CoreAPI.CoreService().Audit().SaveAuditLog(kit.Ctx, log.header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog add cloud sync task audit log failed, err: %s, result: %+v,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog add cloud sync task audit log failed, err: %s, result: %s,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

type AccountAuditLog struct {
	logic       *Logics
	header      http.Header
	ownerID     string
	accountName string
	accountID   int64
	Content     *metadata.BasicContent
}

func (lgc *Logics) NewAccountAuditLog(kit *rest.Kit, ownerID string) *AccountAuditLog {
	return &AccountAuditLog{
		logic:   lgc,
		header:  kit.Header,
		ownerID: ownerID,
		Content: new(metadata.BasicContent),
	}
}

func (log *AccountAuditLog) WithPrevious(kit *rest.Kit, accountID int64) errors.CCError {
	data, err := log.buildLogData(kit, accountID)
	if err != nil {
		return err
	}
	log.Content.PreData = data

	return nil
}

func (log *AccountAuditLog) WithCurrent(kit *rest.Kit, accountID int64) errors.CCError {
	data, err := log.buildLogData(kit, accountID)
	if err != nil {
		return err
	}
	log.Content.CurData = data

	return nil
}

func (log *AccountAuditLog) SaveAuditLog(kit *rest.Kit, action metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudAccountRes,
		Action:       action,
		OperationDetail: &metadata.BasicOpDetail{
			Details:      log.Content,
		},
	}

	auditResult, err := log.logic.CoreAPI.CoreService().Audit().SaveAuditLog(kit.Ctx, log.header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog add cloud account audit log failed, err: %s, result: %+v,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog add cloud account audit log failed, err: %s, result: %s,rid:%s", err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (log *AccountAuditLog) buildLogData(kit *rest.Kit, accountID int64) (map[string]interface{}, errors.CCError) {

	cond := metadata.SearchCloudOption{
		Condition: mapstr.MapStr{common.BKCloudAccountID: accountID},
	}
	res, err := log.logic.CoreAPI.CoreService().Cloud().SearchAccount(kit.Ctx, kit.Header, &cond)
	if err != nil {
		return nil, err
	}
	if len(res.Info) <= 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail)
	}
	data := map[string]interface{}{
		common.BKCloudAccountName: res.Info[0].AccountName,
		common.BKCloudVendor:      res.Info[0].CloudVendor,
		common.BKDescriptionField: res.Info[0].Description,
	}

	log.accountName = res.Info[0].AccountName
	log.accountID = accountID

	return data, nil
}

func (lgc *Logics) GetAddHostLog(kit *rest.Kit, curData map[string]interface{}) (*metadata.AuditLog, error) {

	auditLog := metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       metadata.AuditCreate,
		OperateFrom:  metadata.FromCloudSync,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: &metadata.BasicContent{
					PreData:    nil,
					CurData:    curData,
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
					PreData:    preData,
					CurData:    curData,
				},
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}

	return &auditLog, nil
}

type CloudAreaAuditLog struct {
	logic          *Logics
	kit            *rest.Kit
	MultiCloudArea map[int64]*SingleCloudArea
}

type SingleCloudArea struct {
	CloudName string
	PreData   map[string]interface{}
	CurData   map[string]interface{}
}

func (lgc *Logics) NewCloudAreaLog(kit *rest.Kit) *CloudAreaAuditLog {
	return &CloudAreaAuditLog{
		logic:          lgc,
		kit:            kit,
		MultiCloudArea: make(map[int64]*SingleCloudArea),
	}
}

func (c *CloudAreaAuditLog) WithPrevious(platIDs ...int64) errors.CCError {
	return c.buildAuditLogData(true, false, platIDs...)
}

func (c *CloudAreaAuditLog) WithCurrent(platIDs ...int64) errors.CCError {
	return c.buildAuditLogData(false, true, platIDs...)
}

func (c *CloudAreaAuditLog) buildAuditLogData(withPrevious, withCurrent bool, platIDs ...int64) errors.CCError {
	var err error

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: platIDs}},
	}
	res, err := c.logic.CoreAPI.CoreService().Instance().ReadInstance(c.kit.Ctx, c.kit.Header, common.BKInnerObjIDPlat, query)
	if nil != err {
		return err
	}
	if len(res.Data.Info) <= 0 {
		return errors.New(common.CCErrTopoCloudNotFound, "")
	}

	for _, data := range res.Data.Info {
		cloudID, err := data.Int64(common.BKCloudIDField)
		if err != nil {
			return err
		}

		cloudName, err := data.String(common.BKCloudNameField)
		if err != nil {
			return err
		}

		if c.MultiCloudArea[cloudID] == nil {
			c.MultiCloudArea[cloudID] = new(SingleCloudArea)
		}

		c.MultiCloudArea[cloudID].CloudName = cloudName

		if withPrevious {
			c.MultiCloudArea[cloudID].PreData = data
		}

		if withCurrent {
			c.MultiCloudArea[cloudID].CurData = data
		}
	}

	return nil
}

func (c *CloudAreaAuditLog) SaveAuditLog(action metadata.ActionType) errors.CCError {
	logs := make([]metadata.AuditLog, 0)
	for cloudID, cloudarea := range c.MultiCloudArea {
		logs = append(logs, metadata.AuditLog{
			AuditType:    metadata.CloudResourceType,
			ResourceType: metadata.CloudAreaRes,
			Action:       action,
			ResourceID:   cloudID,
			ResourceName: cloudarea.CloudName,
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{

					Details: &metadata.BasicContent{
						PreData:    cloudarea.PreData,
						CurData:    cloudarea.CurData,
					},
				},
				ModelID: common.BKInnerObjIDPlat,
			},
		})
	}

	auditResult, err := c.logic.CoreAPI.CoreService().Audit().SaveAuditLog(c.kit.Ctx, c.kit.Header, logs...)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog add cloud area audit log failed, err: %s, result: %+v,rid:%s", err, auditResult, c.kit.Rid)
		return c.kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog add cloud area audit log failed, err: %s, result: %s,rid:%s", err, auditResult, c.kit.Rid)
		return c.kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}
