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
	"fmt"
	"strconv"
	"strings"

	cccommon "configcenter/src/common"
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/errors"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ HostInterface = (*host)(nil)

// HostInterface the host interface
type HostInterface interface {
	Maintaince

	Transfer() TransferInterface

	SetBusinessID(bizID int64)
	SetModuleIDS(moduleIDS []int64, isIncrement bool)

	GetBizs() []types.MapStr
	GetSets() []types.MapStr
	GetModules() []types.MapStr

	GetModel() model.Model

	GetInstID() (int64, error)
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)
}

type host struct {
	isIncrement bool
	bizs        []types.MapStr
	sets        []types.MapStr
	modules     []types.MapStr
	bizID       int64
	moduleIDS   []int64
	setIDS      []int64
	target      model.Model
	datas       types.MapStr
}

func (cli *host) GetBizs() []types.MapStr {
	return cli.bizs
}

func (cli *host) GetSets() []types.MapStr {
	return cli.sets
}

func (cli *host) GetModules() []types.MapStr {
	return cli.modules
}

func (cli *host) reset() error {

	// parse business
	datas, err := cli.datas.MapStrArray(Business)
	if nil != err {
		return fmt.Errorf("failed to get biz data , error info is %s", err.Error())
	}
	cli.bizs = datas

	for _, dataVal := range datas {
		id, err := dataVal.Int64(BusinessID)
		if nil != err {
			return fmt.Errorf("failed to get biz id, error info is %s", err.Error())
		}

		cli.bizID = id
	}

	// parse module
	datas, err = cli.datas.MapStrArray(Module)
	if nil != err {
		return fmt.Errorf("failed to get module data , error info is %s", err.Error())
	}
	cli.modules = datas

	for _, dataVal := range datas {
		id, err := dataVal.Int64(ModuleID)
		if nil != err {
			return fmt.Errorf("failed to get module id, error info is %s", err.Error())
		}

		cli.moduleIDS = append(cli.moduleIDS, id)
		//fmt.Println("the module id:", id)
	}

	// parse set
	datas, err = cli.datas.MapStrArray(Set)
	if nil != err {
		return fmt.Errorf("failed to get module data , error info is %s", err.Error())
	}
	cli.sets = datas

	for _, dataVal := range datas {
		id, err := dataVal.Int64(SetID)
		if nil != err {
			return fmt.Errorf("failed to get set id, error info is %s", err.Error())
		}

		cli.setIDS = append(cli.setIDS, id)
	}

	return nil
}

func (cli *host) SetBusinessID(bizID int64) {
	cli.bizID = bizID
}

func (cli *host) SetModuleIDS(moduleIDS []int64, isIncrement bool) {
	cli.moduleIDS = moduleIDS
	cli.isIncrement = isIncrement
}

func (cli *host) GetModel() model.Model {
	return cli.target
}

func (cli *host) GetInstID() (int64, error) {
	return cli.datas.Int64(cccommon.BKHostIDField)
}

func (cli *host) GetInstName() string {
	return cli.datas.String(cccommon.BKHostIDField)
}

func (cli *host) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *host) SetValue(key string, value interface{}) error {
	cli.datas.Set(key, value)
	return nil
}

func (cli *host) Transfer() TransferInterface {
	return &transfer{
		targetHost: cli,
	}
}

func (cli *host) ResetAssociationValue() error {
	attrs, err := cli.target.Attributes()
	if nil != err {
		return err
	}

	for _, attrItem := range attrs {
		if attrItem.GetKey() {

			if model.FieldTypeSingleAsst == attrItem.GetType() {
				asstVals, err := cli.datas.MapStrArray(attrItem.GetID())
				if nil != err {
					return err
				}

				for _, val := range asstVals {

					// the datas is created by find
					valID, exists := val.Get(InstID)
					if exists {
						cli.datas.Set(attrItem.GetID(), valID)
						continue
					}

					// default the datas is new one
					valID, exists = val.Get(attrItem.GetID())
					if exists {
						cli.datas.Set(attrItem.GetID(), valID)
					}
				}

				continue
			}

			if model.FieldTypeMultiAsst == attrItem.GetType() {
				asstVals, err := cli.datas.MapStrArray(attrItem.GetID())
				if nil != err {
					return err
				}

				condVals := make([]string, 0)
				for _, val := range asstVals {

					valID := val.String(InstID)
					if 0 != len(valID) {
						condVals = append(condVals, valID)
					}
				}

				cli.datas.Set(attrItem.GetID(), strings.Join(condVals, ","))
			}
		}
	}

	return nil
}

func (cli *host) search() ([]model.Attribute, []types.MapStr, error) {

	if err := cli.ResetAssociationValue(); nil != err {
		return nil, nil, err
	}

	attrs, err := cli.target.Attributes()
	if nil != err {
		return nil, nil, err
	}

	cond := common.CreateCondition()
	for _, attrItem := range attrs {
		if attrItem.GetKey() {

			attrVal, exists := cli.datas.Get(attrItem.GetID())
			if !exists {
				return nil, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
			}
			cond.Field(attrItem.GetID()).Eq(attrVal)
		}
	}
	//log.Infof("the condition:%#v", cond.ToMapStr())
	// search by condition
	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Host().SearchHost(cond)
	return attrs, items, err
}
func (cli *host) IsExists() (bool, error) {

	attrs, err := cli.target.Attributes()
	if nil != err {
		return false, err
	}

	cond := common.CreateCondition()

	for _, attrItem := range attrs {

		if !attrItem.GetKey() {
			continue
		}

		if !cli.datas.Exists(attrItem.GetID()) {
			return false, errors.New("the key field(" + attrItem.GetID() + ") is not set")
		}

		cond.Field(attrItem.GetID()).Eq(cli.datas[attrItem.GetID()])

	}
	//log.Infof("host search condition:%s %d", string(cond.ToMapStr().ToJSON()), 0)
	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Host().SearchHost(cond)
	if nil != err {
		return false, err
	}

	return 0 != len(items), nil
}
func (cli *host) Create() error {
	//log.Infof("the create host:%#v", cli.datas)
	//log.Infof("create the exists:%s %d", string(cli.datas.ToJSON()), 0)
	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return nil
	}
	//log.Infof("the create host:%#v", cli.datas)
	ids, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Host().CreateHostBatch(cli.bizID, cli.moduleIDS, cli.datas)
	if nil != err {
		return err
	}

	if 0 == len(ids) {
		return fmt.Errorf("the host ids is empty")
	}

	cli.datas.Set(HostID, ids[0])

	return nil

}

func (cli *host) Update() error {

	attrs, existItems, err := cli.search()
	if nil != err {
		return err
	}

	//log.Infof("update the exists:%s %d", string(cli.datas.ToJSON()), 0)

	// clear the invalid data
	cli.datas.ForEach(func(key string, val interface{}) {
		for _, attrItem := range attrs {
			if attrItem.GetID() == key {
				return
			}
		}
		//log.Infof("remove the invalid field:%s", key)
		cli.datas.Remove(key)
	})

	cli.datas.Remove("create_time") //invalid check , need to delete

	for _, existItem := range existItems {

		hostID, err := existItem.Int64(HostID)
		if nil != err {
			return err
		}

		updateCond := common.CreateCondition()
		updateCond.Field(ModuleID).Eq(hostID)

		//log.Infof("the exists:%s %d", string(existItem.ToJSON()), hostID)
		err = client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Host().UpdateHostBatch(cli.datas, strconv.Itoa(int(hostID)))
		if err != nil {
			log.Errorf("failed to update host, error info is %s", err.Error())
			return err
		}

		cli.datas.Set(HostID, hostID)
		if 0 != len(cli.moduleIDS) {
			err = cli.Transfer().MoveToModule(cli.moduleIDS, cli.isIncrement)
			if nil != err {
				log.Errorf("failed to biz(%d) set modules(%#v), error info is %s", cli.bizID, cli.moduleIDS, err.Error())
				return err
			}
		}
	}

	return nil
}
func (cli *host) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()
}
