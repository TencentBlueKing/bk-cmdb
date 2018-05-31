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

package manager

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output"
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/types"

	"fmt"
)

type eventHost struct {
	outputerMgr output.Manager
}

func (cli *eventHost) constructEvent(hostID int, data types.MapStr) (types.MapStr, error) {

	hostModel, err := cli.outputerMgr.GetModel("0", "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}

	fmt.Println("hostid:", hostID, "data:", data, "host model:", hostModel)

	cond := common.CreateCondition()
	cond.Field("bk_host_id").Eq(29)
	hostIter, err := cli.outputerMgr.FindInstsByCondition(hostModel, cond)
	if nil != err {
		return nil, err
	}

	err = hostIter.ForEach(func(item inst.Inst) error {

		vals, err := item.GetValues()
		if nil != err {
			return err
		}

		vals.ForEach(func(key string, value interface{}) {

			if key == "biz" {
				fmt.Println("host:", key, "val:", value)
			}
		})

		return nil
	})

	if nil != err {
		return nil, err
	}

	return nil, nil
}

func (cli *eventHost) parse(data types.MapStr) (types.MapStr, error) {

	dataArr, err := data.MapStrArray("data")
	if nil != err {
		return nil, err
	}

	for _, dataItem := range dataArr {

		curHost, err := dataItem.MapStr("cur_data")

		if nil != err {
			return nil, err
		}

		hostID, err := curHost.Int("bk_host_id")
		if nil != err {
			return nil, err
		}

		eventItem, err := cli.constructEvent(hostID, curHost)
		if nil != err {
			return nil, err
		}

		_ = eventItem

	}

	return nil, nil
}
