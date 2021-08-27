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

package discovery

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/registerdiscover"
)

func newServerMaster(disc *registerdiscover.RegDiscv, path, name string) *master {
	master := &master{
		rd:       disc,
		path:     path,
		name:     name,
		isMaster: false,
	}

	master.run()
	return master
}

type master struct {
	rd       *registerdiscover.RegDiscv
	path     string
	name     string
	isMaster bool
}

// IsMaster judge whether current service is master or slave
func (m *master) IsMaster() bool {
	return m.isMaster
}

func (m *master) run() {
	go m.campaign()
}

// campaign elects to be master
func (m *master) campaign() {
	val := common.GetServerInfo().RegisterIP
	for {
		// block until it is elected, an error occurs, or the context is cancelled
		if err := m.rd.Campaign(m.path, val); err != nil {
			m.isMaster = false
			blog.Errorf("discovery campaign master err: %v", err)
			// try campaign 3 seconds later
			time.Sleep(3 * time.Second)
			continue
		}
		// if campaign return no error, it is elected as master
		m.isMaster = true
	}

}
