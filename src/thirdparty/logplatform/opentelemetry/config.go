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

package opentelemetry

import (
	"errors"
	"fmt"
	"time"

	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
)

var (
	openTelemetryCfg = new(OpenTelemetryConfig)
)

type OpenTelemetryConfig struct {
	// 表示是否开启日志平台openTelemetry跟踪链接入相关功能，布尔值, 默认值为false不开启
	enable bool
	// 日志平台openTelemetry跟踪链功能的自定义上报服务地址
	endpoint string
	// 日志平台openTelemetry跟踪链功能的上报data_id
	bkDataID int64
}

// InitOpenTelemetryConfig init openTelemetry config
func InitOpenTelemetryConfig() error {

	var err error
	maxCnt := 100
	cnt := 0
	for !cc.IsExist("openTelemetry") && cnt < maxCnt {
		blog.V(5).Infof("waiting openTelemetry config to be init")
		cnt++
		time.Sleep(time.Millisecond * 300)
	}

	if cnt == maxCnt {
		return errors.New("no openTelemetry config is found, the config 'openTelemetry' must exist")
	}

	openTelemetryCfg.enable, err = cc.Bool("openTelemetry.enable")
	if err != nil {
		return fmt.Errorf("config openTelemetry.enable err: %v", err)
	}

	// 如果不需要开启OpenTelemetry，那么后续没有必要再检查配置
	if !openTelemetryCfg.enable {
		return nil
	}

	openTelemetryCfg.endpoint, err = cc.String("openTelemetry.endpoint")
	if err != nil {
		return fmt.Errorf("config openTelemetry.endpoint err: %v", err)
	}

	openTelemetryCfg.bkDataID, err = cc.Int64("openTelemetry.bkDataID")
	if err != nil {
		return fmt.Errorf("config openTelemetry.bkDataID err: %v", err)
	}
	return nil
}
