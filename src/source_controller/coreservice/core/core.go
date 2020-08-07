/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/selector"
)

// ModelAttributeGroup model attribute group methods definitions
type ModelAttributeGroup interface {
	CreateModelAttributeGroup(kit *rest.Kit, objID string, inputParam metadata.CreateModelAttributeGroup) (*metadata.CreateOneDataResult, error)
	SetModelAttributeGroup(kit *rest.Kit, objID string, inputParam metadata.SetModelAttributeGroup) (*metadata.SetDataResult, error)
	UpdateModelAttributeGroup(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	UpdateModelAttributeGroupByCondition(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelAttributeGroup(kit *rest.Kit, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error)
	SearchModelAttributeGroupByCondition(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error)
	DeleteModelAttributeGroup(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	DeleteModelAttributeGroupByCondition(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// ModelClassification model classification methods definitions
type ModelClassification interface {
	CreateOneModelClassification(kit *rest.Kit, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error)
	CreateManyModelClassification(kit *rest.Kit, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error)
	SetManyModelClassification(kit *rest.Kit, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error)
	SetOneModelClassification(kit *rest.Kit, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error)
	UpdateModelClassification(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModelClassification(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModelClassification(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelClassificationDataResult, error)
}

// ModelAttribute model attribute methods definitions
type ModelAttribute interface {
	CreateModelAttributes(kit *rest.Kit, objID string, inputParam metadata.CreateModelAttributes) (*metadata.CreateManyDataResult, error)
	SetModelAttributes(kit *rest.Kit, objID string, inputParam metadata.SetModelAttributes) (*metadata.SetDataResult, error)
	UpdateModelAttributes(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	UpdateModelAttributesIndex(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdateAttrIndexData, error)
	UpdateModelAttributesByCondition(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModelAttributes(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModelAttributes(kit *rest.Kit, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error)
	SearchModelAttributesByCondition(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error)
}

// ModelAttrUnique model attribute  unique methods definitions
type ModelAttrUnique interface {
	CreateModelAttrUnique(kit *rest.Kit, objID string, data metadata.CreateModelAttrUnique) (*metadata.CreateOneDataResult, error)
	UpdateModelAttrUnique(kit *rest.Kit, objID string, id uint64, data metadata.UpdateModelAttrUnique) (*metadata.UpdatedCount, error)
	DeleteModelAttrUnique(kit *rest.Kit, objID string, id uint64, meta metadata.DeleteModelAttrUnique) (*metadata.DeletedCount, error)
	SearchModelAttrUnique(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryUniqueResult, error)
}

// ModelOperation model methods
type ModelOperation interface {
	ModelClassification
	ModelAttributeGroup
	ModelAttribute
	ModelAttrUnique

	CreateModel(kit *rest.Kit, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error)
	SetModel(kit *rest.Kit, inputParam metadata.SetModel) (*metadata.SetDataResult, error)
	UpdateModel(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModel(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModel(kit *rest.Kit, modelID int64) (*metadata.DeletedCount, error)
	SearchModel(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelDataResult, error)
	SearchModelWithAttribute(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelWithAttributeDataResult, error)
}

// InstanceOperation instance methods
type InstanceOperation interface {
	CreateModelInstance(kit *rest.Kit, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error)
	CreateManyModelInstance(kit *rest.Kit, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error)
	UpdateModelInstance(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelInstance(kit *rest.Kit, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelInstance(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModelInstance(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// AssociationKind association kind methods
type AssociationKind interface {
	CreateAssociationKind(kit *rest.Kit, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error)
	CreateManyAssociationKind(kit *rest.Kit, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyDataResult, error)
	SetAssociationKind(kit *rest.Kit, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error)
	SetManyAssociationKind(kit *rest.Kit, inputParam metadata.SetManyAssociationKind) (*metadata.SetDataResult, error)
	UpdateAssociationKind(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteAssociationKind(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteAssociationKind(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchAssociationKind(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.SearchAssociationKindResult, error)
}

// ModelAssociation manager model association
type ModelAssociation interface {
	CreateModelAssociation(kit *rest.Kit, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	CreateMainlineModelAssociation(kit *rest.Kit, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	SetModelAssociation(kit *rest.Kit, inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error)
	UpdateModelAssociation(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelAssociation(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelAssociation(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModelAssociation(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// InstanceAssociation manager instance association
type InstanceAssociation interface {
	CreateOneInstanceAssociation(kit *rest.Kit, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error)
	CreateManyInstanceAssociation(kit *rest.Kit, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error)
	SearchInstanceAssociation(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteInstanceAssociation(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// DataSynchronizeOperation manager data synchronize interface
type DataSynchronizeOperation interface {
	SynchronizeInstanceAdapter(kit *rest.Kit, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	SynchronizeModelAdapter(kit *rest.Kit, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	SynchronizeAssociationAdapter(kit *rest.Kit, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	Find(kit *rest.Kit, find *metadata.SynchronizeFindInfoParameter) ([]mapstr.MapStr, uint64, error)
	ClearData(kit *rest.Kit, input *metadata.SynchronizeClearDataParameter) error
	SetIdentifierFlag(kit *rest.Kit, input *metadata.SetIdenifierFlag) ([]metadata.ExceptionResult, error)
}

// TopoOperation methods
type TopoOperation interface {
	SearchMainlineModelTopo(ctx context.Context, header http.Header, withDetail bool) (*metadata.TopoModelNode, error)
	SearchMainlineInstanceTopo(ctx context.Context, header http.Header, objID int64, withDetail bool) (*metadata.TopoInstanceNode, error)
}

// HostOperation methods
type HostOperation interface {
	TransferToInnerModule(kit *rest.Kit, input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, error)
	TransferToNormalModule(kit *rest.Kit, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error)
	TransferToAnotherBusiness(kit *rest.Kit, input *metadata.TransferHostsCrossBusinessRequest) ([]metadata.ExceptionResult, error)
	RemoveFromModule(kit *rest.Kit, input *metadata.RemoveHostsFromModuleOption) ([]metadata.ExceptionResult, error)
	DeleteFromSystem(kit *rest.Kit, input *metadata.DeleteHostRequest) ([]metadata.ExceptionResult, error)
	GetHostModuleRelation(kit *rest.Kit, input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error)
	Identifier(kit *rest.Kit, input *metadata.SearchHostIdentifierParam) ([]metadata.HostIdentifier, error)
	UpdateHostCloudAreaField(kit *rest.Kit, input metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder

	LockHost(kit *rest.Kit, input *metadata.HostLockRequest) errors.CCError
	UnlockHost(kit *rest.Kit, input *metadata.HostLockRequest) errors.CCError
	QueryHostLock(kit *rest.Kit, input *metadata.QueryHostLockRequest) ([]metadata.HostLockData, errors.CCError)

	// host search
	ListHosts(kit *rest.Kit, input metadata.ListHosts) (*metadata.ListHostResult, error)

	// GetDistinctHostIDsByTopoRelation get all  host ids by topology relation condition
	GetDistinctHostIDsByTopoRelation(kit *rest.Kit, input *metadata.DistinctHostIDByTopoRelationRequest) ([]int64, error)
}

// AssociationOperation association methods
type AssociationOperation interface {
	AssociationKind
	ModelAssociation
	InstanceAssociation
}

type AuditOperation interface {
	CreateAuditLog(kit *rest.Kit, logs ...metadata.AuditLog) error
	SearchAuditLog(kit *rest.Kit, param metadata.QueryInput) ([]metadata.AuditLog, uint64, error)
}

type StatisticOperation interface {
	SearchInstCount(kit *rest.Kit, inputParam map[string]interface{}) (uint64, error)
	SearchChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error)
	SearchOperationChart(kit *rest.Kit, inputParam interface{}) (*metadata.ChartClassification, error)
	CreateOperationChart(kit *rest.Kit, inputParam metadata.ChartConfig) (uint64, error)
	UpdateChartPosition(kit *rest.Kit, inputParam interface{}) (interface{}, error)
	DeleteOperationChart(kit *rest.Kit, id int64) (interface{}, error)
	UpdateOperationChart(kit *rest.Kit, inputParam map[string]interface{}) (interface{}, error)
	SearchTimerChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error)
	TimerFreshData(kit *rest.Kit) error
}

// Core core itnerfaces methods
type Core interface {
	ModelOperation() ModelOperation
	InstanceOperation() InstanceOperation
	AssociationOperation() AssociationOperation
	TopoOperation() TopoOperation
	DataSynchronizeOperation() DataSynchronizeOperation
	HostOperation() HostOperation
	AuditOperation() AuditOperation
	StatisticOperation() StatisticOperation
	ProcessOperation() ProcessOperation
	LabelOperation() LabelOperation
	SetTemplateOperation() SetTemplateOperation
	HostApplyRuleOperation() HostApplyRuleOperation
	SystemOperation() SystemOperation
}

// ProcessOperation methods
type ProcessOperation interface {
	// service category
	CreateServiceCategory(kit *rest.Kit, category metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder)
	GetServiceCategory(kit *rest.Kit, categoryID int64) (*metadata.ServiceCategory, errors.CCErrorCoder)
	GetDefaultServiceCategory(kit *rest.Kit) (*metadata.ServiceCategory, errors.CCErrorCoder)
	UpdateServiceCategory(kit *rest.Kit, categoryID int64, category metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder)
	ListServiceCategories(kit *rest.Kit, bizID int64, withStatistics bool) (*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder)
	DeleteServiceCategory(kit *rest.Kit, categoryID int64) errors.CCErrorCoder
	IsServiceCategoryLeafNode(kit *rest.Kit, categoryID int64) (bool, errors.CCErrorCoder)

	// service template
	CreateServiceTemplate(kit *rest.Kit, template metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder)
	GetServiceTemplate(kit *rest.Kit, templateID int64) (*metadata.ServiceTemplate, errors.CCErrorCoder)
	UpdateServiceTemplate(kit *rest.Kit, templateID int64, template metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder)
	ListServiceTemplates(kit *rest.Kit, option metadata.ListServiceTemplateOption) (*metadata.MultipleServiceTemplate, errors.CCErrorCoder)
	DeleteServiceTemplate(kit *rest.Kit, serviceTemplateID int64) errors.CCErrorCoder

	// process template
	CreateProcessTemplate(kit *rest.Kit, template metadata.ProcessTemplate) (*metadata.ProcessTemplate, errors.CCErrorCoder)
	GetProcessTemplate(kit *rest.Kit, templateID int64) (*metadata.ProcessTemplate, errors.CCErrorCoder)
	UpdateProcessTemplate(kit *rest.Kit, templateID int64, property map[string]interface{}) (*metadata.ProcessTemplate, errors.CCErrorCoder)
	ListProcessTemplates(kit *rest.Kit, option metadata.ListProcessTemplatesOption) (*metadata.MultipleProcessTemplate, errors.CCErrorCoder)
	DeleteProcessTemplate(kit *rest.Kit, processTemplateID int64) errors.CCErrorCoder

	// service instance
	CreateServiceInstance(kit *rest.Kit, template metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder)
	GetServiceInstance(kit *rest.Kit, templateID int64) (*metadata.ServiceInstance, errors.CCErrorCoder)
	UpdateServiceInstance(kit *rest.Kit, instanceID int64, instance metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder)
	ListServiceInstance(kit *rest.Kit, option metadata.ListServiceInstanceOption) (*metadata.MultipleServiceInstance, errors.CCErrorCoder)
	ListServiceInstanceDetail(kit *rest.Kit, option metadata.ListServiceInstanceDetailOption) (*metadata.MultipleServiceInstanceDetail, errors.CCErrorCoder)
	DeleteServiceInstance(kit *rest.Kit, serviceInstanceIDs []int64) errors.CCErrorCoder
	AutoCreateServiceInstanceModuleHost(kit *rest.Kit, hostID int64, moduleID int64) (*metadata.ServiceInstance, errors.CCErrorCoder)
	RemoveTemplateBindingOnModule(kit *rest.Kit, moduleID int64) errors.CCErrorCoder
	ReconstructServiceInstanceName(kit *rest.Kit, instanceID int64) errors.CCErrorCoder

	// process instance relation
	CreateProcessInstanceRelation(kit *rest.Kit, relation *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	GetProcessInstanceRelation(kit *rest.Kit, processInstanceID int64) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	UpdateProcessInstanceRelation(kit *rest.Kit, processInstanceID int64, relation metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	ListProcessInstanceRelation(kit *rest.Kit, option metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder)
	ListHostProcessRelation(kit *rest.Kit, option *metadata.ListProcessInstancesWithHostOption) (*metadata.MultipleHostProcessRelation, errors.CCErrorCoder)
	DeleteProcessInstanceRelation(kit *rest.Kit, option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder

	GetBusinessDefaultSetModuleInfo(kit *rest.Kit, bizID int64) (metadata.BusinessDefaultSetModuleInfo, errors.CCErrorCoder)
	GetProc2Module(kit *rest.Kit, option *metadata.GetProc2ModuleOption) ([]metadata.Proc2Module, errors.CCErrorCoder)
}

type LabelOperation interface {
	AddLabel(kit *rest.Kit, tableName string, option selector.LabelAddOption) errors.CCErrorCoder
	RemoveLabel(kit *rest.Kit, tableName string, option selector.LabelRemoveOption) errors.CCErrorCoder
}

type SetTemplateOperation interface {
	CreateSetTemplate(kit *rest.Kit, bizID int64, option metadata.CreateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder)
	UpdateSetTemplate(kit *rest.Kit, setTemplateID int64, option metadata.UpdateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder)
	DeleteSetTemplate(kit *rest.Kit, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder
	GetSetTemplate(kit *rest.Kit, bizID int64, setTemplateID int64) (metadata.SetTemplate, errors.CCErrorCoder)
	ListSetTemplate(kit *rest.Kit, bizID int64, option metadata.ListSetTemplateOption) (metadata.MultipleSetTemplateResult, errors.CCErrorCoder)
	ListSetServiceTemplateRelations(kit *rest.Kit, bizID int64, setTemplateID int64) ([]metadata.SetServiceTemplateRelation, errors.CCErrorCoder)
	ListSetTplRelatedSvcTpl(kit *rest.Kit, bizID, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder)
	UpdateSetTemplateSyncStatus(kit *rest.Kit, setID int64, option metadata.SetTemplateSyncStatus) errors.CCErrorCoder
	ListSetTemplateSyncStatus(kit *rest.Kit, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
	ListSetTemplateSyncHistory(kit *rest.Kit, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
	DeleteSetTemplateSyncStatus(kit *rest.Kit, option metadata.DeleteSetTemplateSyncStatusOption) errors.CCErrorCoder
	ModifySetTemplateSyncStatus(kit *rest.Kit, setID int64, sysncStatus metadata.SyncStatus) errors.CCErrorCoder
}

type HostApplyRuleOperation interface {
	CreateHostApplyRule(kit *rest.Kit, bizID int64, option metadata.CreateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder)
	UpdateHostApplyRule(kit *rest.Kit, bizID int64, ruleID int64, option metadata.UpdateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder)
	DeleteHostApplyRule(kit *rest.Kit, bizID int64, ruleIDs ...int64) errors.CCErrorCoder
	GetHostApplyRule(kit *rest.Kit, bizID int64, ruleID int64) (metadata.HostApplyRule, errors.CCErrorCoder)
	ListHostApplyRule(kit *rest.Kit, bizID int64, option metadata.ListHostApplyRuleOption) (metadata.MultipleHostApplyRuleResult, errors.CCErrorCoder)
	GenerateApplyPlan(kit *rest.Kit, bizID int64, option metadata.HostApplyPlanOption) (metadata.HostApplyPlanResult, errors.CCErrorCoder)
	SearchRuleRelatedModules(kit *rest.Kit, bizID int64, option metadata.SearchRuleRelatedModulesOption) ([]metadata.Module, errors.CCErrorCoder)
	BatchUpdateHostApplyRule(kit *rest.Kit, bizID int64, option metadata.BatchCreateOrUpdateApplyRuleOption) (metadata.BatchCreateOrUpdateHostApplyRuleResult, errors.CCErrorCoder)
	RunHostApplyOnHosts(kit *rest.Kit, bizID int64, option metadata.UpdateHostByHostApplyRuleOption) (metadata.MultipleHostApplyResult, errors.CCErrorCoder)
}

type SystemOperation interface {
	GetSystemUserConfig(kit *rest.Kit) (map[string]interface{}, errors.CCErrorCoder)
	SearchConfigAdmin(kit *rest.Kit) (*metadata.ConfigAdmin, errors.CCErrorCoder)
}

type core struct {
	model           ModelOperation
	instance        InstanceOperation
	association     AssociationOperation
	dataSynchronize DataSynchronizeOperation
	topo            TopoOperation
	host            HostOperation
	audit           AuditOperation
	operation       StatisticOperation
	process         ProcessOperation
	label           LabelOperation
	sys             SystemOperation
	setTemplate     SetTemplateOperation
	hostApplyRule   HostApplyRuleOperation
}

// New create core
func New(
	model ModelOperation,
	instance InstanceOperation,
	association AssociationOperation,
	dataSynchronize DataSynchronizeOperation,
	topo TopoOperation, host HostOperation,
	audit AuditOperation,
	process ProcessOperation,
	label LabelOperation,
	setTemplate SetTemplateOperation,
	operation StatisticOperation,
	hostApplyRule HostApplyRuleOperation,
	sys SystemOperation,
) Core {
	return &core{
		model:           model,
		instance:        instance,
		association:     association,
		dataSynchronize: dataSynchronize,
		topo:            topo,
		host:            host,
		audit:           audit,
		operation:       operation,
		process:         process,
		label:           label,
		sys:             sys,
		setTemplate:     setTemplate,
		hostApplyRule:   hostApplyRule,
	}
}

func (m *core) ModelOperation() ModelOperation {
	return m.model
}

func (m *core) InstanceOperation() InstanceOperation {
	return m.instance
}

func (m *core) AssociationOperation() AssociationOperation {
	return m.association
}

func (m *core) TopoOperation() TopoOperation {
	return m.topo
}

func (m *core) DataSynchronizeOperation() DataSynchronizeOperation {
	return m.dataSynchronize
}

func (m *core) HostOperation() HostOperation {
	return m.host
}

func (m *core) AuditOperation() AuditOperation {
	return m.audit
}

func (m *core) ProcessOperation() ProcessOperation {
	return m.process
}

func (m *core) StatisticOperation() StatisticOperation {
	return m.operation
}

func (m *core) LabelOperation() LabelOperation {
	return m.label
}

func (m *core) SetTemplateOperation() SetTemplateOperation {
	return m.setTemplate
}

func (m *core) SystemOperation() SystemOperation {
	return m.sys
}

func (m *core) HostApplyRuleOperation() HostApplyRuleOperation {
	return m.hostApplyRule
}
