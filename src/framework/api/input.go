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
	"time"

	"configcenter/src/framework/core/input"
)

var (
	inputers = make([]input.InputerParams, 0)
)

// RegisterInputer register a inputer, only execute once
func RegisterInputer(inputer input.Inputer) {

	inputers = append(inputers, input.InputerParams{
		Target: inputer,
		Kind:   input.ExecuteOnce,
		Putter: nil,
	})
}

// RegisterFrequencyInputer execute a non-blocking timing inputer, execute periodically.
func RegisterFrequencyInputer(inputer input.Inputer, frequency time.Duration) {

	if frequency < time.Minute*5 {
		frequency = time.Minute * 5
	}

	inputers = append(inputers, input.InputerParams{
		IsTiming:  true,
		Frequency: frequency,
		Target:    inputer,
		Kind:      input.ExecuteOnce,
		Putter:    nil,
	})
}
