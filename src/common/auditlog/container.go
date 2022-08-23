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

package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// kubeAuditLog provides methods to generate and save kube audit log.
type kubeAuditLog struct {
	audit
}

// mockKube mocked kube related struct, used as a placeholder because the actual structs are not defined yet
// TODO replace it with the actual struct for each generate audit log function, remove it after all is done
type mockKube struct {
	ID    *int64
	Name  *string
	BizID *int64
}

// GenerateClusterAuditLog generate audit log of kube cluster.
func (c *kubeAuditLog) GenerateClusterAuditLog(param *generateAuditCommonParameter, data []mockKube) (
	[]metadata.AuditLog, errors.CCErrorCoder) {

	auditLogs := make([]metadata.AuditLog, len(data))

	for index, d := range data {
		log, err := c.generateAuditLog(param, metadata.KubeCluster, d.ID, d.BizID, d.Name, d)
		if err != nil {
			return nil, err
		}
		auditLogs[index] = log
	}

	return auditLogs, nil
}

// GenerateNodeAuditLog generate audit log of kube node.
func (c *kubeAuditLog) GenerateNodeAuditLog(param *generateAuditCommonParameter, data []mockKube) (
	[]metadata.AuditLog, errors.CCErrorCoder) {

	auditLogs := make([]metadata.AuditLog, len(data))

	for index, d := range data {
		auditLog, err := c.generateAuditLog(param, metadata.KubeNode, d.ID, d.BizID, d.Name, d)
		if err != nil {
			return nil, err
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

// GenerateNamespaceAuditLog generate audit log of kube namespace.
func (c *kubeAuditLog) GenerateNamespaceAuditLog(param *generateAuditCommonParameter, data []mockKube) (
	[]metadata.AuditLog, errors.CCErrorCoder) {

	auditLogs := make([]metadata.AuditLog, len(data))

	for index, d := range data {
		auditLog, err := c.generateAuditLog(param, metadata.KubeNamespace, d.ID, d.BizID, d.Name, d)
		if err != nil {
			return nil, err
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

// GeneratePodAuditLog generate audit log of kube pod.
func (c *kubeAuditLog) GeneratePodAuditLog(param *generateAuditCommonParameter, data []mockKube) (
	[]metadata.AuditLog, errors.CCErrorCoder) {

	auditLogs := make([]metadata.AuditLog, len(data))

	for index, d := range data {
		auditLog, err := c.generateAuditLog(param, metadata.KubePod, d.ID, d.BizID, d.Name, d)
		if err != nil {
			return nil, err
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

// GenerateWorkloadAuditLog generate audit log of kube workload.
func (c *kubeAuditLog) GenerateWorkloadAuditLog(param *generateAuditCommonParameter, data []mockKube) (
	[]metadata.AuditLog, errors.CCErrorCoder) {

	auditLogs := make([]metadata.AuditLog, len(data))

	for index, d := range data {
		auditLog, err := c.generateAuditLog(param, metadata.KubeWorkload, d.ID, d.BizID, d.Name, d)
		if err != nil {
			return nil, err
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

// kubeWorkloadData kube workload audit data struct, including workload type and its actual data
type kubeWorkloadData struct {
	Kind string      `json:"kind" bson:"kind"`
	Data interface{} `json:"data" bson:"data"`
}

// GenerateDeploymentAuditLog generate audit log of kube deployment workload.
// TODO add other workload type generate functions like this
func (c *kubeAuditLog) GenerateDeploymentAuditLog(param *generateAuditCommonParameter, data []mockKube) (
	[]metadata.AuditLog, errors.CCErrorCoder) {

	auditLogs := make([]metadata.AuditLog, len(data))

	for index, d := range data {
		wl := &kubeWorkloadData{
			// TODO replace with actual kind definition
			Kind: "deployment",
			Data: d,
		}
		auditLog, err := c.generateAuditLog(param, metadata.KubeWorkload, d.ID, d.BizID, d.Name, wl)
		if err != nil {
			return nil, err
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

func (c *kubeAuditLog) generateAuditLog(param *generateAuditCommonParameter, typ metadata.ResourceType,
	id, bizID *int64, name *string, data interface{}) (metadata.AuditLog, errors.CCErrorCoder) {

	if id == nil || bizID == nil || name == nil || data == nil {
		return metadata.AuditLog{}, param.kit.CCError.CCError(common.CCErrAuditGenerateLogFailed)
	}

	details := &metadata.GenericOpDetail{
		Data: data,
	}

	switch param.action {
	case metadata.AuditUpdate:
		details.UpdateFields = param.updateFields
	}

	return metadata.AuditLog{
		AuditType:       metadata.KubeType,
		ResourceType:    typ,
		Action:          param.action,
		BusinessID:      *bizID,
		ResourceID:      *id,
		OperateFrom:     param.operateFrom,
		ResourceName:    *name,
		OperationDetail: details,
	}, nil
}

// NewKubeAudit new kube audit log utility struct
func NewKubeAudit(clientSet coreservice.CoreServiceClientInterface) *kubeAuditLog {
	return &kubeAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
