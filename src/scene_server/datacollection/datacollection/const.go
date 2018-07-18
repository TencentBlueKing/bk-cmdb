package datacollection

import (
	"configcenter/src/common"
)

const (
	RedisSnapKeyPrefix        = common.BKCacheKeyV3Prefix + "snapshot:"
	RedisDisKeyPrefix         = common.BKCacheKeyV3Prefix + "discover:"
	MasterProcLockKey         = common.BKCacheKeyV3Prefix + "snapshot:masterlock"
	MasterDisLockKey          = common.BKCacheKeyV3Prefix + "discover:masterlock"
	RedisSnapKeyChannelStatus = common.BKCacheKeyV3Prefix + "snapshot:channelstatus"
)

const (
	MaxSnapSize     = 2000
	MaxDiscoverSize = 1000
)

const (
	DiscoverChan = "discover"
	SnapShotChan = "snapshot"
)
