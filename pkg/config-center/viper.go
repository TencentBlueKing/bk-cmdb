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
	"bytes"
	"context"
	"fmt"
	"maps"
	"reflect"

	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// viperParser is the config files parser that use viper to parse config.
type viperParser struct {
	// parsers are the config files related viper instances.
	parsers map[ConfigType]*viper.Viper
	// prevData are the previous config data.
	prevData map[ConfigType]map[string]any
	// eventHandlers are the config change event handlers.
	eventHandlers map[string]EventHandler[any]
}

// newViperParser create a new viperParser.
func newViperParser() *viperParser {
	return &viperParser{
		parsers:       make(map[ConfigType]*viper.Viper),
		prevData:      make(map[ConfigType]map[string]any),
		eventHandlers: maps.Clone(eventHandlers),
	}
}

// addParser add a viper instance to viperParser.
func (p *viperParser) addParser(conf ConfigType, v *viper.Viper) {
	v.SetConfigName(string(conf))
	v.SetConfigType("yaml")
	p.parsers[conf] = v
	p.prevData[conf] = make(map[string]any)
}

// parseConfigData use viper to parse config data
func (p *viperParser) parseConfigData(ctx context.Context, conf ConfigType, data []byte) error {
	log.Trace(ctx, "start parsing config", "conf", conf, "data", string(data))

	v, ok := p.parsers[conf]
	if !ok {
		log.Error(ctx, "viper cannot parse invalid config", "conf", conf, "data", string(data))
		return fmt.Errorf("config %s is invalid", conf)
	}

	// parse current config file data
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		log.Error(ctx, "viper read config failed", "conf", conf, "data", string(data), log.E(err))
		return err
	}

	// triggers registered config change event handlers for related config keys
	prevData := p.prevData[conf]
	curData := v.AllSettings()
	for key, handler := range p.eventHandlers {
		if !v.IsSet(key) {
			if _, exists := prevData[key]; exists {
				event := &Event[any]{
					Type: DeleteEvent,
				}
				if err := handler(event); err != nil {
					log.Error(ctx, "call config change handler failed", "conf", key, log.E(err), "event", *event)
					return err
				}
			}
			continue
		}

		if reflect.DeepEqual(prevData[key], curData[key]) {
			continue
		}

		event := &Event[any]{
			Type: UpsertEvent,
			Pre:  prevData[key],
			Data: curData[key],
		}
		if err := handler(event); err != nil {
			log.Error(ctx, "call config change handler failed", "conf", key, log.E(err), "event", *event)
			return err
		}
	}
	p.prevData[conf] = curData

	return nil
}
