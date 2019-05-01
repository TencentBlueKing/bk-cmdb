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

package datacollection

import (
	"time"

	"configcenter/src/common"
)

const (
	RedisDisKeyPrefix               = common.BKCacheKeyV3Prefix + "discover:"
	MasterProcLockKey               = common.BKCacheKeyV3Prefix + "snapshot:masterlock"
	MasterDisLockKey                = common.BKCacheKeyV3Prefix + "discover:masterlock"
	MasterNetLockKey                = common.BKCacheKeyV3Prefix + "netcollect:masterlock"
	RedisSnapKeyChannelStatus       = common.BKCacheKeyV3Prefix + "snapshot:channelstatus"
	RedisNetcollectKeyChannelStatus = common.BKCacheKeyV3Prefix + "netcollect:channelstatus"
)

const (
	MaxSnapSize       = 2000
	MaxNetcollectSize = 1000
	MaxDiscoverSize   = 1000
)

var masterProcLockLiveTime = time.Second * 10

const (
	DiscoverChan = "discover"
	SnapShotChan = "snapshot"
)

type Analyzer interface {
	Analyze(mesg string) error
}

type Porter interface {
	Name() string
	Run() error
	Mock(string) error
}
