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

package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// aesCrpytor AES密码器
type aesCrpytor struct {
	// 加解密使用的密钥
	key string
}

// NewAesEncrpytor 生成aesEncrpytor
func NewAesEncrpytor(key string) *aesCrpytor {
	return &aesCrpytor{
		key:    key,
	}
}

// Encrypt AES加密（CBC模式）
func (a *aesCrpytor) Encrypt(plainText string) (string, error) {
	plainData := []byte(plainText)
	k := []byte(a.key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		fmt.Println("err:", err)
		return "", err
	}
	blockSize := block.BlockSize()
	plainData = a.pkcs7Padding(plainData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(plainData))
	// 加密
	blockMode.CryptBlocks(cryted, plainData)

	return base64.StdEncoding.EncodeToString(cryted), nil
}

//Decrypt AES解密（CBC模式）
func (a *aesCrpytor) Decrypt(cryptedText string) (string, error) {
	cryptedByte, err := base64.StdEncoding.DecodeString(cryptedText)
	if err != nil {
		return "", err
	}
	k := []byte(a.key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	plain := make([]byte, len(cryptedByte))
	// 解密
	blockMode.CryptBlocks(plain, cryptedByte)
	plain = a.pkcs7UnPadding(plain)

	return string(plain), nil
}

// pkcs7Padding 补填充码
func (a *aesCrpytor) pkcs7Padding(data []byte, blocksize int) []byte {
	padding := blocksize - len(data)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7UnPadding 去填充码
func (a *aesCrpytor) pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
