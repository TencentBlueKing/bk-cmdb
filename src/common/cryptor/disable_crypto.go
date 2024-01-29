//go:build disable_crypto
// +build disable_crypto

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

// disableCrypto is the special crypto type that do not do encryption and decryption
type disableCrypto struct{}

// NewBkCrypto new cmdb crypto from bk crypto
func NewBkCrypto(conf *Config) (Cryptor, error) {
	return new(disableCrypto), nil
}

// Encrypt plaintext
func (c *disableCrypto) Encrypt(plaintext string) (string, error) {
	return plaintext, nil
}

// Decrypt ciphertext
func (c *disableCrypto) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}
