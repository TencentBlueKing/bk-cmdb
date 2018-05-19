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
 
package output

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/types"
)

func (cli *manager) AddOutputer(target Outputer) OutputerKey {

	cli.outputerLock.Lock()

	key := OutputerKey(common.UUID())
	cli.outputers[key] = target

	cli.outputerLock.Unlock()

	return key
}

func (cli *manager) RemoveOutputer(key OutputerKey) {

	cli.outputerLock.Lock()

	if item, ok := cli.outputers[key]; ok {
		if err := item.Stop(); nil != err {
			log.Errorf("failed to stop the outputer (%s), stop to reove it, error info is %s", item.Name(), err.Error())
		} else {
			log.Infof("remove the outputer(%s)", item.Name())
			delete(cli.outputers, key)
		}
	}

	cli.outputerLock.Unlock()
}
func (cli *manager) FetchOutputer(key OutputerKey) Puter {

	cli.outputerLock.RLock()
	defer func() {
		cli.outputerLock.RUnlock()
	}()

	if item, ok := cli.outputers[key]; ok {
		return item
	}

	return nil
}
func (cli *manager) CreateCustomOutputer(name string, run func(data types.MapStr) error) (OutputerKey, Puter) {

	log.Infof("creater custom outputer:%s", name)
	wrapper := &customWrapper{
		name:    name,
		runFunc: run,
	}

	key := cli.AddOutputer(wrapper)

	return key, wrapper
}
