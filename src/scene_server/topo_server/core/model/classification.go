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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

var _ Classification = (*classification)(nil)

// classification the model classification definition
type classification struct {
	cls       metadata.Classification
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *classification) Parse(data frtypes.MapStr) (*metadata.Classification, error) {
	return cli.cls.Parse(data)
}

func (cli *classification) ToMapStr() (frtypes.MapStr, error) {
	rst := metadata.SetValueToMapStrByTags(cli)
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

func (cli *classification) Create() error {
	return nil
}

func (cli *classification) Update() error {
	return nil
}

func (cli *classification) Delete() error {
	return nil
}

func (cli *classification) IsExists() (bool, error) {
	return false, nil
}

func (cli *classification) Save() error {

	return nil
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
