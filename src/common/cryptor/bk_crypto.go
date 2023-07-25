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
	"encoding/base64"
	"errors"

	"github.com/TencentBlueKing/crypto-golang-sdk"
)

// cmdbBkCrypto blue king crypto
type cmdbBkCrypto struct {
	crypto bkcrypto.Crypto
}

// NewBkCrypto new cmdb crypto from bk crypto
func NewBkCrypto(crypto bkcrypto.Crypto) (Cryptor, error) {
	if crypto == nil {
		return nil, errors.New("bk crypto is nil")
	}

	return &cmdbBkCrypto{
		crypto: crypto,
	}, nil
}

// Encrypt plaintext
func (c *cmdbBkCrypto) Encrypt(plaintext string) (string, error) {
	ciphertext, err := c.crypto.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt ciphertext
func (c *cmdbBkCrypto) Decrypt(ciphertext string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := c.crypto.Decrypt(cipherBytes)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
