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

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(params types.ContextParams, data mapstr.MapStr) (model.AttributeInterface, error)
	DeleteObjectAttribute(params types.ContextParams, cond condition.Condition) error
	FindObjectAttributeWithDetail(params types.ContextParams, cond condition.Condition) ([]*metadata.ObjAttDes, error)
	FindObjectAttribute(params types.ContextParams, cond condition.Condition) ([]model.AttributeInterface, error)
	UpdateObjectAttribute(params types.ContextParams, data mapstr.MapStr, attID int64) error
	UpdateObjectAttributeIndex(params types.ContextParams, objID string, data mapstr.MapStr, attID int64) (*metadata.UpdateAttrIndexData, error)

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

func (a *attribute) CreateObjectAttribute(params types.ContextParams, data mapstr.MapStr) (model.AttributeInterface, error) {

	var err error
	att := a.modelFactory.CreateAttribute(params)
	err = att.Parse(data)
	if nil != err {
		blog.Errorf("[operation-attr] failed to parse the attribute data (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}

	// check if the object is mainline object, if yes. then user can not create required attribute.
	yes, err := a.isMainlineModel(params.Header, att.Attribute().ObjectID)
	if err != nil {
		blog.Warnf("add object attribute, but not allow to add required attribute to mainline object: %+v. rid: %d.", data, params.ReqID)
		return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if yes {
		if att.Attribute().IsRequired {
			return nil, params.Err.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
		}
	}

	// check the object id
	objID := att.Attribute().ObjectID
	err = a.obj.IsValidObject(params, objID)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	// check is the group exist
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(att.Attribute().ObjectID)
	cond.Field(common.BKPropertyGroupIDField).Eq(att.Attribute().PropertyGroup)
	groupResult, err := a.grp.FindObjectGroup(params, cond)
	if nil != err {
		blog.Errorf("[operation-attr] failed to search the attribute group data (%#v), error info is %s, rid: %s", cond.ToMapStr(), err.Error(), params.ReqID)
		return nil, err
	}
	// create the default group
	if 0 == len(groupResult) {
		if nil != params.MetaData {
			cond := condition.CreateCondition()
			cond.Field(common.BKObjIDField).Eq(att.Attribute().ObjectID)
			cond.Field(common.BKPropertyGroupIDField).Eq(common.BKBizDefault)
			groupResult, err := a.grp.FindObjectGroup(params, cond)
			if nil != err {
				blog.Errorf("[operation-attr] failed to search the attribute group data (%#v), error info is %s, rid: %s", cond.ToMapStr(), err.Error(), params.ReqID)
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
				group.Metadata = *params.MetaData
				data := mapstr.NewFromStruct(group, "field")
				if _, err := a.grp.CreateObjectGroup(params, data); nil != err {
					blog.Errorf("[operation-obj] failed to create the default group, err: %s, rid: %s", err.Error(), params.ReqID)
					return nil, params.Err.Error(common.CCErrTopoObjectGroupCreateFailed)
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
		blog.Errorf("[operation-attr] failed to save the attribute data (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, err
	}

	return att, nil
}

func (a *attribute) DeleteObjectAttribute(params types.ContextParams, cond condition.Condition) error {

	attrItems, err := a.FindObjectAttribute(params, cond)
	if nil != err {
		blog.Errorf("[operation-attr] failed to find the attributes by the cond(%v), err: %v, rid: %s", cond.ToMapStr(), err, params.ReqID)
		return params.Err.New(common.CCErrTopoObjectAttributeDeleteFailed, err.Error())
	}

	// auth: check authorization
	// var objID string
	// for idx, attrItem := range attrItems {
	// 	oID := attrItem.Attribute().ObjectID
	// 	if idx == 0 && objID != oID {
	// 		return params.Err.New(common.CCErrTopoObjectAttributeDeleteFailed, "can't attributes of multiple model per request")
	// 	}
	// }
	// if err := a.authManager.AuthorizeByObjectID(params.Context, params.Header, meta.Update, objID); err != nil {
	// 	return params.Err.New(common.CCErrCommAuthorizeFailed, err.Error())
	// }

	for _, attrItem := range attrItems {
		// delete the attribute
		rsp, err := a.clientSet.CoreService().Model().DeleteModelAttr(context.Background(), params.Header, attrItem.Attribute().ObjectID, &metadata.DeleteOption{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[operation-attr] delete object attribute failed, request object controller with err: %v, rid: %s", err, params.ReqID)
			return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-attr] failed to delete the attribute by condition(%v), err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, params.ReqID)
			return params.Err.New(rsp.Code, rsp.ErrMsg)
		}
	}

	return nil
}

func (a *attribute) FindObjectAttributeWithDetail(params types.ContextParams, cond condition.Condition) ([]*metadata.ObjAttDes, error) {
	attrs, err := a.FindObjectAttribute(params, cond)
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
			metadata.GroupFieldGroupID:         attribute.PropertyGroup,
			metadata.GroupFieldObjectID:        attribute.ObjectID,
		})
	}
	grps, err := a.grp.FindObjectGroup(params, grpCond)
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
			return nil, params.Err.CCErrorf(common.CCErrCommParamsIsInvalid, "attribute.bk_property_group: "+attribute.PropertyGroup)
		}
		result.PropertyGroupName = grpName
		results = append(results, result)
	}

	return results, nil
}

func (a *attribute) FindObjectAttribute(params types.ContextParams, cond condition.Condition) ([]model.AttributeInterface, error) {
	limits := cond.GetLimit()
	sort := cond.GetSort()
	start := cond.GetStart()
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	opt := &metadata.QueryCondition{
		Condition: fCond,
		Limit:     metadata.SearchLimit{Limit: limits, Offset: start},
		SortArr:   metadata.NewSearchSortParse().String(sort).ToSearchSortArr(),
	}

	rsp, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(context.Background(), params.Header, opt)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to search attribute by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return model.CreateAttribute(params, a.clientSet, rsp.Data.Info), nil
}

func (a *attribute) UpdateObjectAttribute(params types.ContextParams, data mapstr.MapStr, attID int64) error {

	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(attID).ToMapStr(),
		Data:      data,
	}

	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsByCondition(context.Background(), params.Header, &input)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to update the attribute by the attr-id(%d), error info is %s, rid: %s", attID, rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
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

func (a *attribute) UpdateObjectAttributeIndex(params types.ContextParams, objID string, data mapstr.MapStr, attID int64) (*metadata.UpdateAttrIndexData, error) {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(attID).ToMapStr(),
		Data:      data,
	}

	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsIndex(context.Background(), params.Header, objID, &input)
	if nil != err {
		blog.Errorf("[operation-attr] failed to request object CoreService, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-attr] failed to update the attribute index by the attr-id(%d), error info is %s, rid: %s", attID, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}
