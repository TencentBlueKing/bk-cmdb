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

package condition

import (
	"reflect"

	"configcenter/src/common"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateCondition create a condition object
func CreateCondition() Condition {
	return &condition{}
}

// Condition condition interface
type Condition interface {
	SetPage(page types.MapStr) error
	SetStart(start int64)
	GetStart() int64
	SetLimit(limit int64)
	GetLimit() int64
	SetSort(sort string)
	GetSort() string
	SetFields(fields []string)
	GetFields() []string
	Field(fieldName string) Field
	Parse(data types.MapStr) error
	ToMapStr() types.MapStr
}

// Condition the condition definition
type condition struct {
	start        int64
	limit        int64
	sort         string
	fields       []Field
	filterFields []string
}

// SetPage set the page
func (cli *condition) SetPage(page types.MapStr) error {

	pageInfo := metadata.BasePage{}
	if err := page.MarshalJSONInto(&pageInfo); nil != err {
		return err
	}
	cli.start = int64(pageInfo.Start)
	cli.limit = int64(pageInfo.Limit)
	cli.sort = pageInfo.Sort
	return nil
}

func (cli *condition) SetFields(fields []string) {
	cli.filterFields = fields
}
func (cli *condition) GetFields() []string {
	return cli.filterFields
}

// Parse load the data into condition object
func (cli *condition) Parse(data types.MapStr) error {

	var fieldFunc func(tmpField *field, val interface{})
	fieldFunc = func(tmpField *field, val interface{}) {

		if nil == val {
			return
		}
		valType := reflect.TypeOf(val)

		switch valType.Kind() {
		default:
			tmpField.fieldValue = val

		case reflect.Map:

			tmpMap, err := types.NewFromInterface(val)
			if nil != err {
				panic(err) // very  serious
			}

			tmpMap.ForEach(func(key string, subVal interface{}) {
				switch key {

				default:
					tmp := &field{}
					tmp.fieldName = key
					tmp.opeartor = BKDBEQ
					tmp.condition = tmpField.condition
					fieldFunc(tmp, subVal)
					tmpField.fields = append(tmpField.fields, tmp)
				case BKDBEQ, BKDBGT, BKDBGTE, BKDBIN, BKDBNIN, BKDBLIKE, BKDBLT, BKDBLTE, BKDBNE, BKDBOR:
					tmpField.opeartor = key
					fieldFunc(tmpField, subVal)
				}

			})
		}

	}

	data.ForEach(func(key string, val interface{}) {

		tmpField := &field{}
		tmpField.condition = cli
		tmpField.fieldName = key
		tmpField.opeartor = BKDBEQ
		fieldFunc(tmpField, val)
		cli.fields = append(cli.fields, tmpField)

	})

	return nil
}

// SetStart set the start
func (cli *condition) SetStart(start int64) {
	cli.start = start
}

// GetStart return the start
func (cli *condition) GetStart() int64 {
	return cli.start
}

// SetLimit set the limit num
func (cli *condition) SetLimit(limit int64) {
	cli.limit = limit
}

// GetLimit return the limit num
func (cli *condition) GetLimit() int64 {
	if cli.limit <= 0 {
		return common.BKNoLimit
	}
	return cli.limit
}

// SetSort set the sort field
func (cli *condition) SetSort(sort string) {
	cli.sort = sort
}

// GetSort return the sort field
func (cli *condition) GetSort() string {
	return cli.sort
}

// CreateField create a field
func (cli *condition) Field(fieldName string) Field {
	field := &field{
		fieldName: fieldName,
		condition: cli,
	}
	cli.fields = append(cli.fields, field)
	return field
}

// ToMapStr to MapStr object
func (cli *condition) ToMapStr() types.MapStr {
	tmpResult := types.MapStr{}
	for _, item := range cli.fields {
		tmpResult.Merge(item.ToMapStr())
	}
	return tmpResult
}
