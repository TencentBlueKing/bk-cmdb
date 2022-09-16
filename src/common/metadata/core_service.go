/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/selector"
	"configcenter/src/common/util"
)

// CreateModelAttributeGroup used to create a new group for some attributes
type CreateModelAttributeGroup struct {
	Data Group `json:"data"`
}

// SetModelAttributeGroup used to create a new group for  some attributes, if it is exists, then update it
type SetModelAttributeGroup CreateModelAttributeGroup

// CreateManyModelClassifiaction create many input params
type CreateManyModelClassifiaction struct {
	Data []Classification `json:"datas"`
}

// CreateOneModelClassification create one model classification
type CreateOneModelClassification struct {
	Data Classification `json:"data"`
}

// SetManyModelClassification set many input params
type SetManyModelClassification CreateManyModelClassifiaction

// SetOneModelClassification set one input params
type SetOneModelClassification CreateOneModelClassification

// DeleteModelClassificationResult delete the model classification result
type DeleteModelClassificationResult struct {
	BaseResp `json:",inline"`
	Data     DeletedCount `json:"data"`
}

// CreateModel create model params
type CreateModel struct {
	Spec       Object      `json:"spec"`
	Attributes []Attribute `json:"attributes"`
}

// SetModel define SetMode method input params
type SetModel CreateModel

// SearchModelInfo search  model params
type SearchModelInfo struct {
	Spec       Object      `json:"spec"`
	Attributes []Attribute `json:"attributes"`
}

// CreateModelAttributes create model attributes
type CreateModelAttributes struct {
	Attributes []Attribute `json:"attributes"`
}

// SetModelAttributes TODO
type SetModelAttributes CreateModelAttributes

// CreateModelAttrUnique TODO
type CreateModelAttrUnique struct {
	Data ObjectUnique `json:"data"`
}

// UpdateModelAttrUnique TODO
type UpdateModelAttrUnique struct {
	Data UpdateUniqueRequest `json:"data"`
}

// DeleteModelAttrUnique TODO
type DeleteModelAttrUnique struct {
	BizID int64 `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
}

// CreateModelInstance TODO
type CreateModelInstance struct {
	Data mapstr.MapStr `json:"data"`
}

// CreateManyModelInstance TODO
type CreateManyModelInstance struct {
	Datas []mapstr.MapStr `json:"datas"`
}

// SetModelInstance TODO
type SetModelInstance CreateModelInstance

// SetManyModelInstance TODO
type SetManyModelInstance CreateManyModelInstance

// CreateAssociationKind TODO
type CreateAssociationKind struct {
	Data AssociationKind `json:"data"`
}

// CreateManyAssociationKind TODO
type CreateManyAssociationKind struct {
	Datas []AssociationKind `json:"datas"`
}

// SetAssociationKind TODO
type SetAssociationKind CreateAssociationKind

// SetManyAssociationKind TODO
type SetManyAssociationKind CreateManyAssociationKind

// CreateModelAssociation TODO
type CreateModelAssociation struct {
	Spec Association `json:"spec"`
}

// SetModelAssociation TODO
type SetModelAssociation CreateModelAssociation

// CreateOneInstanceAssociation TODO
type CreateOneInstanceAssociation struct {
	Data InstAsst `json:"data"`
}

// CreateManyInstanceAssociation TODO
type CreateManyInstanceAssociation struct {
	Datas []InstAsst `json:"datas"`
}

// Dimension TODO
type Dimension struct {
	AppID int64 `json:"bk_biz_id"`
}

// SetOneInstanceAssociation TODO
type SetOneInstanceAssociation CreateOneInstanceAssociation

// SetManyInstanceAssociation TODO
type SetManyInstanceAssociation CreateManyInstanceAssociation

// TopoModelNode TODO
type TopoModelNode struct {
	Children []*TopoModelNode
	ObjectID string
}

// SearchTopoModelNodeResult TODO
type SearchTopoModelNodeResult struct {
	BaseResp `json:",inline"`
	Data     TopoModelNode `json:"data"`
}

// LeftestObjectIDList extract leftest node's id of each level, arrange as a list
// it's useful in model mainline topo case, as bk_mainline relationship degenerate to a list.
func (tn *TopoModelNode) LeftestObjectIDList() []string {
	objectIDs := make([]string, 0)
	node := tn
	for {
		objectIDs = append(objectIDs, node.ObjectID)
		if len(node.Children) == 0 {
			break
		}
		node = node.Children[0]
	}
	return objectIDs
}

// TopoInstanceNodeSimplify TODO
type TopoInstanceNodeSimplify struct {
	ObjectID     string `json:"bk_obj_id" field:"bk_obj_id" mapstructure:"bk_obj_id"`
	InstanceID   int64  `json:"bk_inst_id" field:"bk_inst_id" mapstructure:"bk_inst_id"`
	InstanceName string `json:"bk_inst_name" field:"bk_inst_name" mapstructure:"bk_inst_name"`
}

// TopoInstanceNode TODO
type TopoInstanceNode struct {
	Children     []*TopoInstanceNode
	ObjectID     string
	InstanceID   int64
	InstanceName string
	Detail       map[string]interface{}
}

// SearchTopoInstanceNodeResult TODO
type SearchTopoInstanceNodeResult struct {
	BaseResp `json:",inline"`
	Data     TopoInstanceNode `json:"data"`
}

// Name TODO
func (node *TopoInstanceNode) Name() string {
	var name string
	var exist bool
	var val interface{}
	switch node.ObjectID {
	case common.BKInnerObjIDSet:
		val, exist = node.Detail[common.BKSetNameField]
	case common.BKInnerObjIDApp:
		val, exist = node.Detail[common.BKAppNameField]
	case common.BKInnerObjIDModule:
		val, exist = node.Detail[common.BKModuleNameField]
	default:
		val, exist = node.Detail[common.BKInstNameField]
	}

	if exist == true {
		name = util.GetStrByInterface(val)
	} else {
		blog.V(7).Infof("extract topo instance node:%+v name failed", *node)
		name = fmt.Sprintf("%s:%d", node.ObjectID, node.InstanceID)
	}
	return name
}

// TraversalFindModule TODO
func (node *TopoInstanceNode) TraversalFindModule(targetID int64) []*TopoInstanceNode {
	// ex: module1 ==> reverse([bizID, mainline1, ..., mainline2, set1, module1])
	return node.TraversalFindNode(common.BKInnerObjIDModule, targetID)
}

// TraversalFindNode TODO
// common.BKInnerObjIDObject used for matching custom level node
func (node *TopoInstanceNode) TraversalFindNode(objectType string, targetID int64) []*TopoInstanceNode {
	if objectType == common.BKInnerObjIDObject && !common.IsInnerModel(node.ObjectID) && node.InstanceID == targetID {
		return []*TopoInstanceNode{node}
	}
	if node.ObjectID == objectType && node.InstanceID == targetID {
		return []*TopoInstanceNode{node}
	}

	for _, child := range node.Children {
		path := child.TraversalFindNode(objectType, targetID)
		if len(path) > 0 {
			path = append(path, node)
			return path
		}
	}

	return []*TopoInstanceNode{}
}

// DeepFirstTraversal TODO
func (node *TopoInstanceNode) DeepFirstTraversal(f func(node *TopoInstanceNode)) {
	if node == nil {
		return
	}
	for _, child := range node.Children {
		child.DeepFirstTraversal(f)
	}
	f(node)
}

// ToSimplify TODO
func (node *TopoInstanceNode) ToSimplify() *TopoInstanceNodeSimplify {
	if node == nil {
		return nil
	}
	return &TopoInstanceNodeSimplify{
		ObjectID:     node.ObjectID,
		InstanceID:   node.InstanceID,
		InstanceName: node.InstanceName,
	}
}

// TopoInstance TODO
type TopoInstance struct {
	ObjectID         string
	InstanceID       int64
	InstanceName     string
	ParentInstanceID int64
	Detail           map[string]interface{}
	Default          int64
}

// Key generate a unique key for instance(as instances's of different object type maybe conflict)
func (ti *TopoInstance) Key() string {
	return fmt.Sprintf("%s:%d", ti.ObjectID, ti.InstanceID)
}

// TransferHostsCrossBusinessRequest Transfer host across business request parameter
type TransferHostsCrossBusinessRequest struct {
	SrcApplicationIDs []int64 `json:"src_bk_biz_id"`
	DstApplicationID  int64   `json:"dst_bk_biz_id"`
	HostIDArr         []int64 `json:"bk_host_id"`
	DstModuleIDArr    []int64 `json:"bk_module_ids"`
}

// HostModuleRelationRequest gethost module relation request parameter
type HostModuleRelationRequest struct {
	ApplicationID int64    `json:"bk_biz_id" bson:"bk_biz_id" field:"bk_biz_id" mapstructure:"bk_biz_id"`
	SetIDArr      []int64  `json:"bk_set_ids" bson:"bk_set_ids" field:"bk_set_ids" mapstructure:"bk_set_ids"`
	HostIDArr     []int64  `json:"bk_host_ids" bson:"bk_host_ids" field:"bk_host_ids" mapstructure:"bk_host_ids"`
	ModuleIDArr   []int64  `json:"bk_module_ids" bson:"bk_module_ids" field:"bk_module_ids" mapstructure:"bk_module_ids"`
	Page          BasePage `json:"page" bson:"page" field:"page" mapstructure:"page"`
	Fields        []string `json:"field" bson:"field"  field:"field" mapstructure:"field"`
}

// Empty empty struct
func (h *HostModuleRelationRequest) Empty() bool {
	if h.ApplicationID != 0 {
		return false
	}
	if len(h.SetIDArr) != 0 {
		return false
	}
	if len(h.ModuleIDArr) != 0 {
		return false
	}

	if len(h.HostIDArr) != 0 {
		return false
	}
	return true
}

// DeleteHostRequest delete host from application
type DeleteHostRequest struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostIDArr     []int64 `json:"bk_host_ids"`
}

// ListServiceCategoryOption list service category option
type ListServiceCategoryOption struct {
	BizID int64 `json:"bk_biz_id"`
}

// OneServiceCategoryResult TODO
type OneServiceCategoryResult struct {
	BaseResp `json:",inline"`
	Data     ServiceCategory `json:"data"`
}

// OneServiceCategoryWithStatisticsResult TODO
type OneServiceCategoryWithStatisticsResult struct {
	BaseResp `json:",inline"`
	Data     ServiceCategoryWithStatistics `json:"data"`
}

// MultipleServiceCategory TODO
type MultipleServiceCategory struct {
	Count int64             `json:"count"`
	Info  []ServiceCategory `json:"info"`
}

// MultipleServiceCategoryWithStatistics TODO
type MultipleServiceCategoryWithStatistics struct {
	Count int64                           `json:"count"`
	Info  []ServiceCategoryWithStatistics `json:"info"`
}

// MultipleServiceCategoryResult TODO
type MultipleServiceCategoryResult struct {
	BaseResp `json:",inline"`
	Data     MultipleServiceCategory `json:"data"`
}

// MultipleServiceCategoryWithStatisticsResult TODO
type MultipleServiceCategoryWithStatisticsResult struct {
	BaseResp `json:",inline"`
	Data     MultipleServiceCategoryWithStatistics `json:"data"`
}

// ListServiceTemplateOption TODO
type ListServiceTemplateOption struct {
	BusinessID         int64    `json:"bk_biz_id"`
	ServiceCategoryID  *int64   `json:"service_category_id"`
	ServiceTemplateIDs []int64  `json:"service_template_ids"`
	Page               BasePage `json:"page,omitempty"`
	// search service templates by name
	Search string `json:"search"`
	// used with search, means whether search service templates with exact name or not
	IsExact bool `json:"is_exact"`
}

// FindServiceTemplateCountInfoOption TODO
type FindServiceTemplateCountInfoOption struct {
	ServiceTemplateIDs []int64 `json:"service_template_ids"`
}

// Validate TODO
func (o *FindServiceTemplateCountInfoOption) Validate() (rawError errors.RawErrorInfo) {
	maxLimit := 500
	if len(o.ServiceTemplateIDs) == 0 || len(o.ServiceTemplateIDs) > maxLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"service_template_ids", maxLimit},
		}
	}

	return errors.RawErrorInfo{}
}

// FindServiceTemplateCountInfoResult TODO
type FindServiceTemplateCountInfoResult struct {
	ServiceTemplateID    int64 `json:"service_template_id"`
	ProcessTemplateCount int64 `json:"process_template_count"`
	ServiceInstanceCount int64 `json:"service_instance_count"`
	ModuleCount          int64 `json:"module_count"`
}

// OneServiceTemplateResult TODO
type OneServiceTemplateResult struct {
	BaseResp `json:",inline"`
	Data     ServiceTemplate `json:"data"`
}

// OneServiceTemplateWithStatisticsResult TODO
type OneServiceTemplateWithStatisticsResult struct {
	BaseResp `json:",inline"`
	Data     ServiceTemplateWithStatistics `json:"data"`
}

// MultipleServiceTemplateDetailResult TODO
type MultipleServiceTemplateDetailResult struct {
	BaseResp `json:",inline"`
	Data     MultipleServiceTemplateDetail `json:"data"`
}

// MultipleServiceTemplateDetail TODO
type MultipleServiceTemplateDetail struct {
	Count uint64                  `json:"count"`
	Info  []ServiceTemplateDetail `json:"info"`
}

// MultipleServiceTemplate TODO
type MultipleServiceTemplate struct {
	Count uint64            `json:"count"`
	Info  []ServiceTemplate `json:"info"`
}

// ListServiceInstanceOption TODO
type ListServiceInstanceOption struct {
	BusinessID         int64              `json:"bk_biz_id"`
	ServiceTemplateID  int64              `json:"service_template_id"`
	HostIDs            []int64            `json:"bk_host_ids"`
	ModuleIDs          []int64            `json:"bk_module_ids"`
	SearchKey          *string            `json:"search_key"`
	ServiceInstanceIDs []int64            `json:"service_instance_ids"`
	Selectors          selector.Selectors `json:"selectors"`
	Fields             []string           `json:"fields"`
	Page               BasePage           `json:"page"`
}

// ListServiceInstanceDetailOption TODO
type ListServiceInstanceDetailOption struct {
	BusinessID int64 `json:"bk_biz_id"`
	ModuleID   int64 `json:"bk_module_id"`
	// HostID is legacy code, compatible with this field, but not recommended.
	HostID             int64              `json:"bk_host_id"`
	HostList           []int64            `json:"bk_host_list"`
	ServiceInstanceIDs []int64            `json:"service_instance_ids"`
	Selectors          selector.Selectors `json:"selectors,omitempty"`
	Page               BasePage           `json:"page,omitempty"`
}

// ListProcessInstanceRelationOption TODO
type ListProcessInstanceRelationOption struct {
	BusinessID         int64    `json:"bk_biz_id"`
	ProcessIDs         []int64  `json:"process_ids,omitempty"`
	ServiceInstanceIDs []int64  `json:"service_instance_id,omitempty"`
	ProcessTemplateID  int64    `json:"process_template_id,omitempty"`
	HostID             int64    `json:"host_id,omitempty"`
	Page               BasePage `json:"page" field:"page"`
}

// MultipleServiceTemplateResult TODO
type MultipleServiceTemplateResult struct {
	BaseResp `json:",inline"`
	Data     MultipleServiceTemplate `json:"data"`
}

// OneProcessTemplateResult TODO
type OneProcessTemplateResult struct {
	BaseResp `json:",inline"`
	Data     ProcessTemplate `json:"data"`
}

// MultipleProcessTemplate TODO
type MultipleProcessTemplate struct {
	Count uint64            `json:"count"`
	Info  []ProcessTemplate `json:"info"`
}

// MultipleProcessTemplateResult TODO
type MultipleProcessTemplateResult struct {
	BaseResp `json:",inline"`
	Data     MultipleProcessTemplate `json:"data"`
}

// DeleteProcessInstanceRelationOption TODO
type DeleteProcessInstanceRelationOption struct {
	BusinessID         *int64  `json:"bk_biz_id"`
	ProcessIDs         []int64 `json:"bk_process_id,omitempty"`
	ServiceInstanceIDs []int64 `json:"service_instance_id,omitempty"`
	ProcessTemplateIDs []int64 `json:"process_template_id,omitempty"`
	HostIDs            []int64 `json:"bk_host_id,omitempty"`
	ModuleIDs          []int64 `json:"bk_module_id,omitempty"`
}

// ListProcessTemplatesOption TODO
type ListProcessTemplatesOption struct {
	BusinessID         int64    `json:"bk_biz_id" bson:"bk_biz_id"`
	ProcessTemplateIDs []int64  `json:"process_template_ids,omitempty" bson:"process_template_ids"`
	ServiceTemplateIDs []int64  `json:"service_template_ids,omitempty" bson:"service_template_ids"`
	Page               BasePage `json:"page" field:"page" bson:"page"`
}

// ListServiceCategoriesOption TODO
type ListServiceCategoriesOption struct {
	BusinessID         int64   `json:"bk_biz_id" bson:"bk_biz_id"`
	ServiceCategoryIDs []int64 `json:"service_category_ids,omitempty" bson:"service_category_ids"`
	WithStatistics     bool    `json:"with_statistics" bson:"with_statistics"`
}

// OneServiceInstanceResult TODO
type OneServiceInstanceResult struct {
	BaseResp `json:",inline"`
	Data     ServiceInstance `json:"data"`
}

// ManyServiceInstanceResult TODO
type ManyServiceInstanceResult struct {
	BaseResp `json:",inline"`
	Data     []*ServiceInstance `json:"data"`
}

// MultipleServiceInstance 增加指定返回field字段的result briefInfo。
type MultipleServiceInstance struct {
	Count uint64            `json:"count"`
	Info  []ServiceInstance `json:"info"`
}

// MultipleServiceInstanceDetail TODO
type MultipleServiceInstanceDetail struct {
	Count uint64                  `json:"count"`
	Info  []ServiceInstanceDetail `json:"info"`
}

// MultipleServiceInstanceResult TODO
type MultipleServiceInstanceResult struct {
	BaseResp `json:",inline"`
	Data     MultipleServiceInstance `json:"data"`
}

// MultipleServiceInstanceDetailResult TODO
type MultipleServiceInstanceDetailResult struct {
	BaseResp `json:",inline"`
	Data     MultipleServiceInstanceDetail `json:"data"`
}

// OneProcessInstanceRelationResult TODO
type OneProcessInstanceRelationResult struct {
	BaseResp `json:",inline"`
	Data     ProcessInstanceRelation `json:"data"`
}

// ManyProcessInstanceRelationResult TODO
type ManyProcessInstanceRelationResult struct {
	BaseResp `json:",inline"`
	Data     []*ProcessInstanceRelation `json:"data"`
}

// MultipleProcessInstanceRelation TODO
type MultipleProcessInstanceRelation struct {
	Count uint64                    `json:"count"`
	Info  []ProcessInstanceRelation `json:"info"`
}

// MultipleProcessInstanceRelationResult TODO
type MultipleProcessInstanceRelationResult struct {
	BaseResp `json:",inline"`
	Data     MultipleProcessInstanceRelation `json:"data"`
}

// MultipleHostProcessRelation TODO
type MultipleHostProcessRelation struct {
	Count uint64                `json:"count"`
	Info  []HostProcessRelation `json:"info"`
}

// MultipleHostProcessRelationResult TODO
type MultipleHostProcessRelationResult struct {
	BaseResp `json:",inline"`
	Data     MultipleHostProcessRelation `json:"data"`
}

// BusinessDefaultSetModuleInfo TODO
type BusinessDefaultSetModuleInfo struct {
	IdleSetID       int64 `json:"idle_set_id"`
	IdleModuleID    int64 `json:"idle_module_id"`
	FaultModuleID   int64 `json:"fault_module_id"`
	RecycleModuleID int64 `json:"recycle_module_id"`
}

// IsInternalModule TODO
func (b BusinessDefaultSetModuleInfo) IsInternalModule(moduleID int64) bool {
	if moduleID == b.IdleModuleID ||
		moduleID == b.FaultModuleID ||
		moduleID == b.RecycleModuleID {
		return true
	}
	return false
}

// BusinessDefaultSetModuleInfoResult TODO
type BusinessDefaultSetModuleInfoResult struct {
	BaseResp `json:",inline"`
	Data     BusinessDefaultSetModuleInfo `json:"data"`
}

// RemoveTemplateBoundOnModuleResult TODO
type RemoveTemplateBoundOnModuleResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		ServiceTemplateID int64 `json:"service_template_id" bson:"service_template_id" field:"service_template_id"`
	} `json:"data"`
}

// MultipleMap TODO
type MultipleMap struct {
	Count uint64                   `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

// DistinctHostIDByTopoRelationRequest  distinct host id by topology request
type DistinctHostIDByTopoRelationRequest struct {
	ApplicationIDArr []int64 `json:"bk_biz_ids" bson:"bk_biz_ids" field:"bk_biz_ids" mapstructure:"bk_biz_ids"`
	SetIDArr         []int64 `json:"bk_set_ids" bson:"bk_set_ids" field:"bk_set_ids" mapstructure:"bk_set_ids"`
	HostIDArr        []int64 `json:"bk_host_ids" bson:"bk_host_ids" field:"bk_host_ids" mapstructure:"bk_host_ids"`
	ModuleIDArr      []int64 `json:"bk_module_ids" bson:"bk_module_ids" field:"bk_module_ids" mapstructure:"bk_module_ids"`
}

// Empty empty struct
func (h *DistinctHostIDByTopoRelationRequest) Empty() bool {
	if len(h.ApplicationIDArr) != 0 {
		return false
	}
	if len(h.SetIDArr) != 0 {
		return false
	}
	if len(h.ModuleIDArr) != 0 {
		return false
	}

	if len(h.HostIDArr) != 0 {
		return false
	}
	return true
}

// CloudAccountResult TODO
type CloudAccountResult struct {
	BaseResp `json:",inline"`
	Data     CloudAccount `json:"data"`
}

// MultipleCloudAccountResult TODO
type MultipleCloudAccountResult struct {
	BaseResp `json:",inline"`
	Data     MultipleCloudAccount `json:"data"`
}

// TransferHostResourceDirectory TODO
type TransferHostResourceDirectory struct {
	ModuleID int64   `json:"bk_module_id"`
	HostID   []int64 `json:"bk_host_id"`
}

// MultipleCloudAccountConfResult TODO
type MultipleCloudAccountConfResult struct {
	BaseResp `json:",inline"`
	Data     MultipleCloudAccountConf `json:"data"`
}

// CreateSyncTaskResult TODO
type CreateSyncTaskResult struct {
	BaseResp `json:",inline"`
	Data     CloudSyncTask `json:"data"`
}

// CreateSyncHistoryesult TODO
type CreateSyncHistoryesult struct {
	BaseResp `json:",inline"`
	Data     SyncHistory `json:"data"`
}

// MultipleCloudSyncTaskResult TODO
type MultipleCloudSyncTaskResult struct {
	BaseResp `json:",inline"`
	Data     MultipleCloudSyncTask `json:"data"`
}

// MultipleSyncHistoryResult TODO
type MultipleSyncHistoryResult struct {
	BaseResp `json:",inline"`
	Data     MultipleSyncHistory `json:"data"`
}

// MultipleSyncRegionResult TODO
type MultipleSyncRegionResult struct {
	BaseResp `json:",inline"`
	Data     []*Region `json:"data"`
}

// DistinctFieldOption TODO
type DistinctFieldOption struct {
	TableName string                 `json:"table_name"`
	Field     string                 `json:"field"`
	Filter    map[string]interface{} `json:"filter"`
}

// Validate TODO
func (d *DistinctFieldOption) Validate() (rawError errors.RawErrorInfo) {
	if d.TableName == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"table_name"},
		}
	}

	if d.Field == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"field"},
		}
	}

	return errors.RawErrorInfo{}
}

// CreateModelTable create model table params
type CreateModelTable struct {
	ObjectIDs  []string `json:"bk_object_ids"`
	IsMainLine bool     `json:"is_main_line"`
}

// ServiceTemplateAttributeResult service template attribute result
type ServiceTemplateAttributeResult struct {
	BaseResp `json:",inline"`
	Data     *ServTempAttrData `json:"data"`
}

// SetTemplateAttributeResult set template attribute result
type SetTemplateAttributeResult struct {
	BaseResp `json:",inline"`
	Data     *SetTempAttrData `json:"data"`
}
