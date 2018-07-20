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

package inst

import (
	"io"

	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ HostIterator = (*hostIterator)(nil)

// HostIterator the iterator interface for the host
type HostIterator interface {
	Next() (HostInterface, error)
	ForEach(callbackItem func(item HostInterface) error) error
}

type hostIterator struct {
	targetModel model.Model
	cond        common.Condition
	buffer      []types.MapStr
	bufIdx      int
}

func newHostIterator(target model.Model, cond common.Condition) (*hostIterator, error) {
	grpIterator := &hostIterator{
		targetModel: target,
		cond:        cond,
		buffer:      make([]types.MapStr, 0),
	}

	grpIterator.cond.SetLimit(DefaultLimit)
	grpIterator.cond.SetStart(grpIterator.bufIdx)

	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: target.GetSupplierAccount()}).Host().SearchHost(cond)
	if nil != err {
		return nil, err
	}

	grpIterator.buffer = items
	grpIterator.bufIdx = 0

	grpIterator.buffer = append(grpIterator.buffer, items...)

	return grpIterator, nil
}

func (cli *hostIterator) ForEach(itemCallback func(item HostInterface) error) (err error) {
	var item HostInterface
	for {

		item, err = cli.Next()
		if nil != err {
			if io.EOF == err {
				err = nil
			}
			break
		}

		if nil == item {
			break
		}

		err = itemCallback(item)
		if nil != err {
			break
		}
	}
	cli.bufIdx = 0
	return err
}

func (cli *hostIterator) Next() (HostInterface, error) {
	if len(cli.buffer) == cli.bufIdx {

		cli.cond.SetStart(cli.bufIdx)

		existItems, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.targetModel.GetSupplierAccount()}).Host().SearchHost(cli.cond)
		if nil != err {
			return nil, err
		}
		//fmt.Println("the err:", err)
		if 0 == len(existItems) {
			cli.bufIdx = 0
			return nil, io.EOF
		}

		cli.buffer = append(cli.buffer, existItems...)
	}

	tmpItem := cli.buffer[cli.bufIdx]
	cli.bufIdx++

	returnItem := &host{
		target: cli.targetModel,
		datas:  tmpItem,
	}

	if err := returnItem.reset(); nil != err {
		return nil, err
	}

	return returnItem, nil
}
