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
	"encoding/json/v2"
	"fmt"
	"maps"
	"reflect"

	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
)

// viperParser is the config files parser that use viper to parse config.
type viperParser struct {
	// parsers are the config files related viper instances.
	parsers map[string]*viper.Viper
	// prevData are the previous config data.
	prevData map[string]map[string]any
	// eventHandlers are the config change event handlers.
	eventHandlers map[string]EventHandler
}

// newViperParser create a new viperParser.
func newViperParser() *viperParser {
	return &viperParser{
		parsers:       make(map[string]*viper.Viper),
		prevData:      make(map[string]map[string]any),
		eventHandlers: maps.Clone(eventHandlers),
	}
}

// addParser add a viper instance to viperParser.
func (p *viperParser) addParser(fileName string, v *viper.Viper) {
	p.parsers[fileName] = v
	p.prevData[fileName] = make(map[string]any)
}

// parseConfigData use viper to parse config data
func (p *viperParser) parseConfigData(ctx context.Context, fileName string, data []byte) error {
	logger.Trace(ctx, "start parsing config", "conf", fileName, "data", string(data))

	v, ok := p.parsers[fileName]
	if !ok {
		logger.Error(ctx, "viper cannot parse invalid config", "conf", fileName, "data", string(data))
		return fmt.Errorf("config %s is invalid", fileName)
	}

	// parse current config file data
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		logger.Error(ctx, "viper read config failed", "conf", fileName, "data", string(data), "err", err)
		return err
	}

	// triggers registered config change event handlers for related config keys
	prevData := p.prevData[fileName]
	curData := v.AllSettings()
	for key, handler := range p.eventHandlers {
		if !v.IsSet(key) {
			if _, exists := prevData[key]; exists {
				event := &Event{
					Type: DeleteEvent,
				}
				if err := handler(event); err != nil {
					logger.Error(ctx, "call config change handler failed", "conf", key, "err", err, "event", *event)
					return err
				}
			}
			continue
		}

		if reflect.DeepEqual(prevData[key], curData[key]) {
			continue
		}

		event := &Event{
			Type: UpsertEvent,
			Pre:  prevData[key],
			Data: curData[key],
		}
		if err := handler(event); err != nil {
			logger.Error(ctx, "call config change handler failed", "conf", key, "err", err, "event", *event)
			return err
		}
	}
	p.prevData[fileName] = curData

	return nil
}

// ConvertBasic converts config value to specified basic type value.
func ConvertBasic[T cast.Basic](data any) (T, error) {
	return cast.ToE[T](data)
}

// Convert converts config value to pointer of specified type.
func Convert[T any](data any) (*T, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %v", err)
	}

	result := new(T)
	if err = json.Unmarshal(marshal, result); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	return result, nil
}
