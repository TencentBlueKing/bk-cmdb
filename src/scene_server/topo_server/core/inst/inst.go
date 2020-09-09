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

package inst

import (
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/model"
)

// Inst the inst interface
type Inst interface {
	model.Operation
	GetObject() model.Object

	GetMainlineParentInst() (Inst, error)
	GetMainlineChildInst() ([]Inst, error)

	GetParentObjectWithInsts() ([]*ObjectWithInsts, error)
	GetChildObjectWithInsts() ([]*ObjectWithInsts, error)

	SetMainlineParentInst(instID int64) error
	SetMainlineChildInst(targetInst Inst) error

	GetInstID() (int64, error)
	GetParentID() (int64, error)
	GetInstName() (string, error)

	SetAssoID(id int64)
	GetAssoID() int64

	SetValue(key string, value interface{})

	SetValues(values mapstr.MapStr)

	GetValues() mapstr.MapStr

	ToMapStr() mapstr.MapStr

	IsDefault() bool

	GetBizID() (int64, error)

	CheckInstanceExists(nonInnerAttributes []model.AttributeInterface) (exist bool, filter condition.Condition, err error)
	UpdateInstance(filter condition.Condition, data mapstr.MapStr, nonInnerAttributes []model.AttributeInterface) error
}

var _ Inst = (*inst)(nil)

type inst struct {
	clientSet apimachinery.ClientSetInterface
	kit       *rest.Kit
	datas     mapstr.MapStr
	target    model.Object
	// this instance associate with object id, as is InstAsst table "id" filed.
	assoID int64
}

func (cli *inst) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.datas)
}

func (cli *inst) SetAssoID(id int64) {
	cli.assoID = id
}

func (cli *inst) GetAssoID() int64 {
	return cli.assoID
}

func (cli *inst) searchInsts(targetModel model.Object, cond condition.Condition) ([]Inst, error) {

	queryInput := &metadata.QueryInput{}
	queryInput.Condition = cond.ToMapStr()

	if targetModel.Object().ObjectID != common.BKInnerObjIDHost {
		rsp, err := cli.clientSet.CoreService().Instance().ReadInstance(cli.kit.Ctx, cli.kit.Header, targetModel.GetObjectID(), &metadata.QueryCondition{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[inst-inst] failed to request the object controller , error info is %s, rid: %s", err.Error(), cli.kit.Rid)
			return nil, cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[inst-inst] failed to search the inst, error info is %s, rid: %s", rsp.ErrMsg, cli.kit.Rid)
			return nil, cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}

		return CreateInst(cli.kit, cli.clientSet, targetModel, rsp.Data.Info), nil
	}

	// search hosts
	rsp, err := cli.clientSet.CoreService().Host().GetHosts(cli.kit.Ctx, cli.kit.Header, queryInput)
	if nil != err {
		blog.Errorf("[inst-inst] failed to request the object controller , error info is %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[inst-inst] failed to search the inst, error info is %s, rid: %s", rsp.ErrMsg, cli.kit.Rid)
		return nil, cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return CreateInst(cli.kit, cli.clientSet, targetModel, mapstr.NewArrayFromMapStr(rsp.Data.Info)), nil

}

func (cli *inst) Create() error {
	rid := cli.kit.Rid
	objID := cli.target.GetObjectID()

	// 暂停使用的model不允许创建实例
	if cli.target.Object().IsPaused {
		return cli.kit.CCError.Error(common.CCErrorTopoModelStopped)
	}
	// app/set/module等非通用模型没有 bk_obj_id 字段
	if cli.target.IsCommon() {
		cli.datas.Set(common.BKObjIDField, objID)
	}

	cli.datas.Set(common.BKOwnerIDField, cli.kit.SupplierAccount)

	data := &metadata.CreateModelInstance{Data: cli.datas}
	rsp, err := cli.clientSet.CoreService().Instance().CreateInstance(cli.kit.Ctx, cli.kit.Header, objID, data)
	if nil != err {
		blog.Errorf("failed to create object instance, error info is %s, rid: %s", err.Error(), rid)
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to create object instance ,error info is %v, rid: %s", rsp.ErrMsg, rid)
		return cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	cli.datas.Set(cli.target.GetInstIDFieldName(), rsp.Data.Created.ID)

	return nil
}

func (cli *inst) Update(data mapstr.MapStr) error {
	exist, filter, err := cli.CheckInstanceExists(nil)
	if err != nil {
		return err
	}
	if exist == true {
		return cli.UpdateInstance(filter, data, nil)
	}
	return cli.kit.CCError.CCError(common.CCErrCommNotFound)
}

func (cli *inst) UpdateInstance(filter condition.Condition, data mapstr.MapStr, nonInnerAttributes []model.AttributeInterface) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)
	rid := cli.kit.Rid
	tObj := cli.target.Object()
	objID := tObj.ObjectID
	if tObj.IsPaused {
		return cli.kit.CCError.Error(common.CCErrorTopoModelStopped)
	}

	// execute update action
	updateOption := metadata.UpdateOption{
		Condition: filter.ToMapStr(),
		Data:      data,
	}
	rsp, err := cli.clientSet.CoreService().Instance().UpdateInstance(cli.kit.Ctx, cli.kit.Header, objID, &updateOption)
	if nil != err {
		blog.Errorf("failed to update object(%s)'s instances, err: %s, rid: %s", objID, err.Error(), rid)
		return cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to update object(%s)'s instances, err: %s, rid: %s", objID, rsp.ErrMsg, rid)
		return cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	// read the updated data
	instItems, err := cli.searchInsts(cli.target, filter)
	if nil != err {
		blog.ErrorJSON("[inst-inst] failed to search updated data, cond: %s, err: %s, rid: %s", filter.ToMapStr(), err.Error(), rid)
		return err
	}

	// TODO: 这种实现方案非常不安全
	for _, item := range instItems { // should be only one item
		cli.datas = item.GetValues()
	}

	return nil
}

func (cli *inst) IsExists() (bool, error) {
	exist, _, err := cli.CheckInstanceExists(nil)
	return exist, err
}

func (cli *inst) CheckInstanceExists(nonInnerAttributes []model.AttributeInterface) (exist bool, filter condition.Condition, err error) {
	rid := cli.kit.Rid
	objID := cli.target.GetObjectID()
	instIDField := cli.target.GetInstIDFieldName()
	instNameField := cli.target.GetInstNameFieldName()
	if nonInnerAttributes == nil {
		var err error
		nonInnerAttributes, err = cli.target.GetNonInnerAttributes()
		if nil != err {
			blog.Errorf("failed to get attributes for the object(%s), err: %s, rid: %s", objID, err.Error(), rid)
			return false, nil, err
		}
	}

	cond := condition.CreateCondition()
	// if the inst id already exist, query it with id directly,
	// otherwise, when import a object instance, the other field may be changed.
	if id, exist := cli.datas[instIDField]; exist {
		cond.Field(instIDField).Eq(id)
	} else {
		if cli.target.IsCommon() {
			cond.Field(common.BKObjIDField).Eq(objID)
		}
		val, exists := cli.datas.Get(common.BKInstParentStr)
		if exists {
			cond.Field(common.BKInstParentStr).Eq(val)
		}

		for _, attrItem := range nonInnerAttributes {
			// check the inst
			attr := attrItem.Attribute()
			if attr.IsOnly || attr.PropertyID == instNameField {

				val, exists := cli.datas.Get(attr.PropertyID)
				if !exists {
					return false, nil, cli.kit.CCError.Errorf(common.CCErrCommParamsLostField, attr.PropertyID)
				}
				cond.Field(attr.PropertyID).Eq(val)
			}
		}
	}

	queryCond := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}
	rsp, err := cli.clientSet.CoreService().Instance().ReadInstance(cli.kit.Ctx, cli.kit.Header, objID, queryCond)
	if nil != err {
		blog.Errorf("failed to search object(%s) instances, err: %s, rid: %s", objID, err.Error(), rid)
		return false, nil, cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s) instances, err: %s, rid: %s", objID, rsp.ErrMsg, rid)
		return false, nil, cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != rsp.Data.Count, cond, nil
}

func (cli *inst) Save(data mapstr.MapStr) error {

	if nil != data {
		cli.SetValues(data)
	}
	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		if nil == data {
			return cli.Update(cli.datas)
		}
		return cli.Update(data)
	}

	return cli.Create()
}

func (cli *inst) GetObject() model.Object {
	return cli.target
}

func (cli *inst) GetInstID() (int64, error) {
	return cli.datas.Int64(cli.target.GetInstIDFieldName())
}
func (cli *inst) GetParentID() (int64, error) {
	return cli.datas.Int64(common.BKInstParentStr)
}

func (cli *inst) GetInstName() (string, error) {

	return cli.datas.String(cli.target.GetInstNameFieldName())
}

func (cli *inst) ToMapStr() mapstr.MapStr {
	return cli.datas
}

func (cli *inst) SetValue(key string, value interface{}) {
	cli.datas.Set(key, value)
}

func (cli *inst) SetValues(values mapstr.MapStr) {
	cli.datas.Merge(values)
}

func (cli *inst) GetValues() mapstr.MapStr {
	return cli.datas
}

func (cli *inst) IsDefault() bool {
	if cli.datas.Exists(common.BKDefaultField) {
		defaultVal, err := cli.datas.Int64(common.BKDefaultField)
		if nil != err {
			defaultFieldValue := cli.datas[common.BKDefaultField]
			blog.Errorf("[operation-inst]the `default` field's value(%#v) invalid, err: %s, rid: %s", defaultFieldValue, err.Error(), cli.kit.Rid)
			return false
		}

		if defaultVal == int64(common.DefaultAppFlag) {
			return false
		}
	}

	return false
}

func (cli *inst) GetBizID() (int64, error) {
	switch cli.target.Object().ObjectID {
	case common.BKInnerObjIDApp:
		return cli.GetInstID()
	case common.BKInnerObjIDSet:
		return util.GetInt64ByInterface(cli.datas[common.BKAppIDField])
	case common.BKInnerObjIDModule:
		return util.GetInt64ByInterface(cli.datas[common.BKAppIDField])
	default:
		return metadata.ParseBizIDFromData(cli.datas)
	}
}
