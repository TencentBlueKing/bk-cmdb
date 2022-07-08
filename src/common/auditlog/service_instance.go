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
	"sync"

	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// svcInstLog provides methods to generate and save svcInst audit log.
type svcInstLog struct {
	audit
	*svcInstLogGenerator
}

// svcInstLogGenerator provides methods to generate svcInst audit log.
type svcInstLogGenerator struct {
	data         []metadata.ServiceInstance
	procAuditMap map[int64][]metadata.SvcInstProOpDetail
	sync.Mutex
}

// WithServiceInstance set service instances data in audit log
func (s *svcInstLogGenerator) WithServiceInstance(data []metadata.ServiceInstance) {
	s.data = data
}

// WithServiceInstanceByIDs set service instances data in audit log by ids, get data from core service
func (s *svcInstLog) WithServiceInstanceByIDs(kit *rest.Kit, bizID int64, svcInstIDs []int64,
	fields []string) errors.CCErrorCoder {

	data, err := s.GetSvcInstByIDs(kit, bizID, svcInstIDs, fields)
	if err != nil {
		return err
	}
	s.WithServiceInstance(data)
	return nil
}

// WithProc set process data audit detail for service instance audit log
func (s *svcInstLogGenerator) WithProc(param *generateAuditCommonParameter, processes []mapstr.MapStr,
	relations []metadata.ProcessInstanceRelation) errors.CCErrorCoder {

	if len(processes) == 0 || len(relations) == 0 {
		return nil
	}

	procMap := make(map[int64]mapstr.MapStr)
	for _, proc := range processes {
		procID, err := util.GetInt64ByInterface(proc[common.BKProcessIDField])
		if err != nil {
			blog.Errorf("get process(%+v) id failed, err: %v, rid: %s", proc, err, param.kit.Rid)
			return param.kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
		}
		procMap[procID] = proc
	}

	for _, relation := range relations {
		proc, exists := procMap[relation.ProcessID]
		if !exists {
			continue
		}

		details := metadata.BasicOpDetail{}

		switch param.action {
		case metadata.AuditCreate:
			details.Details = &metadata.BasicContent{CurData: proc}
		case metadata.AuditDelete:
			details.Details = &metadata.BasicContent{PreData: proc}
		case metadata.AuditUpdate:
			details.Details = &metadata.BasicContent{PreData: proc, UpdateFields: param.updateFields}
		}

		s.Lock()
		s.procAuditMap[relation.ServiceInstanceID] = append(s.procAuditMap[relation.ServiceInstanceID],
			metadata.SvcInstProOpDetail{
				Action:        param.action,
				ProcessIDs:    relation.ProcessID,
				ProcessNames:  util.GetStrByInterface(proc[common.BKProcessNameField]),
				BasicOpDetail: details,
			})
		s.Unlock()
	}
	return nil
}

// WithProcByRelations set process data in audit log by process relations, get proc data from core service
func (s *svcInstLog) WithProcByRelations(param *generateAuditCommonParameter,
	relations []metadata.ProcessInstanceRelation, fields []string) errors.CCErrorCoder {

	kit := param.kit

	if len(relations) == 0 {
		return nil
	}

	processIDs := make([]int64, 0)
	for _, relation := range relations {
		processIDs = append(processIDs, relation.ProcessID)
	}

	// we need to search id, name and host id field for even if it is not updated
	if len(fields) > 0 {
		fields = append(fields, common.BKFieldID, common.BKFieldName, common.BKHostIDField)
	}

	reqParam := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{
				common.BKDBIN: processIDs,
			},
		},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Fields: fields,
	}
	processRes, rawErr := s.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc, reqParam)
	if rawErr != nil {
		blog.Errorf("get processes failed, request: %+v, err: %v, rid: %s", reqParam, rawErr, kit.Rid)
		return kit.CCError.CCError(common.CCErrProcGetProcessInstanceFailed)
	}
	return s.WithProc(param, processRes.Info, relations)
}

// WithProcBySvcInstIDs set process data in audit log by service instance ids, get data from core service
func (s *svcInstLog) WithProcBySvcInstIDs(param *generateAuditCommonParameter, bizID int64, svcInstIDs []int64,
	fields []string) errors.CCErrorCoder {

	kit := param.kit

	if len(svcInstIDs) == 0 {
		blog.Errorf("get process relations by empty service instance ids for audit log, rid: %s", kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "service_instance_ids")

	}

	// get process instance relations in service instances
	relOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: svcInstIDs,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := s.clientSet.Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, relOpt)
	if err != nil {
		blog.Errorf("get process relations failed, option: %+v, err: %v, rid: %s", relOpt, err, kit.Rid)
		return err
	}

	return s.WithProcByRelations(param, relRes.Info, fields)
}

// GenerateAuditLog generate audit log of service instances.
func (s *svcInstLogGenerator) GenerateAuditLog(param *generateAuditCommonParameter) []metadata.AuditLog {
	auditLogs := make([]metadata.AuditLog, len(s.data))

	for index, svcInst := range s.data {
		action := param.action

		details := &metadata.ServiceInstanceOpDetail{
			HostID: svcInst.HostID,
		}
		if processes, exists := s.procAuditMap[svcInst.ID]; exists {
			details.Processes = processes
		}

		svcInstMapStr := mapstr.SetValueToMapStrByTags(svcInst)

		switch action {
		case metadata.AuditCreate:
			details.BasicOpDetail = metadata.BasicOpDetail{Details: &metadata.BasicContent{CurData: svcInstMapStr}}
		case metadata.AuditDelete:
			details.BasicOpDetail = metadata.BasicOpDetail{Details: &metadata.BasicContent{PreData: svcInstMapStr}}
		case metadata.AuditUpdate:
			if len(param.updateFields) == 0 && len(details.Processes) == 0 {
				continue
			}

			if len(param.updateFields) > 0 {
				preData := make(mapstr.MapStr)
				for key := range param.updateFields {
					preData[key] = svcInstMapStr[key]
				}

				details.BasicOpDetail = metadata.BasicOpDetail{
					Details: &metadata.BasicContent{
						PreData:      preData,
						UpdateFields: param.updateFields,
					},
				}
			}
		}

		auditLog := metadata.AuditLog{
			AuditType:       metadata.BusinessResourceType,
			ResourceType:    metadata.ServiceInstanceRes,
			Action:          action,
			BusinessID:      svcInst.BizID,
			ResourceID:      svcInst.ID,
			OperateFrom:     param.operateFrom,
			ResourceName:    svcInst.Name,
			OperationDetail: details,
		}
		auditLogs[index] = auditLog
	}

	return auditLogs
}

// GetSvcInstByIDs get service instances data in audit log by ids
func (s svcInstLog) GetSvcInstByIDs(kit *rest.Kit, bizID int64, svcInstIDs []int64, fields []string) (
	[]metadata.ServiceInstance, errors.CCErrorCoder) {

	if len(svcInstIDs) == 0 {
		blog.Errorf("get service instance data by empty ids for audit log, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "service_instance_ids")
	}

	// we need to search id, name, biz and host id field for display even if it is not updated
	if len(fields) > 0 {
		fields = append(fields, common.BKFieldID, common.BKFieldName, common.BKHostIDField, common.BKAppIDField)
	}

	option := &metadata.ListServiceInstanceOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: svcInstIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: fields,
	}

	resp, err := s.clientSet.Process().ListServiceInstance(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("get service instance by ids(%+v) failed, err: %v, rid: %s", svcInstIDs, err, kit.Rid)
		return nil, err
	}

	return resp.Info, nil
}

// NewSvcInstAudit create a new service instance audit log operator
func NewSvcInstAudit(clientSet coreservice.CoreServiceClientInterface) *svcInstLog {
	return &svcInstLog{
		audit: audit{
			clientSet: clientSet,
		},
		svcInstLogGenerator: NewSvcInstAuditGenerator(),
	}
}

// NewSvcInstAuditGenerator create a new service instance audit log generator
func NewSvcInstAuditGenerator() *svcInstLogGenerator {
	return &svcInstLogGenerator{
		data:         make([]metadata.ServiceInstance, 0),
		procAuditMap: make(map[int64][]metadata.SvcInstProOpDetail),
	}
}
