/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package findopt

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// ConvertToMongoOptions convert find many opt into mongo options array
func (m *Many) ConvertToMongoOptions() *options.FindOptions {

	option := &options.FindOptions{}

	if 0 != m.Limit {
		option.Limit = &m.Limit
	}

	option.Skip = &m.Skip

	sortD := primitive.D{}
	for _, sortItem := range m.Sort {

		if sortItem.Descending {
			sortD = append(sortD, primitive.E{Key: sortItem.Name, Value: -1})
			continue
		}

		sortD = append(sortD, primitive.E{Key: sortItem.Name, Value: 1})

	}
	if 0 != len(sortD) {
		option.Sort = sortD
	}

	fieldD := primitive.D{}
	for _, fieldItem := range m.Fields {
		if fieldItem.Hide {
			fieldD = append(fieldD, primitive.E{Key: fieldItem.Name, Value: 0})
			continue
		}

		fieldD = append(fieldD, primitive.E{Key: fieldItem.Name, Value: 1})
	}
	if 0 != len(fieldD) {
		option.Projection = fieldD
	}

	return option
}
