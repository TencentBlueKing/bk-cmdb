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

package filter

type OperatorType string

const (
	// Matches values that are equal to a specified value.
	EQ OperatorType = "$eq"

	// Matches values that are greater than a specified value.
	GT OperatorType = "$gt"

	// 	Matches values that are greater than or equal to a specified value.
	GTE OperatorType = "$gte"

	// Matches any of the values specified in an array.
	IN OperatorType = "$in"

	// Matches values that are less than a specified value.
	LT OperatorType = "$lt"

	// 	Matches values that are less than or equal to a specified value.
	LTE OperatorType = "$lte"

	// Matches all values that are not equal to a specified value.
	NE OperatorType = "$ne"

	// Matches none of the values specified in an array.
	NIN OperatorType = "$nin"

	// Matches documents that have the specified field.
	// Exists OperatorType = "$exists"

)
