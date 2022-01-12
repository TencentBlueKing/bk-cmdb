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

const (
	// windowOSType window os type
	windowOSType = "windows"
	// linuxOSType linux os type
	linuxOSType = "linux"
	// defaultHostIdentifierFileName default host identifier fileName
	defaultHostIdentifierFileName = "hostid"
	// defaultHostIdentifierLinuxFilePath default host identifier linux filePath
	defaultHostIdentifierLinuxFilePath = "/var/lib/gse/host"
	// defaultHostIdentifierLinuxFileOwner default host identifier linux fileOwner
	defaultHostIdentifierLinuxFileOwner = "root"
	// defaultHostIdentifierLinuxFileRight default host identifier linux fileRight
	defaultHostIdentifierLinuxFileRight = 644
	// defaultHostIdentifierWindowsFilePath default host identifier windows filePath
	defaultHostIdentifierWindowsFilePath = "c:/gse/data/host"
	// defaultHostIdentifierWindowsFileOwner default host identifier windows fileOwner
	defaultHostIdentifierWindowsFileOwner = "root"
	// defaultHostIdentifierWindowsFileRight default host identifier windows fileRight
	defaultHostIdentifierWindowsFileRight = 644
	// defaultRateLimiterQPS default rate limiter QPS
	defaultRateLimiterQPS = 200
	// defaultRateLimiterBurst default rate limiter burst
	defaultRateLimiterBurst = 200
)

// FileConf host identifier file config struct
type FileConf struct {
	FileName  string
	FilePath  string
	FileOwner string
	FileRight int32
}

func getHostIdentifierFileConf(osType string) *FileConf {
	var prefix string
	switch osType {
	case common.HostOSTypeEnumWindows:
		prefix = windowOSType
	default:
		prefix = linuxOSType
	}

	fileName, err := cc.String("eventServer.hostIdentifier.fileName")
	if err == nil {
		fileName = defaultHostIdentifierFileName
	}

	filePath, err := cc.String("eventServer.hostIdentifier." + prefix + ".filePath")
	if err != nil {
		blog.Errorf("get host identifier filePath error, err: %v", err)
		filePath = getDefaultHostIdentifierFilePath(osType)
	}
	fileOwner, err := cc.String("eventServer.hostIdentifier." + prefix + ".fileOwner")
	if err != nil {
		blog.Errorf("get host identifier fileOwner error, err: %v", err)
		fileOwner = getDefaultHostIdentifierFileOwner(osType)
	}
	fileRight, err := cc.Int("eventServer.hostIdentifier." + prefix + ".fileRight")
	if err != nil {
		blog.Errorf("get host identifier fileRight error, err: %v", err)
		fileRight = getDefaultHostIdentifierFileRight(osType)
	}

	return &FileConf{
		FileName:  fileName,
		FilePath:  filePath,
		FileOwner: fileOwner,
		FileRight: int32(fileRight),
	}
}

func getDefaultHostIdentifierFilePath(osType string) string {
	switch osType {
	case common.HostOSTypeEnumWindows:
		return defaultHostIdentifierWindowsFilePath
	default:
		return defaultHostIdentifierLinuxFilePath
	}
}

func getDefaultHostIdentifierFileOwner(osType string) string {
	switch osType {
	case common.HostOSTypeEnumWindows:
		return defaultHostIdentifierWindowsFileOwner
	default:
		return defaultHostIdentifierLinuxFileOwner
	}
}

func getDefaultHostIdentifierFileRight(osType string) int {
	switch osType {
	case common.HostOSTypeEnumWindows:
		return defaultHostIdentifierWindowsFileRight
	default:
		return defaultHostIdentifierLinuxFileRight
	}
}

func getRateLimiterConfig() (int64, int64) {
	qps, err := cc.Int64("eventServer.hostIdentifier.rateLimiter.qps")
	if err != nil {
		blog.Errorf("can't find the value of eventServer.hostIdentifier.rateLimiter.qps settings, "+
			"set the default value: %s", defaultRateLimiterQPS)
		qps = defaultRateLimiterQPS
	}
	burst, err := cc.Int64("eventServer.hostIdentifier.rateLimiter.burst")
	if err != nil {
		blog.Errorf("can't find the value of eventServer.hostIdentifier.rateLimiter.burst setting,"+
			"set the default value: %s", defaultRateLimiterBurst)
		burst = defaultRateLimiterBurst
	}
	return qps, burst
}
