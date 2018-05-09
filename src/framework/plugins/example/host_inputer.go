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

package example

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/output/module/model"

	"fmt"
	"time"
)

func init() {

	// api.RegisterInputer(host)
	api.RegisterFrequencyInputer(host, time.Minute*5)
}

var host = &hostInputer{}

type hostInputer struct {
}

// Name the Inputer name.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *hostInputer) Name() string {
	return "host_inputer"
}

// Run the input should not be blocked
func (cli *hostInputer) Run(ctx input.InputerContext) *input.InputerResult {

	host, err := api.CreateHost("0")
	if nil != err {
		fmt.Println("err:", err.Error())
	}

	// set the inner field
	host.SetInnerIP("192.168.1.1")
	host.SetOsBit("64")
	host.SetOsName("os-test")
	host.SetOsType(api.HostOSTypeLinux)
	host.SetSLA(api.HostSLALevel1)
	host.SetAssetID("host2122")
	host.SetInnerMac("1d2-3d-d-d")
	host.SetOperator("test_user")
	host.SetBakOperator("test_bak_user")
	host.SetCPU(5)
	host.SetCPUMhz(12)
	host.SetDisk(3456)
	host.SetMem(12334)
	host.SetCPUModule("cpu-xxx")

	// create a new field
	hostModel := host.GetModel()
	hostAttr := hostModel.CreateAttribute()
	hostAttr.SetID("host_field_id")
	hostAttr.SetName("host_field_id(test)")
	hostAttr.SetType(model.FieldTypeLongChar)
	err = hostAttr.Save()
	if nil != err {
		fmt.Println("err attr:", err)
	}

	// set the custom field
	host.SetValue("host_field_id", "test_custom_d")

	// save the host
	err = host.Save()
	if nil != err {
		fmt.Println("err:", err)
	}
	return nil

}

func (cli *hostInputer) Stop() error {
	return nil
}
