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

// Package cryptor defines a package that supports crypto related operations
package cryptor

import (
	"errors"
	"fmt"
)

// Cryptor 密码器
type Cryptor interface {
	// Encrypt 加密方法
	Encrypt(plaintext string) (string, error)

	// Decrypt 解密方法
	Decrypt(ciphertext string) (string, error)
}

// NewCrypto new crypto by config
func NewCrypto(conf *Config) (Cryptor, error) {
	if conf == nil {
		return nil, errors.New("crypto config is nil")
	}

	err := conf.Validate()
	if err != nil {
		return nil, fmt.Errorf("validate crypto config failed, err: %v", err)
	}

	return NewBkCrypto(conf)
}
