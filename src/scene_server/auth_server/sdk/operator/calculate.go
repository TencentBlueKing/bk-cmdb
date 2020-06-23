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

package operator

// Calculator is used to calculate the logic value of a list/array of
// Operator with different of "Calculator Type" instance.
type Calculator interface {
	// name of the calculator
	Name() string

	// the calculated resulted with multiple Policy.
	Result(p []*Policy) (bool, error)
}

const And = OperType("AND")

type AndOper OperType

func (a *AndOper) Name() string {
	return "AND"
}

func (a *AndOper) Result(p []*Policy) (bool, error) {
	return true, nil
}

const Or = OperType("OR")

type OrOper OperType

func (o *OrOper) Name() string {
	return "OR"
}
