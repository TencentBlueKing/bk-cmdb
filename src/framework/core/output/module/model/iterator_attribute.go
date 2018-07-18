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

package model

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/types"
	//"fmt"
	"io"
)

var _ AttributeIterator = (*attributeIterator)(nil)

type attributeIterator struct {
	cond   common.Condition
	buffer []types.MapStr
	bufIdx int
}

func newAttributeIterator(supplierAccount string, cond common.Condition) (AttributeIterator, error) {

	attrIterator := &attributeIterator{
		cond:   cond,
		buffer: make([]types.MapStr, 0),
	}

	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: supplierAccount}).Attribute().SearchObjectAttributes(cond)
	if nil != err {
		return nil, err
	}

	attrIterator.buffer = items
	attrIterator.bufIdx = 0
	if 0 == len(attrIterator.buffer) {
		return nil, io.EOF
	}

	return attrIterator, nil
}

func (cli *attributeIterator) ForEach(itemCallback func(item Attribute) error) error {

	for {

		item, err := cli.Next()
		if nil != err {
			if io.EOF == err {
				return nil
			}
			return err
		}

		if nil == item {
			return nil
		}
		err = itemCallback(item)
		if nil != err {
			return err
		}
	}

}

func (cli *attributeIterator) Next() (Attribute, error) {

	if len(cli.buffer) == cli.bufIdx {
		cli.bufIdx = 0
		return nil, io.EOF
	}

	tmpItem := cli.buffer[cli.bufIdx]
	cli.bufIdx++
	returnItem := &attribute{}
	common.SetValueToStructByTags(returnItem, tmpItem)

	return returnItem, nil
}
