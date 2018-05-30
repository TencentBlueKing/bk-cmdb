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
	"configcenter/src/framework/core/types"
	"github.com/tidwall/gjson"

	"fmt"
)

type eventSubscription struct {
}

func (cli *eventSubscription) run() error {
	// TODO：启动CMDB 3.0 事件订阅，对读取到的数据做加工整理成真是的时间对象并投递
	return nil
}

func (cli *eventSubscription) puts(data gjson.Result) (types.MapStr, error) {

	fmt.Println("puts:", data.String())
	return nil, nil
}

func (cli *eventSubscription) register(eventType types.EventType, eventFunc types.EventCallbackFunc) types.EventKey {
	return types.EventKey("")
}

func (cli *eventSubscription) unregister(eventKey types.EventKey) {

}
