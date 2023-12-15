/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package common

// DeploymentMethod is the deployment method
type DeploymentMethod string

// String get string value
func (d *DeploymentMethod) String() string {
	return string(*d)
}

// Set value
func (d *DeploymentMethod) Set(s string) error {
	*d = DeploymentMethod(s)
	return nil
}

// Type returns value type
func (d *DeploymentMethod) Type() string {
	return "DeploymentMethod"
}

const (
	// OpenSourceDeployment is the open-source deployment method, do not rely on api gateway
	OpenSourceDeployment DeploymentMethod = "open_source"
	// BluekingDeployment is the deployment method for blueking, using api gateway
	BluekingDeployment DeploymentMethod = "blueking"
)
