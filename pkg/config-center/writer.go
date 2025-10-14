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
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
)

// Writer is the config writer that reads config from config file and writes to config center.
type Writer struct {
	// registry is the config registry that writes config to config center.
	registry Registry
	// directory is the config file directory.
	directory string
	// parser is the viper parser that parses the config file.
	parser *viperParser
}

// NewWriter creates a new config writer and starts writing config to config center.
func NewWriter(registry Registry, directory string) *Writer {
	return &Writer{
		registry:  registry,
		directory: directory,
		parser:    newViperParser(),
	}
}

// RunConfigWrite starts writing config to config center.
func (w *Writer) RunConfigWrite(ctx context.Context) error {
	for _, config := range allConfTypes {
		// initialize viper parser
		v := viper.New()
		v.SetConfigName(config)
		v.SetConfigType("yaml")
		v.AddConfigPath(w.directory)
		w.parser.addParser(config, v)

		// write config to config center
		if err := w.WriteConfig(ctx, config); err != nil {
			return err
		}

		// watch config file change and trigger config update
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			if err := w.WriteConfig(ctx, config); err != nil {
				logger.Error(ctx, "watch config file change failed", "file", config, "err", err)
				return
			}
		})
	}

	return nil
}

// WriteConfig writes config from config file to config center.
func (w *Writer) WriteConfig(ctx context.Context, fileName string) error {
	// read config file
	file, err := os.OpenInRoot(w.directory, fileName+".yaml")
	if err != nil {
		logger.Error(ctx, "open config file failed", "err", err, "file", fileName, "dir", w.directory)
		return err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			logger.Error(ctx, "close config file failed", "file", fileName, "err", err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error(ctx, "read config file failed", "file", fileName, "dir", w.directory, "err", err)
		return err
	}

	// write config to config center
	path := getConfigRegisterPath(fileName)
	if err = w.registry.Write(ctx, path, data); err != nil {
		logger.Error(ctx, "write config to config center failed", "err", err, "path", path, "data", string(data))
		return err
	}

	// parse config file data
	if err = w.parser.parseConfigData(ctx, fileName, data); err != nil {
		return err
	}

	return nil
}
