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

package metadata

type AddDeviceResult struct {
	DeviceID uint64 `json:"device_id"`
}

type BatchAddDevice struct {
	Data []NetcollectDevice `json:"data"`
}

type BatchAddDeviceResult struct {
	Result   bool   `json:"result"`
	ErrMsg   string `json:"error_msg"`
	DeviceID uint64 `json:"device_id"`
}

type SearchNetDevice struct {
	Count uint64             `json:"count"`
	Info  []NetcollectDevice `json:"info"`
}

type SearchNetDeviceResult struct {
	BaseResp `json:",inline"`
	Data     SearchNetDevice `json:"data"`
}

type NetCollSearchParams struct {
	Page      BasePage        `json:"page,omitempty"`
	Fields    []string        `json:"fields,omitempty"`
	Condition []ConditionItem `json:"condition,omitempty"`
}

type DeleteNetDeviceBatchOpt struct {
	DeviceIDs []uint64 `json:"device_id"`
}

type AddNetPropertyResult struct {
	NetcollectPropertyID uint64 `json:"netcollect_property_id"`
}

type BatchAddNetPropertyResult struct {
	Result               bool   `json:"result"`
	ErrMsg               string `json:"error_msg"`
	NetcollectPropertyID uint64 `json:"netcollect_property_id"`
}

type BatchAddNetProperty struct {
	Data []NetcollectProperty `json:"data"`
}

type SearchNetProperty struct {
	Count uint64               `json:"count"`
	Info  []NetcollectProperty `json:"info"`
}

type SearchNetPropertyResult struct {
	BaseResp `json:",inline"`
	Data     SearchNetProperty `json:"data"`
}

type DeleteNetPropertyBatchOpt struct {
	NetcollectPropertyIDs []uint64 `json:"netcollect_property_id"`
}
