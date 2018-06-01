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
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output"
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/types"
)

type eventHost struct {
	outputerMgr output.Manager
}

func (cli *eventHost) constructEvent(hostID int, data types.MapStr) (types.MapStr, error) {

	hostModel, err := cli.outputerMgr.GetModel("0", "bk_host_manage", "host")
	if nil != err {
		log.Errorf("failed to serch the host model, %s", err.Error())
		return nil, err
	}

	cond := common.CreateCondition()
	cond.Field("bk_host_id").Eq(hostID)
	hostIter, err := cli.outputerMgr.FindInstsByCondition(hostModel, cond)
	if nil != err {
		log.Errorf("failed to serch the host, %s", err.Error())
		return nil, err
	}

	result := types.MapStr{}
	err = hostIter.ForEach(func(item inst.Inst) error {

		vals, err := item.GetValues()
		if nil != err {
			log.Errorf("failed to get host info, %s", err.Error())
			return err
		}

		//fmt.Println("vals:", vals)

		biz, err := vals.MapStrArray("biz")
		if nil != err {
			log.Errorf("failed to get biz, %s", err.Error())
			return err
		}

		set, err := vals.MapStrArray("set")
		if nil != err {
			log.Errorf("failed to get set, %s", err.Error())
			return err
		}
		vals.Remove("set")

		module, err := vals.MapStrArray("module")
		if nil != err {
			log.Errorf("failed to get module, %s", err.Error())
			return err
		}
		vals.Remove("module")

		for _, bizItem := range biz {

			bizID, err := bizItem.Int("bk_biz_id")
			if nil != err {
				log.Errorf("failed to get biz id, %s", err.Error())
				return err
			}
			setArr := make([]types.MapStr, 0)
			for _, setItem := range set {

				setBizID, err := setItem.Int("bk_biz_id")
				if nil != err {
					log.Errorf("failed to get biz id, %s", err.Error())
					return err
				}

				if setBizID != bizID {
					continue
				}

				setArr = append(setArr, setItem)
				moduleArr := make([]types.MapStr, 0)

				setID, err := setItem.Int("bk_set_id")
				if nil != err {
					log.Errorf("failed to get set id, %s", err.Error())
					return err
				}
				for _, moduleItem := range module {

					moduleSetID, err := moduleItem.Int("bk_set_id")
					if nil != err {
						log.Errorf("failed to get set id, %s", err.Error())
						return err
					}

					if setID != moduleSetID {
						continue
					}

					moduleArr = append(moduleArr, moduleItem)
				}
				setItem.Set("module", moduleArr)

			}
			bizItem.Set("set", setArr)

		}

		result = vals

		return nil
	})

	if nil != err {
		log.Errorf("failed to for each the data, %s", err.Error())
		return nil, err
	}

	return result, nil
}

func (cli *eventHost) parse(data types.MapStr) (*types.Event, error) {

	dataArr, err := data.MapStrArray("data")
	if nil != err {
		return nil, err
	}

	eve := &types.Event{}
	for _, dataItem := range dataArr {

		curHost, err := dataItem.MapStr("cur_data")

		if nil != err {
			log.Errorf("failed to get the curr data, %s", err.Error())
			return nil, err
		}

		hostID, err := curHost.Int("bk_host_id")
		if nil != err {
			log.Errorf("failed to get the host id, %s", err.Error())
			return nil, err
		}

		eventItem, err := cli.constructEvent(hostID, curHost)
		if nil != err {
			log.Errorf("failed to construct the host event, %s", err.Error())
			return nil, err
		}

		eve.AddData(eventItem)

	}

	return eve, nil
}
