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
    "configcenter/src/common/mapstr"
    "configcenter/src/common/metadata"
    "configcenter/src/common/selector"
)

// ModelAttributeGroup model attribute group methods definitions
type ModelAttributeGroup interface {
	CreateModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.CreateModelAttributeGroup) (*metadata.CreateOneDataResult, error)
	SetModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.SetModelAttributeGroup) (*metadata.SetDataResult, error)
	UpdateModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	UpdateModelAttributeGroupByCondition(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error)
	SearchModelAttributeGroupByCondition(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error)
	DeleteModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	DeleteModelAttributeGroupByCondition(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// ModelClassification model classification methods definitions
type ModelClassification interface {
	CreateOneModelClassification(ctx ContextParams, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error)
	CreateManyModelClassification(ctx ContextParams, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error)
	SetManyModelClassification(ctx ContextParams, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error)
	SetOneModelClassification(ctx ContextParams, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error)
	UpdateModelClassification(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModelClassification(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModelClassification(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelClassificationDataResult, error)
}

// ModelAttribute model attribute methods definitions
type ModelAttribute interface {
	CreateModelAttributes(ctx ContextParams, objID string, inputParam metadata.CreateModelAttributes) (*metadata.CreateManyDataResult, error)
	SetModelAttributes(ctx ContextParams, objID string, inputParam metadata.SetModelAttributes) (*metadata.SetDataResult, error)
	UpdateModelAttributes(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	UpdateModelAttributesByCondition(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModelAttributes(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModelAttributes(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error)
	SearchModelAttributesByCondition(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error)
}

// ModelAttrUnique model attribute  unique methods definitions
type ModelAttrUnique interface {
	CreateModelAttrUnique(ctx ContextParams, objID string, data metadata.CreateModelAttrUnique) (*metadata.CreateOneDataResult, error)
	UpdateModelAttrUnique(ctx ContextParams, objID string, id uint64, data metadata.UpdateModelAttrUnique) (*metadata.UpdatedCount, error)
	DeleteModelAttrUnique(ctx ContextParams, objID string, id uint64, meta metadata.DeleteModelAttrUnique) (*metadata.DeletedCount, error)
	SearchModelAttrUnique(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryUniqueResult, error)
}

// ModelOperation model methods
type ModelOperation interface {
	ModelClassification
	ModelAttributeGroup
	ModelAttribute
	ModelAttrUnique

	CreateModel(ctx ContextParams, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error)
	SetModel(ctx ContextParams, inputParam metadata.SetModel) (*metadata.SetDataResult, error)
	UpdateModel(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModel(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModel(ctx ContextParams, modelID int64) (*metadata.DeletedCount, error)
	SearchModel(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelDataResult, error)
	SearchModelWithAttribute(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelWithAttributeDataResult, error)
}

// InstanceOperation instance methods
type InstanceOperation interface {
	CreateModelInstance(ctx ContextParams, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error)
	CreateManyModelInstance(ctx ContextParams, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error)
	UpdateModelInstance(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelInstance(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelInstance(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModelInstance(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// AssociationKind association kind methods
type AssociationKind interface {
	CreateAssociationKind(ctx ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error)
	CreateManyAssociationKind(ctx ContextParams, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyDataResult, error)
	SetAssociationKind(ctx ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error)
	SetManyAssociationKind(ctx ContextParams, inputParam metadata.SetManyAssociationKind) (*metadata.SetDataResult, error)
	UpdateAssociationKind(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteAssociationKind(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteAssociationKind(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchAssociationKind(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
}

// ModelAssociation manager model association
type ModelAssociation interface {
	CreateModelAssociation(ctx ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	CreateMainlineModelAssociation(ctx ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	SetModelAssociation(ctx ContextParams, inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error)
	UpdateModelAssociation(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelAssociation(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModelAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// InstanceAssociation manager instance association
type InstanceAssociation interface {
	CreateOneInstanceAssociation(ctx ContextParams, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error)
	CreateManyInstanceAssociation(ctx ContextParams, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error)
	SearchInstanceAssociation(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteInstanceAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// DataSynchronizeOperation manager data synchronize interface
type DataSynchronizeOperation interface {
	SynchronizeInstanceAdapter(ctx ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	SynchronizeModelAdapter(ctx ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	SynchronizeAssociationAdapter(ctx ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	Find(ctx ContextParams, find *metadata.SynchronizeFindInfoParameter) ([]mapstr.MapStr, uint64, error)
	ClearData(ctx ContextParams, input *metadata.SynchronizeClearDataParameter) error
	SetIdentifierFlag(ctx ContextParams, input *metadata.SetIdenifierFlag) ([]metadata.ExceptionResult, error)
}

// TopoOperation methods
type TopoOperation interface {
	SearchMainlineModelTopo(ctx context.Context, header http.Header, withDetail bool) (*metadata.TopoModelNode, error)
	SearchMainlineInstanceTopo(ctx context.Context, header http.Header, objID int64, withDetail bool) (*metadata.TopoInstanceNode, error)
}

// HostOperation methods
type HostOperation interface {
	TransferToInnerModule(ctx ContextParams, input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, error)
	TransferToNormalModule(ctx ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error)
	TransferToAnotherBusiness(ctx ContextParams, input *metadata.TransferHostsCrossBusinessRequest) ([]metadata.ExceptionResult, error)
	RemoveFromModule(ctx ContextParams, input *metadata.RemoveHostsFromModuleOption) ([]metadata.ExceptionResult, error)
	DeleteFromSystem(ctx ContextParams, input *metadata.DeleteHostRequest) ([]metadata.ExceptionResult, error)
	GetHostModuleRelation(ctx ContextParams, input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error)
	Identifier(ctx ContextParams, input *metadata.SearchHostIdentifierParam) ([]metadata.HostIdentifier, error)
	UpdateHostCloudAreaField(ctx ContextParams, input metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder

	LockHost(params ContextParams, input *metadata.HostLockRequest) errors.CCError
	UnlockHost(params ContextParams, input *metadata.HostLockRequest) errors.CCError
	QueryHostLock(params ContextParams, input *metadata.QueryHostLockRequest) ([]metadata.HostLockData, errors.CCError)

	// cloud sync
	CreateCloudSyncTask(ctx ContextParams, input *metadata.CloudTaskList) (uint64, error)
	CreateResourceConfirm(ctx ContextParams, input *metadata.ResourceConfirm) (uint64, error)
	CreateCloudSyncHistory(ctx ContextParams, input *metadata.CloudHistory) (uint64, error)
	CreateConfirmHistory(ctx ContextParams, input mapstr.MapStr) (uint64, error)

	// host search
	ListHosts(ctx ContextParams, input metadata.ListHosts) (*metadata.ListHostResult, error)
}

// AssociationOperation association methods
type AssociationOperation interface {
	AssociationKind
	ModelAssociation
	InstanceAssociation
}

type AuditOperation interface {
	CreateAuditLog(ctx ContextParams, logs ...metadata.SaveAuditLogParams) error
	SearchAuditLog(ctx ContextParams, param metadata.QueryInput) ([]metadata.OperationLog, uint64, error)
}

type StatisticOperation interface {
	SearchInstCount(ctx ContextParams, inputParam mapstr.MapStr) (uint64, error)
	SearchChartDataCommon(ctx ContextParams, inputParam metadata.ChartConfig) (interface{}, error)
	SearchOperationChart(ctx ContextParams, inputParam interface{}) (*metadata.ChartClassification, error)
	CreateOperationChart(ctx ContextParams, inputParam metadata.ChartConfig) (uint64, error)
	UpdateChartPosition(ctx ContextParams, inputParam interface{}) (interface{}, error)
	DeleteOperationChart(ctx ContextParams, id int64) (interface{}, error)
	UpdateOperationChart(ctx ContextParams, inputParam mapstr.MapStr) (interface{}, error)
	SearchTimerChartData(ctx ContextParams, inputParam metadata.ChartConfig) (interface{}, error)
	TimerFreshData(params ContextParams) error
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
}

// ProcessOperation methods
type ProcessOperation interface {
	// service category
	CreateServiceCategory(ctx ContextParams, category metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder)
	GetServiceCategory(ctx ContextParams, categoryID int64) (*metadata.ServiceCategory, errors.CCErrorCoder)
	GetDefaultServiceCategory(ctx ContextParams) (*metadata.ServiceCategory, errors.CCErrorCoder)
	UpdateServiceCategory(ctx ContextParams, categoryID int64, category metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder)
	ListServiceCategories(ctx ContextParams, bizID int64, withStatistics bool) (*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder)
	DeleteServiceCategory(ctx ContextParams, categoryID int64) errors.CCErrorCoder
	IsServiceCategoryLeafNode(ctx ContextParams, categoryID int64) (bool, errors.CCErrorCoder)

	// service template
	CreateServiceTemplate(ctx ContextParams, template metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder)
	GetServiceTemplate(ctx ContextParams, templateID int64) (*metadata.ServiceTemplate, errors.CCErrorCoder)
	UpdateServiceTemplate(ctx ContextParams, templateID int64, template metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder)
	ListServiceTemplates(ctx ContextParams, option metadata.ListServiceTemplateOption) (*metadata.MultipleServiceTemplate, errors.CCErrorCoder)
	DeleteServiceTemplate(ctx ContextParams, serviceTemplateID int64) errors.CCErrorCoder

	// process template
	CreateProcessTemplate(ctx ContextParams, template metadata.ProcessTemplate) (*metadata.ProcessTemplate, errors.CCErrorCoder)
	GetProcessTemplate(ctx ContextParams, templateID int64) (*metadata.ProcessTemplate, errors.CCErrorCoder)
	UpdateProcessTemplate(ctx ContextParams, templateID int64, property map[string]interface{}) (*metadata.ProcessTemplate, errors.CCErrorCoder)
	ListProcessTemplates(ctx ContextParams, option metadata.ListProcessTemplatesOption) (*metadata.MultipleProcessTemplate, errors.CCErrorCoder)
	DeleteProcessTemplate(ctx ContextParams, processTemplateID int64) errors.CCErrorCoder

	// service instance
	CreateServiceInstance(ctx ContextParams, template metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder)
	GetServiceInstance(ctx ContextParams, templateID int64) (*metadata.ServiceInstance, errors.CCErrorCoder)
	UpdateServiceInstance(ctx ContextParams, instanceID int64, instance metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder)
	ListServiceInstance(ctx ContextParams, option metadata.ListServiceInstanceOption) (*metadata.MultipleServiceInstance, errors.CCErrorCoder)
	ListServiceInstanceDetail(ctx ContextParams, option metadata.ListServiceInstanceDetailOption) (*metadata.MultipleServiceInstanceDetail, errors.CCErrorCoder)
	DeleteServiceInstance(ctx ContextParams, serviceInstanceIDs []int64) errors.CCErrorCoder
	AutoCreateServiceInstanceModuleHost(ctx ContextParams, hostID int64, moduleID int64) (*metadata.ServiceInstance, errors.CCErrorCoder)
	RemoveTemplateBindingOnModule(ctx ContextParams, moduleID int64) errors.CCErrorCoder
	ReconstructServiceInstanceName(ctx ContextParams, instanceID int64) errors.CCErrorCoder

	// process instance relation
	CreateProcessInstanceRelation(ctx ContextParams, relation *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	GetProcessInstanceRelation(ctx ContextParams, processInstanceID int64) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	UpdateProcessInstanceRelation(ctx ContextParams, processInstanceID int64, relation metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder)
	ListProcessInstanceRelation(ctx ContextParams, option metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder)
	DeleteProcessInstanceRelation(ctx ContextParams, option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder

	GetBusinessDefaultSetModuleInfo(ctx ContextParams, bizID int64) (metadata.BusinessDefaultSetModuleInfo, errors.CCErrorCoder)
	GetProc2Module(ctx ContextParams, option *metadata.GetProc2ModuleOption) ([]metadata.Proc2Module, errors.CCErrorCoder)
}

type LabelOperation interface {
	AddLabel(ctx ContextParams, tableName string, option selector.LabelAddOption) errors.CCErrorCoder
	RemoveLabel(ctx ContextParams, tableName string, option selector.LabelRemoveOption) errors.CCErrorCoder
}

type SetTemplateOperation interface {
	CreateSetTemplate(ctx ContextParams, bizID int64, option metadata.CreateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder)
	UpdateSetTemplate(ctx ContextParams, setTemplateID int64, option metadata.UpdateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder)
	DeleteSetTemplate(ctx ContextParams, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder
	GetSetTemplate(ctx ContextParams, bizID int64, setTemplateID int64) (metadata.SetTemplate, errors.CCErrorCoder)
	ListSetTemplate(ctx ContextParams, bizID int64, option metadata.ListSetTemplateOption) (metadata.MultipleSetTemplateResult, errors.CCErrorCoder)
	ListSetServiceTemplateRelations(ctx ContextParams, bizID int64, setTemplateID int64) ([]metadata.SetServiceTemplateRelation, errors.CCErrorCoder)
	ListSetTplRelatedSvcTpl(ctx ContextParams, bizID, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder)
	UpdateSetTemplateSyncStatus(ctx ContextParams, setID int64, option metadata.SetTemplateSyncStatus) errors.CCErrorCoder
	ListSetTemplateSyncStatus(ctx ContextParams, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
	ListSetTemplateSyncHistory(ctx ContextParams, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
	DeleteSetTemplateSyncStatus(ctx ContextParams, option metadata.DeleteSetTemplateSyncStatusOption) errors.CCErrorCoder
}

type HostApplyRuleOperation interface {
	CreateHostApplyRule(ctx ContextParams, bizID int64, option metadata.CreateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder)
	UpdateHostApplyRule(ctx ContextParams, bizID int64, ruleID int64, option metadata.UpdateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder)
	DeleteHostApplyRule(ctx ContextParams, bizID int64, ruleIDs ...int64) errors.CCErrorCoder
	GetHostApplyRule(ctx ContextParams, bizID int64, ruleID int64) (metadata.HostApplyRule, errors.CCErrorCoder)
	ListHostApplyRule(ctx ContextParams, bizID int64, option metadata.ListHostApplyRuleOption) (metadata.MultipleHostApplyRuleResult, errors.CCErrorCoder)
	GenerateApplyPlan(ctx ContextParams, bizID int64, option metadata.HostApplyPlanOption) (metadata.HostApplyPlanResult, errors.CCErrorCoder)
	SearchRuleRelatedModules(ctx ContextParams, bizID int64, option metadata.SearchRuleRelatedModulesOption) ([]metadata.Module, errors.CCErrorCoder)
	BatchUpdateHostApplyRule(ctx ContextParams, bizID int64, option metadata.BatchCreateOrUpdateApplyRuleOption) (metadata.BatchCreateOrUpdateHostApplyRuleResult, errors.CCErrorCoder)
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

func (m *core) HostApplyRuleOperation() HostApplyRuleOperation {
	return m.hostApplyRule
}
