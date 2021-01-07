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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
)

// Classification classification operation interface declaration
type Classification interface {
	Operation
	Parse(data mapstr.MapStr) (*metadata.Classification, error)
	GetObjects() ([]Object, error)
	Classify() metadata.Classification
	ToMapStr() mapstr.MapStr
}

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	FieldValid
	cls       metadata.Classification
	kit       *rest.Kit
	clientSet apimachinery.ClientSetInterface
}

func (cli *classification) Classify() metadata.Classification {
	return cli.cls
}

func (cli *classification) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.cls)
}

func (cli *classification) Parse(data mapstr.MapStr) (*metadata.Classification, error) {
	return cli.cls.Parse(data)
}

func (cli *classification) ToMapStr() mapstr.MapStr {
	rst := mapstr.SetValueToMapStrByTags(&cli.cls)
	return rst
}

func (cli *classification) GetObjects() ([]Object, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.ModelFieldObjCls).Eq(cli.cls.ClassificationID)

	rsp, err := cli.clientSet.CoreService().Model().ReadModel(context.Background(), cli.kit.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the classification(%s) object, error info is %s, rid: %s", cli.cls.ClassificationID, rsp.ErrMsg, cli.kit.Rid)
		return nil, cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	rstItems := make([]Object, 0)
	for _, item := range rsp.Data.Info {
		tmpObj := &object{
			obj:   item.Spec,
			isNew: false,
		}
		rstItems = append(rstItems, tmpObj)
	}

	return rstItems, nil
}

func (cli *classification) IsValid(isUpdate bool, data mapstr.MapStr) error {
	if !isUpdate || data.Exists(metadata.ClassFieldClassificationID) {
		val, err := cli.FieldValid.Valid(cli.kit, data, metadata.ClassFieldClassificationID)
		if nil != err {
			return err
		}
		if err = cli.FieldValid.ValidID(cli.kit, val); nil != err {
			return err
		}
	}
	if !isUpdate || data.Exists(metadata.ClassFieldClassificationName) {
		val, err := cli.FieldValid.Valid(cli.kit, data, metadata.ClassFieldClassificationName)
		if nil != err {
			return err
		}
		if err = cli.FieldValid.ValidName(cli.kit, val); nil != err {
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
		return cli.kit.CCError.Errorf(common.CCErrCommDuplicateItem, cli.GetID()+"/"+cli.GetName())
	}

	input := metadata.CreateOneModelClassification{Data: cli.cls}
	rsp, err := cli.clientSet.CoreService().Model().CreateModelClassification(context.Background(), cli.kit.Header, &input)
	if nil != err {
		blog.Errorf("failed to request object controller, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to create classification(%s), error info is %s, rid: %s", cli.cls.ClassificationID, rsp.ErrMsg, cli.kit.Rid)
		return cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	cli.cls.ID = int64(rsp.Data.Created.ID)
	return nil
}

func (cli *classification) Update(data mapstr.MapStr) error {

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
		return cli.kit.CCError.Errorf(common.CCErrCommDuplicateItem, cli.GetName())
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

		//TODO : if update with kit cli.kit.MetaData
		input := metadata.UpdateOption{
			Condition: cond.ToMapStr(),
			Data:      data,
		}
		rsp, err := cli.clientSet.CoreService().Model().UpdateModelClassification(context.Background(), cli.kit.Header, &input)
		if nil != err {
			blog.Errorf("failed to request object controller, err: %s, rid: %s", err.Error(), cli.kit.Rid)
			return err
		}

		if !rsp.Result {
			blog.Errorf("failed to update the classificaiotn(%s), error info is %s, rid: %s", cli.cls.ClassificationID, rsp.ErrMsg, cli.kit.Rid)
			return cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}

		cli.cls = item
	}

	return nil
}

func (cli *classification) search(cond condition.Condition) ([]metadata.Classification, error) {
	rsp, err := cli.clientSet.CoreService().Model().ReadModelClassification(context.Background(), cli.kit.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("failed to search the classification, error info is %s, rid: %s", rsp.ErrMsg, cli.kit.Rid)
		return nil, cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
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
	cond.Field(metadata.ClassFieldClassificationName).Eq(cli.cls.ClassificationName)
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

func (cli *classification) Save(data mapstr.MapStr) error {

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
	cli.cls.OwnerID = supplierAccount
}

func (cli *classification) GetSupplierAccount() string {
	return cli.cls.OwnerID
}

func (cli *classification) SetIcon(classificationIcon string) {
	cli.cls.ClassificationIcon = classificationIcon
}

func (cli *classification) GetIcon() string {
	return cli.cls.ClassificationIcon
}
