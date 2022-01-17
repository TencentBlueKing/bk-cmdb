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

package logplatform

import (
	"errors"
	"fmt"
	"time"

	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
)

var (
	OpenTelemetryCfg = new(OpenTelemetryConfig)
)

type OpenTelemetryConfig struct {
	// 表示是否开启日志平台openTelemetry跟踪链接入相关功能，布尔值, 默认值为false不开启
	Enable bool
	// 日志平台openTelemetry跟踪链功能的自定义上报服务地址
	EndPoint string
	// 日志平台openTelemetry跟踪链功能的上报data_id
	BkDataID int64
}

// InitOpenTelemetryConfig init openTelemetry config
func InitOpenTelemetryConfig() error {
	var err error
	maxCnt := 100
	cnt := 0
	for !cc.IsExist("logPlatform.openTelemetry") && cnt < maxCnt {
		blog.Infof("waiting logPlatform.openTelemetry config to be init")
		cnt++
		time.Sleep(time.Millisecond * 300)
	}

	if cnt == maxCnt {
		blog.Infof("init openTelemetry failed, no openTelemetry config is found, " +
			"the config 'logPlatform.openTelemetry' must exist")
		return fmt.Errorf("init openTelemetry failed, " +
			"no openTelemetry config is found, the config 'logPlatform.openTelemetry' must exist")
	}

	OpenTelemetryCfg.Enable, err = cc.Bool("logPlatform.openTelemetry.enable")
	if err != nil {
		blog.Errorf("init openTelemetry failed, openTelemetry.enable err: %v", err)
		return errors.New("config logPlatform.openTelemetry is not found")
	}

	//如果不需要开启OpenTelemetry，那么后续没有必要再检查配置
	if !OpenTelemetryCfg.Enable {
		return nil
	}

	OpenTelemetryCfg.EndPoint, err = cc.String("logPlatform.openTelemetry.endpoint")
	if err != nil {
		blog.Errorf("init openTelemetry failed, err: %v", err)
		return errors.New("logPlatform.openTelemetry.endpoint is not found")
	}

	OpenTelemetryCfg.BkDataID, err = cc.Int64("logPlatform.openTelemetry.bkDataID")
	if err != nil {
		blog.Errorf("init openTelemetry failed, err: %v", err)
		return errors.New("config logPlatform.openTelemetry.bkDataID is not found")
	}
	return nil
}

