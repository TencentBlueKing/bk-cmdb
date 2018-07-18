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
	"io"

	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/types"
)

var _ ClassificationIterator = (*classificationIterator)(nil)

type classificationIterator struct {
	cond   common.Condition
	buffer []types.MapStr
	bufIdx int
}

func newClassificationIterator(supplierAccount string, cond common.Condition) (ClassificationIterator, error) {

	clsIterator := &classificationIterator{
		cond:   cond,
		buffer: make([]types.MapStr, 0),
	}

	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: supplierAccount}).Classification().SearchClassifications(cond)
	if nil != err {
		//fmt.Println("err:", err.Error(), items)
		return nil, err
	}

	clsIterator.buffer = items
	clsIterator.bufIdx = 0
	if 0 == len(clsIterator.buffer) {
		return nil, io.EOF
	}

	return clsIterator, nil
}

func (cli *classificationIterator) ForEach(itemCallback func(item Classification) error) error {

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

func (cli *classificationIterator) Next() (Classification, error) {

	if len(cli.buffer) == cli.bufIdx {
		cli.bufIdx = 0
		return nil, io.EOF
	}

	tmpItem := cli.buffer[cli.bufIdx]
	cli.bufIdx++
	returnItem := &classification{}
	common.SetValueToStructByTags(returnItem, tmpItem)

	return returnItem, nil
}
