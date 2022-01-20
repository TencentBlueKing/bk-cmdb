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
