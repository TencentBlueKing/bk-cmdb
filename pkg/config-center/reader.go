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

// Package cc is cmdb's config center.
package cc

import (
	"context"

	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
)

// Reader is the config reader that reads and watches config from config center.
type Reader struct {
	// discovery is the config discovery that reads config from config center.
	discovery Discovery
	// parser is the viper parser that parses the config file.
	parser *viperParser
}

// NewReader creates a new config reader.
func NewReader(discovery Discovery) *Reader {
	return &Reader{
		discovery: discovery,
		parser:    newViperParser(),
	}
}

// RunConfigRead reads config from config center and watches config change events.
func (r *Reader) RunConfigRead(ctx context.Context) error {
	for _, config := range allConfTypes {
		// initialize viper parser
		r.parser.addParser(config, viper.New())

		// read config from config center
		if err := r.ReadConfig(ctx, config); err != nil {
			return err
		}
	}

	// watch config change events from config center and triggers config update
	watchChan, err := r.discovery.Watch(ctx, configPath)
	if err != nil {
		logger.Error(ctx, "watch config change events failed", logger.E(err))
		return err
	}

	go func() {
		for event := range watchChan {
			conf := getConfigTypeByRegisterPath(event.Key)

			data := event.Data
			if event.Type == DeleteEvent {
				data = make([]byte, 0)
			}

			if err = r.parser.parseConfigData(ctx, conf, data); err != nil {
				logger.Error(ctx, "parse config change event failed", "key", conf, "event", event, logger.E(err))
				continue
			}
		}
	}()

	return nil
}

// ReadConfig reads config file data from config center.
func (r *Reader) ReadConfig(ctx context.Context, config ConfigType) error {
	// read config file
	key := getConfigRegisterPath(config)
	data, err := r.discovery.Read(ctx, key)
	if err != nil {
		logger.Error(ctx, "read config from discovery failed", "config", config, "key", key, logger.E(err))
		return err
	}

	// parse config file data
	if err = r.parser.parseConfigData(ctx, config, data); err != nil {
		logger.Error(ctx, "parse config data failed", "key", key, "data", data, logger.E(err))
		return err
	}

	return nil
}
