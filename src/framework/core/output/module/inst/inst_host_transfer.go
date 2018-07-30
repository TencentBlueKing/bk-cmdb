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

	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/client/v3"
)

// TransferInterface host trnasfer operation
type TransferInterface interface {
	// MoveToModule transfer host to business module
	MoveToModule(newModuleIDS []int64, isIncrement bool) error

	// MoveToFaultModule transfer host module to fault module
	MoveToFaultModule() error

	// MoveToIdleModule transfer host module to idle module
	MoveToIdleModule() error

	// MoveToResourcePools transfer host module to resource pools
	MoveToResourcePools() error

	// MoveToBusinessIdleModuleFromResourcePools transfer host to business module
	MoveToBusinessIdleModuleFromResourcePools(bizID int64) error

	// MoveToAnotherBusinessModules transfer host to another business modules
	MoveToAnotherBusinessModules(bizID int64, moduleID int64) error

	// ResetBusinessHosts transfer the hosts in set or module to the idle module
	ResetBusinessHosts(setID, moduleID int64) error
}

type transfer struct {
	targetHost *host
}

// MoveToModule transfer host to business module
func (t *transfer) MoveToModule(newModuleIDS []int64, isIncrement bool) error {

	hostID, err := t.targetHost.GetInstID()
	if nil != err {
		log.Errorf("failed to get the host id, error info is %s", err.Error())
		return err
	}

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().TransferHostToBusinessModule(t.targetHost.bizID, []int64{hostID}, newModuleIDS, isIncrement)
}

// MoveToFaultModule transfer host module to fault module
func (t *transfer) MoveToFaultModule() error {

	hostID, err := t.targetHost.GetInstID()
	if nil != err {
		log.Errorf("failed to get the host id, error info is %s", err.Error())
		return err
	}

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().TransferHostToBusinessFaultModule(t.targetHost.bizID, []int64{hostID})
}

// MoveToIdleModule transfer host module to idle module
func (t *transfer) MoveToIdleModule() error {
	hostID, err := t.targetHost.GetInstID()
	if nil != err {
		log.Errorf("failed to get the host id, error info is %s", err.Error())
		return err
	}

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().TransferHostToBusinessIdleModule(t.targetHost.bizID, []int64{hostID})
}

// MoveToResourcePools transfer host module to resource pools
func (t *transfer) MoveToResourcePools() error {
	hostID, err := t.targetHost.GetInstID()
	if nil != err {
		log.Errorf("failed to get the host id, error info is %s", err.Error())
		return err
	}

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().TransferHostToResourcePools(t.targetHost.bizID, []int64{hostID})
}

// MoveToBusinessIdleModuleFromResourcePools transfer host to business module
func (t *transfer) MoveToBusinessIdleModuleFromResourcePools(bizID int64) error {
	hostID, err := t.targetHost.GetInstID()
	if nil != err {
		log.Errorf("failed to get the host id, error info is %s", err.Error())
		return err
	}

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().TransferHostFromResourcePoolsToBusinessIdleModule(bizID, []int64{hostID})
}

// MoveToAnotherBusinessModules transfer host to another business modules
func (t *transfer) MoveToAnotherBusinessModules(bizID int64, moduleID int64) error {

	hostInfo := v3.HostInfo{}
	cloudID, err := t.targetHost.datas.Int64(PlatID)
	if nil != err {
		log.Errorf("failed to get the host cloud id, error info is %s", err.Error())
		return fmt.Errorf("failed to get the host cloud id, error info is %s", err.Error())
	}

	innerIP := t.targetHost.datas.String(HostInnerIP)
	if nil != err {
		log.Errorf("failed to get the host innper ip, error info is %s", err.Error())
		return err
	}

	hostInfo.CloudID = cloudID
	hostInfo.HostInnerIP = innerIP

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().TransferHostToAnotherBusinessModules(bizID, moduleID, []*v3.HostInfo{})
}

// ResetBusinessHosts transfer the hosts in set or module to the idle module
func (t *transfer) ResetBusinessHosts(setID, moduleID int64) error {

	return client.GetClient().CCV3(client.Params{SupplierAccount: t.targetHost.GetModel().GetSupplierAccount()}).Host().ResetBusinessHosts(t.targetHost.bizID, moduleID, setID)
}
