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
	"io"
	//"fmt"
)

var _ Iterator = (*iterator)(nil)

type iterator struct {
	cond   common.Condition
	buffer []types.MapStr
	bufIdx int
}

func newModelIterator(cond common.Condition) (Iterator, error) {

	objIterator := &iterator{
		cond:   cond,
		buffer: make([]types.MapStr, 0),
	}

	items, err := client.GetClient().CCV3().Model().SearchObjects(cond)
	if nil != err {
		return nil, err
	}

	objIterator.buffer = items
	objIterator.bufIdx = 0
	if 0 == len(objIterator.buffer) {
		return nil, nil
	}

	return objIterator, nil
}

func (cli *iterator) ForEach(itemCallback func(item Model) error) error {

	for {

		item, err := cli.Next()
		if nil != err {
			if io.EOF == err{
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

func (cli *iterator) Next() (Model, error) {

	if len(cli.buffer) == cli.bufIdx {
		cli.bufIdx = 0
		return nil, io.EOF
	}

	tmpItem := cli.buffer[cli.bufIdx]
	cli.bufIdx++
	returnItem := &model{}
	common.SetValueToStructByTags(returnItem, tmpItem)

	return returnItem, nil
}
