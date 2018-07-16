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
	"strconv"
	"strings"

	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

// OperationInterface inst operation interface
type OperationInterface interface {
	CreateCommonInst(target model.Model) CommonInstInterface
	CreateSetInst(target model.Model) SetInterface
	CreateModuleInst(target model.Model) ModuleInterface
	CreateBusinessInst(target model.Model) BusinessInterface
	CreatePlatInst(target model.Model) CommonInstInterface
	CreateHostInst(target model.Model) HostInterface

	DeleteHosts(supplierAccount string, hostIDS []int64) error

	FindCommonInstLikeName(target model.Model, instName string) (Iterator, error)
	FindCommonInstByCondition(target model.Model, cond common.Condition) (Iterator, error)
	FindBusinessLikeName(target model.Model, businessName string) (BusinessIterator, error)
	FindBusinessByCondition(target model.Model, condition common.Condition) (BusinessIterator, error)
	FindHostsLikeName(target model.Model, hostname string) (HostIterator, error)
	FindHostsByCondition(target model.Model, condition common.Condition) (HostIterator, error)
	FindModulesLikeName(target model.Model, moduleName string) (ModuleIterator, error)
	FindModulesByCondition(target model.Model, cond common.Condition) (ModuleIterator, error)
	FindPlatsLikeName(target model.Model, platName string) (Iterator, error)
	FindPlatsByCondition(target model.Model, cond common.Condition) (Iterator, error)
	FindSetsLikeName(target model.Model, setName string) (SetIterator, error)
	FindSetsByCondition(target model.Model, cond common.Condition) (SetIterator, error)
}

// Operation inst operation interface
func Operation() OperationInterface {
	return &operation{}
}

type operation struct {
}

func (o *operation) CreateCommonInst(target model.Model) CommonInstInterface {
	return &inst{target: target, datas: types.MapStr{}}
}

func (o *operation) CreateBusinessInst(target model.Model) BusinessInterface {
	return &business{target: target, datas: types.MapStr{}}
}

func (o *operation) CreateSetInst(target model.Model) SetInterface {
	return &set{target: target, datas: types.MapStr{}}
}

func (o *operation) CreateModuleInst(target model.Model) ModuleInterface {
	return &module{target: target, datas: types.MapStr{}}
}

func (o *operation) CreatePlatInst(target model.Model) CommonInstInterface {
	return &inst{target: target, datas: types.MapStr{}}
}
func (o *operation) CreateHostInst(target model.Model) HostInterface {
	return &host{target: target, datas: types.MapStr{}}
}

func (o *operation) DeleteHosts(supplierAccount string, hostIDS []int64) error {

	hostIDArr := make([]string, 0)
	for hostID := range hostIDS {
		hostIDArr = append(hostIDArr, strconv.Itoa(hostID))
	}

	return client.GetClient().CCV3(client.Params{SupplierAccount: supplierAccount}).Host().DeleteHostBatch(strings.Join(hostIDArr, ","))
}

func (o *operation) FindCommonInstLikeName(target model.Model, instName string) (Iterator, error) {
	cond := common.CreateCondition().Field(InstName).Like(instName)
	return newIteratorInst(target, cond)
}

func (o *operation) FindCommonInstByCondition(target model.Model, cond common.Condition) (Iterator, error) {
	return newIteratorInst(target, cond)
}

func (o *operation) FindBusinessLikeName(target model.Model, businessName string) (BusinessIterator, error) {
	cond := common.CreateCondition().Field(BusinessNameField).Like(businessName)
	return newIteratorInstBusiness(target, cond)
}

func (o *operation) FindBusinessByCondition(target model.Model, condition common.Condition) (BusinessIterator, error) {
	return newIteratorInstBusiness(target, condition)
}

func (o *operation) FindHostsLikeName(target model.Model, hostname string) (HostIterator, error) {
	cond := common.CreateCondition().Field(HostNameField).Like(hostname)
	return newHostIterator(target, cond)
}

func (o *operation) FindHostsByCondition(target model.Model, condition common.Condition) (HostIterator, error) {
	return newHostIterator(target, condition)
}

func (o *operation) FindModulesLikeName(target model.Model, moduleName string) (ModuleIterator, error) {
	cond := common.CreateCondition().Field(ModuleName).Like(moduleName)
	return newIteratorInstModule(target, cond)
}

func (o *operation) FindModulesByCondition(target model.Model, cond common.Condition) (ModuleIterator, error) {
	return newIteratorInstModule(target, cond)
}

func (o *operation) FindPlatsLikeName(target model.Model, platName string) (Iterator, error) {
	cond := common.CreateCondition().Field(PlatName).Like(platName)
	return newIteratorInst(target, cond)
}

func (o *operation) FindPlatsByCondition(target model.Model, cond common.Condition) (Iterator, error) {
	return newIteratorInst(target, cond)
}

func (o *operation) FindSetsLikeName(target model.Model, setName string) (SetIterator, error) {
	cond := common.CreateCondition().Field(SetName).Like(setName)
	return newIteratorInstSet(target, cond)
}

func (o *operation) FindSetsByCondition(target model.Model, cond common.Condition) (SetIterator, error) {
	return newIteratorInstSet(target, cond)
}
