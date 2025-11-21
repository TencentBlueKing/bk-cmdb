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
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// RegistryWriter is the config writer that reads config from config file and writes to config center.
type RegistryWriter struct {
	// registry is the config registry that writes config to config center.
	registry Registry
	// directory is the config file directory.
	directory string
	// viperParser is the viper parser that parses the config file.
	*viperParser
	// neededConfigs are the needed config types.
	neededConfigs []ConfigType
}

// NewRegistryWriter creates a new config writer and starts writing config to config center.
func NewRegistryWriter(registry Registry, directory string, neededConfigs []ConfigType) *RegistryWriter {
	return &RegistryWriter{
		registry:      registry,
		directory:     directory,
		viperParser:   newViperParser(),
		neededConfigs: neededConfigs,
	}
}

// RunConfigWrite starts writing config to config center.
func (w *RegistryWriter) RunConfigWrite(ctx context.Context) error {
	for _, config := range w.neededConfigs {
		// initialize viper parser
		v := viper.New()
		v.AddConfigPath(w.directory)
		w.addParser(config, v)

		// write config to config center
		if err := w.WriteConfig(ctx, config); err != nil {
			return err
		}

		// watch config file change and trigger config update
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			if err := w.WriteConfig(ctx, config); err != nil {
				log.Error(ctx, "watch config file change failed", "file", config, log.E(err))
				return
			}
		})
	}

	return nil
}

// WriteConfig writes config from config file to config center.
func (w *RegistryWriter) WriteConfig(ctx context.Context, conf ConfigType) error {
	// read config file
	file, err := os.OpenInRoot(w.directory, fmt.Sprintf("%s.yaml", conf))
	if err != nil {
		log.Error(ctx, "open config file failed", log.E(err), "file", conf, "dir", w.directory)
		return err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Error(ctx, "close config file failed", "file", conf, log.E(err))
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		log.Error(ctx, "read config file failed", "file", conf, "dir", w.directory, log.E(err))
		return err
	}

	// write config to config center
	path := getConfigRegisterPath(conf)
	if err = w.registry.Write(ctx, path, data); err != nil {
		log.Error(ctx, "write config to config center failed", log.E(err), "path", path, "data", string(data))
		return err
	}

	// parse config file data
	if err = w.parseConfigData(ctx, conf, data); err != nil {
		log.Error(ctx, "parse config data failed", "conf", conf, "data", data, log.E(err))
		return err
	}

	return nil
}
