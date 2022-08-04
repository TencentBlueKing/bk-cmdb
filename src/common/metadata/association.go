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

package metadata

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/querybuilder"
)

const (
	// AssociationFieldObjectID the association data field definition
	AssociationFieldObjectID = "bk_obj_id"
	// AssociationFieldAsstID the association data field bk_obj_asst_id
	AssociationFieldAsstID = "bk_obj_asst_id"
	// AssociationFieldSupplierAccount TODO
	// AssociationFieldObjectAttributeID the association data field definition
	// AssociationFieldObjectAttributeID = "bk_object_att_id"
	// AssociationFieldSupplierAccount the association data field definition
	AssociationFieldSupplierAccount = "bk_supplier_account"
	// AssociationFieldAssociationObjectID the association data field definition
	// AssociationFieldAssociationForward = "bk_asst_forward"
	// AssociationFieldAssociationObjectID the association data field definition
	AssociationFieldAssociationObjectID = "bk_asst_obj_id"
	// AssociationFieldAssociationId the association data field definition
	// AssociationFieldAssociationName = "bk_asst_name"
	// AssociationFieldAssociationId auto incr id
	AssociationFieldAssociationId = "id"
	// AssociationFieldAssociationKind TODO
	AssociationFieldAssociationKind = "bk_asst_id"
)

// SearchAssociationTypeRequest TODO
type SearchAssociationTypeRequest struct {
	BasePage  `json:"page"`
	Condition map[string]interface{} `json:"condition"`
}

// SearchAssociationType struct for search association type
type SearchAssociationType struct {
	Count int                `json:"count"`
	Info  []*AssociationKind `json:"info"`
}

// SearchAssociationTypeResult TODO
type SearchAssociationTypeResult struct {
	BaseResp `json:",inline"`
	Data     SearchAssociationType `json:"data"`
}

// CreateAssociationTypeResult TODO
type CreateAssociationTypeResult struct {
	BaseResp `json:",inline"`
	Data     RspID `json:"data"`
}

// UpdateAssociationTypeRequest TODO
type UpdateAssociationTypeRequest struct {
	AsstName  string `field:"bk_asst_name" json:"bk_asst_name" bson:"bk_asst_name"`
	SrcDes    string `field:"src_des" json:"src_des" bson:"src_des"`
	DestDes   string `field:"dest_des" json:"dest_des" bson:"dest_des"`
	Direction string `field:"direction" json:"direction" bson:"direction"`
}

// UpdateManyAssociationTypeRequest params of update many association type
type UpdateManyAssociationTypeRequest struct {
	Data map[int64]UpdateAssociationTypeRequest `json:"data"`
}

// UpdateAssociationTypeResult TODO
type UpdateAssociationTypeResult struct {
	BaseResp `json:",inline"`
	Data     string `json:"data"`
}

// DeleteAssociationTypeResult TODO
type DeleteAssociationTypeResult struct {
	BaseResp `json:",inline"`
	Data     string `json:"data"`
}

// SearchAssociationObjectRequest TODO
type SearchAssociationObjectRequest struct {
	Condition mapstr.MapStr `json:"condition"`
}

// SearchAssociationObjectResult TODO
type SearchAssociationObjectResult struct {
	BaseResp `json:",inline"`
	Data     []*Association `json:"data"`
}

// CreateAssociationObjectResult TODO
type CreateAssociationObjectResult struct {
	BaseResp `json:",inline"`
	Data     RspID `json:"data"`
}

// UpdateAssociationObjectRequest TODO
type UpdateAssociationObjectRequest struct {
	AsstName string `field:"bk_asst_name" json:"bk_asst_name" bson:"bk_asst_name"`
}

// UpdateAssociationObjectResult TODO
type UpdateAssociationObjectResult struct {
	BaseResp `json:",inline"`
	Data     string `json:"data"`
}

// DeleteAssociationObjectResult TODO
type DeleteAssociationObjectResult struct {
	BaseResp `json:",inline"`
	Data     string `json:"data"`
}

// SearchAssociationRelatedInstRequestCond TODO
type SearchAssociationRelatedInstRequestCond struct {
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id,omitempty" bson:"bk_obj_id,omitempty"`
	InstID   int64  `field:"bk_inst_id" json:"bk_inst_id,omitempty" bson:"bk_inst_id,omitempty"`
}

// SearchAssociationInstRequest TODO
type SearchAssociationInstRequest struct {
	Condition mapstr.MapStr `json:"condition"` // construct condition mapstr by condition.Condition
	ObjID     string        `json:"bk_obj_id"`
}

// SearchAssociationRelatedInstRequest TODO
type SearchAssociationRelatedInstRequest struct {
	Fields    []string                                `json:"fields"`
	Page      BasePage                                `json:"page"`
	Condition SearchAssociationRelatedInstRequestCond `json:"condition"`
}

// SearchAssociationInstResult TODO
type SearchAssociationInstResult struct {
	BaseResp `json:",inline"`
	Data     []*InstAsst `json:"data"`
}

// AsstResult model item result
type AsstResult struct {
	Count int           `json:"count"`
	Info  []Association `json:"info"`
}

// SearchAsstModelResp query association model result
type SearchAsstModelResp struct {
	BaseResp `json:",inline"`
	Data     AsstResult `json:"data"`
}

// SearchInstAssociationListResult the struct of list instance association result
type SearchInstAssociationListResult struct {
	Association struct {
		Src []InstAsst `json:"src"`
		Dst []InstAsst `json:"dst"`
	} `json:"association"`
	Inst map[string][]mapstr.MapStr `json:"instance"`
}

// InstAndAssocDetailResult search inst and association detail result
type InstAndAssocDetailResult struct {
	BaseResp `json:",inline"`
	Data     InstAndAssocDetailData `json:"data"`
}

// InstAndAssocDetailData search inst and association detail return data
type InstAndAssocDetailData struct {
	Asst []InstAsst      `field:"association" json:"association"`
	Src  []mapstr.MapStr `field:"src" json:"src"`
	Dst  []mapstr.MapStr `field:"dst" json:"dst"`
}

// InstAndAssocRequest search inst and association detail request
type InstAndAssocRequest struct {
	Condition struct {
		AsstFilter *querybuilder.QueryFilter `field:"asst_filter" json:"asst_filter"`
		AsstFields []string                  `field:"asst_fields" json:"asst_fields"`
		SrcFields  []string                  `field:"src_fields" json:"src_fields"`
		DstFields  []string                  `field:"dst_fields" json:"dst_fields"`
		SrcDetail  bool                      `field:"src_detail" json:"src_detail"`
		DstDetail  bool                      `field:"dst_detail" json:"dst_detail"`
	} `field:"condition" json:"condition"`
	Page BasePage `field:"page" json:"page"`
}

// Validate validate InstAndAssocDetailData
func (assoc *InstAndAssocRequest) Validate() errors.RawErrorInfo {

	if assoc.Condition.AsstFilter == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"condition.asst_filter"},
		}
	}

	filterOption := querybuilder.RuleOption{NeedSameSliceElementType: true}
	if key, err := assoc.Condition.AsstFilter.Validate(&filterOption); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{fmt.Sprintf("condition.asst_filter.%s", key)},
		}
	}

	if assoc.Page.Limit > 200 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommPageLimitIsExceeded,
		}
	}

	if assoc.Page.Limit == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"page.limit"},
		}
	}

	return errors.RawErrorInfo{}
}

// CreateAssociationInstRequest TODO
type CreateAssociationInstRequest struct {
	ObjectAsstID string `field:"bk_obj_asst_id" json:"bk_obj_asst_id,omitempty" bson:"bk_obj_asst_id,omitempty"`
	InstID       int64  `field:"bk_inst_id" json:"bk_inst_id,omitempty" bson:"bk_inst_id,omitempty"`
	AsstInstID   int64  `field:"bk_asst_inst_id" json:"bk_asst_inst_id,omitempty" bson:"bk_asst_inst_id,omitempty"`
}

// CreateAssociationInstResult TODO
type CreateAssociationInstResult struct {
	BaseResp `json:",inline"`
	Data     RspID `json:"data"`
}

// CreateManyInstAsstRequest parameter structure for creating multiple instances associations
type CreateManyInstAsstRequest struct {
	ObjectID     string     `field:"bk_obj_id" json:"bk_obj_id,omitempty" bson:"bk_obj_id,omitempty"`
	AsstObjectID string     `field:"bk_asst_obj_id" json:"bk_asst_obj_id,omitempty" bson:"bk_asst_obj_id,omitempty"`
	ObjectAsstID string     `field:"bk_obj_asst_id" json:"bk_obj_asst_id,omitempty" bson:"bk_obj_asst_id,omitempty"`
	Details      []InstAsst `field:"details" json:"details,omitempty" bson:"details,omitempty"`
}

// Validate TODO
func (assoc *CreateManyInstAsstRequest) Validate() errors.RawErrorInfo {
	if len(assoc.ObjectAsstID) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{common.AssociationObjAsstIDField},
		}
	}

	if len(assoc.ObjectID) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if len(assoc.AsstObjectID) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{common.BKAsstObjIDField},
		}
	}

	if len(assoc.Details) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommInstDataNil,
			Args:    []interface{}{"details"},
		}
	}

	if len(assoc.Details) > 200 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"details", 200},
		}
	}

	// NOTE: if bk_obj_asst_id changes, the logic here needs to be modified
	if assoc.ObjectAsstID[:len(assoc.ObjectID)] != assoc.ObjectID {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if assoc.ObjectAsstID[len(assoc.ObjectAsstID)-len(assoc.AsstObjectID):] != assoc.AsstObjectID {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKAsstObjIDField},
		}
	}

	return errors.RawErrorInfo{}
}

// AssociationInstDetails source and target instances of associations
type AssociationInstDetails struct {
	InstID     int64 `field:"bk_inst_id" json:"bk_inst_id,omitempty" bson:"bk_inst_id,omitempty"`
	AsstInstID int64 `field:"bk_asst_inst_id" json:"bk_asst_inst_id,omitempty" bson:"bk_asst_inst_id,omitempty"`
}

// CreateManyInstAsstResultDetail details of creating instance association result
type CreateManyInstAsstResultDetail struct {
	SuccessCreated map[int64]int64  `json:"success_created"`
	Error          map[int64]string `json:"error_msg"`
}

// CreateManyInstAsstResult  result of creating instance association
type CreateManyInstAsstResult struct {
	BaseResp `json:",inline"`
	Data     CreateManyInstAsstResultDetail `json:"data"`
}

// NewManyInstAsstResultDetail TODO
func NewManyInstAsstResultDetail() *CreateManyInstAsstResultDetail {
	return &CreateManyInstAsstResultDetail{
		SuccessCreated: make(map[int64]int64, 0),
		Error:          make(map[int64]string, 0),
	}
}

// DeleteAssociationInstRequest TODO
type DeleteAssociationInstRequest struct {
	Condition mapstr.MapStr `json:"condition"`
}

// DeleteAssociationInstBatchRequest TODO
type DeleteAssociationInstBatchRequest struct {
	ObjectID string  `json:"bk_obj_id"`
	ID       []int64 `json:"id"`
}

// DeleteAssociationInstResult TODO
type DeleteAssociationInstResult struct {
	BaseResp `json:",inline"`
	Data     string `json:"data"`
}

// DeleteAssociationInstBatchResult TODO
type DeleteAssociationInstBatchResult struct {
	BaseResp `json:",inline"`
	Data     int `json:"data"`
}

// AssociationKindIDs TODO
type AssociationKindIDs struct {
	// the association kind ids.
	AsstIDs []string `json:"asst_ids"`
}

// ListAssociationsWithAssociationKindResult TODO
type ListAssociationsWithAssociationKindResult struct {
	BaseResp `json:",inline"`
	Data     AssociationList `json:"data"`
}

// AssociationList TODO
type AssociationList struct {
	Associations []AssociationDetail `json:"associations"`
}

// AssociationDetail TODO
type AssociationDetail struct {
	// the ID of the association kind.
	AssociationKindID string        `json:"bk_asst_id"`
	Associations      []Association `json:"assts"`
}

// AssociationDirection 关联类型
type AssociationDirection string

const (
	// NoneDirection TODO
	NoneDirection AssociationDirection = "none"
	// DestinationToSource TODO
	DestinationToSource AssociationDirection = "src_to_dest"
	// SourceToDestination TODO
	SourceToDestination AssociationDirection = "dest_to_src"
	// Bidirectional TODO
	Bidirectional AssociationDirection = "bidirectional"
)

// AssociationKind TODO
type AssociationKind struct {
	ID int64 `field:"id" json:"id" bson:"id"`
	// a unique association id created by user.
	AssociationKindID string `field:"bk_asst_id" json:"bk_asst_id" bson:"bk_asst_id"`
	// a memorable name for this association kind, could be a chinese name, a english name etc.
	AssociationKindName string `field:"bk_asst_name" json:"bk_asst_name" bson:"bk_asst_name"`
	// the owner that this association type belongs to.
	OwnerID string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	// the describe for the relationship from source object to the target(destination) object, which will be displayed
	// when the topology is constructed between objects.
	SourceToDestinationNote string `field:"src_des" json:"src_des" bson:"src_des"`
	// the describe for the relationship from the target(destination) object to source object, which will be displayed
	// when the topology is constructed between objects.
	DestinationToSourceNote string `field:"dest_des" json:"dest_des" bson:"dest_des"`
	// the association direction between two objects.
	Direction AssociationDirection `field:"direction" json:"direction" bson:"direction"`
	// whether this is a pre-defined kind.
	IsPre *bool `field:"ispre" json:"ispre" bson:"ispre"`
}

// Parse TODO
func (cli *AssociationKind) Parse(data mapstr.MapStr) (*AssociationKind, error) {
	// TODO support parse metadata params
	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *AssociationKind) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(cli)
}

// AssociationOnDeleteAction TODO
type AssociationOnDeleteAction string

// AssociationMapping TODO
type AssociationMapping string

const (
	// NoAction TODO
	// this is a default action, which is do nothing when a association between object is deleted.
	NoAction AssociationOnDeleteAction = "none"
	// DeleteSource TODO
	// delete related source object instances when the association is deleted.
	DeleteSource AssociationOnDeleteAction = "delete_src"
	// DeleteDestinatioin TODO
	// delete related destination object instances when the association is deleted.
	DeleteDestinatioin AssociationOnDeleteAction = "delete_dest"

	// OneToOneMapping TODO
	// the source object can be related with only one destination object
	OneToOneMapping AssociationMapping = "1:1"
	// OneToManyMapping TODO
	// the source object can be related with multiple destination objects
	OneToManyMapping AssociationMapping = "1:n"
	// ManyToManyMapping TODO
	// multiple source object can be related with multiple destination objects
	ManyToManyMapping AssociationMapping = "n:n"
)

// Association defines the association between two objects.
type Association struct {
	ID      int64  `field:"id" json:"id" bson:"id"`
	OwnerID string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`

	// the unique id belongs to  this association, should be generated with rules as follows:
	// "$ObjectID"_"$AsstID"_"$AsstObjID"
	AssociationName string `field:"bk_obj_asst_id" json:"bk_obj_asst_id" bson:"bk_obj_asst_id"`
	// the alias name of this association, which is a substitute name in the association kind $AsstKindID
	AssociationAliasName string `field:"bk_obj_asst_name" json:"bk_obj_asst_name" bson:"bk_obj_asst_name"`

	// describe which object this association is defined for.
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	// describe where the Object associate with.
	AsstObjID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	// the association kind used by this association.
	AsstKindID string `field:"bk_asst_id" json:"bk_asst_id" bson:"bk_asst_id"`

	// defined which kind of association can be used between the source object and destination object.
	Mapping AssociationMapping `field:"mapping" json:"mapping" bson:"mapping"`
	// describe the action when this association is deleted.
	OnDelete AssociationOnDeleteAction `field:"on_delete" json:"on_delete" bson:"on_delete"`
	// describe whether this association is a pre-defined association or not,
	// if true, it means this association is used by cmdb itself.
	IsPre *bool `field:"ispre" json:"ispre" bson:"ispre"`
}

// CanUpdate TODO
// return field means which filed is set but is forbidden to update.
func (a *Association) CanUpdate() (field string, can bool) {
	if a.ID != 0 {
		return "id", false
	}

	if len(a.OwnerID) != 0 {
		return "bk_supplier_account", false
	}

	if len(a.AssociationName) != 0 {
		return "bk_obj_asst_id", false
	}

	if len(a.ObjectID) != 0 {
		return "bk_obj_id", false
	}

	if len(a.AsstObjID) != 0 {
		return "bk_asst_obj_id", false
	}

	if len(a.Mapping) != 0 {
		return "mapping", false
	}

	if a.IsPre != nil {
		return "ispre", false
	}

	// only on delete, association kind id, alias name can be update.
	return "", true
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *Association) Parse(data mapstr.MapStr) (*Association, error) {
	// TODO support parse metadata params
	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Association) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(cli)
}

// MainlineAssociation defines the mainline association between two objects.
type MainlineAssociation struct {
	Association `json:",inline"`

	ClassificationID string `json:"bk_classification_id,omitempty" bson:"-"`
	ObjectIcon       string `json:"bk_obj_icon,omitempty" bson:"-"`
	ObjectName       string `json:"bk_obj_name,omitempty" bson:"-"`
}

// InstAsst an association definition between instances.
type InstAsst struct {
	// sequence ID
	ID int64 `field:"id" json:"id,omitempty"`
	// inst id associate to ObjectID
	InstID int64 `field:"bk_inst_id" json:"bk_inst_id,omitempty" bson:"bk_inst_id"`
	// association source ObjectID
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id,omitempty" bson:"bk_obj_id"`
	// inst id associate to AsstObjectID
	AsstInstID int64 `field:"bk_asst_inst_id" json:"bk_asst_inst_id,omitempty"  bson:"bk_asst_inst_id"`
	// association target ObjectID
	AsstObjectID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id,omitempty" bson:"bk_asst_obj_id"`
	// bk_supplier_account
	OwnerID string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// association id between two object
	ObjectAsstID string `field:"bk_obj_asst_id" json:"bk_obj_asst_id,omitempty" bson:"bk_obj_asst_id"`
	// association kind id
	AssociationKindID string `field:"bk_asst_id" json:"bk_asst_id,omitempty" bson:"bk_asst_id"`

	// BizID the business ID
	BizID int64 `field:"bk_biz_id" json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
}

// GetInstID TODO
func (asst InstAsst) GetInstID(objID string) (instID int64, ok bool) {
	switch objID {
	case asst.ObjectID:
		return asst.InstID, true
	case asst.AsstObjectID:
		return asst.AsstInstID, true
	default:
		return 0, false
	}
}

// InstNameAsst TODO
type InstNameAsst struct {
	ID         string `json:"id"`
	ObjID      string `json:"bk_obj_id"`
	ObjIcon    string `json:"bk_obj_icon"`
	InstID     int64  `json:"bk_inst_id"`
	ObjectName string `json:"bk_obj_name"`
	InstName   string `json:"bk_inst_name"`
	AssoID     int64  `json:"asso_id"`
	// AsstName   string                 `json:"bk_asst_name"`
	// AsstID   string                 `json:"bk_asst_id"`
	InstInfo map[string]interface{} `json:"inst_info,omitempty"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *InstAsst) Parse(data mapstr.MapStr) (*InstAsst, error) {

	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *InstAsst) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(cli)
}

// MainlineObjectTopo the mainline object topo
type MainlineObjectTopo struct {
	ObjID      string `field:"bk_obj_id" json:"bk_obj_id"`
	ObjName    string `field:"bk_obj_name" json:"bk_obj_name"`
	OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account"`
	NextObj    string `field:"bk_next_obj" json:"bk_next_obj"`
	NextName   string `field:"bk_next_name" json:"bk_next_name"`
	PreObjID   string `field:"bk_pre_obj_id" json:"bk_pre_obj_id"`
	PreObjName string `field:"bk_pre_obj_name" json:"bk_pre_obj_name"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *MainlineObjectTopo) Parse(data mapstr.MapStr) (*MainlineObjectTopo, error) {

	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *MainlineObjectTopo) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(cli)
}

// TopoInst 实例拓扑结构
type TopoInst struct {
	InstID             int64  `json:"bk_inst_id"`
	InstName           string `json:"bk_inst_name"`
	ObjID              string `json:"bk_obj_id"`
	ObjName            string `json:"bk_obj_name"`
	Default            int    `json:"default"`
	ServiceTemplateID  int64  `json:"service_template_id,omitempty"`
	SetTemplateID      int64  `json:"set_template_id,omitempty"`
	HostApplyEnabled   *bool  `json:"host_apply_enabled,omitempty"`
	HostApplyRuleCount *int64 `json:"host_apply_rule_count,omitempty"`
}

// TopoInstRst 拓扑实例
type TopoInstRst struct {
	TopoInst `json:",inline"`
	Child    []*TopoInstRst `json:"child"`
}

// TopoNodeHostAndSerInstCount topo节点主机/服务实例数量
type TopoNodeHostAndSerInstCount struct {
	ObjID                string `json:"bk_obj_id"`
	InstID               int64  `json:"bk_inst_id"`
	HostCount            int64  `json:"host_count"`
	ServiceInstanceCount int64  `json:"service_instance_count"`
}

// HostAndSerInstCountOption 获取主机/服务实例查询参数结构
type HostAndSerInstCountOption struct {
	Condition []CountOptions `json:"condition"`
}

// CountOptions 获取主机/服务实例入参条件
type CountOptions struct {
	ObjID  string `json:"bk_obj_id"`
	InstID int64  `json:"bk_inst_id"`
}

// TopoInstRstVisitor TODO
type TopoInstRstVisitor func(tir *TopoInstRst)

// DeepFirstTraverse TODO
func (tir *TopoInstRst) DeepFirstTraverse(visitor TopoInstRstVisitor) {
	for _, child := range tir.Child {
		child.DeepFirstTraverse(visitor)
	}
	visitor(tir)
}

// ConditionItem subcondition
type ConditionItem struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// AssociationParams  association params
type AssociationParams struct {
	Page      BasePage                   `json:"page,omitempty"`
	Fields    map[string][]string        `json:"fields,omitempty"`
	Condition map[string][]ConditionItem `json:"condition,omitempty"`
}

// ResponeImportAssociation  import association result
type ResponeImportAssociation struct {
	BaseResp `json:",inline"`
	Data     ResponeImportAssociationData `json:"data"`
}

// RowMsgData TODO
type RowMsgData struct {
	Row int    `json:"row"`
	Msg string `json:"message"`
}

// ResponeImportAssociationData TODO
type ResponeImportAssociationData struct {
	ErrMsgMap []RowMsgData `json:"err_msg"`
}

// RequestImportAssociation  import association result
type RequestImportAssociation struct {
	AssociationInfoMap    map[int]ExcelAssociation `json:"association_info"`
	AsstObjectUniqueIDMap map[string]int64         `json:"asst_object_unique_id_info"`
	ObjectUniqueID        int64                    `json:"object_unique_id"`
}

// RequestInstAssociationObjectID 要求根据实例信息（实例的模型ID，实例ID）和模型ID（关联关系中的源，目的模型ID）, 返回关联关系的请求参数
type RequestInstAssociationObjectID struct {
	Condition RequestInstAssociationObjectIDCondition `json:"condition"`
	Page      BasePage                                `json:"page"`
}

// RequestInstAssociationObjectIDCondition  query condition
type RequestInstAssociationObjectIDCondition struct {
	// 实例得模型ID
	ObjectID string `json:"bk_obj_id"`
	// 实例ID
	InstID int64 `json:"bk_inst_id"`
	// ObjectID是否为目标模型， 默认false， 关联关系中的源模型，否则是目标模型
	IsTargetObject bool `json:"is_target_object"`

	// 关联对象的模型ID
	AssociationObjectID string `json:"association_obj_id"`
}

// InstBaseInfo instance base info
type InstBaseInfo struct {
	ID   int64  `json:"bk_inst_id"`
	Name string `json:"bk_inst_name"`
}

// FindTopoPathRequest TODO
type FindTopoPathRequest struct {
	Nodes []TopoNode `json:"topo_nodes" mapstructure:"topo_nodes"`
}

// TopoPathResult TODO
type TopoPathResult struct {
	Nodes []NodeTopoPath `json:"nodes" mapstructure:"nodes"`
}

// NodeTopoPath TODO
type NodeTopoPath struct {
	BizID int64                       `json:"bk_biz_id" mapstructure:"bk_biz_id"`
	Node  TopoNode                    `json:"topo_node" mapstructure:"topo_node"`
	Path  []*TopoInstanceNodeSimplify `json:"topo_path" mapstructure:"topo_path"`
}

// InstAsstQueryCondition TODO
type InstAsstQueryCondition struct {
	Cond  QueryCondition `json:"cond"`
	ObjID string         `json:"bk_obj_id"`
}

// InstAsstDeleteOption TODO
type InstAsstDeleteOption struct {
	Opt   DeleteOption `json:"opt"`
	ObjID string       `json:"bk_obj_id"`
}

// FindAssociationByObjectAssociationIDRequest 专用接口， 为excel 导入使用
type FindAssociationByObjectAssociationIDRequest struct {
	ObjAsstIDArr []string `json:"bk_obj_asst_ids"`
}

// FindAssociationByObjectAssociationIDResponse 专用接口， 为excel 导入使用
type FindAssociationByObjectAssociationIDResponse struct {
	BaseResp
	Data []Association `json:"data"`
}

// AssociationKindYaml yaml's field about association kind
type AssociationKindYaml struct {
	YamlHeader `json:",inline" yaml:",inline"`
	AsstKind   []AssociationKind `json:"asst_kind" yaml:"asst_kind"`
}

// Validate validate total yaml of association
func (asst *AssociationKindYaml) Validate() errors.RawErrorInfo {
	if err := asst.YamlHeader.Validate(); err.ErrCode != 0 {
		return err
	}

	for _, item := range asst.AsstKind {
		if len(item.AssociationKindID) == 0 {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrWebVerifyYamlFail,
				Args:    []interface{}{common.AssociationKindIDField + " not found"},
			}
		}

		if item.AssociationKindID == common.AssociationKindMainline {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrWebVerifyYamlFail,
				Args:    []interface{}{common.AssociationKindMainline + " forbidden operations"},
			}
		}
	}

	return errors.RawErrorInfo{}
}
