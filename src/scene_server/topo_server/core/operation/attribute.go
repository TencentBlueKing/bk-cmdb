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

package operation

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, metaData *metadata.Metadata) (model.AttributeInterface, error)
	DeleteObjectAttribute(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) error
	FindObjectAttributeWithDetail(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]*metadata.ObjAttDes, error)
	FindObjectAttribute(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.AttributeInterface, error)
	UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, bizID int64) error
	UpdateObjectAttributeIndex(kit *rest.Kit, objID string, data mapstr.MapStr, attID int64) (*metadata.UpdateAttrIndexData, error)

	FindBusinessAttribute(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Attribute, error)

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface, asst AssociationOperationInterface, grp GroupOperationInterface)
}

// NewAttributeOperation create a new attribute operation instance
func NewAttributeOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) AttributeOperationInterface {
	return &attribute{
		clientSet:   client,
		authManager: authManager,
	}
}

type attribute struct {
	clientSet    apimachinery.ClientSetInterface
	authManager  *extensions.AuthManager
	modelFactory model.Factory
	instFactory  inst.Factory
	obj          ObjectOperationInterface
	asst         AssociationOperationInterface
	grp          GroupOperationInterface
}

func (a *attribute) SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface, asst AssociationOperationInterface, grp GroupOperationInterface) {
	a.modelFactory = modelFactory
	a.instFactory = instFactory
	a.obj = obj
	a.asst = asst
	a.grp = grp
}

func (a *attribute) CreateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, metaData *metadata.Metadata) (model.AttributeInterface, error) {

	var err error
	att := a.modelFactory.CreateAttribute(kit)
	err = att.Parse(data)
	if nil != err {
		blog.Errorf("[operation-attr] failed to parse the attribute data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	// check if the object is mainline object, if yes. then user can not create required attribute.
	yes, err := a.isMainlineModel(kit.Header, att.Attribute().ObjectID)
	if err != nil {
		blog.Warnf("add object attribute, but not allow to add required attribute to mainline object: %+v. rid: %d.", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if yes {
		if att.Attribute().IsRequired {
			return nil, kit.CCError.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
		}
	}

	// check the object id
	objID := att.Attribute().ObjectID
	err = a.obj.IsValidObject(kit, objID, metaData)
	if nil != err {
		return nil, kit.CCError.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	// check is the group exist
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(att.Attribute().ObjectID)
	cond.Field(common.BKPropertyGroupIDField).Eq(att.Attribute().PropertyGroup)
	groupResult, err := a.grp.FindObjectGroup(kit, cond, metaData)
	if nil != err {
		blog.Errorf("[operation-attr] failed to search the attribute group data (%#v), error info is %s, rid: %s", cond.ToMapStr(), err.Error(), kit.Rid)
		return nil, err
	}
	// create the default group
	if 0 == len(groupResult) {
		if nil != metaData {
			cond := condition.CreateCondition()
			cond.Field(common.BKObjIDField).Eq(att.Attribute().ObjectID)
			cond.Field(common.BKPropertyGroupIDField).Eq(common.BKBizDefault)
			groupResult, err := a.grp.FindObjectGroup(kit, cond, metaData)
			if nil != err {
				blog.Errorf("[operation-attr] failed to search the attribute group data (%#v), error info is %s, rid: %s", cond.ToMapStr(), err.Error(), kit.Rid)
				return nil, err
			}
			if 0 == len(groupResult) {
				group := metadata.Group{
					IsDefault:  true,
					GroupIndex: -1,
					GroupName:  common.BKBizDefault,
					GroupID:    common.BKBizDefault,
					ObjectID:   att.Attribute().ObjectID,
					OwnerID:    att.Attribute().OwnerID,
				}
				group.Metadata = *metaData
				data := mapstr.NewFromStruct(group, "field")
				grp, err := a.grp.CreateObjectGroup(kit, data, metaData)
				if nil != err {
					blog.Errorf("[operation-obj] failed to create the default group, err: %s, rid: %s", err.Error(), kit.Rid)
					return nil, kit.CCError.Error(common.CCErrTopoObjectGroupCreateFailed)
				}
				//audit the CreateObjectGroup action
				err = NewObjectAttrGroupAudit(kit, a.clientSet, grp.Group().ID).buildSnapshotForPre().SaveAuditLog(metadata.AuditCreate)
				if err != nil {
					blog.Errorf("create object attribute group %s success, but update to auditLog failed, err: %v, rid: %s", grp.Group().GroupName, err, kit.Rid)
					return nil, err
				}
			}
			att.Attribute().PropertyGroup = common.BKBizDefault
		} else {
			att.Attribute().PropertyGroup = common.BKDefaultField
		}
	}

	// create a new one
	err = att.Create()
	if nil != err {
		blog.Errorf("[operation-attr] failed to save the attribute data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	//package audit response
	err = NewObjectAttrAudit(kit, a.clientSet, att.Attribute().ID).buildSnapshotForCur().SaveAuditLog(metadata.AuditCreate)
	if err != nil {
		blog.Errorf("create object attribute %s success, but update to auditLog failed, err: %v, rid: %s", att.Attribute().PropertyName, err, kit.Rid)
		return nil, err
	}

	return att, nil
}

func (a *attribute) DeleteObjectAttribute(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) error {

	attrItems, err := a.FindObjectAttribute(kit, cond, metaData)
	if nil != err {
		blog.Errorf("[operation-attr] failed to find the attributes by the cond(%v), err: %v, rid: %s", cond.ToMapStr(), err, kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectAttributeDeleteFailed, err.Error())
	}

	// auth: check authorization
	// var objID string
	// for idx, attrItem := range attrItems {
	// 	oID := attrItem.Attribute().ObjectID
	// 	if idx == 0 && objID != oID {
	// 		return kit.CCError.New(common.CCErrTopoObjectAttributeDeleteFailed, "can't attributes of multiple model per request")
	// 	}
	// }
	// if err := a.authManager.AuthorizeByObjectID(kit.Ctx, kit.Header, meta.Update, objID); err != nil {
	// 	return kit.CCError.New(common.CCErrCommAuthorizeFailed, err.Error())
	// }

	for _, attrItem := range attrItems {
		//get PreData
		objAudit := NewObjectAttrAudit(kit, a.clientSet, attrItem.Attribute().ID).buildSnapshotForPre()

		// delete the attribute
		rsp, err := a.clientSet.CoreService().Model().DeleteModelAttr(context.Background(), kit.Header, attrItem.Attribute().ObjectID, &metadata.DeleteOption{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[operation-attr] delete object attribute failed, request object controller with err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-attr] failed to delete the attribute by condition(%v), err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, kit.Rid)
			return kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		//saveAuditLog
		err = objAudit.SaveAuditLog(metadata.AuditDelete)
		if err != nil {
			blog.Errorf("Delete object attribute success, but update to auditLog failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	return nil
}

func (a *attribute) FindObjectAttributeWithDetail(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]*metadata.ObjAttDes, error) {
	attrs, err := a.FindObjectAttribute(kit, cond, metaData)
	if nil != err {
		blog.ErrorJSON("FindObjectAttribute failed, err: %s, cond: %s", err, cond.ToMapStr())
		return nil, err
	}
	results := make([]*metadata.ObjAttDes, 0)
	// if can't find any attribute of a obj, to return, for example, when the obj is not exist
	if len(attrs) == 0 {
		return results, nil
	}
	grpCond := condition.CreateCondition()
	grpOrCond := grpCond.NewOR()
	for _, attr := range attrs {
		attribute := attr.Attribute()
		grpOrCond.Item(map[string]interface{}{
			metadata.GroupFieldGroupID:  attribute.PropertyGroup,
			metadata.GroupFieldObjectID: attribute.ObjectID,
		})
	}
	grps, err := a.grp.FindObjectGroup(kit, grpCond, metaData)
	if nil != err {
		blog.ErrorJSON("FindObjectGroup failed, err: %s, grpCond: %s", err, grpCond.ToMapStr())
		return nil, err
	}
	grpMap := make(map[string]string)
	for _, grp := range grps {
		grpMap[grp.Group().GroupID] = grp.Group().GroupName
	}
	for _, attr := range attrs {
		attribute := attr.Attribute()
		result := &metadata.ObjAttDes{Attribute: *attribute}
		grpName, ok := grpMap[attribute.PropertyGroup]
		if !ok {
			blog.ErrorJSON("attribute [%s] has an invalid bk_property_group %s", *attribute, attribute.PropertyGroup)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attribute.bk_property_group: "+attribute.PropertyGroup)
		}
		result.PropertyGroupName = grpName
		results = append(results, result)
	}

	return results, nil
}

func (a *attribute) FindBusinessAttribute(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Attribute, error) {
	opt := &metadata.QueryCondition{
		Condition: cond,
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDApp, opt)
	if err != nil {
		blog.Errorf("find business attributes failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf("find business attributes failed, err: %s rid: %s", resp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data.Info, nil
}

func (a *attribute) FindObjectAttribute(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.AttributeInterface, error) {
	limits := cond.GetLimit()
	sort := cond.GetSort()
	start := cond.GetStart()
	fCond := cond.ToMapStr()
	if nil != metaData {
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	opt := &metadata.QueryCondition{
		Condition: fCond,
		Page:      metadata.BasePage{Limit: int(limits), Start: int(start), Sort: sort},
	}

	rsp, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(context.Background(), kit.Header, opt)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to search attribute by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return model.CreateAttribute(kit, a.clientSet, rsp.Data.Info), nil
}

func (a *attribute) UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, bizID int64) error {
	// TODO replace this logic with cond := metadata.NewPublicOrBizConditionByBizID(bizID) when old interface can't operate biz custom field
	var cond map[string]interface{}
	if bizID == 0 {
		cond = make(map[string]interface{}, 0)
	} else {
		cond = metadata.NewPublicOrBizConditionByBizID(bizID)
	}
	cond[common.BKFieldID] = attID
	input := metadata.UpdateOption{
		Condition: cond,
		Data:      data,
	}

	//get PreData
	objAudit := NewObjectAttrAudit(kit, a.clientSet, attID).buildSnapshotForPre()

	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsByCondition(context.Background(), kit.Header, &input)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to update the attribute by the attr-id(%d), error info is %s, rid: %s", attID, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	//get CurData and saveAuditLog
	err = objAudit.buildSnapshotForCur().SaveAuditLog(metadata.AuditUpdate)
	if err != nil {
		blog.Errorf("update object attribute-id %s success, but update to auditLog failed, err: %v, rid: %s", attID, err, kit.Rid)
		return err
	}

	return nil
}

func (a *attribute) isMainlineModel(head http.Header, modelID string) (bool, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := a.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), head,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return false, err
	}

	if !asst.Result {
		return false, errors.New(asst.Code, asst.ErrMsg)
	}

	if len(asst.Data.Info) <= 0 {
		return false, fmt.Errorf("model association [%+v] not found", cond)
	}

	for _, mainline := range asst.Data.Info {
		if mainline.ObjectID == modelID {
			return true, nil
		}
	}

	return false, nil
}

func (a *attribute) UpdateObjectAttributeIndex(kit *rest.Kit, objID string, data mapstr.MapStr, attID int64) (*metadata.UpdateAttrIndexData, error) {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(attID).ToMapStr(),
		Data:      data,
	}

	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsIndex(context.Background(), kit.Header, objID, &input)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object CoreService, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to update the attribute index by the attr-id(%d), error info is %s, rid: %s", attID, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}
