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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type CloudAreaAuditLog struct {
	kit            *rest.Kit
	logic          *Logics
	header         http.Header
	ownerID        string
	MultiCloudArea map[int64]*SingleCloudArea
}

type SingleCloudArea struct {
	CloudName string
	PreData   map[string]interface{}
	CurData   map[string]interface{}
}

func (lgc *Logics) NewCloudAreaLog(kit *rest.Kit) *CloudAreaAuditLog {
	return &CloudAreaAuditLog{
		kit:            kit,
		logic:          lgc,
		header:         kit.Header,
		ownerID:        kit.SupplierAccount,
		MultiCloudArea: make(map[int64]*SingleCloudArea),
	}
}

func (c *CloudAreaAuditLog) WithPrevious(ctx context.Context, platIDs ...int64) errors.CCError {
	return c.buildAuditLogData(ctx, true, false, platIDs...)
}

func (c *CloudAreaAuditLog) WithCurrent(ctx context.Context, platIDs ...int64) errors.CCError {
	return c.buildAuditLogData(ctx, false, true, platIDs...)
}

func (c *CloudAreaAuditLog) buildAuditLogData(ctx context.Context, withPrevious, withCurrent bool, platIDs ...int64) errors.CCError {
	var err error
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: platIDs}},
	}

	res, err := c.logic.CoreAPI.CoreService().Instance().ReadInstance(ctx, c.header, common.BKInnerObjIDPlat, query)
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

func (c *CloudAreaAuditLog) SaveAuditLog(ctx context.Context, action metadata.ActionType) errors.CCError {
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
						PreData: cloudarea.PreData,
						CurData: cloudarea.CurData,
					},
				},
				ModelID: common.BKInnerObjIDPlat,
			},
		})
	}

	auditResult, err := c.logic.CoreAPI.CoreService().Audit().SaveAuditLog(ctx, c.header, logs...)
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
