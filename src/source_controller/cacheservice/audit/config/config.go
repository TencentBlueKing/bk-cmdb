/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// Package config is the audit data reporting configuration
package config

import "errors"

// Config is the audit center related configuration
type Config struct {
	Enabled   bool   `mapstructure:"enabled"`
	AppCode   string `mapstructure:"appCode"`
	AppSecret string `mapstructure:"appSecret"`
	Endpoint  string `mapstructure:"endpoint"`
	Token     string `mapstructure:"token"`
}

// Validate Config
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.AppCode == "" {
		return errors.New("appCode is not set")
	}

	if c.AppSecret == "" {
		return errors.New("appSecret is not set")
	}

	if c.Endpoint == "" {
		return errors.New("endpoint is not set")
	}

	if c.Token == "" {
		return errors.New("token is not set")
	}

	return nil
}
