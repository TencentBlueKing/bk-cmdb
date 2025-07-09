/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package collections

// Analyzer is common collection analyzer interface.
type Analyzer interface {
	// Analyze analyzes message from collectors or delay queue.
	Analyze(message *string, sourceType string) (bool, error)

	// Hash returns a hash value of the input message string.
	Hash(cloudid, ip string) (string, error)

	// Mock returns mock message that could be analyzed by the Analyzer.
	Mock() string
}

// Porter is common porter interface. It handles
// message from collectors base on Analyzer.
type Porter interface {
	// Name returns name of the Porter.
	Name() string

	// Run runs the Porter.
	Run() error

	// Mock supports mock service in Porter.
	Mock() error
}
