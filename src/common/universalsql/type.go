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
	// GT TODO
	// Comparison Operator
	GT string = "$gt"
	// LT TODO
	LT string = "$lt"
	// GTE TODO
	GTE string = "$gte"
	// LTE TODO
	LTE string = "$lte"
	// IN TODO
	IN string = "$in"
	// NIN TODO
	NIN string = "$nin"
	// EQ TODO
	EQ string = "$eq"
	// NEQ TODO
	NEQ string = "$ne"
	// REGEX TODO
	REGEX string = "$regex"

	// AND TODO
	// Logic Operator
	AND string = "$and"
	// OR TODO
	OR string = "$or"
	// NOT TODO
	NOT string = "$not"
	// NOR TODO
	NOR string = "$nor"

	// EXISTS TODO
	// TODO:
	// Elements Operator
	EXISTS string = "$exists"
	// TYPE TODO
	TYPE string = "$type"

	// ALL TODO
	// Array Operator
	ALL string = "$all"
	// ELEMMATCH TODO
	ELEMMATCH string = "$elemMatch"
	// SIZE TODO
	SIZE string = "$size"
)
