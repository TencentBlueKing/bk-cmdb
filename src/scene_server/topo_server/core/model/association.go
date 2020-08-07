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
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

var _ Association = (*association)(nil)

type association struct {
	asst       metadata.Association
	isMainLine bool
	kit        *rest.Kit
	clientSet  apimachinery.ClientSetInterface
}

func (cli *association) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.asst)
}

func (cli *association) GetType() AssociationType {
	return CommonAssociation
}

func (cli *association) IsExists() (bool, error) {
	return false, nil
}

func (cli *association) SetSupplierAccount(supplierAccount string) {
	cli.asst.OwnerID = supplierAccount
}

func (cli *association) GetSupplierAccount() string {
	return cli.asst.OwnerID
}

func (cli *association) Create() error {
	return nil
}
func (cli *association) Delete() error {
	return nil
}
func (cli *association) Update(data mapstr.MapStr) error {
	return nil
}
func (cli *association) Save(data mapstr.MapStr) error {
	return nil
}

func (cli *association) SetTopo(parent, child Object) error {
	return nil
}

func (cli *association) GetTopo(obj Object) (Topo, error) {
	return nil, nil
}

func (cli *association) ToMapStr() (mapstr.MapStr, error) {
	rst := mapstr.SetValueToMapStrByTags(&cli.asst)
	return rst, nil
}

func (cli *association) Parse(data mapstr.MapStr) (*metadata.Association, error) {
	return cli.asst.Parse(data)
}
