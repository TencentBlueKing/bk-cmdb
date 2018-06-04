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
	"configcenter/src/framework/core/types"

	"context"
	"errors"
)

type eventRegister struct {
	eventKey  types.EventKey
	eventType types.EventType
	callback  types.EventCallbackFunc
	datas     types.MapStr
}

type eventSubscription struct {
	datas             chan types.MapStr
	registers         map[types.EventType][]*eventRegister
	setMgr            *eventSet
	moduleMgr         *eventModule
	moduleTransferMgr *eventModuleTransfer
	hostMgr           *eventHost
	hostIdentifierMgr *eventHostIdentifier
	businessMgr       *eventBusiness
	instMgr           *eventInst
}

func (cli *eventSubscription) setOutputer(output output.Manager) {

	cli.hostMgr.outputerMgr = output
	cli.setMgr.outputerMgr = output
	cli.moduleMgr.outputerMgr = output
	cli.hostIdentifierMgr.outputerMgr = output
	cli.businessMgr.outputerMgr = output

}

func (cli *eventSubscription) sendEvent(eveType types.EventType, eveData []*types.Event) error {
	if items, ok := cli.registers[eveType]; ok {

		for _, eveItem := range items {
			if nil != eveItem.callback {
				if err := eveItem.callback(eveData); nil != err {
					log.Errorf("failed to send the event, %s", err.Error())
				}
			}
		}

		return nil
	}
	log.Infof("not support the event type %s", eveType)
	return nil
}

func (cli *eventSubscription) run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-cli.datas:
			{
				objID := msg.String("obj_type")

				switch objID {
				case EventHost:
					if dataEve, err := cli.hostMgr.parse(msg); nil != err {
						log.Errorf("failed to parse the host event, error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventHostType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}
				case EventBusiness:
					if dataEve, err := cli.businessMgr.parse(msg); nil != err {
						log.Errorf("failed to parse the business event, error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventBusinessType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}
				case EventModule:
					if dataEve, err := cli.moduleMgr.parse(msg); nil != err {
						log.Errorf("failed to parse the module event, error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventBusinessType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}
				case EventModuleTransfer:
					if dataEve, err := cli.moduleTransferMgr.parse(msg); nil != err {
						log.Errorf("failed to parse the module transfer event, error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventBusinessType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}
				case EventSet:
					if dataEve, err := cli.setMgr.parse(msg); nil != err {
						log.Errorf("failed to parse the set , error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventSetType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}

				case EventHostIdentifier:
					if dataEve, err := cli.hostIdentifierMgr.parse(msg); nil != err {
						log.Errorf("failed to parse hostindentifier, error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventHostIdentifierType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}

				default:
					if dataEve, err := cli.instMgr.parse(msg); nil != err {
						log.Errorf("failed to parse inst event, error %s", err.Error())
					} else {
						if err := cli.sendEvent(types.EventInstType, dataEve); nil != err {
							log.Errorf("failed to send event, %s", err.Error())
						}
					}
				}
			}
		}
	}
}

func (cli *eventSubscription) puts(data types.MapStr) (types.MapStr, error) {

	select {
	default:
		return nil, errors.New("the event queue is full")
	case cli.datas <- data:
		return nil, nil
	}
}

func (cli *eventSubscription) register(key types.EventKey, eventType types.EventType, eventFunc types.EventCallbackFunc) types.EventKey {

	if 0 == len(key) {
		key = types.EventKey(common.UUID())
	}

	regEve := &eventRegister{
		eventType: eventType,
		eventKey:  types.EventKey(key),
		callback:  eventFunc,
	}

	if items, ok := cli.registers[eventType]; ok {

		cli.registers[eventType] = append(items, regEve)

		return types.EventKey(key)
	}

	regs := make([]*eventRegister, 0)
	regs = append(regs, regEve)
	cli.registers[eventType] = regs

	return key
}

func (cli *eventSubscription) unregister(eventKey types.EventKey) {

	delSlice := func(s []*eventRegister, i int) []*eventRegister {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}

	for eveType, items := range cli.registers {
		for idx, item := range items {
			if 0 == item.eventKey.Compare(eventKey) {
				cli.registers[eveType] = delSlice(items, idx)
			}
		}
	}
}
