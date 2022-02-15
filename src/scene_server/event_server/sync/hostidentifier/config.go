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

package hostidentifier

import (
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
)

// RateLimiter push identifier rate limiter
type RateLimiter struct {
	Qps   int64
	Burst int64
}

// HostIdentifierConf host identifier config
type HostIdentifierConf struct {
	StartUp                bool
	BatchSyncIntervalHours int
	LinuxFileConf          *FileConf
	WinFileConf            *FileConf
	RateLimiter            *RateLimiter
}

// ParseIdentifierConf parser host identifier config
func ParseIdentifierConf() (*HostIdentifierConf, error) {
	startUp, err := cc.Bool("eventServer.hostIdentifier.startUp")
	if err != nil {
		blog.Errorf("get eventServer.hostIdentifier.startUp error, err: %v", err)
		return nil, err
	}

	if !startUp {
		blog.Warnf("eventServer.hostIdentifier.startUp is false, will not start sync host identifier")
		return &HostIdentifierConf{
			StartUp: startUp,
		}, nil
	}

	batchSyncIntervalHours, err := cc.Int("eventServer.hostIdentifier.batchSyncIntervalHours")
	if err != nil {
		blog.Errorf("get eventServer.hostIdentifier.batchSyncIntervalHours error, err: %v", err)
		return nil, err
	}

	winFileConfig, err := newHostIdentifierFileConf("windows")
	if err != nil {
		blog.Errorf("get eventServer hostIdentifier windows config error, err: %v", err)
		return nil, err
	}

	linuxFileConfig, err := newHostIdentifierFileConf("linux")
	if err != nil {
		blog.Errorf("get eventServer hostIdentifier linux config error, err: %v", err)
		return nil, err
	}

	qps, burst, err := getRateLimiterConfig()
	if err != nil {
		blog.Errorf("get evenServer hostIdentifier rate limiter config error, err: %v", err)
		return nil, err
	}

	rateLimiter := &RateLimiter{
		Qps:   qps,
		Burst: burst,
	}

	return &HostIdentifierConf{
		StartUp:                startUp,
		BatchSyncIntervalHours: batchSyncIntervalHours,
		LinuxFileConf:          linuxFileConfig,
		WinFileConf:            winFileConfig,
		RateLimiter:            rateLimiter,
	}, nil
}

func getRateLimiterConfig() (int64, int64, error) {
	qps, err := cc.Int64("eventServer.hostIdentifier.rateLimiter.qps")
	if err != nil {
		return 0, 0, err
	}

	burst, err := cc.Int64("eventServer.hostIdentifier.rateLimiter.burst")
	if err != nil {
		return 0, 0, err
	}
	return qps, burst, nil
}

// FileConf host identifier file config struct
type FileConf struct {
	FileName      string
	FilePath      string
	FileOwner     string
	FilePrivilege int32
}

// newHostIdentifierFileConf new host identifier file config
func newHostIdentifierFileConf(prefix string) (*FileConf, error) {
	fileName, err := cc.String("eventServer.hostIdentifier.fileName")
	if err != nil {
		blog.Errorf("get host identifier fileName error, err: %v", err)
		return nil, err
	}

	filePath, err := cc.String("eventServer.hostIdentifier." + prefix + ".filePath")
	if err != nil {
		blog.Errorf("get host identifier filePath error, err: %v", err)
		return nil, err
	}

	fileOwner, err := cc.String("eventServer.hostIdentifier." + prefix + ".fileOwner")
	if err != nil {
		blog.Errorf("get host identifier fileOwner error, err: %v", err)
		return nil, err
	}

	filePrivilege, err := cc.Int("eventServer.hostIdentifier." + prefix + ".filePrivilege")
	if err != nil {
		blog.Errorf("get host identifier filePrivilege error, err: %v", err)
		return nil, err
	}

	return &FileConf{
		FileName:      fileName,
		FilePath:      filePath,
		FileOwner:     fileOwner,
		FilePrivilege: int32(filePrivilege),
	}, nil
}

func (h *HostIdentifier) getHostIdentifierFileConf(osType string) *FileConf {
	switch osType {
	case common.HostOSTypeEnumWindows:
		return h.winFileConfig
	default:
		return h.linuxFileConfig
	}
}
