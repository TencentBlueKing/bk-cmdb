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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Classification classification operation interface declaration
type Classification interface {
	Operation
	Parse(data frtypes.MapStr) (*metadata.Classification, error)

	GetObjects() ([]Object, error)

	SetID(classificationID string)
	GetID() string

	SetName(classificationName string)
	GetName() string

	SetType(classificationType string)
	GetType() string

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetIcon(classificationIcon string)
	GetIcon() string

	ToMapStr() (frtypes.MapStr, error)
}

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	FieldValid
	cls       metadata.Classification
	params    types.ContextParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *classification) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.cls)
}

func (cli *classification) Parse(data frtypes.MapStr) (*metadata.Classification, error) {
	return cli.cls.Parse(data)
}

func (cli *classification) ToMapStr() (frtypes.MapStr, error) {
	rst := metadata.SetValueToMapStrByTags(&cli.cls)
	return rst, nil
}

func (cli *classification) GetObjects() ([]Object, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.ModelFieldObjCls).Eq(cli.cls.ClassificationID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjects(context.Background(), cli.params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the classification(%s) object, error info is %s", cli.cls.ClassificationID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Object, 0)
	for _, item := range rsp.Data {

		tmpObj := &object{
			isNew: false,
		}

		err := metadata.SetValueToStructByTags(tmpObj.obj, item.ToMapStr())
		if nil != err {
			return nil, err
		}

		rstItems = append(rstItems, tmpObj)
	}

	return rstItems, nil
}

func (cli *classification) IsValid(isUpdate bool, data frtypes.MapStr) error {

	if !isUpdate || data.Exists(metadata.ClassFieldClassificationID) {
		if _, err := cli.FieldValid.Valid(cli.params, data, metadata.ClassFieldClassificationID); nil != err {
			return err
		}
	}
	if !isUpdate || data.Exists(metadata.ClassFieldClassificationName) {
		if _, err := cli.FieldValid.Valid(cli.params, data, metadata.ClassFieldClassificationName); nil != err {
			return err
		}
	}
	return nil
}

func (cli *classification) Create() error {

	if err := cli.IsValid(false, cli.cls.ToMapStr()); nil != err {
		return err
	}
	exists, err := cli.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return cli.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	rsp, err := cli.clientSet.ObjectController().Meta().CreateClassification(context.Background(), cli.params.Header, &cli.cls)
	if nil != err {
		blog.Errorf("failed to request object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("faield to create classification(%s), error info is %s", cli.cls.ClassificationID, rsp.ErrMsg)
		return cli.params.Err.Error(rsp.Code)
	}

	cli.cls.ID = rsp.Data.ID
	return nil
}

func (cli *classification) Update(data frtypes.MapStr) error {

	data.Remove(metadata.ClassFieldClassificationID)
	data.Remove(metadata.ClassificationFieldID)

	if err := cli.IsValid(true, data); nil != err {
		return err
	}

	exists, err := cli.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return cli.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	cond := condition.CreateCondition()
	if 0 == len(cli.cls.ClassificationID) {
		cond.Field(metadata.ClassificationFieldID).Eq(cli.cls.ID)
	} else {
		cond.Field(metadata.ClassFieldClassificationID).Eq(cli.cls.ClassificationID)
	}

	updateItems, err := cli.search(cond)
	if nil != err {
		return err
	}

	for _, item := range updateItems { // only one item

		rsp, err := cli.clientSet.ObjectController().Meta().UpdateClassification(context.Background(), item.ID, cli.params.Header, data)
		if nil != err {
			blog.Errorf("failed to resuest object controller, error info is %s", err.Error())
			return err
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("faile to update the classificaiotn(%s), error info is %s", cli.cls.ClassificationID, rsp.ErrMsg)
			return cli.params.Err.Error(rsp.Code)
		}

		cli.cls = item
	}

	return nil
}

func (cli *classification) search(cond condition.Condition) ([]metadata.Classification, error) {

	rsp, err := cli.clientSet.ObjectController().Meta().SelectClassifications(context.Background(), cli.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the classificaiont, error info is %s", rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	return rsp.Data, nil
}

func (cli *classification) IsExists() (bool, error) {

	// check id
	cond := condition.CreateCondition()
	cond.Field(metadata.ClassFieldClassificationID).Eq(cli.cls.ClassificationID)
	cond.Field(metadata.ClassificationFieldID).NotIn([]int64{cli.cls.ID})
	items, err := cli.search(cond)
	if nil != err {
		return false, err
	}
	if 0 != len(items) {
		return true, err
	}

	// check name
	cond = condition.CreateCondition()
	cond.Field(metadata.ClassFieldClassificationID).Eq(cli.cls.ClassificationName)
	cond.Field(metadata.ClassificationFieldID).NotIn([]int64{cli.cls.ID})
	items, err = cli.search(cond)
	if nil != err {
		return false, err
	}
	if 0 != len(items) {
		return true, err
	}

	return false, nil
}

func (cli *classification) Save(data frtypes.MapStr) error {

	if nil != data {
		if _, err := cli.cls.Parse(data); nil != err {
			return err
		}
	}

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if !exists {
		return cli.Create()
	}

	if nil != data {
		return cli.Update(data)
	}

	return cli.Update(cli.cls.ToMapStr())
}

func (cli *classification) SetID(classificationID string) {
	cli.cls.ClassificationID = classificationID
}

func (cli *classification) GetID() string {
	return cli.cls.ClassificationID
}

func (cli *classification) SetName(classificationName string) {
	cli.cls.ClassificationName = classificationName
}

func (cli *classification) GetName() string {
	return cli.cls.ClassificationName
}

func (cli *classification) SetType(classificationType string) {
	cli.cls.ClassificationType = classificationType
}

func (cli *classification) GetType() string {
	return cli.cls.ClassificationType
}

func (cli *classification) SetSupplierAccount(supplierAccount string) {
	// TODO: need to add owner field
}

func (cli *classification) GetSupplierAccount() string {
	// TODO: need to add owner field
	return ""
}

func (cli *classification) SetIcon(classificationIcon string) {
	cli.cls.ClassificationIcon = classificationIcon
}

func (cli *classification) GetIcon() string {
	return cli.cls.ClassificationIcon
}
