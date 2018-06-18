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
	"context"
	"encoding/json"
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"

	"configcenter/src/scene_server/topo_server/core/types"
)

var _ Object = (*object)(nil)

type object struct {
	obj       meta.Object
	isNew     bool
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *object) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.obj)
}

func (cli *object) GetInstIDFieldName() string {

	switch cli.obj.ObjectID {
	case common.BKInnerObjIDApp:
		return common.BKAppIDField
	case common.BKInnerObjIDSet:
		return common.BKSetIDField
	case common.BKInnerObjIDModule:
		return common.BKModuleIDField
	case common.BKINnerObjIDObject:
		return common.BKInstIDField
	case common.BKInnerObjIDHost:
		return common.BKHostIDField
	case common.BKInnerObjIDProc:
		return common.BKProcIDField
	case common.BKInnerObjIDPlat:
		return common.BKCloudIDField
	default:
		return common.BKInstIDField
	}

}

func (cli *object) GetInstNameFieldName() string {
	switch cli.obj.ObjectID {
	case common.BKInnerObjIDApp:
		return common.BKAppNameField
	case common.BKInnerObjIDSet:
		return common.BKSetNameField
	case common.BKInnerObjIDModule:
		return common.BKModuleNameField
	case common.BKInnerObjIDHost:
		return common.BKHostInnerIPField
	case common.BKInnerObjIDProc:
		return common.BKProcNameField
	case common.BKInnerObjIDPlat:
		return common.BKCloudNameField
	default:
		return common.BKInstNameField
	}
}

func (cli *object) GetObjectType() string {
	switch cli.obj.ObjectID {
	case common.BKInnerObjIDApp:
		return cli.obj.ObjectID
	case common.BKInnerObjIDSet:
		return cli.obj.ObjectID
	case common.BKInnerObjIDModule:
		return cli.obj.ObjectID
	case common.BKInnerObjIDHost:
		return cli.obj.ObjectID
	case common.BKInnerObjIDProc:
		return cli.obj.ObjectID
	case common.BKInnerObjIDPlat:
		return cli.obj.ObjectID
	default:
		return common.BKINnerObjIDObject
	}
}

func (cli *object) search() ([]meta.Object, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(cli.params.Header.OwnerID).Field(common.BKObjIDField).Eq(cli.obj.ObjectID)

	condStr, err := cond.ToMapStr().ToJSON()
	if nil != err {
		return nil, err
	}
	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjects(context.Background(), cli.params.Header.ToHeader(), condStr)

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.obj.ObjectID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	return rsp.Data, nil

}

func (cli *object) IsExists() ([]meta.Object, bool, error) {

	items, err := cli.search()
	if nil != err {
		return nil, false, err
	}

	return items, 0 != len(items), nil
}

func (cli *object) Create() error {

	rsp, err := cli.clientSet.ObjectController().Meta().CreateObject(context.Background(), cli.params.Header.ToHeader(), &cli.obj)

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.obj.ObjectID, rsp.ErrMsg)
		return cli.params.Err.Error(rsp.Code)
	}

	cli.obj.ID = rsp.Data.ID

	return nil
}

func (cli *object) Delete() error {
	rsp, err := cli.clientSet.ObjectController().Meta().DeleteObject(context.Background(), cli.obj.ID, cli.params.Header.ToHeader(), nil)

	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[opration-obj] failed to delete the object by the id(%d)", cli.obj.ID)
		return cli.params.Err.Error(rsp.Code)
	}
	return nil
}

func (cli *object) Update() error {

	data := meta.SetValueToMapStrByTags(cli)

	items, err := cli.search()
	if nil != err {
		return err
	}

	for _, item := range items {

		rsp, err := cli.clientSet.ObjectController().Meta().UpdateObject(context.Background(), item.ID, cli.params.Header.ToHeader(), data)

		if nil != err {
			blog.Errorf("failed to request the object controller, error info is %s", err.Error())
			return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("failed to search the object(%s), error info is %s", cli.obj.ObjectID, rsp.ErrMsg)
			return cli.params.Err.Error(rsp.Code)
		}
	}
	return nil
}

func (cli *object) Parse(data frtypes.MapStr) (*meta.Object, error) {

	err := meta.SetValueToStructByTags(&cli.obj, data)
	if nil != err {
		return nil, err
	}

	if 0 == len(cli.obj.ObjectID) {
		return nil, cli.params.Err.Errorf(common.CCErrCommParamsNeedSet, meta.ModelFieldObjectID)
	}

	if 0 == len(cli.obj.ObjCls) {
		return nil, cli.params.Err.Errorf(common.CCErrCommParamsNeedSet, meta.ModelFieldObjCls)
	}

	return nil, err
}

func (cli *object) ToMapStr() (frtypes.MapStr, error) {
	rst := meta.SetValueToMapStrByTags(&cli.obj)
	return rst, nil
}

func (cli *object) Save() error {

	if _, exists, err := cli.IsExists(); nil != err {
		return err
	} else if !exists {
		return cli.Create()
	}

	return cli.Update()

}

func (cli *object) CreateGroup() Group {
	return &group{
		grp: meta.Group{
			OwnerID:  cli.obj.OwnerID,
			ObjectID: cli.obj.ObjectID,
		},
	}
}

func (cli *object) CreateAttribute() Attribute {
	return &attribute{
		attr: meta.Attribute{
			OwnerID:  cli.obj.OwnerID,
			ObjectID: cli.obj.ObjectID,
		},
	}
}

func (cli *object) GetAttributes() ([]Attribute, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AttributeFieldObjectID).Eq(cli.obj.ObjectID).Field(meta.AttributeFieldSupplierAccount).Eq(cli.params.Header.OwnerID)
	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.obj.ObjectID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Attribute, 0)
	for _, item := range rsp.Data {

		attr := &attribute{
			attr:      item,
			params:    cli.params,
			clientSet: cli.clientSet,
		}

		rstItems = append(rstItems, attr)
	}

	return rstItems, nil
}

func (cli *object) GetGroups() ([]Group, error) {

	cond := condition.CreateCondition()

	cond.Field(meta.GroupFieldObjectID).Eq(cli.obj.ObjectID).Field(meta.GroupFieldSupplierAccount).Eq(cli.params.Header.OwnerID)
	rsp, err := cli.clientSet.ObjectController().Meta().SelectGroup(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.obj.ObjectID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Group, 0)
	for _, item := range rsp.Data {
		grp := &group{
			grp:       item,
			params:    cli.params,
			clientSet: cli.clientSet,
		}
		rstItems = append(rstItems, grp)
	}

	return rstItems, nil
}

func (cli *object) SetClassification(class Classification) {
	cli.obj.ObjCls = class.GetID()
}

func (cli *object) GetClassification() (Classification, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.ClassFieldClassificationID).Eq(cli.obj.ObjCls)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectClassifications(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", cli.obj.ObjectID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	for _, item := range rsp.Data {

		return &classification{
			cls:       item,
			params:    cli.params,
			clientSet: cli.clientSet,
		}, nil // only one classification
	}

	return nil, fmt.Errorf("invalid classification(%s) for the object(%s)", cli.obj.ObjCls, cli.obj.ObjectID)
}

func (cli *object) SetIcon(objectIcon string) {
	cli.obj.ObjIcon = objectIcon
}

func (cli *object) GetIcon() string {
	return cli.obj.ObjIcon
}

func (cli *object) SetID(objectID string) {
	cli.obj.ObjectID = objectID
}

func (cli *object) GetID() string {
	return cli.obj.ObjectID
}

func (cli *object) SetName(objectName string) {
	cli.obj.ObjectName = objectName
}

func (cli *object) GetName() string {
	return cli.obj.ObjectName
}

func (cli *object) SetIsPre(isPre bool) {
	cli.obj.IsPre = isPre
}

func (cli *object) GetIsPre() bool {
	return cli.obj.IsPre
}

func (cli *object) SetIsPaused(isPaused bool) {
	cli.obj.IsPaused = isPaused
}

func (cli *object) GetIsPaused() bool {
	return cli.obj.IsPaused
}

func (cli *object) SetPosition(position string) {
	cli.obj.Position = position
}

func (cli *object) GetPosition() string {
	return cli.obj.Position
}

func (cli *object) SetSupplierAccount(supplierAccount string) {
	cli.obj.OwnerID = supplierAccount
}

func (cli *object) GetSupplierAccount() string {
	return cli.obj.OwnerID
}

func (cli *object) SetDescription(description string) {
	cli.obj.Description = description
}

func (cli *object) GetDescription() string {
	return cli.obj.Description
}

func (cli *object) SetCreator(creator string) {
	cli.obj.Creator = creator
}

func (cli *object) GetCreator() string {
	return cli.obj.Creator
}

func (cli *object) SetModifier(modifier string) {
	cli.obj.Modifier = modifier
}

func (cli *object) GetModifier() string {
	return cli.obj.Modifier
}
