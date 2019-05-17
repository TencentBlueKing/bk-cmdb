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

package input

import (
	"configcenter/src/framework/common"
)

// create a new inputer key
func makeInputerKey() InputerKey {
	return InputerKey(common.UUID())
}

// checkWorkerExists check whether the inputer exists
func inputerExists(target MapInputer, key InputerKey) bool {
	_, ok := target[key]
	return ok
}

// deleteInputer delete a inputer from MapInputer
func deleteInputer(target MapInputer, key InputerKey) bool {

	if inputerExists(target, key) {
		delete(target, key)
	}

	return true
}
