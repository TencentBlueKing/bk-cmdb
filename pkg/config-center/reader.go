/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package cc

import (
	"context"
	"time"

	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// Reader is the config reader that reads and watches config from config center.
type Reader struct {
	// discovery is the config discovery that reads config from config center.
	discovery Discovery
	// viperParser is the viper parser that parses the config file.
	*viperParser
	// neededConfigs are the needed config types.
	neededConfigs []ConfigType
}

// NewReader creates a new config reader.
func NewReader(discovery Discovery, neededConfigs []ConfigType) *Reader {
	return &Reader{
		discovery:     discovery,
		viperParser:   newViperParser(),
		neededConfigs: neededConfigs,
	}
}

// RunConfigRead reads config from config center and watches config change events.
func (r *Reader) RunConfigRead(ctx context.Context) error {
	for _, config := range r.neededConfigs {
		// initialize viper parser
		r.addParser(config, viper.New())

		// read config from config center
		if err := r.ReadConfig(ctx, config); err != nil {
			return err
		}
	}

	// watch config change events from config center and triggers config update
	watchChan, err := r.discovery.Watch(ctx, configPath)
	if err != nil {
		log.Error(ctx, "watch config change events failed", log.E(err))
		return err
	}

	go func() {
		for event := range watchChan {
			conf := getConfigTypeByRegisterPath(event.Key)

			data := event.Data
			if event.Type == DeleteEvent {
				data = make([]byte, 0)
			}

			if err = r.parseConfigData(ctx, conf, data); err != nil {
				log.Error(ctx, "parse config change event failed", "key", conf, "event", event, log.E(err))
				continue
			}
		}
	}()

	return nil
}

// ReadConfig reads config file data from config center.
func (r *Reader) ReadConfig(ctx context.Context, config ConfigType) error {
	// read config from discovery, wait until config is ready
	key := getConfigRegisterPath(config)

	var data []byte
	var err error
	for range ConfigWaitTime {
		data, err = r.discovery.Read(ctx, key)
		if err != nil {
			log.Error(ctx, "read config from discovery failed", "config", config, "key", key, log.E(err))
			return err
		}

		if len(data) > 0 {
			break
		}

		log.Info(ctx, "waiting for config to be ready ...", "conf", config)
		time.Sleep(time.Second)
	}

	// parse config file data
	if err = r.parseConfigData(ctx, config, data); err != nil {
		log.Error(ctx, "parse config data failed", "key", key, "data", data, log.E(err))
		return err
	}

	return nil
}
