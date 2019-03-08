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

package types

import (
	"time"
)

// Event the cmdb event definition
type Event struct {
	action     string
	actionTime time.Time
	curDatas   MapStr
	preDatas   MapStr
}

// SetAction set the event action
func (cli *Event) SetAction(action string) {
	cli.action = action
}

// GetAction return the event action
func (cli *Event) GetAction() string {
	return cli.action
}

// SetActionTime set the action time
func (cli *Event) SetActionTime(tm time.Time) {
	cli.actionTime = tm
}

// GetActionTime get the action time
func (cli *Event) GetActionTime() time.Time {
	return cli.actionTime
}

// SetCurrData set the event data
func (cli *Event) SetCurrData(data MapStr) {
	cli.curDatas = data
}

// GetCurrData get the event data
func (cli *Event) GetCurrData() MapStr {
	return cli.curDatas
}

// SetPreData set the event data
func (cli *Event) SetPreData(data MapStr) {
	cli.preDatas = data
}

// GetPreData get the event data
func (cli *Event) GetPreData() MapStr {
	return cli.preDatas
}
