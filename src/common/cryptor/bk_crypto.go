//go:build !disable_crypto
// +build !disable_crypto

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
	"fmt"
	"strings"

	"configcenter/src/test/util"

	"github.com/TencentBlueKing/crypto-golang-sdk"
)

// cmdbBkCrypto blueking crypto
type cmdbBkCrypto struct {
	conf *Config
}

// NewBkCrypto new cmdb crypto from bk crypto
func NewBkCrypto(conf *Config) (Cryptor, error) {
	return &cmdbBkCrypto{
		conf: conf,
	}, nil
}

// Encrypt plaintext
func (c *cmdbBkCrypto) Encrypt(plaintext string) (string, error) {
	randVal := []byte(util.RandSeq(algorithmRandLen[c.conf.Algorithm]))

	crypto, err := newBkCrypto(c.conf, randVal)
	if err != nil {
		return "", err
	}

	ciphertext, err := crypto.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	// append randVal to the front of the ciphertext
	ciphertext = append(randVal, ciphertext...)
	cipherStr := string(algorithmToPrefixMap[c.conf.Algorithm]) + base64.StdEncoding.EncodeToString(ciphertext)
	return cipherStr, nil
}

// Decrypt ciphertext
func (c *cmdbBkCrypto) Decrypt(ciphertext string) (string, error) {
	splitCiphertext := strings.Split(ciphertext, "]")

	// ciphertext has no prefix, treat it as plaintext
	if len(splitCiphertext) != 2 {
		return ciphertext, nil
	}

	// get encrypt algorithm by prefix, decode by this algorithm
	prefix := AlgorithmPrefix(splitCiphertext[0] + "]")
	algorithm := prefixToAlgorithmMap[prefix]
	cryptoAlgorithm := c.conf.Algorithm
	c.conf.Algorithm = algorithm
	defer func() { c.conf.Algorithm = cryptoAlgorithm }()
	if err := c.conf.Validate(); err != nil {
		return "", fmt.Errorf("prefix crypto config failed, err: %v", err)
	}

	ciphertext = splitCiphertext[1]
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// get randVal from the front of the ciphertext
	randValLen := algorithmRandLen[algorithm]
	if len(cipherBytes) <= randValLen {
		return "", fmt.Errorf("ciphertext length(%d) <= random value length", len(ciphertext))
	}
	randVal := cipherBytes[:randValLen]

	crypto, err := newBkCrypto(c.conf, randVal)
	c.conf.Algorithm = cryptoAlgorithm
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.Decrypt(cipherBytes[randValLen:])
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func newBkCrypto(conf *Config, randVal []byte) (bkcrypto.Crypto, error) {
	var bkCrypto bkcrypto.Crypto
	var err error

	if len(randVal) != algorithmRandLen[conf.Algorithm] {
		return nil, fmt.Errorf("init %s crypto failed, random value is empty", conf.Algorithm)
	}

	switch conf.Algorithm {
	case AesGcm:
		bkCrypto, err = bkcrypto.NewAesGcm([]byte(conf.AesGcm.Key), randVal)
	default:
		return nil, fmt.Errorf("crypto algorithm %s is invalid", conf.Algorithm)
	}

	if err != nil {
		return nil, fmt.Errorf("init %s crypto failed, err: %v", conf.Algorithm, err)
	}

	return bkCrypto, nil

}
