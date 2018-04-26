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
 
package output

import "configcenter/src/framework/core/types"

type customWrapper struct {
	name    string
	runFunc func(data types.MapStr) error
}

// Name the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *customWrapper) Name() string {
	return cli.name
}

// Run the output main loop. This should block until singnalled to stop by invocation of the Stop() method.
func (cli *customWrapper) Put(data types.MapStr) error {
	return cli.runFunc(data)
}

// Stop is the invoked to signal that the Run() method should its execution.
// It will be invoked at most once.
func (cli *customWrapper) Stop() error {
	// only compatible with the Outputer interface
	return nil
}
