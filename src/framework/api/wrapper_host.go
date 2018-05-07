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

package api

import (
	"configcenter/src/framework/core/output/module/inst"
)

// HostWrapper the host wrapper
type HostWrapper struct {
	host inst.Inst
}

// SetValue set the key value
func (cli *HostWrapper) SetValue(key string, val interface{}) error {
	return cli.host.SetValue(key, val)
}

// Save save the data
func (cli *HostWrapper) Save() error {
	if err := cli.host.SetValue(fieldImportFrom, HostImportFromAPI); nil != err {
		return err
	}
	return cli.host.Save()
}

// SetBakOperator set the bak operator
func (cli *HostWrapper) SetBakOperator(bakOperator string) error {
	return cli.host.SetValue(fieldBakOperator, bakOperator)
}

// SetOsBit set the os bit
func (cli *HostWrapper) SetOsBit(osbit string) error {
	return cli.host.SetValue(fieldOsBit, osbit)
}

// SetSLA set the sla
func (cli *HostWrapper) SetSLA(sla string) error {
	return cli.host.SetValue(fieldSLA, sla)
}

// SetCloudID set the cloudid for the host
func (cli *HostWrapper) SetCloudID(cloudID int64) error {
	return cli.host.SetValue(fieldCloudID, cloudID)
}

// SetInnerIP set the inner ip
func (cli *HostWrapper) SetInnerIP(innerIP string) error {
	return cli.host.SetValue(fieldHostInnerIP, innerIP)
}

// SetOperator set the operator for the host
func (cli *HostWrapper) SetOperator(operator string) error {
	return cli.host.SetValue(fieldHostOperator, operator)
}

// SetStateName set the state name for the host
func (cli *HostWrapper) SetStateName(stateName string) error {
	return cli.host.SetValue(fieldStateName, stateName)
}

// SetCPU set the cpu core num  for the host
func (cli *HostWrapper) SetCPU(cpu int64) error {
	return cli.host.SetValue(fieldCPU, cpu)
}

// SetCPUMhz set the cpu mhz
func (cli *HostWrapper) SetCPUMhz(cpuMhz float64) error {
	return cli.host.SetValue(fieldCPUMhz, cpuMhz)
}

// SetOsType set the os type for the host
func (cli *HostWrapper) SetOsType(osType string) error {
	return cli.host.SetValue(fieldOsType, osType)
}

// SetOuterIP set the outer ip for the host
func (cli *HostWrapper) SetOuterIP(outerIP string) error {
	return cli.host.SetValue(fieldHostOuterIP, outerIP)
}

// SetAssetID set the assetid for the host
func (cli *HostWrapper) SetAssetID(assetID string) error {
	return cli.host.SetValue(fieldAssetID, assetID)
}

// SetMac set the mac for the host
func (cli *HostWrapper) SetMac(mac string) error {
	return cli.host.SetValue(fieldMac, mac)
}

// SetProvinceName set the province name for the host
func (cli *HostWrapper) SetProvinceName(provinceName string) error {
	return cli.host.SetValue(fieldProvinceName, provinceName)
}

// SetSN set the sn for the host
func (cli *HostWrapper) SetSN(sn string) error {
	return cli.host.SetValue(fieldSN, sn)
}

// SetCPUModule set the cpu module for the host
func (cli *HostWrapper) SetCPUModule(cpuModule string) error {
	return cli.host.SetValue(fieldCPUModule, cpuModule)
}

// SetName set the host name for the host
func (cli *HostWrapper) SetName(hostName string) error {
	return cli.host.SetValue(fieldHostName, hostName)
}

// SetISPName set the isp name for the host
func (cli *HostWrapper) SetISPName(ispName string) error {
	return cli.host.SetValue(fieldISPName, ispName)
}

// SetServiceTerm set the service term for the host
func (cli *HostWrapper) SetServiceTerm(serviceTerm int64) error {
	return cli.host.SetValue(fieldServiceTerm, serviceTerm)
}

// SetComment set the comment for the host
func (cli *HostWrapper) SetComment(comment string) error {
	return cli.host.SetValue(fieldComment, comment)
}

// SetMem set the mem for the host
func (cli *HostWrapper) SetMem(mem int64) error {
	return cli.host.SetValue(fieldMem, mem)
}

// SetDisk set the capacity of the disk for the host
func (cli *HostWrapper) SetDisk(disk int64) error {
	return cli.host.SetValue(fieldDisk, disk)
}

// SetOsName set the os name for the host
func (cli *HostWrapper) SetOsName(osName string) error {
	return cli.host.SetValue(fieldOsName, osName)
}

// SetOsVersion set the os version for the host
func (cli *HostWrapper) SetOsVersion(osVersion string) error {
	return cli.host.SetValue(fieldOsVersion, osVersion)
}
