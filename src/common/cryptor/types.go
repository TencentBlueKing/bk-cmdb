/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package cryptor

import (
	"errors"
	"fmt"
)

// Config defines cmdb crypto configuration
type Config struct {
	Enabled   bool        `mapstructure:"enabled"`
	Algorithm Algorithm   `mapstructure:"algorithm"`
	Sm4       *Sm4Conf    `mapstructure:"sm4"`
	AesGcm    *AesGcmConf `mapstructure:"aesGcm"`
}

// Validate Config
func (conf Config) Validate() error {
	if !conf.Enabled {
		return nil
	}

	switch conf.Algorithm {
	case Sm4:
		if conf.Sm4 == nil {
			return errors.New("sm4 config is not set")
		}

		if conf.Sm4.Key == "" {
			return errors.New("sm4 key is not set")
		}
	case AesGcm:
		if conf.AesGcm == nil {
			return errors.New("aes-gcm config is not set")
		}

		if conf.AesGcm.Key == "" {
			return errors.New("aes-gcm key is not set")
		}
	default:
		return fmt.Errorf("crypto algorithm %s is invalid", conf.Algorithm)
	}

	return nil
}

// Algorithm defines cryptography algorithm types
type Algorithm string

const (
	Sm4    Algorithm = "SHANGMI"
	AesGcm Algorithm = "CLASSIC"
)

var algorithmRandLen = map[Algorithm]int{
	Sm4:    16,
	AesGcm: 12,
}

// Sm4Conf defines SM4 cryptography algorithm configuration
type Sm4Conf struct {
	Key string `json:"key"`
}

// AesGcmConf defines AES-GCM cryptography algorithm configuration
type AesGcmConf struct {
	Key string `json:"key"`
}

// AlgorithmPrefix defines cryptography algorithm types
type AlgorithmPrefix string

const (
	Sm4Prefix    AlgorithmPrefix = "[Cipher:::SM4]"
	AesGcmPrefix AlgorithmPrefix = "[Cipher:::AesGcm]"
)

var algorithmToPrefixMap = map[Algorithm]AlgorithmPrefix{
	Sm4:    Sm4Prefix,
	AesGcm: AesGcmPrefix,
}

var prefixToAlgorithmMap = make(map[AlgorithmPrefix]Algorithm)

func init() {
	for algorithm, prefix := range algorithmToPrefixMap {
		prefixToAlgorithmMap[prefix] = algorithm
	}
}
