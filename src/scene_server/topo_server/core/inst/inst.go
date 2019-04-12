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
	"context"
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
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

	SetValue(key string, value interface{}) error

	SetValues(values mapstr.MapStr)

	GetValues() mapstr.MapStr

	ToMapStr() mapstr.MapStr

	IsDefault() bool
}

var _ Inst = (*inst)(nil)

type inst struct {
	clientSet apimachinery.ClientSetInterface
	params    types.ContextParams
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
		rsp, err := cli.clientSet.CoreService().Instance().ReadInstance(context.Background(), cli.params.Header, targetModel.GetObjectID(), &metadata.QueryCondition{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[inst-inst] failed to request the object controller , error info is %s", err.Error())
			return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[inst-inst] failed to search the inst, error info is %s", rsp.ErrMsg)
			return nil, cli.params.Err.New(rsp.Code, rsp.ErrMsg)
		}

		return CreateInst(cli.params, cli.clientSet, targetModel, rsp.Data.Info), nil
	}

	// search hosts
	rsp, err := cli.clientSet.HostController().Host().GetHosts(context.Background(), cli.params.Header, queryInput)
	if nil != err {
		blog.Errorf("[inst-inst] failed to request the object controller , error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[inst-inst] failed to search the inst, error info is %s", rsp.ErrMsg)
		return nil, cli.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return CreateInst(cli.params, cli.clientSet, targetModel, mapstr.NewArrayFromMapStr(rsp.Data.Info)), nil

}

func (cli *inst) Create() error {
	if cli.target.Object().IsPaused {
		return cli.params.Err.Error(common.CCErrorTopoModleStopped)
	}
	if cli.target.IsCommon() {
		cli.datas.Set(common.BKObjIDField, cli.target.Object().ObjectID)
	}

	cli.datas.Set(common.BKOwnerIDField, cli.params.SupplierAccount)

    rsp, err := cli.clientSet.CoreService().Instance().CreateInstance(context.Background(), cli.params.Header, cli.target.GetObjectID(), &metadata.CreateModelInstance{Data: cli.datas})
	if nil != err {
		blog.Errorf("failed to create object instance, error info is %s", err.Error())
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to create object instance ,error info is %v", rsp.ErrMsg)
		return cli.params.Err.Error(common.CCErrTopoInstCreateFailed)
	}

	cli.datas.Set(cli.target.GetInstIDFieldName(), rsp.Data.Created.ID)

	return nil
}

func (cli *inst) Update(data mapstr.MapStr) error {
	if cli.target.Object().IsPaused {
		return cli.params.Err.Error(common.CCErrorTopoModleStopped)
	}
	instIDName := cli.target.GetInstIDFieldName()
	instID, exists := cli.datas.Get(instIDName)

	tObj := cli.target.Object()
	cond := condition.CreateCondition()
	if cli.target.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(tObj.ObjectID)
	}

	if exists {
		// construct the update condition by the instid
		cond.Field(instIDName).Eq(instID)
	} else {
		// construct the update condition by the only key

		attrs, err := cli.target.GetAttributesExceptInnerFields()
		if nil != err {
			blog.Errorf("failed to get attributes for the object(%s), error info is is %s", tObj.ObjectID, err.Error())
			return err
		}

		for _, attrItem := range attrs {
			// check the inst
			att := attrItem.Attribute()
			if att.IsOnly || att.PropertyID == cli.target.GetInstNameFieldName() {
				val, exists := cli.datas.Get(att.PropertyID)
				if !exists {
					return cli.params.Err.Errorf(common.CCErrCommParamsLostField, att.PropertyID)
				}
				cond.Field(att.PropertyID).Eq(val)
			}
		}

	}

	// execute update action
	updateCond := metadata.UpdateOption{}
	updateCond.Data = data
	updateCond.Condition = cond.ToMapStr()
	rsp, err := cli.clientSet.CoreService().Instance().UpdateInstance(context.Background(), cli.params.Header, cli.target.GetObjectID(), &updateCond)
	if nil != err {
		blog.Errorf("failed to update the object(%s) instances, error info is %s", tObj.ObjectID, err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to update the object(%s) instances, error info is %s", tObj.ObjectID, rsp.ErrMsg)
		return cli.params.Err.Error(common.CCErrTopoInstUpdateFailed)
	}

	// read the new data
	instItems, err := cli.searchInsts(cli.target, cond)
	if nil != err {
		blog.Errorf("[inst-inst] failed to search the new insts data, error info is %s", err.Error())
		return err
	}

	for _, item := range instItems { // should be only one item
		cli.datas = item.GetValues()
	}

	return nil
}

func (cli *inst) IsExists() (bool, error) {

	tObj := cli.target.Object()
	attrs, err := cli.target.GetAttributesExceptInnerFields()
	if nil != err {
		blog.Errorf("failed to get attributes for the object(%s), error info is is %s", tObj.ObjectID, err.Error())
		return false, err
	}

	cond := condition.CreateCondition()
	// if the inst id already exist, query it with id directly,
	// otherwise, when import a object instance, the other field may be changed.
	if id, exist := cli.datas[cli.target.GetInstIDFieldName()]; exist {
		cond.Field(cli.target.GetInstIDFieldName()).Eq(id)
	} else {
		if cli.target.IsCommon() {
			cond.Field(common.BKObjIDField).Eq(tObj.ObjectID)
		}
		val, exists := cli.datas.Get(common.BKInstParentStr)
		if exists {
			cond.Field(common.BKInstParentStr).Eq(val)
		}

		for _, attrItem := range attrs {

			// check the inst
			attr := attrItem.Attribute()
			if attr.IsOnly || attr.PropertyID == cli.target.GetInstNameFieldName() {

				val, exists := cli.datas.Get(attr.PropertyID)
				if !exists {
					return false, cli.params.Err.Errorf(common.CCErrCommParamsLostField, attr.PropertyID)
				}
				cond.Field(attr.PropertyID).Eq(val)
			}

		}
	}

	queryCond := metadata.QueryInput{}
	queryCond.Condition = cond.ToMapStr()

	rsp, err := cli.clientSet.CoreService().Instance().ReadInstance(
		context.Background(), cli.params.Header, cli.target.GetObjectID(), &metadata.QueryCondition{Condition: cond.ToMapStr()},
	)
	if nil != err {
		blog.Errorf("failed to search object(%s) instances  , error info is %s", tObj.ObjectID, err.Error())
		return false, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object (%s) instances, error info is %s", tObj.ObjectID, rsp.ErrMsg)
		return false, cli.params.Err.Error(common.CCErrTopoInstSelectFailed)
	}

	return 0 != rsp.Data.Count, nil
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

func (cli *inst) SetValue(key string, value interface{}) error {
	cli.datas.Set(key, value)
	return nil
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
			blog.Errorf("[operation-inst]the default value(%#v) is invalid, error info is %s", cli.datas[common.BKDefaultField], err.Error())
			return false
		}

		if defaultVal == int64(common.DefaultAppFlag) {
			return false
		}
	}

	return false
}
