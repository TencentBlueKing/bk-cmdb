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

// Package process TODO
package process

// ValidateProcessBindIPEmptyHook validate process bind ip value which is empty, right now empty bind ip is valid
func ValidateProcessBindIPEmptyHook() error {
	return nil
}

// ValidateProcessBindIPHook validate if process bind ip value is valid, the value is not empty
func ValidateProcessBindIPHook(bindIP string) error {
	return nil
}

// ValidateProcessBindProtocolHook validate if process bind protocol value is valid, the value is not empty
func ValidateProcessBindProtocolHook(bindProtocol string) error {
	return nil
}

// NeedIPv6OptionsHook returns if process ipv6 options needs to be supported
func NeedIPv6OptionsHook() bool {
	return true
}
