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

	"errors"
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
	hostIter, err := cli.outputerMgr.InstOperation().FindHostsByCondition(hostModel, cond)
	if nil != err {
		log.Errorf("failed to serch the host, %s", err.Error())
		return nil, err
	}

	err = hostIter.ForEach(func(item inst.HostInterface) error {

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

		data.Set("biz", biz)

		return nil
	})

	if nil != err {
		log.Errorf("failed to for each the data, %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (cli *eventHost) getEvent(dataKey string, dataItem types.MapStr) (types.MapStr, error) {

	curHost, err := dataItem.MapStr(dataKey)

	if nil != err {
		log.Errorf("failed to get the curr data, %s", err.Error())
		return nil, err
	}

	hostID, err := curHost.Int("bk_host_id")
	if nil != err {
		log.Errorf("failed to get the host id, %s", err.Error())
		return nil, err
	}

	return cli.constructEvent(hostID, curHost)
}

func (cli *eventHost) parse(data types.MapStr) ([]*types.Event, error) {

	dataArr, err := data.MapStrArray("data")
	if nil != err {
		return nil, err
	}

	tm, err := data.Time("action_time")
	if nil != err {
		log.Error("failed to get action time")
		return nil, err
	}

	action := data.String("action")
	if 0 == len(action) {
		log.Error("the event action is not set")
		return nil, errors.New("the event action is not set")
	}

	eves := make([]*types.Event, 0)
	for _, dataItem := range dataArr {

		currEvent, err := cli.getEvent("cur_data", dataItem)
		if nil != err {
			log.Errorf("failed to get the current host event,%s", err.Error())
			return nil, err
		}

		prevEvent, err := cli.getEvent("pre_data", dataItem)
		if nil != err {
			log.Errorf("failed to get the prev host event,%s", err.Error())
			return nil, err
		}

		ev := &types.Event{}
		ev.SetCurrData(currEvent)
		ev.SetPreData(prevEvent)
		ev.SetActionTime(*tm)
		ev.SetAction(action)
		eves = append(eves, ev)

	}

	return eves, nil
}
