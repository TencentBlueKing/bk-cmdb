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
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ Iterator = (*hostIterator)(nil)

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

	items, err := client.GetClient().CCV3().Host().SearchHost(cond)
	if nil != err {
		return nil, err
	}

	grpIterator.buffer = items
	grpIterator.bufIdx = 0
	if 0 == len(grpIterator.buffer) {
		return nil, nil
	}

	return grpIterator, nil
}

func (cli *hostIterator) ForEach(itemCallback func(item Inst) error) (err error) {
	var item Inst
	for {

		item, err = cli.Next()
		if nil != err {
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

func (cli *hostIterator) Next() (Inst, error) {
	if len(cli.buffer) == cli.bufIdx {

		cli.cond.SetStart(cli.bufIdx)

		existItems, err := client.GetClient().CCV3().Module().SearchModules(cli.cond)
		if nil != err {
			return nil, err
		}

		if 0 == len(existItems) {
			cli.bufIdx = 0
			return nil, nil
		}

		cli.buffer = append(cli.buffer, existItems...)
	}

	tmpItem := cli.buffer[cli.bufIdx]
	cli.bufIdx++

	returnItem := &host{
		target: cli.targetModel,
		datas:  tmpItem,
	}
	return returnItem, nil
}
