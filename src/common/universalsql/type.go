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

package universalsql

const (
	//Comparison Operator
	GT    string = "$gt"
	LT    string = "$lt"
	GTE   string = "$gte"
	LTE   string = "$lte"
	IN    string = "$in"
	NIN   string = "$nin"
	EQ    string = "$eq"
	NEQ   string = "$ne"
	REGEX string = "$regex"

	//Logic Operator
	AND string = "$and"
	OR  string = "$or"
	NOT string = "$not"
	NOR string = "$nor"

	//TODO:
	//Elements Operator
	EXISTS string = "$exists"
	TYPE   string = "$type"

	//Array Operator
	ALL       string = "$all"
	ELEMMATCH string = "$elemMatch"
	SIZE      string = "$size"
)
