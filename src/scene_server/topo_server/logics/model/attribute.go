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

package model

import (
	"fmt"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, modelBizID int64) (*metadata.Attribute, error)
	DeleteObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) error
	UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64) error
	UpdateObjectAttributeIndex(kit *rest.Kit, objID string, data mapstr.MapStr,
		attID int64) (*metadata.UpdateAttrIndexData, error)
	FindObjectAttributeWithDetail(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) ([]*metadata.ObjAttDes, error)
	FindObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) ([]metadata.Attribute, error)
	FindBusinessAttribute(kit *rest.Kit, cond mapstr.MapStr, objID string) ([]metadata.Attribute, error)
}

// NewAttributeOperation create a new attribute operation instance
func NewAttributeOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) AttributeOperationInterface {
	return &attribute{
		clientSet:   client,
		authManager: authManager,
	}
}

type attribute struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// FindObject  find object by condition
func (a *attribute) FindObject(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Object, error) {
	rsp, err := a.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	return rsp.Data.Info, nil
}

// IsValidObject check whether objID is a real model's bk_obj_id field in backend
func (a *attribute) IsValidObject(kit *rest.Kit, objID string) error {

	checkObjCond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}

	objItems, err := a.FindObject(kit, checkObjCond)
	if err != nil {
		blog.Errorf("failed to check the object repeated, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, err.Error())
	}

	if len(objItems) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return nil
}

// FindObjectGroup find object group data by condition
func (a *attribute) FindObjectGroup(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) ([]metadata.Group, error) {
	util.AddModelBizIDCondition(cond, modelBizID)

	rsp, err := a.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header,
		metadata.QueryCondition{Condition: cond})
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[opeartion-grp] failed to search the group by the condition(%#v), error info is %s, rid: %s", cond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

// CreateObjectGroup create object group
func (a *attribute) CreateObjectGroup(kit *rest.Kit, data mapstr.MapStr) (*metadata.Group, error) {
	grp := metadata.Group{}

	err := mapstr.SetValueToStructByTags(&grp, data)
	if nil != err {
		blog.Errorf("[operation-grp] failed to parse the group data(%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	//  check the object
	if err = a.isObjectInGroupValidObject(kit, grp.ObjectID); nil != err {
		blog.Errorf("[operation-grp] the group (%#v) is in valid, rid: %s", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// create a new group
	rsp, err := a.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header, grp.ObjectID,
		metadata.CreateModelAttributeGroup{Data: grp})
	if nil != err {
		blog.Errorf("[operation-grp] failed to save the group data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}
	if !rsp.Result {
		blog.Errorf("[model-grp] failed to create the group(%s), err: is %s, rid: %s", grp.GroupID, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectGroupCreateFailed)
	}

	grp.ID = int64(rsp.Data.Created.ID)

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, grp.ID, &grp)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but generate audit log failed, err: %v, rid: %s", grp.GroupName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute group %s success, but save audit log failed, err: %v, rid: %s", grp.GroupName, err, kit.Rid)
		return nil, err
	}

	return &grp, nil
}

// isObjectInGroupValidObject is object in group valid object
func (a *attribute) isObjectInGroupValidObject(kit *rest.Kit, objectID string) error {
	// check source object exists
	objRsp, err := a.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: objectID},
	})
	if err != nil {
		blog.Errorf("read the object(%s) failed, err: %s, rid: %s", objectID, err.Error(), kit.Rid)
		return err
	}

	if !objRsp.Result {
		blog.Errorf("read the object(%s) failed, err: %s, rid: %s", objectID, objRsp.ErrMsg, kit.Rid)
		return err
	}

	if len(objRsp.Data.Info) == 0 {
		blog.Errorf("the object(%s) is invalid, return is empty, rid: %s", objectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return nil
}

// CreateObjectAttribute create object attribute
func (a *attribute) CreateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, modelBizID int64) (*metadata.Attribute,
	error) {
	attr := &metadata.Attribute{}
	err := mapstr.SetValueToStructByTags(&attr, data)
	if nil != err {
		blog.Errorf("[operation-attr] failed to parse the attr data(%#v), error info is %s, rid: %s", data,
			err.Error(), kit.Rid)
		return nil, err
	}

	// check if the object is mainline object, if yes. then user can not create required attribute.
	yes, err := a.isMainlineModel(kit, attr.ObjectID)
	if err != nil {
		blog.Warnf("add object attribute, but not allow to add required attribute to mainline object: %+v. rid: %d.", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if yes {
		if attr.IsRequired {
			return nil, kit.CCError.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
		}
	}

	// check the object id
	objID := attr.ObjectID
	err = a.IsValidObject(kit, objID)
	if nil != err {
		return nil, kit.CCError.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	cond := mapstr.MapStr{common.BKObjIDField: attr.ObjectID, common.BKPropertyGroupIDField: attr.PropertyGroup}
	groupResult, err := a.FindObjectGroup(kit, cond, modelBizID)
	if nil != err {
		blog.Errorf("[operation-attr] failed to search the attribute group data (%#v), error info is %s, rid: %s", cond, err.Error(), kit.Rid)
		return nil, err
	}
	// create the default group
	if 0 == len(groupResult) {
		if modelBizID > 0 {
			cond := mapstr.MapStr{common.BKObjIDField: attr.ObjectID, common.BKPropertyGroupIDField: common.BKBizDefault}
			groupResult, err := a.FindObjectGroup(kit, cond, modelBizID)
			if nil != err {
				blog.Errorf("[operation-attr] failed to search the attribute group data (%#v), error info is %s, rid: %s", cond, err.Error(), kit.Rid)
				return nil, err
			}
			if 0 == len(groupResult) {
				group := metadata.Group{
					IsDefault:  true,
					GroupIndex: -1,
					GroupName:  common.BKBizDefault,
					GroupID:    common.BKBizDefault,
					ObjectID:   attr.ObjectID,
					OwnerID:    attr.OwnerID,
					BizID:      modelBizID,
				}

				data := mapstr.NewFromStruct(group, "field")
				_, err := a.CreateObjectGroup(kit, data)
				if nil != err {
					blog.Errorf("[operation-obj] failed to create the default group, err: %s, rid: %s", err.Error(), kit.Rid)
					return nil, kit.CCError.Error(common.CCErrTopoObjectGroupCreateFailed)
				}
			}
			attr.PropertyGroup = common.BKBizDefault
		} else {
			attr.PropertyGroup = common.BKDefaultField
		}
	}

	input := metadata.CreateModelAttributes{Attributes: []metadata.Attribute{*attr}}
	rsp, err := a.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, attr.ObjectID, &input)
	if nil != err {
		blog.ErrorJSON("failed to request coreService to create model attrs, the err: %s, ObjectID: %s, input: %s, rid: %s", err.Error(), attr.ObjectID, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.ErrorJSON("create model attrs failed, ObjectID: %s, input: %s, rid: %s", attr.ObjectID, input, kit.Rid)
		return nil, rsp.CCError()
	}

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attr.ID, nil)
	if err != nil {
		blog.Errorf("create object attribute %s success, but generate audit log failed, err: %v, rid: %s",
			attr.PropertyName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute %s success, but save audit log failed, err: %v, rid: %s",
			attr.PropertyName, err, kit.Rid)
		return nil, err
	}

	return attr, nil
}

// DeleteObjectAttribute delete object attribute
func (a *attribute) DeleteObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) error {

	attrItems, err := a.FindObjectAttribute(kit, cond, modelBizID)
	if nil != err {
		blog.Errorf("[operation-attr] failed to find the attributes by the cond(%v), err: %v, rid: %s", cond, err, kit.Rid)
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

	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)

	for _, attrItem := range attrItems {
		// generate audit log of model attribute.
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attrItem.ID, &attrItem)
		if err != nil {
			blog.Errorf("generate audit log failed before delete model attribute %s, err: %v, rid: %s",
				attrItem.PropertyName, err, kit.Rid)
			return err
		}

		// delete the attribute.
		rsp, err := a.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, attrItem.ObjectID, &metadata.DeleteOption{Condition: cond})
		if nil != err {
			blog.Errorf("[operation-attr] delete object attribute failed, request object controller with err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rsp.Result {
			blog.Errorf("[operation-attr] failed to delete the attribute by condition(%v), err: %s, rid: %s", cond, rsp.ErrMsg, kit.Rid)
			return kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}

		// save audit log.
		if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
			blog.Errorf("delete object attribute %s success, but save audit log failed, err: %v, rid: %s",
				attrItem.PropertyName, err, kit.Rid)
			return err
		}
	}

	return nil
}

// UpdateObjectAttribute update object attribute
func (a *attribute) UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64) error {
	cond := make(map[string]interface{})
	util.AddModelBizIDCondition(cond, modelBizID)

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update model attribute, attID: %d, err: %v, rid: %s",
			attID, err, kit.Rid)
		return err
	}

	// to update.
	cond[common.BKFieldID] = attID
	input := metadata.UpdateOption{
		Condition: cond,
		Data:      data,
	}
	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsByCondition(kit.Ctx, kit.Header, &input)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to update the attribute by the attr-id(%d), error info is %s, rid: %s", attID, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object attribute success, but save audit log failed, attID: %d, err: %v, rid: %s",
			attID, err, kit.Rid)
		return err
	}

	return nil
}

// UpdateObjectAttributeIndex update object attribute index by obj id
func (a *attribute) UpdateObjectAttributeIndex(kit *rest.Kit, objID string, data mapstr.MapStr,
	attID int64) (*metadata.UpdateAttrIndexData, error) {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(attID).ToMapStr(),
		Data:      data,
	}

	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsIndex(kit.Ctx, kit.Header, objID, &input)
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

// isMainlineModel check is mainline model by module id
func (a *attribute) isMainlineModel(kit *rest.Kit, modelID string) (bool, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := a.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
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

// FindObjectAttributeWithDetail find object attribute detail by condition
func (a *attribute) FindObjectAttributeWithDetail(kit *rest.Kit, cond mapstr.MapStr,
	modelBizID int64) ([]*metadata.ObjAttDes, error) {
	attrs, err := a.FindObjectAttribute(kit, cond, modelBizID)
	if nil != err {
		blog.ErrorJSON("FindObjectAttribute failed, err: %s, cond: %s", err, cond)
		return nil, err
	}
	results := make([]*metadata.ObjAttDes, 0)
	// if can't find any attribute of a obj, to return, for example, when the obj is not exist
	if len(attrs) == 0 {
		return results, nil
	}

	grps := make([]metadata.Group, 0)
	for _, attr := range attrs {
		cond := mapstr.MapStr{metadata.GroupFieldGroupID: attr.PropertyGroup}
		grp, err := a.FindObjectGroup(kit, cond, modelBizID)
		if nil != err {
			blog.ErrorJSON("FindObjectGroup failed, err: %s, grpCond: %s", err, cond)
			return nil, err
		}
		grps = append(grps, grp...)
	}

	grpMap := make(map[string]string)
	for _, grp := range grps {
		grpMap[grp.GroupID] = grp.GroupName
	}
	for _, attr := range attrs {
		result := &metadata.ObjAttDes{Attribute: attr}
		grpName, ok := grpMap[attr.PropertyGroup]
		if !ok {
			blog.ErrorJSON("attribute [%s] has an invalid bk_property_group %s", attr, attr.PropertyGroup)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attribute.bk_property_group: "+attr.PropertyGroup)
		}
		result.PropertyGroupName = grpName
		results = append(results, result)
	}

	return results, nil
}

// FindObjectAttribute find object attribute
func (a *attribute) FindObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) ([]metadata.Attribute, error) {
	util.AddModelBizIDCondition(cond, modelBizID)

	opt := &metadata.QueryCondition{
		Condition: cond,
	}

	rsp, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, opt)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to search attribute by the condition(%#v), error info is %s, rid: %s",
			cond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

// FindBusinessAttribute find attribute by obj id
func (a *attribute) FindBusinessAttribute(kit *rest.Kit, cond mapstr.MapStr, objID string) ([]metadata.Attribute,
	error) {
	opt := &metadata.QueryCondition{
		Condition: cond,
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, opt)
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
