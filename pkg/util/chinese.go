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

package util

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// gb2312Delimiter gb2312 pinyin delimiters
var gb2312Delimiter = []int64{0xB0A1, 0xB0C5, 0xB2C1, 0xB4EE, 0xB6EA, 0xB7A2, 0xB8C1, 0xB9FE, 0xBBF7, 0xBFA6, 0xC0AC,
	0xC2E8, 0xC4C3, 0xC5B6, 0xC5BE, 0xC6DA, 0xC8BB, 0xC8F6, 0xCBFA, 0xCDDA, 0xCEF4, 0xD1B9, 0xD4D1, 0xD7F9}

// gb2312DelimiterInitials gb2312 pinyin delimiters corresponding initials
var gb2312DelimiterInitials = []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "O", "P", "Q",
	"R", "S", "T", "W", "X", "Y", "Z"}

// GetInitials get the initials of the input, if the input is chinese, returns the initials of its pinyin
func GetInitials(input string) string {
	initials := string([]rune(input)[:1])

	// transform the initials into gbk format
	gbkInitials, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(initials)),
		simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return initials
	}

	// get the integer format of the initials
	intInitials, err := strconv.ParseInt(hex.EncodeToString(gbkInitials), 16, 0)
	if err != nil {
		return initials
	}

	// use gb2312 delimiters to get the first letter of the pinyin since gb2312 is sorted by pinyin
	if intInitials < gb2312Delimiter[0] {
		return initials
	}

	for i := 1; i < len(gb2312Delimiter); i++ {
		if intInitials < gb2312Delimiter[i] {
			return gb2312DelimiterInitials[i-1]
		}
	}

	return initials
}
