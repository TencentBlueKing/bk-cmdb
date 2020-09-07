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
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type cloudAreaAuditLog struct {
	audit
}

// GenerateAuditLog batch generate audit log for cloud area, auto get cloud area data by cloud ID.
func (h *cloudAreaAuditLog) GenerateAuditLog(kit *rest.Kit, action metadata.ActionType, platIDs []int64,
	OperateFrom metadata.OperateFromType, updateFields map[string]interface{}) ([]metadata.AuditLog, error) {
	if len(platIDs) == 0 {
		return make([]metadata.AuditLog, 0), nil
	}

	// build query condition for read cloud.
	var err error
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: platIDs}},
	}

	// to query plat.
	res, err := h.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDPlat, query)
	if nil != err {
		return nil, err
	}

	// mapping from cloudID to cloudData.
	var mutilCloudArea = make(map[int64]mapstr.MapStr)
	for _, data := range res.Data.Info {
		cloudID, err := data.Int64(common.BKCloudIDField)
		if err != nil {
			return nil, err
		}

		mutilCloudArea[cloudID] = data
	}

	// to generate audit log.
	var logs = make([]metadata.AuditLog, 0)
	for cloudID, cloudData := range mutilCloudArea {
		cloudName, _ := cloudData.String(common.BKCloudNameField)

		var basicDetail *metadata.BasicContent
		switch action {
		case metadata.AuditCreate:
			basicDetail = &metadata.BasicContent{
				CurData: cloudData,
			}
		case metadata.AuditDelete:
			basicDetail = &metadata.BasicContent{
				PreData: cloudData,
			}
		case metadata.AuditUpdate:
			basicDetail = &metadata.BasicContent{
				PreData:      cloudData,
				UpdateFields: updateFields,
			}
		}

		logs = append(logs, metadata.AuditLog{
			AuditType:    metadata.CloudResourceType,
			ResourceType: metadata.CloudAreaRes,
			Action:       action,
			ResourceID:   cloudID,
			ResourceName: cloudName,
			OperateFrom:  OperateFrom,
			OperationDetail: &metadata.BasicOpDetail{
				Details: basicDetail,
			},
		})
	}

	return logs, err
}

func NewCloudAreaAuditLog(clientSet coreservice.CoreServiceClientInterface) *cloudAreaAuditLog {
	return &cloudAreaAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
